package http

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"tgm/errors"
	"tgm/users"
)

var (
	NonModifiableFieldsForNonAdmin = []string{"Username", "Scope", "LockPassword", "Perm", "Commands", "Rules"}
)

type requestUserData struct {
	What  string   `json:"what"`  // Answer to: what data type?
	Which []string `json:"which"` // Answer to: which fields?
	Data  string   `json:"data"`  // Answer to: which fields?
}

type modifyUserRequest struct {
	modifyRequest
	Data *users.User `json:"data"`
}

func getUserParameter(_ http.ResponseWriter, r *http.Request) (*requestUserData, error) {
	if r.Body == nil {
		return nil, errors.ErrEmptyRequest
	}

	req := &requestUserData{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		return nil, err
	}

	if req.What != "user" {
		return nil, errors.ErrInvalidDataType
	}

	return req, nil
}

func getUserID(r *http.Request) (uint, error) {
	vars := mux.Vars(r)
	i, err := strconv.ParseUint(vars["id"], 10, 0) //nolint:gomnd
	if err != nil {
		return 0, err
	}
	return uint(i), err
}

func getUser(_ http.ResponseWriter, r *http.Request) (*modifyUserRequest, error) {
	if r.Body == nil {
		return nil, errors.ErrEmptyRequest
	}

	req := &modifyUserRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		return nil, err
	}

	if req.What != "user" {
		return nil, errors.ErrInvalidDataType
	}

	return req, nil
}

func withSelfOrAdmin(fn handleFunc) handleFunc {
	return withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
		id, err := getUserID(r)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		if d.user.ID != id && !d.user.Perm.Admin {
			return http.StatusForbidden, nil
		}

		d.raw = id
		return fn(w, r, d)
	})
}

var usersGetHandler = withAdmin(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	req, err := getUserParameter(w, r)
	users, err := d.store.Users.Gets(d.server.Root, req.Data)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	for _, u := range users {
		u.Password = ""
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].ID < users[j].ID
	})

	return renderJSON(w, r, users)
})

var userGetHandler = withSelfOrAdmin(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	u, err := d.store.Users.Get(d.server.Root, d.raw.(uint))
	if err == errors.ErrNotExist {
		return http.StatusNotFound, err
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}

	u.Password = ""
	return renderJSON(w, r, u)
})

var userDeleteHandler = withSelfOrAdmin(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	err := d.store.Users.Delete(d.raw.(uint))
	if err != nil {
		return errToStatus(err), err
	}

	return http.StatusOK, nil
})

var userPostHandler = withAdmin(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	req, err := getUser(w, r)
	if err != nil {
		return http.StatusBadRequest, err
	}

	if len(req.Which) != 0 {
		return http.StatusBadRequest, nil
	}

	if req.Data.Password == "" {
		return http.StatusBadRequest, errors.ErrEmptyPassword
	}

	req.Data.PasswordHint = users.HintPwd(req.Data.Password)
	req.Data.Password, err = users.HashPwd(req.Data.Password)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	userHome, err := d.settings.MakeUserDir(req.Data.Username, req.Data.Scope, d.server.Root)
	if err != nil {
		log.Printf("create user: failed to mkdir user home dir: [%s]", userHome)
		return http.StatusInternalServerError, err
	}
	req.Data.Scope = userHome
	log.Printf("user: %s, home dir: [%s].", req.Data.Username, userHome)

	err = d.store.Users.Save(req.Data)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	w.Header().Set("Location", "/settings/users/"+strconv.FormatUint(uint64(req.Data.ID), 10)) //nolint:gomnd
	return http.StatusCreated, nil
})

var userPutHandler = withSelfOrAdmin(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	req, err := getUser(w, r)
	if err != nil {
		return http.StatusBadRequest, err
	}

	if req.Data.ID != d.raw.(uint) {
		return http.StatusBadRequest, nil
	}

	if len(req.Which) == 0 || (len(req.Which) == 1 && req.Which[0] == "all") {
		if !d.user.Perm.Admin {
			return http.StatusForbidden, nil
		}

		if req.Data.Password != "" {
			req.Data.PasswordHint = users.HintPwd(req.Data.Password)
			req.Data.Password, err = users.HashPwd(req.Data.Password)
		} else {
			var suser *users.User
			suser, err = d.store.Users.Get(d.server.Root, d.raw.(uint))
			req.Data.Password = suser.Password
			req.Data.PasswordHint = suser.PasswordHint
		}

		if err != nil {
			return http.StatusInternalServerError, err
		}

		req.Which = []string{}
	}

	for k, v := range req.Which {
		v = strings.Title(v)
		req.Which[k] = v

		if v == "Password" {
			if !d.user.Perm.Admin && d.user.LockPassword {
				return http.StatusForbidden, nil
			}

			req.Data.Password = users.HintPwd(req.Data.Password)
			req.Data.Password, err = users.HashPwd(req.Data.Password)
			if err != nil {
				return http.StatusInternalServerError, err
			}
		}

		for _, f := range NonModifiableFieldsForNonAdmin {
			if !d.user.Perm.Admin && v == f {
				return http.StatusForbidden, nil
			}
		}
	}

	err = d.store.Users.Update(req.Data, req.Which...)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
})

var userGetShellsHandler = withAdmin(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	shells, err := getShells()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return renderJSON(w, r, shells)
})

func getShells() (map[string]string, error) {
	var m map[string]string
	m = make(map[string]string)
	// this is for Linux/Unix machines
	file, err := os.Open("/etc/shells")
	if err != nil {
		log.Print(err)
		return m, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	re := regexp.MustCompile(`\r?\n`)

	for {
		line, err := reader.ReadString('\n')
		line = re.ReplaceAllString(line, "")

		if strings.HasPrefix(line, "/bin/bash") || strings.HasPrefix(line, "/bin/tcsh") || strings.HasPrefix(line, "/bin/zsh") {
			m[line] = line
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Print(err)
			return m, err
		}
	}
	return m, nil
}

var userGetGroupsHandler = withAdmin(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	shells, err := getGroups()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return renderJSON(w, r, shells)
})

func getGroups() (map[string]string, error) {
	var m map[string]string
	m = make(map[string]string)
	// this is for Linux/Unix machines
	file, err := os.Open("/etc/group")
	if err != nil {
		log.Print(err)
		return m, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')

		if equal := strings.Index(line, "#"); equal < 0 {
			lineSlice := strings.FieldsFunc(line, func(divide rune) bool {
				return divide == ':'
			})

			if len(lineSlice) > 0 {
				gid, err := strconv.Atoi(lineSlice[2])
				if err == nil {
					if gid >= 1000 && gid <= 65500 {
						m[lineSlice[0]] = lineSlice[0]
					}
				}
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Print(err)
			return m, err
		}
	}
	return m, nil
}

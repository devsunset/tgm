package http

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
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

type User struct {
	ID    string `json:"id"`
	Uid   string `json:"uid"`
	Gid   string `json:"gid"`
	Home  string `json:"home"`
	Shell string `json:"shell"`
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

	// DB에 저장된 정보는 Sync가 안맞는 경우 발생 할 수 있음으로 최대한 리눅스 계정과
	// Sync  처리 하기 위해 리눅스 계정 정보 조회 하여 리턴 값 대체 처리
	if strings.Compare(u.Username, "admin") != 0 {
		// /etc/passwd
		// u.Shell = ""
		users, _ := getUsers()
		for _, user := range users {
			if user.ID == u.Username {
				u.Shell = strings.TrimSpace(user.Shell)
				break
			}
		}
		// groups $USER
		// u.Group = ""
		groupsCmd := exec.Command("groups", u.Username)
		if out, err := groupsCmd.Output(); err != nil {
			log.Println("get groups error", err)
		} else {
			outStr := string(out)
			outStr = strings.TrimSpace(outStr)
			outStr = outStr[strings.Index(outStr, ":")+1:]
			slice := strings.Split(outStr, " ")
			realGroup := ""
			for _, sgroup := range slice {
				if sgroup != u.Username && sgroup != "" {
					realGroup = realGroup + strings.TrimSpace(sgroup) + ","
				}
			}
			if realGroup != "" {
				realGroup = realGroup[:len(realGroup)-1]
			}
			u.Group = realGroup
		}
		//chage -l $USER | grep "Account expires"
		// u.ExpireDay = ""
		expiresCmd := exec.Command("sh", "-c", "chage -l  "+u.Username+" | grep 'Account expires'")
		if out, err := expiresCmd.Output(); err != nil {
			log.Println("get chage info error", err)
		} else {
			outStr := string(out)
			outStr = strings.TrimSpace(outStr)

			if strings.Index(outStr, "never") > -1 {
				u.ExpireDay = "9999-12-31"
			} else {
				slice := strings.Split(outStr, " ")
				month := "01"
				if slice[2] == "Jan" {
					month = "01"
				} else if slice[2] == "Feb" {
					month = "02"
				} else if slice[2] == "Mar" {
					month = "03"
				} else if slice[2] == "Apr" {
					month = "04"
				} else if slice[2] == "May" {
					month = "05"
				} else if slice[2] == "Jun" {
					month = "06"
				} else if slice[2] == "Jul" {
					month = "07"
				} else if slice[2] == "Aug" {
					month = "08"
				} else if slice[2] == "Sep" {
					month = "09"
				} else if slice[2] == "Oct" {
					month = "10"
				} else if slice[2] == "Nov" {
					month = "11"
				} else if slice[2] == "Dec" {
					month = "12"
				}
				u.ExpireDay = slice[4] + "-" + month + "-" + slice[3][:len(slice[3])-1]
			}
		}
		// passwd -S $USER
		// u.PasswrodExpireDay = ""
		// u.PasswordExpireWarningDay = ""
		// u.LockAccount = true
		passwdCmd := exec.Command("passwd", "-S", u.Username)
		if out, err := passwdCmd.Output(); err != nil {
			log.Println("get passwd info error", err)
		} else {
			outStr := string(out)
			outStr = strings.TrimSpace(outStr)
			slice := strings.Split(outStr, " ")

			u.PasswrodExpireDay = slice[4]
			u.PasswordExpireWarningDay = slice[5]

			if strings.Compare(slice[1], "LK") == 0 {
				u.LockAccount = true
			} else {
				u.LockAccount = false
			}
		}
	}

	u.Password = ""
	return renderJSON(w, r, u)
})

var userDeleteHandler = withSelfOrAdmin(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	u, _ := d.store.Users.Get(d.server.Root, d.raw.(uint))
	err := d.store.Users.Delete(d.raw.(uint))
	if err != nil {
		return errToStatus(err), err
	}

	/*
			Usage: userdel [options] LOGIN

		Options:
		  -f, --force                   force some actions that would fail otherwise
		                                e.g. removal of user still logged in
		                                or files, even if not owned by the user
		  -h, --help                    display this help message and exit
		  -r, --remove                  remove home directory and mail spool
		  -R, --root CHROOT_DIR         directory to chroot into
		  -P, --prefix PREFIX_DIR       prefix directory where are located the /etc/* files
		  -Z, --selinux-user            remove any SELinux user mapping for the user
	*/
	//USER DEL
	userDel := exec.Command("userdel", u.Username, "-r")
	if _, err := userDel.Output(); err != nil {
		log.Println(err, "There was an error by user del", u.Username, err)
	} else {
		log.Println("user del", u.Username, "successfully")
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

	var password = req.Data.Password

	req.Data.PasswordHint = users.HintPwd(req.Data.Password)
	req.Data.Password, err = users.HashPwd(req.Data.Password)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	/*
		userHome, err := d.settings.MakeUserDir(req.Data.Username, req.Data.Scope, d.server.Root)
		if err != nil {
			log.Printf("create user: failed to mkdir user home dir: [%s]", userHome)
			return http.StatusInternalServerError, err
		}
	*/

	req.Data.Scope = strings.TrimSpace(req.Data.Scope)
	if req.Data.Scope == "" || req.Data.Scope == "./" || req.Data.Scope == "/" {
		req.Data.Scope = "/home"
	}

	if l := len(req.Data.Scope); l > 0 && strings.HasSuffix(req.Data.Scope, "/") {
		req.Data.Scope = req.Data.Scope[:l-1]
	}

	req.Data.Scope = req.Data.Scope + "/" + req.Data.Username

	//req.Data.Scope = userHome
	//log.Printf("user: %s, home dir: [%s].", req.Data.Username, userHome)

	err = d.store.Users.Save(req.Data)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	log.Println("@@@@@@@@@@@@@ =================>>> USER ADD", req.Data.Username)
	//////////////////////////////////////////////////////////////////////////////
	// Linux user connection
	// fmt.Println("=======================================")
	// fmt.Println(req.Data)
	// fmt.Println("=======================================")
	// fmt.Println(req.Data.Username)
	// fmt.Println(password)
	// fmt.Println(req.Data.Shell)
	// fmt.Println(req.Data.Group)
	// fmt.Println(req.Data.ExpireDay)
	// fmt.Println(req.Data.PasswrodExpireDay)
	// fmt.Println(req.Data.PasswordExpireWarningDay)
	// fmt.Println(req.Data.Scope)
	// fmt.Println("=======================================")

	/*
			Usage: useradd [options] LOGIN
		       useradd -D
		       useradd -D [options]

					Options:
				  -b, --base-dir BASE_DIR       base directory for the home directory of the
				                                new account
				  -c, --comment COMMENT         GECOS field of the new account
				  -d, --home-dir HOME_DIR       home directory of the new account
				  -D, --defaults                print or change default useradd configuration
				  -e, --expiredate EXPIRE_DATE  expiration date of the new account
				  -f, --inactive INACTIVE       password inactivity period of the new account
				  -g, --gid GROUP               name or ID of the primary group of the new
				                                account
				  -G, --groups GROUPS           list of supplementary groups of the new
				                                account
				  -h, --help                    display this help message and exit
				  -k, --skel SKEL_DIR           use this alternative skeleton directory
				  -K, --key KEY=VALUE           override /etc/login.defs defaults
				  -l, --no-log-init             do not add the user to the lastlog and
				                                faillog databases
				  -m, --create-home             create the user's home directory
				  -M, --no-create-home          do not create the user's home directory
				  -N, --no-user-group           do not create a group with the same name as
				                                the user
				  -o, --non-unique              allow to create users with duplicate
				                                (non-unique) UID
				  -p, --password PASSWORD       encrypted password of the new account
				  -r, --system                  create a system account
				  -R, --root CHROOT_DIR         directory to chroot into
				  -P, --prefix PREFIX_DIR       prefix directory where are located the /etc/* files
				  -s, --shell SHELL             login shell of the new account
				  -u, --uid UID                 user ID of the new account
				  -U, --user-group              create a group with the same name as the user
				  -Z, --selinux-user SEUSER     use a specific SEUSER for the SELinux user mapping
	*/
	//USER CREATE
	// -p 옵션 사용 하면 shadow 파일에 암호가 평문으로 저장됨 ㅠㅠ - 로그인 불가  옵션 사용 금지
	// 하단 옵션으로 처리 하면 useradd 명령어에서 처리 가능 하나 암호화 알고리즘이 일반 암호화 알고리즘과 다르게 적용되어 암호화 되어 저장됨
	// sudo useradd -m -p $(perl -e 'print crypt($ARGV[0], "password_value")' 'password_value') username
	// sudo useradd -p $(openssl passwd -crypt password_value) username
	// userAdd := exec.Command("useradd", req.Data.Username, "-m", "-d", req.Data.Scope, "-s", req.Data.Shell, "-e", req.Data.ExpireDay, "-p", password)
	userAdd := exec.Command("useradd", req.Data.Username, "-m", "-d", req.Data.Scope, "-s", req.Data.Shell, "-e", req.Data.ExpireDay)
	if _, err := userAdd.Output(); err != nil {
		log.Println(err, "There was an error by adding user", req.Data.Username, err)
		//d.store.Users.Delete(ID)
	} else {
		log.Println("useradd", req.Data.Username, "successfully")
	}
	/*
		Usage: usermod [options] LOGIN

		Options:
		  -c, --comment COMMENT         new value of the GECOS field
		  -d, --home HOME_DIR           new home directory for the user account
		  -e, --expiredate EXPIRE_DATE  set account expiration date to EXPIRE_DATE
		  -f, --inactive INACTIVE       set password inactive after expiration
		                                to INACTIVE
		  -g, --gid GROUP               force use GROUP as new primary group
		  -G, --groups GROUPS           new list of supplementary GROUPS
		  -a, --append                  append the user to the supplemental GROUPS
		                                mentioned by the -G option without removing
		                                the user from other groups
		  -h, --help                    display this help message and exit
		  -l, --login NEW_LOGIN         new value of the login name
		  -L, --lock                    lock the user account
		  -m, --move-home               move contents of the home directory to the
		                                new location (use only with -d)
		  -o, --non-unique              allow using duplicate (non-unique) UID
		  -p, --password PASSWORD       use encrypted password for the new password
		  -R, --root CHROOT_DIR         directory to chroot into
		  -P, --prefix PREFIX_DIR       prefix directory where are located the /etc/* files
		  -s, --shell SHELL             new login shell for the user account
		  -u, --uid UID                 new UID for the user account
		  -U, --unlock                  unlock the user account
		  -v, --add-subuids FIRST-LAST  add range of subordinate uids
		  -V, --del-subuids FIRST-LAST  remove range of subordinate uids
		  -w, --add-subgids FIRST-LAST  add range of subordinate gids
		  -W, --del-subgids FIRST-LAST  remove range of subordinate gids
		  -Z, --selinux-user SEUSER     new SELinux user mapping for the user account
	*/
	//USER MOD GROUPS
	if req.Data.Group != "" {
		strtmp := req.Data.Group
		slice := strings.Split(strtmp, ",")
		for _, sgroup := range slice {
			if sgroup != "" {
				userMod := exec.Command("usermod", req.Data.Username, "-a", "-G", sgroup)
				if _, err := userMod.Output(); err != nil {
					log.Println(err, "There was an error by user add group", sgroup, req.Data.Username, err)
				} else {
					log.Println("user add group", sgroup, req.Data.Username, "successfully")
				}
			}
		}
	}

	/*
			Usage: passwd [OPTION...] <accountName>
		  -k, --keep-tokens       keep non-expired authentication tokens
		  -d, --delete            delete the password for the named account (root only); also removes password lock if any
		  -l, --lock              lock the password for the named account (root only)
		  -u, --unlock            unlock the password for the named account (root only)
		  -e, --expire            expire the password for the named account (root only)
		  -f, --force             force operation
		  -x, --maximum=DAYS      maximum password lifetime (root only)
		  -n, --minimum=DAYS      minimum password lifetime (root only)
		  -w, --warning=DAYS      number of days warning users receives before password expiration (root only)
		  -i, --inactive=DAYS     number of days after password expiration when an account becomes disabled (root only)
		  -S, --status            report password status on the named account (root only)
		      --stdin             read new tokens from stdin (root only)

		Help options:
		  -?, --help              Show this help message
		      --usage             Display brief usage message
	*/
	//USER MOD PASSWD
	passwdMod := exec.Command("sh", "-c", "echo -n  "+password+" |  passwd "+req.Data.Username+" --stdin")
	stdoutStderr, err := passwdMod.CombinedOutput()
	if err != nil {
		log.Println(err, "There was an error by user mod passwd", req.Data.Username, err)
	} else {
		log.Println("passwd mod", req.Data.Username, "successfully")
		log.Println("passwd mod", string(stdoutStderr))
	}

	passwdPeriodMod := exec.Command("passwd", req.Data.Username, "-x", req.Data.PasswrodExpireDay, "-w", req.Data.PasswordExpireWarningDay)
	if _, err := passwdPeriodMod.Output(); err != nil {
		log.Println(err, "There was an error by user mod passwd period", req.Data.Username, err)
	} else {
		log.Println("passwd period mod", req.Data.Username, "successfully")
	}
	//////////////////////////////////////////////////////////////////////////////

	w.Header().Set("Location", "/settings/users/"+strconv.FormatUint(uint64(req.Data.ID), 10)) //nolint:gomnd
	return http.StatusCreated, nil
})

var userPutHandler = withSelfOrAdmin(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	var password = ""
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
			password = req.Data.Password
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
	} // if len(req.Which) == 0 || (len(req.Which) == 1 && req.Which[0] == "all") {

	for k, v := range req.Which {
		v = strings.Title(v)
		req.Which[k] = v

		if v == "Password" {
			if !d.user.Perm.Admin && d.user.LockPassword {
				return http.StatusForbidden, nil
			}
			password = req.Data.Password
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
	} //for k, v := range req.Which {

	err = d.store.Users.Update(req.Data, req.Which...)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	// log.Println("@@@@@@@@@@@@@ =================>>>", req.Data)
	log.Println("@@@@@@@@@@@@@ =================>>> USER MODIFY", req.Data.Username)

	if req.Data.Username != "admin" {
		var suser *users.User
		suser, err = d.store.Users.Get(d.server.Root, d.raw.(uint))

		// PASSWORD CHANGE
		if password != "" {
			passwdMod := exec.Command("sh", "-c", "echo -n  "+password+" |  passwd "+suser.Username+" --stdin")
			stdoutStderr, err := passwdMod.CombinedOutput()
			if err != nil {
				log.Println(err, "There was an error by user mod passwd", suser.Username, err)
			} else {
				log.Println("passwd mod", suser.Username, "successfully")
				log.Println("passwd mod", string(stdoutStderr))
			}
		}

		if req.Data.Username != "" {
			// Shell
			// 계정 유효 일자
			// Group
			userMod := exec.Command("usermod", suser.Username, "-s", suser.Shell, "-e", suser.ExpireDay)
			if _, err := userMod.Output(); err != nil {
				log.Println(err, "There was an error by user mod shell, expireday", suser.Username, err)
			} else {
				log.Println("user mod  shell, expireday", suser.Username, "successfully")
			}

			groupsCmd := exec.Command("groups", suser.Username)
			if out, err := groupsCmd.Output(); err != nil {
				log.Println("get groups error", err)
			} else {
				outStr := string(out)
				outStr = strings.TrimSpace(outStr)
				outStr = outStr[strings.Index(outStr, ":")+1:]
				slice := strings.Split(outStr, " ")
				for _, sgroup := range slice {
					if sgroup != suser.Username && sgroup != "" {
						gpasswdDel := exec.Command("gpasswd", "-d", suser.Username, sgroup)
						if _, err := gpasswdDel.Output(); err != nil {
							log.Println(err, "There was an error by user delete group", sgroup, err)
						} else {
							log.Println("user delete group", sgroup, "successfully")
						}
					}
				}

				if suser.Group != "" {
					slice := strings.Split(suser.Group, ",")
					for _, sgroup := range slice {
						if sgroup != "" {
							gpasswdAdd := exec.Command("gpasswd", "-a", suser.Username, sgroup)
							if _, err := gpasswdAdd.Output(); err != nil {
								log.Println(err, "There was an error by user add group", sgroup, err)
							} else {
								log.Println("user add group", sgroup, "successfully")
							}
						}
					}
				}
			}

			// 암호 기간 만료일
			// 암호 기간 만료 경고일
			passwdPeriodMod := exec.Command("passwd", suser.Username, "-x", suser.PasswrodExpireDay, "-w", suser.PasswordExpireWarningDay)
			if _, err := passwdPeriodMod.Output(); err != nil {
				log.Println(err, "There was an error by user mod passwd period", suser.Username, err)
			} else {
				log.Println("passwd period mod", suser.Username, "successfully")
			}

			// 계정 잠금
			if suser.LockAccount {
				userLock := exec.Command("usermod", suser.Username, "-L")
				if _, err := userLock.Output(); err != nil {
					log.Println(err, "There was an error by user lock", suser.Username, err)
				} else {
					log.Println("user lock", suser.Username, "successfully")
				}
			} else {
				userLock := exec.Command("usermod", suser.Username, "-U")
				if _, err := userLock.Output(); err != nil {
					log.Println(err, "There was an error by user lock", suser.Username, err)
				} else {
					log.Println("user lock", suser.Username, "successfully")
				}
			}
		}
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
	users, _ := getUsers()
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
						if checkGroup(lineSlice[0], users) {
							m[lineSlice[0]] = lineSlice[0]
						}
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

func getUsers() ([]User, error) {
	var LinuxUsers [][]string
	users := []User{}

	// this is for Linux/Unix machines
	file, err := os.Open("/etc/passwd")
	if err != nil {
		log.Print(err)
		return users, err
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
				uid, err := strconv.Atoi(lineSlice[2])
				if err == nil {
					if uid >= 1000 && uid <= 65500 {
						LinuxUsers = append(LinuxUsers, lineSlice)
					}
				}
			}
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Print(err)
			return users, err
		}
	}

	for _, data := range LinuxUsers {
		user := User{}
		user.ID = data[0]
		user.Uid = data[2]
		user.Gid = data[3]
		if len(data) == 6 {
			user.Home = data[4]
			user.Shell = data[5]
		} else if len(data) == 7 {
			user.Home = data[5]
			user.Shell = data[6]
		} else {
			user.Home = ""
			user.Shell = ""
		}

		users = append(users, user)
	}

	return users, nil
}

func checkGroup(group string, users []User) bool {
	for _, user := range users {
		if user.ID == group {
			return false
		}
	}
	return true
}

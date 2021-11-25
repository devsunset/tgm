package users

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"tgm/errors"
)

// StorageBackend is the interface to implement for a users storage.
type StorageBackend interface {
	GetBy(interface{}) (*User, error)
	Gets() ([]*User, error)
	Save(u *User) error
	Update(u *User, fields ...string) error
	DeleteByID(uint) error
	DeleteByUsername(string) error
}

type Store interface {
	Get(baseScope string, id interface{}) (user *User, err error)
	Gets(baseScope string, userId string) ([]*User, error)
	Update(user *User, fields ...string) error
	Save(user *User) error
	Delete(id interface{}) error
	LastUpdate(id uint) int64
}

// Storage is a users storage.
type Storage struct {
	back    StorageBackend
	updated map[uint]int64
	mux     sync.RWMutex
}

// Group describes a Group.
type Group struct {
	ID      string `json:"id"`
	Gid     string `json:"gid"`
	Members string `json:"members"`
	Primary string `json:"primary"`
}
type Account struct {
	ID    string `json:"id"`
	Uid   string `json:"uid"`
	Gid   string `json:"gid"`
	Home  string `json:"home"`
	Shell string `json:"shell"`
}

// NewStorage creates a users storage from a backend.
func NewStorage(back StorageBackend) *Storage {
	return &Storage{
		back:    back,
		updated: map[uint]int64{},
	}
}

// Get allows you to get a user by its name or username. The provided
// id must be a string for username lookup or a uint for id lookup. If id
// is neither, a ErrInvalidDataType will be returned.
func (s *Storage) Get(baseScope string, id interface{}) (user *User, err error) {
	user, err = s.back.GetBy(id)
	if err != nil {
		return
	}
	if err := user.Clean(baseScope); err != nil {
		return nil, err
	}

	if strings.Compare(user.Username, "admin") != 0 {
		fmt.Println(">>>==================== : user : ", user.Username)
		// DB  저장된 값이 아닌 LINUX 시스템에서 값을 불러와야 할지  고민중 ....
		// user.Shell = "/bin/zsh"
		// user.Group = "testGroup"
		// user.ExpireDay = "2021-12-24"
		// user.PasswrodExpireDay = "90"
		// user.PasswordExpireWarningDay = "7"
		// user.LockAccount = true
	}

	return
}

func getAccounts() ([]Account, error) {
	var LinuxUsers [][]string
	accounts := []Account{}

	// this is for Linux/Unix machines
	file, err := os.Open("/etc/passwd")
	if err != nil {
		log.Print(err)
		return accounts, err
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
			return accounts, err
		}
	}

	for _, data := range LinuxUsers {
		account := Account{}
		account.ID = data[0]
		account.Uid = data[2]
		account.Gid = data[3]
		if len(data) == 6 {
			account.Home = data[4]
			account.Shell = data[5]
		} else if len(data) == 7 {
			account.Home = data[5]
			account.Shell = data[6]
		} else {
			account.Home = ""
			account.Shell = ""
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

func getGroups() ([]Group, error) {
	accounts, _ := getAccounts()

	var LinuxGroups [][]string
	groups := []Group{}

	// this is for Linux/Unix machines
	file, err := os.Open("/etc/group")
	if err != nil {
		log.Print(err)
		return groups, err
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
						LinuxGroups = append(LinuxGroups, lineSlice)
					}
				}
			}
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Print(err)
			return groups, err
		}
	}

	for _, data := range LinuxGroups {
		group := Group{}
		group.ID = data[0]
		group.Gid = data[2]
		if checkPrimary(data[2], accounts) {
			group.Primary = "P"
			if len(data[3]) == 1 {
				group.Members = "User Primary Group"
			} else {
				group.Members = data[3]
			}
		} else {
			group.Primary = ""
			group.Members = data[3]
		}
		groups = append(groups, group)
	}

	return groups, nil
}

func checkPrimary(gid string, accounts []Account) bool {
	for _, account := range accounts {
		if account.Gid == gid {
			return true
		}
	}
	return false
}

func getUserShell(userId string, accounts []Account) string {
	for _, account := range accounts {
		if account.ID == userId {
			return account.Shell
		}
	}
	return ""
}

func getUserGroup(userId string, accounts []Account, groups []Group) string {
	groupid := ""
	gid := ""
	for _, account := range accounts {
		if account.ID == userId {
			gid = account.Gid
		}
	}

	for _, group := range groups {
		if group.Gid == gid {
			groupid = group.ID
		}
	}

	for _, group := range groups {
		members := group.Members
		slice := strings.Split(members, ",")
		for _, str := range slice {
			str = strings.Trim(str, " ")
			str = strings.Trim(str, " \n")
			if str == userId {
				if groupid == "" {
					groupid = group.ID
				} else {
					groupid += "," + group.ID
				}
			}
		}
	}
	return groupid
}

func getUserPasswdStatus(userId string) string {
	argGroup := []string{"-S", userId}
	cmd := exec.Command("passwd", argGroup...)
	if out, err := cmd.Output(); err != nil {
		//log.Println(err, "There was an error by account status check", userId)
		return ""
	} else {
		//log.Println(string(out))
		return string(out)
	}
}

// Gets gets a list of all users.
func (s *Storage) Gets(baseScope string, userId string) ([]*User, error) {
	accounts, _ := getAccounts()
	groups, _ := getGroups()
	users, err := s.back.Gets()
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if err := user.Clean(baseScope); err != nil { //nolint:govet
			return nil, err
		}
		user.Shell = getUserShell(user.Username, accounts)
		user.Group = getUserGroup(user.Username, accounts, groups)
		passwdStatus := getUserPasswdStatus(user.Username)
		if passwdStatus == "" {
			user.Lock = ""
		} else {
			slice := strings.Split(passwdStatus, " ")
			if len(slice) > 2 {
				user.Lock = slice[1]
			} else {
				user.Lock = ""
			}
		}
	}

	userId = strings.Trim(userId, " ")
	if userId != "" {
		searchUsers := []*User{}

		for _, user := range users {
			if strings.Contains(user.Username, userId) {
				searchUsers = append(searchUsers, user)
			}
		}
		return searchUsers, nil
	}

	return users, err
}

// Update updates a user in the database.
func (s *Storage) Update(user *User, fields ...string) error {
	err := user.Clean("", fields...)
	if err != nil {
		return err
	}

	err = s.back.Update(user, fields...)
	if err != nil {
		return err
	}

	s.mux.Lock()
	s.updated[user.ID] = time.Now().Unix()
	s.mux.Unlock()
	return nil
}

// Save saves the user in a storage.
func (s *Storage) Save(user *User) error {
	if err := user.Clean(""); err != nil {
		return err
	}

	return s.back.Save(user)
}

// Delete allows you to delete a user by its name or username. The provided
// id must be a string for username lookup or a uint for id lookup. If id
// is neither, a ErrInvalidDataType will be returned.
func (s *Storage) Delete(id interface{}) error {
	switch id := id.(type) {
	case string:
		user, err := s.back.GetBy(id)
		if err != nil {
			return err
		}
		if user.ID == 1 {
			return errors.ErrRootUserDeletion
		}
		return s.back.DeleteByUsername(id)
	case uint:
		if id == 1 {
			return errors.ErrRootUserDeletion
		}
		return s.back.DeleteByID(id)
	default:
		return errors.ErrInvalidDataType
	}
}

// LastUpdate gets the timestamp for the last update of an user.
func (s *Storage) LastUpdate(id uint) int64 {
	s.mux.RLock()
	defer s.mux.RUnlock()
	if val, ok := s.updated[id]; ok {
		return val
	}
	return 0
}

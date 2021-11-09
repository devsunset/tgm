package groups

import (
	"bufio"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"tgm/errors"
)

// Group describes a Group.
type Group struct {
	ID      string `json:"id"`
	Gid     string `json:"gid"`
	Members string `json:"members"`
	Primary string `json:"primary"`
}
type User struct {
	ID    string `json:"id"`
	Uid   string `json:"uid"`
	Gid   string `json:"gid"`
	Home  string `json:"home"`
	Shell string `json:"shell"`
}

type Store interface {
	Gets() ([]Group, error)
	Save(ID string) error
	Delete(ID string) error
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

		if len(data) == 5 {
			user.Home = data[4]
			user.Shell = data[5]
		} else if len(data) == 6 {
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

func checkPrimary(gid string, users []User) bool {
	for _, user := range users {
		if user.Gid == gid {
			return true
		}
	}
	return false
}

func getGroups() ([]Group, error) {
	users, _ := getUsers()

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
		group.Members = data[3]
		if checkPrimary(data[2], users) {
			group.Primary = "P"
		} else {
			group.Primary = ""
		}
		groups = append(groups, group)
	}

	return groups, nil
}

func Gets() ([]Group, error) {
	return getGroups()
}

func Save(ID string) error {
	linuxGroups, err := getGroups()

	if err != nil {
		return err
	}

	for _, data := range linuxGroups {
		if data.ID == ID {
			log.Print("Group already exists")
			return errors.ErrExistsGroupID
		}
	}

	argGroup := []string{ID}

	groupCmd := exec.Command("groupadd", argGroup...)

	if out, err := groupCmd.Output(); err != nil {
		log.Println(err, "There was an error by adding group", ID)
		return errors.ErrCreateGroupID
	} else {
		log.Println(string(out))
	}

	return nil
}

func Delete(ID string) error {
	argGroup := []string{ID}

	groupCmd := exec.Command("groupdel", argGroup...)

	if out, err := groupCmd.Output(); err != nil {
		log.Println(err, "There was an error by delete group", ID)
		return errors.ErrDeleteGroupID
	} else {
		log.Println(string(out))
	}

	return nil
}

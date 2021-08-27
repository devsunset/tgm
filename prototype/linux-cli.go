package main

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

const (
	userFile string = "/etc/passwd"
)

type Users struct {
	Users []User `json:"users"`
}

type User struct {
	Name      string `json:"name"`
	Directory string `json:"directory"`
	Group     string `json:group`
	Shell     string `json:shell`
}

// Read json file and return slice of byte.
func ReadUsers(f string) []byte {
	jsonFile, err := os.Open(f)

	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	data, _ := ioutil.ReadAll(jsonFile)
	return data
}

// Read file /etc/passwd and return slice of users
func ReadEtcPasswd(f string) (list []string) {
	file, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	r := bufio.NewScanner(file)

	for r.Scan() {
		lines := r.Text()
		parts := strings.Split(lines, ":")
		list = append(list, parts[0])
	}
	return list
}

// Check if user on the host
func check(s []string, u string) bool {
	for _, w := range s {
		if u == w {
			return true
		}
	}
	return false
}

// Return securely generated random bytes
func CreateRandom(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println(err)
		//os.Exit(1)
	}
	return string(b)
}

// User is created by executing shell command useradd
func AddNewUser(u *User) (bool, string) {
	encrypt := base64.StdEncoding.EncodeToString([]byte(CreateRandom(9)))

	argUser := []string{"-m", "-d", u.Directory, "-G", u.Group, "-s", u.Shell, u.Name}
	argPass := []string{"-c", fmt.Sprintf("echo %s:%s | chpasswd", u.Name, encrypt)}

	userCmd := exec.Command("useradd", argUser...)
	passCmd := exec.Command("/bin/sh", argPass...)

	if out, err := userCmd.Output(); err != nil {
		fmt.Println(err, "There was an error by adding user", u.Name)
		return false, ""
	} else {
		fmt.Printf("Output: %s\n", out)

		if _, err := passCmd.Output(); err != nil {
			fmt.Println(err)
			return false, ""
		}
		return true, encrypt
	}
}

func groupadd() {
	fmt.Println("groupadd  groupadd.json execute ....")
}

func groupdel() {
	fmt.Println("groupdel groupdel.json execute ....")
}

func useradd() {
	fmt.Println("useradd useradd.json execute ....")

	NameOfFile := os.Args[1]
	data := ReadUsers(NameOfFile)

	var u Users
	json.Unmarshal(data, &u)

	userList := ReadEtcPasswd(userFile)

	for i := range u.Users {
		c := check(userList, u.Users[i].Name)
		if c == false {
			if info, passwd := AddNewUser(&u.Users[i]); info == true {
				fmt.Println("User was added:>", u.Users[i].Name, "=>", "Password:>", passwd)
			}
		} else {
			fmt.Println("The user already exists:>", u.Users[i].Name)
		}
	}
}

func userdel() {
	fmt.Println("userdel userdel.json execute ....")
}

func main() {

	encrypt := base64.StdEncoding.EncodeToString([]byte(CreateRandom(3)))
	fmt.Println(encrypt)

	if len(os.Args) == 1 {
		fmt.Println("====================================")
		fmt.Println("### Usage ###")
		fmt.Println("")
		fmt.Println("go run linux-cli.go <<command>>")
		fmt.Println("ex) go run linux-cli.go useradd")
		fmt.Println("")

		fmt.Println("====================================")
		fmt.Println("--- command list -------------------")
		fmt.Println("====================================")
		fmt.Println("groupadd")
		fmt.Println("groupdel")
		fmt.Println("useradd")
		fmt.Println("userdel")
		fmt.Println("====================================")
		os.Exit(1)
	}

	if os.Args[1] == "groupadd" {
		groupadd()
	} else if os.Args[1] == "groupdel" {
		groupdel()
	} else if os.Args[1] == "useradd" {
		useradd()
	} else if os.Args[1] == "userdel" {
		userdel()
	} else {
		fmt.Println("Invalid command args")
		os.Exit(1)
	}
}

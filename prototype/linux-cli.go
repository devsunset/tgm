package main

import (
	"bufio"
	"crypto/rand"
	// "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"os/exec"
	"strings"
)

////////////////////////////////////////
const (
	groupFile string = "/etc/group"
)

const (
	userFile string = "/etc/passwd"
)

type Groups struct {
	Groups []Group `json:"groups"`
}

type Users struct {
	Users []User `json:"users"`
}

type Group struct {
	Group     string `json:group`
}

type User struct {
	Name      string `json:"name"`
	Directory string `json:"directory"`
	Group     string `json:group`
	Shell     string `json:shell`
}
////////////////////////////////////////

// Read json file and return slice of byte.
func ReadJsonFile(f string) []byte {
	jsonFile, err := os.Open(f)

	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	data, _ := ioutil.ReadAll(jsonFile)
	return data
}

// Read file /etc/group and return slice of groups
func ReadEtcGroup(f string) (list []string) {
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
func checkAccount(s []string, u string) bool {
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

////////////////////////////////////////

// Group is created by executing shell command groupadd
func AddGroup(g *Group) (bool) {
	argGroup := []string{g.Group}
	
	groupCmd := exec.Command("groupadd", argGroup...)

	fmt.Println(groupCmd)

	if out, err := groupCmd.Output(); err != nil {		
		fmt.Println(err, "There was an error by adding group", g.Group)
		return false
	} else {
		fmt.Printf("Output: %s\n", out)
		return true
	}
}

// Group is deleted by executing shell command groupdel
func DeleteGroup(g *Group) (bool) {
	argGroup := []string{g.Group}
	
	groupCmd := exec.Command("groupdel", argGroup...)

	fmt.Println(groupCmd)

	if out, err := groupCmd.Output(); err != nil {		
		fmt.Println(err, "There was an error by deleting group", g.Group)
		return false
	} else {
		fmt.Printf("Output: %s\n", out)
		return true
	}
}

// User is created by executing shell command useradd
func AddNewUser(u *User) (bool, string) {	
	// encrypt := base64.StdEncoding.EncodeToString([]byte(CreateRandom(9)))
	encrypt := "test1qazxsw2"

	argUser := []string{"-m", "-d", u.Directory, "-G", u.Group, "-s", u.Shell, u.Name}
	argPass := []string{"-c", fmt.Sprintf("echo %s:%s | chpasswd", u.Name, encrypt)}

	userCmd := exec.Command("useradd", argUser...)
	passCmd := exec.Command("/bin/sh", argPass...)

	fmt.Println(userCmd)
	fmt.Println(passCmd)

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

// User is deleted by executing shell command userdel
func DeleteUser(u *User) (bool) {	

	argUser := []string{u.Name}
	// argUser := []string{"-r", u.Name}


	userCmd := exec.Command("userdel", argUser...)

	fmt.Println(userCmd)

	if out, err := userCmd.Output(); err != nil {		
		fmt.Println(err, "There was an error by deleting user", u.Name)
		return false
	} else {
		fmt.Printf("Output: %s\n", out)
		return true
	}
}

////////////////////////////////////////

func groupadd() {
	fmt.Println("groupadd jsonfile/group.json execute ....")

	NameOfFile := "jsonfile/group.json"
	data := ReadJsonFile(NameOfFile)

	var g Groups
	json.Unmarshal(data, &g)

	groupList := ReadEtcGroup(groupFile)

	for i := range g.Groups {
		c := checkAccount(groupList, g.Groups[i].Group)
		if c == false {
			if info := AddGroup(&g.Groups[i]); info == true {
				fmt.Println("Group was added:>", g.Groups[i].Group)
			}
		} else {
			fmt.Println("The group already exists:>", g.Groups[i].Group)
		}
	}
}

func groupdel() {
	fmt.Println("groupdel jsonfile/group.json execute ....")

	NameOfFile := "jsonfile/group.json"
	data := ReadJsonFile(NameOfFile)

	var g Groups
	json.Unmarshal(data, &g)

	groupList := ReadEtcGroup(groupFile)

	for i := range g.Groups {
		c := checkAccount(groupList, g.Groups[i].Group)
		if c == false {
			fmt.Println("The group is not exists:>", g.Groups[i].Group)
		} else {
			
			if info := DeleteGroup(&g.Groups[i]); info == true {
				fmt.Println("Group was deleted:>", g.Groups[i].Group)
			}
		}
	}
}

func getgroups() {
	fmt.Println("getgroups execute ....")

	var LinuxGroups []string

	// this is for Linux/Unix machines
	file, err := os.Open(groupFile)

	if err != nil {
			fmt.Println(err)
			os.Exit(1)
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
						LinuxGroups = append(LinuxGroups, lineSlice[0])
					}
			}

			if err == io.EOF {
					break
			}
			if err != nil {
					fmt.Println(err)
					os.Exit(1)
			}

	}

	for _, name := range LinuxGroups {			
			usr, err := user.LookupGroup(name)
			if err != nil {
					panic(err)
			}
			// see https://pkg.go.dev/os/user#Group
			fmt.Printf("NAME:%s\n", usr.Name)
			fmt.Printf("GID:%s\n", usr.Gid)
			fmt.Println("*********************************")
	}
}

////////////////////////////////////////

func useradd() {
	fmt.Println("useradd user.json execute ....")

	NameOfFile := "jsonfile/user.json"
	data := ReadJsonFile(NameOfFile)

	var u Users
	json.Unmarshal(data, &u)

	userList := ReadEtcPasswd(userFile)

	for i := range u.Users {
		c := checkAccount(userList, u.Users[i].Name)
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
	fmt.Println("userdel user.json execute ....")

	NameOfFile := "jsonfile/user.json"
	data := ReadJsonFile(NameOfFile)

	var u Users
	json.Unmarshal(data, &u)

	userList := ReadEtcPasswd(userFile)

	for i := range u.Users {
		c := checkAccount(userList, u.Users[i].Name)
		if c == false {
			fmt.Println("The user not exists:>", u.Users[i].Name)
		} else {
			
			if info := DeleteUser(&u.Users[i]); info == true {
				fmt.Println("User was deleted:>", u.Users[i].Name)
			}
		}
	}
}

func getusers() {
	fmt.Println("getusers execute ....")

	var LinuxUsers []string

	// this is for Linux/Unix machines
	file, err := os.Open(userFile)

	if err != nil {
			fmt.Println(err)
			os.Exit(1)
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
							LinuxUsers = append(LinuxUsers, lineSlice[0])
					}

			}

			if err == io.EOF {
					break
			}
			if err != nil {
					fmt.Println(err)
					os.Exit(1)
			}

	}

	for _, name := range LinuxUsers {			
			usr, err := user.Lookup(name)
			if err != nil {
					panic(err)
			}
			// see https://golang.org/pkg/os/user/#User			
			fmt.Printf("username:%s\n", usr.Username)
			fmt.Printf("homedir:%s\n", usr.HomeDir)
			fmt.Printf("UID:%s\n", usr.Uid)
			fmt.Printf("GID:%s\n", usr.Gid)
			fmt.Printf("DisplayName:%s\n", usr.Name)
			fmt.Println("*********************************")
	}
}

////////////////////////////////////////

func main() {
	if len(os.Args) == 1 {
		fmt.Println("====================================")
		fmt.Println("### Usage ###")
		fmt.Println("")
		fmt.Println("    go run linux-cli.go <<command>>")
		fmt.Println("ex) go run linux-cli.go useradd")
		fmt.Println("")

		fmt.Println("====================================")
		fmt.Println("--- command list -------------------")
		fmt.Println("====================================")
		fmt.Println("groupadd")
		fmt.Println("groupdel")
		fmt.Println("getgroups")
		fmt.Println("====================================")
		fmt.Println("useradd")
		fmt.Println("userdel")
		fmt.Println("getusers")
		fmt.Println("====================================")
		os.Exit(1)
	}

	if os.Args[1] == "groupadd" {
		groupadd()
	} else if os.Args[1] == "groupdel" {
		groupdel()
	} else if os.Args[1] == "getgroups" {
		getgroups()
	} else if os.Args[1] == "useradd" {
		useradd()
	} else if os.Args[1] == "userdel" {
		userdel()
	} else if os.Args[1] == "getusers" {
		getusers()
	} else {
		fmt.Println("Invalid command args")
		os.Exit(1)
	}
}

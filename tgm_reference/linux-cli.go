package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

////////////////////////////////////////
const (
	groupFile string = "/etc/group"
	userFile  string = "/etc/passwd"
	shellFile string = "/etc/shells"
)

type Groups struct {
	Groups []Group `json:"groups"`
}

type Users struct {
	Users []User `json:"users"`
}

type Group struct {
	Group string `json:group`
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

////////////////////////////////////////
// Group is created by executing shell command groupadd
func AddGroup(g *Group) bool {
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
func DeleteGroup(g *Group) bool {
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
	encrypt := "test1qazxsw2"

	argUser := []string{"-m", "-d", u.Directory, "-G", u.Group, "-s", u.Shell, u.Name, "-e", "2021-12-31"}
	argPass := []string{"-c", fmt.Sprintf("echo %s:%s | chpasswd", u.Name, encrypt)}
	argChage := []string{"-M", "90", "-W", "9", u.Name}

	// Password expires										: never
	// Password inactive									: never				chage  -I
	// Minimum number of days between password change		: 0					chage  -m
	// Maximum number of days between password change		: 99999				chage  -M
	// Number of days of warning before password expires	: 7   				chage  -W

	// ⚡ root@localhost  /workspace/devwork/tgm/prototype   main ±  chage -l -i tgm1
	// Last password change									: 2021-09-08
	// Password expires										: 2021-12-07
	// Password inactive									: never
	// Account expires										: 2021-12-31
	// Minimum number of days between password change		: 0
	// Maximum number of days between password change		: 90
	// Number of days of warning before password expires	: 9

	userCmd := exec.Command("useradd", argUser...)
	passCmd := exec.Command("/bin/sh", argPass...)
	chageCmd := exec.Command("chage", argChage...)

	fmt.Println(userCmd)
	fmt.Println(passCmd)
	fmt.Println(chageCmd)

	if out, err := userCmd.Output(); err != nil {
		fmt.Println(err, "There was an error by adding user", u.Name)
		return false, ""
	} else {
		fmt.Printf("Output: %s\n", out)
		if _, err := passCmd.Output(); err != nil {

			fmt.Println(err)
			return false, ""
		} else {
			if _, err := chageCmd.Output(); err != nil {

				fmt.Println(err)
				return false, ""
			}
		}
		return true, encrypt
	}
}

// User is deleted by executing shell command userdel
func DeleteUser(u *User) bool {

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
	fmt.Println("groupadd jsonfile/group.json execute ...\n")

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
	fmt.Println("groupdel jsonfile/group.json execute ...\n")

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
	fmt.Println("getgroups execute ...\n")

	var LinuxGroups [][]string

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
			fmt.Println(err)
			os.Exit(1)
		}
	}

	for _, groups := range LinuxGroups {
		if err != nil {
			panic(err)
		}
		fmt.Printf("Name:%s\n", groups[0])
		fmt.Printf("Gid:%s\n", groups[2])
		fmt.Printf("Members:%s\n", groups[3])
		fmt.Println("*********************************")
	}
}

////////////////////////////////////////
func useradd() {
	fmt.Println("useradd user.json execute ...\n")

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
	fmt.Println("userdel user.json execute ...\n")

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
	fmt.Println("getusers execute ...\n")

	var LinuxUsers [][]string

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
			fmt.Println(err)
			os.Exit(1)
		}

	}

	for _, users := range LinuxUsers {
		if err != nil {
			panic(err)
		}
		fmt.Printf("UserName:%s\n", users[0])
		fmt.Printf("Uid:%s\n", users[2])
		fmt.Printf("Gid:%s\n", users[3])
		fmt.Printf("Home Dir:%s\n", users[4])
		fmt.Printf("Shell:%s\n", users[5])
		fmt.Println("*********************************")
	}
}

func getshells() {
	fmt.Println("getshells execute ...\n")

	var Shelles []string

	file, err := os.Open(shellFile)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')

		if err == io.EOF {
			break
		}

		if equal := strings.Index(line, "#"); equal < 0 {
			Shelles = append(Shelles, line)
		}

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	for _, shell := range Shelles {
		if err != nil {
			panic(err)
		}
		fmt.Printf("SHELL:%s", shell)
	}
}

func getaccountinfo() {
	fmt.Println("getaccountinfo execute ...\n")

	argsChage := []string{"-l", "-i", "tgm1"}

	cmdChage := exec.Command("chage", argsChage...)

	fmt.Println(cmdChage)

	if outChage, errChage := cmdChage.Output(); errChage != nil {
		fmt.Printf("err: %s\n", errChage)
	} else {
		fmt.Printf("Output: %s\n", outChage)
	}

	argsPasswd := []string{"-S", "tgm1"}

	cmdPasswd := exec.Command("passwd", argsPasswd...)

	fmt.Println(cmdPasswd)

	if outPasswd, errPasswd := cmdPasswd.Output(); errPasswd != nil {
		fmt.Printf("err: %s\n", errPasswd)
	} else {
		fmt.Printf("Output: %s\n", outPasswd)
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
		fmt.Println("---------- command list ------------")
		fmt.Println("====================================")
		fmt.Println("groupadd")
		fmt.Println("groupdel")
		fmt.Println("getgroups")
		fmt.Println("====================================")
		fmt.Println("useradd")
		fmt.Println("userdel")
		fmt.Println("getusers")
		fmt.Println("====================================")
		fmt.Println("getshells")
		fmt.Println("getaccountinfo")
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
	} else if os.Args[1] == "getshells" {
		getshells()
	} else if os.Args[1] == "getaccountinfo" {
		getaccountinfo()
	} else {
		fmt.Println("Invalid command args")
		os.Exit(1)
	}

	/*
		/etc/passwd : 사용자 정보가 담긴 파일
		root:x:0:0:root:/root:/bin/bash
		# 사용자명:패스워드:UID:GID:사용자정보:홈디렉토리:쉘

		> useradd [options] 사용자명
		# options
		# -c [텍스트] : 사용자정보
		# -m : 홈디렉토리 생성
		# -M : 홈디렉토리 생성 안함
		# -d [폴더] : 홈디렉토리 지정
		# -N : 사용자 개인 그룹 생성하지 않음. default : 생성
		# -u [UID] : UID 직접 지정
		# -g [GID] : GID 직접 지정
		# -s [Shell] : shell 지정

		> usermod [options] 사용자명
		# options
		# -c [텍스트] : 사용자정보 수정
		# -d [폴더] : 홈디렉토리 변경
		# -u [UID] : UID 변경
		# -s [Shell] : shell 지정
		# -L : 계정 락킹
		# -U : 계정 언락킹
		# -g [group] : 사용자 기본 그룹 변경
		> usermod -g user2 user1
		# user1의 기본 그룹을 user2로 변경
		# -G [groups] : 사용자 그룹 추가,변경(제거). 기본 그룹은 영향을 받지 않는다.
		> usermod -a -G group1,group2 user1
		# user1에 group1,group2를 추가. -a 옵션은 기존그룹에 추가할지 안할지 여부이다.
		> usermod -G group1 user1
		# 그룹을 제거하는 방법. -a를 주지 않아 기존 그룹을 유지하지 않았다.

		> passwd [options] 사용자명
		# options
		# -d : 패스워드 삭제
		# -e : 패스워드 강제 만료
		# -l, -u : 패스워드 락킹/언락킹. usermod에 있는 -L, -U 옵션을 쓰는것이 더 좋다고 한다.

		Usage: passwd [options] [LOGIN]
		Options:
		-a, --all                     report password status on all accounts
		-d, --delete                  delete the password for the named account
		-e, --expire                  force expire the password for the named account
		-h, --help                    display this help message and exit
		-k, --keep-tokens             change password only if expired
		-i, --inactive INACTIVE       set password inactive after expiration
										to INACTIVE
		-l, --lock                    lock the password of the named account
		-n, --mindays MIN_DAYS        set minimum number of days before password
										change to MIN_DAYS
		-q, --quiet                   quiet mode
		-r, --repository REPOSITORY   change password in REPOSITORY repository
		-R, --root CHROOT_DIR         directory to chroot into
		-S, --status                  report password status on the named account
		-u, --unlock                  unlock the password of the named account
		-w, --warndays WARN_DAYS      set expiration warning days to WARN_DAYS
		-x, --maxdays MAX_DAYS        set maximum number of days before password
										change to MAX_DAYS

		> passwd -S 사용자명
		tgm1 P 08/30/2021 0 99999 7 -1
			P : passwd -u
			L : passwd -l

		> chage -l tgm1
		Last password change								: Aug 30, 2021
		Password expires									: never
		Password inactive									: never				chage  -I
		Account expires										: never				passwd -E
		Minimum number of days between password change		: 0					chage  -m
		Maximum number of days between password change		: 99999				chage  -M
		Number of days of warning before password expires	: 7   				chage  -W

		userdel [options] 사용자명
		# options
		# -r : 사용자의 홈디렉토리, 메일박스, 임시디렉토리 까지 같이 삭제.
			그룹은 기본적으로 속한 사용자가 없으면 자동 삭제된다

		------------------------------------

		/etc/group : 그룹 정보가 담긴 파일
		sudo:x:27:user1,user2
		# 그룹명:패스워드:GID:사용자리스트

		> groupadd [options] 그룹명
		# options
		# -g [GID] : GID 지정

		> groupmod [options] 그룹명
		# options
		# -n [이름] : 그룹명 변경

		> groupdel 그룹명
	*/
}
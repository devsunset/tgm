package groups

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Group describes a Group.
type Group struct {
	ID      string `json:"id"`
	Gid     string `json:"gid"`
	Members string `json:"members"`
}
type Store interface {
	Gets() ([]Group, error)
	Save(ID string) error
	Delete(ID string) error
}

func Gets() ([]Group, error) {
	var LinuxGroups [][]string
	groups := []Group{}

	// this is for Linux/Unix machines
	file, err := os.Open("/etc/group")
	if err != nil {
		fmt.Println(err)
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
			fmt.Println(err)
			return groups, err
		}
	}

	for _, data := range LinuxGroups {
		if err != nil {
			panic(err)
		}
		group := Group{}
		group.ID = data[0]
		group.Gid = data[2]
		group.Members = data[3]
		groups = append(groups, group)
	}

	return groups, nil
}

func Save(ID string) error {
	return nil
}

func Delete(ID string) error {
	return nil
}

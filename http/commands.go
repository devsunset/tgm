package http

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/gorilla/websocket"

	"tgm/runner"
)

const (
	WSWriteDeadline = 10 * time.Second
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var (
	cmdNotAllowed = []byte("Command not allowed.")
)

func wsErr(ws *websocket.Conn, r *http.Request, status int, err error) { //nolint:unparam
	txt := http.StatusText(status)
	if err != nil || status >= 400 {
		log.Printf("%s: %v %s %v", r.URL.Path, status, r.RemoteAddr, err)
	}
	if err := ws.WriteControl(websocket.CloseInternalServerErr, []byte(txt), time.Now().Add(WSWriteDeadline)); err != nil {
		log.Print(err)
	}
}

var commandsHandler = withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer conn.Close()

	var raw string

	for {
		_, msg, err := conn.ReadMessage() //nolint:govet
		if err != nil {
			wsErr(conn, r, http.StatusInternalServerError, err)
			return 0, nil
		}

		raw = strings.TrimSpace(string(msg))
		if raw != "" {
			break
		}
	}

	if d.server.EnableCmdLimit {
		if !d.server.EnableExec || !d.user.CanExecute(strings.Split(raw, " ")[0]) {
			if err := conn.WriteMessage(websocket.TextMessage, cmdNotAllowed); err != nil { //nolint:govet
				wsErr(conn, r, http.StatusInternalServerError, err)
			}

			return 0, nil
		}
	} else {
		if !d.server.EnableExec {
			if err := conn.WriteMessage(websocket.TextMessage, cmdNotAllowed); err != nil { //nolint:govet
				wsErr(conn, r, http.StatusInternalServerError, err)
			}

			return 0, nil
		}
	}

	// if !d.server.EnableExec || !d.user.CanExecute(strings.Split(raw, " ")[0]) {
	// 	if err := conn.WriteMessage(websocket.TextMessage, cmdNotAllowed); err != nil { //nolint:govet
	// 		wsErr(conn, r, http.StatusInternalServerError, err)
	// 	}

	// 	return 0, nil
	// }

	command, err := runner.ParseCommand(d.settings, raw)
	if err != nil {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(err.Error())); err != nil { //nolint:govet
			wsErr(conn, r, http.StatusInternalServerError, err)
		}
		return 0, nil
	}

	// TGM root 권한으로 실행 되지만 console 실행시는 사용자 권한으로 실행되게 처리
	// $ sudo -u 아이디 "명령어"

	// original source
	// cmd := exec.Command(command[0], command[1:]...) //nolint:gosec

	var cmds []string
	if d.user.Username != "admin" {
		cmds = append(cmds, "sudo")
		cmds = append(cmds, "-u")
		cmds = append(cmds, d.user.Username)
		cmds = append(cmds, command...)
	} else {
		cmds = append(cmds, command...)
	}

	cmd := exec.Command(cmds[0], cmds[1:]...) //nolint:gosec

	cmd.Dir = d.user.FullPath(r.URL.Path)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		wsErr(conn, r, http.StatusInternalServerError, err)
		return 0, nil
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		wsErr(conn, r, http.StatusInternalServerError, err)
		return 0, nil
	}

	if err := cmd.Start(); err != nil {
		wsErr(conn, r, http.StatusInternalServerError, err)
		return 0, nil
	}

	s := bufio.NewScanner(io.MultiReader(stdout, stderr))
	for s.Scan() {
		if err := conn.WriteMessage(websocket.TextMessage, s.Bytes()); err != nil {
			log.Print(err)
		}
	}

	if err := cmd.Wait(); err != nil {
		wsErr(conn, r, http.StatusInternalServerError, err)
	}

	return 0, nil
})

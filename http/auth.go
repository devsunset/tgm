package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"

	"tgm/errors"
	"tgm/users"
)

const (
	TokenExpirationTime    = time.Hour * 1
	TokenMaxExpirationTime = time.Hour * 24 * 30
)

type userInfo struct {
	ID           uint              `json:"id"`
	Locale       string            `json:"locale"`
	ViewMode     users.ViewMode    `json:"viewMode"`
	SingleClick  bool              `json:"singleClick"`
	Perm         users.Permissions `json:"perm"`
	Commands     []string          `json:"commands"`
	LockPassword bool              `json:"lockPassword"`
	HideDotfiles bool              `json:"hideDotfiles"`
}

type authToken struct {
	User userInfo `json:"user"`
	jwt.StandardClaims
}

type extractor []string

func (e extractor) ExtractToken(r *http.Request) (string, error) {
	token, _ := request.HeaderExtractor{"X-Auth"}.ExtractToken(r)

	// Checks if the token isn't empty and if it contains two dots.
	// The former prevents incompatibility with URLs that previously
	// used basic auth.
	if token != "" && strings.Count(token, ".") == 2 {
		return token, nil
	}

	auth := r.URL.Query().Get("auth")
	if auth != "" && strings.Count(auth, ".") == 2 {
		return auth, nil
	}

	cookie, _ := r.Cookie("auth")
	if cookie != nil && strings.Count(cookie.Value, ".") == 2 {
		return cookie.Value, nil
	}

	return "", request.ErrNoTokenInRequest
}

func withUser(fn handleFunc) handleFunc {
	return func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
		keyFunc := func(token *jwt.Token) (interface{}, error) {
			return d.settings.Key, nil
		}

		var tk authToken
		token, err := request.ParseFromRequest(r, &extractor{}, keyFunc, request.WithClaims(&tk))

		if err != nil || !token.Valid {
			return http.StatusForbidden, nil
		}

		expired := !tk.VerifyExpiresAt(time.Now().Add(time.Hour).Unix(), true)
		updated := d.store.Users.LastUpdate(tk.User.ID) > tk.IssuedAt

		if expired || updated {
			w.Header().Add("X-Renew-Token", "true")
		}

		d.user, err = d.store.Users.Get(d.server.Root, tk.User.ID)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		return fn(w, r, d)
	}
}

func withAdmin(fn handleFunc) handleFunc {
	return withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
		if !d.user.Perm.Admin {
			return http.StatusForbidden, nil
		}

		return fn(w, r, d)
	})
}

var loginHandler = func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	auther, err := d.store.Auth.Get(d.settings.AuthMethod)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	user, err := auther.Auth(r, d.store.Users, d.server.Root)
	if err == os.ErrPermission {
		return http.StatusForbidden, nil
	} else if err != nil {
		return http.StatusInternalServerError, err
	} else {
		////////////////////////////////////////////////////////////////////////////
		// LINUX 계정 정보 유효성 체크
		////////////////////////////////////////////////////////////////////////////
		// DB에 저장된 정보는 Sync가 안맞는 경우 발생 할 수 있음으로 최대한 리눅스 계정과
		// Sync  처리 하기 위해 리눅스 계정 정보 조회 하여 리턴 값 대체 처리
		if strings.Compare(user.Username, "admin") != 0 {
			//chage -l $USER | grep "Account expires"
			// user.ExpireDay = ""
			expiresCmd := exec.Command("sh", "-c", "chage -l  "+user.Username+" | grep 'Account expires'")
			if out, err := expiresCmd.Output(); err != nil {
				log.Println("get chage info error", err)
			} else {
				outStr := string(out)
				outStr = strings.TrimSpace(outStr)
				if strings.Index(outStr, "never") > -1 {
					user.ExpireDay = "9999-12-31"
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
					user.ExpireDay = slice[4] + "-" + month + "-" + slice[3][:len(slice[3])-1]
				}
			}
			// passwd -S $USER
			// u.PasswrodExpireDay = ""
			// u.PasswordExpireWarningDay = ""
			// u.LockAccount = true
			passwdCmd := exec.Command("passwd", "-S", user.Username)
			if out, err := passwdCmd.Output(); err != nil {
				log.Println("get passwd info error", err)
			} else {
				outStr := string(out)
				outStr = strings.TrimSpace(outStr)
				slice := strings.Split(outStr, " ")

				user.PasswrodExpireDay = slice[4]
				user.PasswordExpireWarningDay = slice[5]

				if strings.Compare(slice[1], "LK") == 0 {
					user.LockAccount = true
				} else {
					user.LockAccount = false
				}
			}

			// 계정 유효 일자 마감
			expireDay, _ := strconv.Atoi(strings.Replace(user.ExpireDay, "-", "", 2))

			t := time.Now()
			formatted := fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day())
			today, _ := strconv.Atoi(strings.Replace(formatted, "-", "", 2))

			if expireDay < today {
				return http.StatusForbidden, nil
			}

			// 계정 잠김
			if user.LockAccount == true {
				return http.StatusForbidden, nil
			}

		}
		////////////////////////////////////////////////////////////////////////////
		return printToken(w, r, d, user)
	}
}

var passwdValidPeriodCheckHandler = func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	auther, err := d.store.Auth.Get(d.settings.AuthMethod)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	user, err := auther.Auth(r, d.store.Users, d.server.Root)
	if err == os.ErrPermission {
		return http.StatusForbidden, nil
	} else if err != nil {
		return http.StatusInternalServerError, err
	} else {
		////////////////////////////////////////////////////////////////////////////
		// LINUX 계정 정보 유효성 체크
		////////////////////////////////////////////////////////////////////////////
		result := "S"
		expireDay := ""
		passwordExpireWarningDay := ""
		if strings.Compare(user.Username, "admin") != 0 {
			//chage -l $USER | grep "Password expires"
			expiresCmd := exec.Command("sh", "-c", "chage -l  "+user.Username+" | grep 'Password expires'")
			if out, err := expiresCmd.Output(); err != nil {
				log.Println("get chage info error", err)
			} else {
				outStr := string(out)
				outStr = strings.TrimSpace(outStr)
				if strings.Index(outStr, "never") > -1 {
					expireDay = "never"
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
					expireDay = slice[4] + "-" + month + "-" + slice[3][:len(slice[3])-1]
				}

			}
			// passwd -S $USER
			passwdCmd := exec.Command("passwd", "-S", user.Username)
			if out, err := passwdCmd.Output(); err != nil {
				log.Println("get passwd info error", err)
			} else {
				outStr := string(out)
				outStr = strings.TrimSpace(outStr)
				slice := strings.Split(outStr, " ")
				passwordExpireWarningDay = slice[5]
			}

			if expireDay != "never" {

				// 패스 워드 유효 일자 마감
				expireDayx, _ := strconv.Atoi(strings.Replace(expireDay, "-", "", 2))
				t := time.Now()
				formatted := fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day())
				today, _ := strconv.Atoi(strings.Replace(formatted, "-", "", 2))

				if expireDayx < today {
					result = "E"
				} else {
					if passwordExpireWarningDay != "-1" {
						passwordExpireWarningDayx, _ := strconv.Atoi(passwordExpireWarningDay)
						passExpireDay, _ := time.Parse("2006-01-02", expireDay)
						time := passExpireDay.AddDate(0, 0, -passwordExpireWarningDayx)
						formattedx := fmt.Sprintf("%d-%02d-%02d", time.Year(), time.Month(), time.Day())
						warningDay, _ := strconv.Atoi(strings.Replace(formattedx, "-", "", 2))
						if warningDay <= today {
							result = expireDay
						}
					}
				}
			}
		}

		w.Header().Set("Content-Type", "text/plain")
		if _, err := w.Write([]byte(result)); err != nil {
			return http.StatusInternalServerError, err
		}
		return 0, nil
	}
}

type signupBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var signupHandler = func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	if !d.settings.Signup {
		return http.StatusMethodNotAllowed, nil
	}

	if r.Body == nil {
		return http.StatusBadRequest, nil
	}

	info := &signupBody{}
	err := json.NewDecoder(r.Body).Decode(info)
	if err != nil {
		return http.StatusBadRequest, err
	}

	if info.Password == "" || info.Username == "" {
		return http.StatusBadRequest, nil
	}

	user := &users.User{
		Username: info.Username,
	}

	d.settings.Defaults.Apply(user)

	pwd, err := users.HashPwd(info.Password)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	user.Password = pwd
	user.PasswordHint = users.HintPwd(info.Password)

	userHome, err := d.settings.MakeUserDir(user.Username, user.Scope, d.server.Root)
	if err != nil {
		log.Printf("create user: failed to mkdir user home dir: [%s]", userHome)
		return http.StatusInternalServerError, err
	}
	user.Scope = userHome
	log.Printf("new user: %s, home dir: [%s].", user.Username, userHome)

	err = d.store.Users.Save(user)
	if err == errors.ErrExist {
		return http.StatusConflict, err
	} else if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

var renewHandler = withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	return printToken(w, r, d, d.user)
})

func printToken(w http.ResponseWriter, r *http.Request, d *data, user *users.User) (int, error) {
	tokenMaxPeriod := false

	rememberme := r.URL.Query().Get("rememberme")
	if rememberme == "true" {
		tokenMaxPeriod = true
	}

	if !tokenMaxPeriod {
		c1, err := r.Cookie("rememberme")
		if err != nil {

		} else {
			if c1.Value == "true" {
				tokenMaxPeriod = true
			}
		}
	}

	TokenExpirationTime_Value := time.Hour * 1

	if tokenMaxPeriod {
		TokenExpirationTime_Value = TokenMaxExpirationTime
	} else {
		TokenExpirationTime_Value = TokenExpirationTime
	}

	claims := &authToken{
		User: userInfo{
			ID:           user.ID,
			Locale:       user.Locale,
			ViewMode:     user.ViewMode,
			SingleClick:  user.SingleClick,
			Perm:         user.Perm,
			LockPassword: user.LockPassword,
			Commands:     user.Commands,
			HideDotfiles: user.HideDotfiles,
		},
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(TokenExpirationTime_Value).Unix(),
			Issuer:    "tgmJwtAuthManagerIssuer",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(d.settings.Key)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	w.Header().Set("Content-Type", "text/plain")
	if _, err := w.Write([]byte(signed)); err != nil {
		return http.StatusInternalServerError, err
	}
	return 0, nil
}

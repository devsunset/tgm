package auth

import (
	"net/http"

	"tgm/settings"
	"tgm/users"
)

// MethodNoAuth is used to identify no auth.
const MethodNoAuth settings.AuthMethod = "noauth"

// NoAuth is no auth implementation of auther.
type NoAuth struct{}

// Auth uses authenticates user 1.
func (a NoAuth) Auth(r *http.Request, sto users.Store, root string) (*users.User, error) {
	return sto.Get(root, uint(1))
}

// LoginPage tells that no auth doesn't require a login page.
func (a NoAuth) LoginPage() bool {
	return false
}

package auth

import "net/http"

// define Authenticated interface
type Authenticated interface {
	GetID() int64
	GetName() string
}

// define func to autenticate
type authenticateFunc func(string, string) (Authenticated, error)

// create handler to hold the authenticate func
type handler struct {
	authenticate authenticateFunc
}

// create exported func to receive the authenticate func
func HandleAuth(fn authenticateFunc) func(http.ResponseWriter, *http.Request) {
	h := handler{fn}

	return h.auth
}

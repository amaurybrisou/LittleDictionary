package middlewares

import (  
  "net/http"
)

type BasicAuth struct {
    Login    string
    Password string
    Realm    string
}

func NewBasicAuth(login, pass string) *BasicAuth {
    return &BasicAuth{Login: login, Password: pass}
}

func (a *BasicAuth) ValidAuth(r *http.Request) bool {
    username, password, _ := r.BasicAuth()
    return username == a.Login && password == a.Password
}

func (a *BasicAuth) BasicAuthHandler(h http.Handler) http.Handler {
    f := func(w http.ResponseWriter, r *http.Request) {
        if !a.ValidAuth(r) {
            a.Authenticate(w, r)
        } else {
            h.ServeHTTP(w, r)
        }
    }

    return http.HandlerFunc(f)
}

func (a *BasicAuth) Authenticate(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("WWW-Authenticate", `Basic realm="`+a.Realm+`"`)
    http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
}
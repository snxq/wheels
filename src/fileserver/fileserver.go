package main

import (
	"flag"
	"fmt"
	"net/http"
)

var (
	port   = flag.String("port", ":80", "server port")
	path   = flag.String("path", ".", "file path")
	user   = flag.String("user", "root", "auth username")
	pass   = flag.String("pass", "", "auth password")
	domain = flag.String("domain", "", "domain for auth")
)

func main() {
	flag.Parse()

	http.ListenAndServe(*port, http.HandlerFunc(auth))
}

func auth(w http.ResponseWriter, r *http.Request) {
	u, p, ok := r.BasicAuth()
	w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, *domain))
	if !ok || u != *user || p != *pass {
		http.Error(w, "auth failed", http.StatusUnauthorized)
		return
	}

	http.FileServer(http.Dir(*path)).ServeHTTP(w, r)
}

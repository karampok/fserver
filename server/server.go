package server

import (
	"net/http"

	"github.com/karampok/fserver/filesystem"
)

// Server ...
type Server struct {
	dir, user, pass string
	rtr             *http.ServeMux
}

// NewServer ...
func NewServer(d, p string) *Server {
	s := &Server{dir: d, pass: p, rtr: http.NewServeMux()}
	s.routes()
	return s
}

func (s *Server) auth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.pass != "" { //!currentUser(r).IsAdmin {
			http.NotFound(w, r)
			return
		}
		h(w, r)
	}
}

func (s *Server) routes() {
	sd := filesystem.FS{http.Dir(s.dir)}
	fs := http.FileServer(sd)
	s.rtr.HandleFunc("/", s.auth(fs.ServeHTTP))
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("server", ".me")
	s.rtr.ServeHTTP(w, r)
}

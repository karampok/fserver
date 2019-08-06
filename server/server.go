package server

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"net/http"
	"time"

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
		_, pass, _ := r.BasicAuth()

		if pass == s.pass {
			h(w, r)
			return
		}
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, "Unauthorized.", 401)
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

// GenTLSConfig ...
func GenTLSConfig() *tls.Config {
	ret := &tls.Config{}
	ret.Certificates = make([]tls.Certificate, 1)

	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization: []string{"ME"},
			Country:      []string{"CH"},
			Province:     []string{"ZRH"},
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(0, 0, 1),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
		DNSNames:     []string{"fserver"},
	}

	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &priv.PublicKey, priv)
	if err != nil {
		panic(err)
	}

	ret.Certificates[0] = tls.Certificate{
		Certificate: [][]byte{certBytes},
		PrivateKey:  priv,
	}

	return ret
}

package main

import (
	"context"
	"flag"
	"log"
	rdr "math/rand"

	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/karampok/fserver/server"
)

func main() {
	lg := log.New(os.Stdout, "", 0)

	addr := flag.String("l", ":8443", "address:port to listen at (eg. 0.0.0.0:8443)")
	flag.Parse()

	dir := "."
	if flag.NArg() > 0 {
		dir = flag.Arg(0)
	}

	pass := randomString(16)
	h := &http.Server{
		Addr:      *addr,
		Handler:   server.NewServer(dir, pass),
		TLSConfig: server.GenTLSConfig(),
	}

	go func() {
		lg.Printf("Serving dir[%s] on [%s] protected by user:< :%s >\n", dir, *addr, pass)
		if err := h.ListenAndServeTLS("", ""); err != nil {
			lg.Fatal(err)
		}
		return
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	lg.Println("Shutting down the server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := h.Shutdown(ctx); err != nil {
		lg.Fatal(err)
	}
	lg.Println("Server stopped gracefully ...")
}

func randomString(length int) string {
	seededRand := rdr.New(rdr.NewSource(time.Now().UnixNano()))
	charset := "abcdefghijklmnopqrstuvwxyz" + "0123456789"

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/karampok/fserver/server"
)

func main() {
	lg := log.New(os.Stdout, "", 0)

	addr := flag.String("l", ":8443", "address:port to listen at (eg. 0.0.0.0:8443)")
	https := flag.Bool("https", false, "serve https")
	protected := flag.Bool("p", false, "generate protected password")
	flag.Parse()

	dir := "."
	if flag.NArg() > 0 {
		dir = flag.Arg(0)
	}
	pass := ""
	if *protected {
		pass = randomString(8)
	}

	h := &http.Server{Addr: *addr, Handler: server.NewServer(dir, pass)}
	go func() {
		lg.Printf("Serving dir[%s] on [%s] protected by pass[%s]\n", dir, *addr, pass)
		if *https {
			crt, key := "certs/server.crt", "certs/server.key"
			if err := h.ListenAndServeTLS(crt, key); err != nil {
				lg.Fatal(err)
			}
			return
		}
		if err := h.ListenAndServe(); err != nil {
			lg.Fatal(err)
		}
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
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	charset := "abcdefghijklmnopqrstuvwxyz" + "0123456789"

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/karampok/fserver/server"
)

func main() {
	lg := log.New(os.Stdout, "", 0)

	addr := flag.String("l", ":8443", "Address:Port to listen at (eg. 0.0.0.0:8443)")
	protected := flag.Bool("p", false, "generate protected password")
	flag.Parse()

	crt := "certs/server.crt"
	key := "certs/server.key"

	dir := "."
	if flag.NArg() > 0 {
		dir = flag.Arg(0)
	}
	pass := ""
	if *protected {
		pass = "random"
	}

	h := &http.Server{Addr: *addr, Handler: server.NewServer(dir, pass)}
	go func() {
		lg.Printf("Serving [%s] on %s protected by [%s]\n", dir, *addr, pass)
		if err := h.ListenAndServeTLS(crt, key); err != nil {
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

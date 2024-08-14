package main

import (
	"flag"
	check "forum/internal"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	cfg := new(check.Config)
	flag.StringVar(&cfg.Addr, "addr", ":8081", "HTTP network address")
	flag.StringVar(&cfg.StaticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.Parse()
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	srv := check.CreateServer(infoLog, errorLog, cfg)

	infoLog.Printf("Starting server on https://localhost%s", cfg.Addr)
	err := srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}

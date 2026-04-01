package main

import (
	"log"

	"github.com/iainjreid/source/db"
	"github.com/iainjreid/source/ssh"
	"github.com/iainjreid/source/web"
	storage "github.com/iainjreid/go-git-sql"
	"golang.org/x/sync/errgroup"
)

func main() {
	if err := db.DB.Ping(); err != nil {
		log.Fatal("DB unreachable: ", err)
	}

	storage := storage.NewStorage(db.DB)
	if err := db.HardReset(); err != nil {
		panic(err)
	}

	wg := new(errgroup.Group)

	wg.Go(func() error {
		log.Println("Starting Web server")
		return web.NewServer(storage)
	})

	wg.Go(func() error {
		log.Println("Starting SSH server")
		return ssh.NewServer(storage)
	})

	err := wg.Wait()
	if err != nil {
		log.Println(err)
	}
}

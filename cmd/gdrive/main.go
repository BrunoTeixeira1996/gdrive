package main

import (
	"log"

	"github.com/BrunoTeixeira1996/gdrive/internal/auth"
	"github.com/BrunoTeixeira1996/gdrive/internal/handles"
)

func run() error {
	server, err := auth.GetDriveService()
	if err != nil {
		return err
	}

	if err := handles.Init(server); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalf(err.Error())
	}
}

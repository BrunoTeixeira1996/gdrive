package main

import (
	"flag"
	"log"

	"github.com/BrunoTeixeira1996/gdrive/internal/auth"
	"github.com/BrunoTeixeira1996/gdrive/internal/handles"
)

func run() error {
	gokrazyFlag := flag.Bool("gokrazy", false, "use this if you want to use gokrazy")
	flag.Parse()

	if *gokrazyFlag {
		log.Println("[run info] ok lets do this on gokrazy then ...")
	}

	server, err := auth.GetDriveService(*gokrazyFlag)
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

// TODO: make work in gokrazy

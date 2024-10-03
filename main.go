package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/BrunoTeixeira1996/gdrive/internal/action"
	"github.com/BrunoTeixeira1996/gdrive/internal/auth"
)

func run() error {
	var (
		driveFolderFlag = flag.String("drivefolder", "", "Faturas/2024/i_Setembro")
	)

	flag.Parse()

	if *driveFolderFlag == "" {
		return fmt.Errorf("[run error] please provide the drivefolder flag")
	}
	server, err := auth.GetDriveService()
	if err != nil {
		return err
	}

	folderId, err := action.GetPathId(server, *driveFolderFlag)
	if err != nil {
		return err
	}

	log.Printf("[run info] gathered folder id '%s' for folder '%s'\n", folderId, *driveFolderFlag)

	if err = action.OutputCSV(server, folderId); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

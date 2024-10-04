package handles

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/BrunoTeixeira1996/gdrive/internal/action"
	"google.golang.org/api/drive/v3"
)

type CSV struct {
	FullPathFolder string
	Content        string
	NOfFiles       string
}

type UI struct {
	Server   *drive.Service
	AllPaths string
	CSV      CSV
	Tmpl     *template.Template
}

// Queries a google drive path and outputs in csv mode
func (ui *UI) queryInvoices(w http.ResponseWriter, r *http.Request) {
	var (
		outputcsv string
		folderId  string
		nOfFiles  string
		err       error
	)

	switch r.Method {
	case http.MethodGet:
		if err = ui.Tmpl.ExecuteTemplate(w, "gdrive.html.tmpl", ui); err != nil {
			log.Printf("[queryInvoices error] executing template (get): %s\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
		path := r.FormValue("path")
		log.Printf("[queryInvoices info] querying for %s\n", path)
		folderId, err = action.GetPathId(ui.Server, path)
		if err != nil {
			log.Printf("[queryInvoices error] get path id: %s\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if outputcsv, nOfFiles, err = action.OutputCSV(ui.Server, folderId); err != nil {
			log.Printf("[queryInvoices error] output csv: %s\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ui.CSV.Content = outputcsv
		ui.CSV.FullPathFolder = path
		ui.CSV.NOfFiles = nOfFiles

		if err = ui.Tmpl.ExecuteTemplate(w, "gdrive.html.tmpl", ui); err != nil {
			log.Printf("[queryInvoices error] executing template (post): %s\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

//go:embed assets/*
var assetsDir embed.FS

func Init(server *drive.Service) error {

	var err error

	tmpl, err := template.ParseFS(assetsDir, "assets/*.tmpl")
	if err != nil {
		return fmt.Errorf("[handles error] could not parse template: %s\n", err)
	}

	ui := &UI{
		Server: server,
		Tmpl:   tmpl,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", ui.queryInvoices)

	log.Printf("[handles info] listen at port 9393 \n")

	err = http.ListenAndServe(":9393", mux)
	if err != nil && err != http.ErrServerClosed {
		panic("[handles error] trying to start http server: " + err.Error())
	}

	return nil
}

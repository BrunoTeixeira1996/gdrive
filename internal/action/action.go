package action

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"google.golang.org/api/drive/v3"
)

type File struct {
	Id          string
	Name        string
	Category    string
	Price       string
	Owner       string
	Description string
}

func (f *File) extractParentheses() (string, string) {
	// Regex to capture both first and second sets of parentheses
	re := regexp.MustCompile(`\(([^)]*)\)\(([^)]*)\)`)
	match := re.FindStringSubmatch(f.Name)

	if len(match) > 2 {
		return match[1], match[2] // Return both first and second parentheses
	}
	return "", "" // Return empty strings if not found
}

func (f *File) fixFile() {
	parts := strings.Split(f.Name, "_")
	if len(parts) < 2 {
		// there's no underscore
		log.Printf("[action error] unexpected file name format: %s\n", f.Name)
		return
	}

	f.Category = parts[0]

	pricePart := strings.Split(parts[1], "(")[0]
	f.Price = strings.Replace(pricePart, "-", ".", 1) // Replace the first hyphen with a period

	f.Owner, f.Description = f.extractParentheses()
}

// Modify the files
func fixFiles(files *[]File) {
	for i := range *files {
		(*files)[i].fixFile()
	}
}

func OutputCSV(server *drive.Service, folderId string) (string, string, error) {
	var csvRows []string

	files, err := getFilesFromFolder(server, folderId)
	if err != nil {
		return "", "", err
	}

	// 	Mix_21-12()(PrendaRegina).pdf
	//  Mix_21-12()().pdf
	//  Mix_21-12(B)(PC).pdf

	// output to terminal csv format like this
	// 12.1,Carro,
	// 22.1,Supermercado,
	// 31,Veterinario,Alex
	// 12.21,Mix,Bruno

	for _, f := range files {
		row := f.Price + "," + f.Category + "," + f.Owner + "," + f.Description
		csvRows = append(csvRows, row)
	}

	return strings.Join(csvRows, "\n"), strconv.Itoa(len(files)), nil
}

// findFolderID traverses a given folder path in Google Drive and returns the folder ID
func GetPathId(server *drive.Service, driveFolder string) (string, error) {
	parentID := "root"
	folders := strings.Split(driveFolder, "/")

	for _, folderName := range folders {
		query := fmt.Sprintf("name = '%s' and mimeType = 'application/vnd.google-apps.folder' and '%s' in parents and trashed = false", folderName, parentID)
		r, err := server.Files.List().Q(query).Fields("files(id, name)").Do()
		if err != nil {
			return "", fmt.Errorf("[action error] could not perform query to search for path id: %s\n", err)
		}

		if len(r.Files) == 0 {
			return "", fmt.Errorf("[action error] folder '%s' not found in '%s'\n", folderName, parentID)
		}

		// Assume the first folder match is the correct one
		parentID = r.Files[0].Id
	}

	return parentID, nil
}

func getFilesFromFolder(server *drive.Service, folderId string) ([]File, error) {
	var files []File

	query := fmt.Sprintf("'%s' in parents", folderId)
	pageToken := ""
	for {
		// List files in the folder, handling pagination
		r, err := server.Files.List().
			Q(query).
			Fields("nextPageToken, files(id, name)").
			PageToken(pageToken).
			PageSize(1000). // You can set up to 1000 per page
			Do()
		if err != nil {
			return []File{}, fmt.Errorf("[action error] unable to retrieve files: %v", err)
		}

		// Print out the file names and IDs
		if len(r.Files) == 0 {
			fmt.Println("[auth info] no files found.")
		} else {
			log.Printf("[auth info] this folder has %d files\n", len(r.Files))
			for _, f := range r.Files {
				// ignore xlsx files
				if strings.Contains(f.Name, "xlsx") {
					continue
				}

				var file File
				file.Id = f.Id
				file.Name = f.Name
				files = append(files, file)
			}
		}

		// If there's a next page, continue retrieving
		if r.NextPageToken == "" {
			break
		}
		pageToken = r.NextPageToken
	}

	// this fixes the category and price from the file name
	fixFiles(&files)

	return files, nil
}

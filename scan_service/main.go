package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"

	"github.com/Azure/azure-storage-blob-go/2018-03-28/azblob"
)

const (
	downloadDirectory       = "./downloads"
	quarantineContainerName = "quarantine"
	cleanContainerName      = "clean"
	virusContainerName      = "virus"
	scanCommand             = "c:\\Program Files\\Windows Defender\\MpCmdRun.exe"
)

func main() {

	accountName := os.Getenv("AZ_ACC_NAME")
	accountKey := os.Getenv("AZ_ACC_KEY")
	host := os.Getenv("HOST")

	log.Printf("Found settings accountName = %s", accountName)

	// Create Downloads Directory if does not exists
	if _, err := os.Stat(downloadDirectory); os.IsNotExist(err) {
		os.Mkdir(downloadDirectory, os.ModePerm)
	}

	http.HandleFunc("/scan/", func(w http.ResponseWriter, r *http.Request) {

		// Split url and get parts
		scanurl := r.URL.Path
		parts := strings.Split(scanurl, "/")

		if len(parts) > 2 {

			fileName := parts[2]
			log.Printf("Request to download file %s", fileName)

			ctx := context.Background()

			// Download Blob from Container
			c, err := azblob.NewSharedKeyCredential(accountName, accountKey)
			if err != nil {
				sendError(w, err.Error())
				return
			}

			p := azblob.NewPipeline(c, azblob.PipelineOptions{})

			u, err := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net", accountName))
			if err != nil {
				sendError(w, err.Error())
				return
			}

			service := azblob.NewServiceURL(*u, p)
			quarantineContainer := service.NewContainerURL(quarantineContainerName)
			b := quarantineContainer.NewBlobURL(fileName)

			f, err := os.Create(path.Join("downloads", fileName))
			if err != nil {
				sendError(w, err.Error())
			}
			defer f.Close()

			// Write Blob to File
			log.Printf("Downloading Blob %s from url %s", fileName, u)
			err = azblob.DownloadBlobToFile(ctx, b, 0, 0, f, azblob.DownloadFromBlobOptions{})
			if err != nil {
				sendError(w, err.Error())
				return
			}

			// Scan File
			dir, err := os.Getwd()
			if err != nil {
				sendError(w, err.Error())
				return
			}

			result := scanFile(fmt.Sprintf("%s\\downloads\\%s", dir, fileName))

			// Publish result to new container
			var containerToMoveFileTo azblob.ContainerURL

			if result.Status == "OK" {
				containerToMoveFileTo = service.NewContainerURL(cleanContainerName)
			} else {
				containerToMoveFileTo = service.NewContainerURL(virusContainerName)
			}

			bURL := containerToMoveFileTo.NewBlockBlobURL(fileName)
			_, err = azblob.UploadFileToBlockBlob(ctx, f, bURL, azblob.UploadToBlockBlobOptions{})

			if err != nil {
				sendError(w, err.Error())
				return
			}

			// Send Response
			js, err := json.Marshal(result)
			if err != nil {
				sendError(w, err.Error())
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
		} else {
			// Send an Error back
			sendError(w, "You need to pass an object name")
		}
	})

	log.Print("Server is starting")
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:80", host), nil))
}

func sendError(w http.ResponseWriter, err string) {
	log.Printf("ERROR: %s", err)
	http.Error(w, err, http.StatusInternalServerError)
}

func scanFile(fileName string) ScanResult {

	log.Printf("File being scanned = %s", fileName)

	result := ScanResult{
		Status:      "OK",
		Description: "",
	}

	scanArgs := []string{"-Scan", "-ScanType", "3", "-File", fileName, "-DisableRemediation"}
	out, _ := exec.Command(scanCommand, scanArgs...).Output()

	r, _ := regexp.Compile("Threat                  :\\s(.+)")

	log.Printf("SCANNER OUTPUT", string(out))

	if r.MatchString(string(out)) {
		description := r.FindStringSubmatch(string(out))
		result.Status = "VIRUS FOUND"
		result.Description = description[1]
		log.Printf("ERROR FOUND", string(out))
	}

	return result
}

type ScanResult struct {
	Status      string
	Description string
}

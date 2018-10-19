package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/Azure/azure-storage-blob-go/2018-03-28/azblob"
	"github.com/dutchcoders/go-clamd"
)

const (
	downloadDirectory       = "./downloads"
	quarantineContainerName = "quarantine"
	cleanContainerName      = "clean"
)

func main() {

	accountName := os.Getenv("AZ_ACC_NAME")
	accountKey := os.Getenv("AZ_ACC_KEY")

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
			}

			p := azblob.NewPipeline(c, azblob.PipelineOptions{})

			u, err := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net", accountName))
			if err != nil {
				sendError(w, err.Error())
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
			}

			// Scan File
			clam := clamd.NewClamd("tcp://localhost:3310")
			response, err := clam.ScanFile(path.Join("downloads", fileName))

			if err != nil {
				sendError(w, err.Error())
			}

			result := ScanResult{
				Status:      "",
				Description: "",
			}
			for s := range response {
				result.Status = s.Status
				result.Description = s.Description
			}

			// Publish result to new container

			// Send Response
			js, err := json.Marshal(result)
			if err != nil {
				sendError(w, err.Error())
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
		} else {
			// Send an Error back
			sendError(w, "You need to pass an object name")
		}
	})

	log.Print("Server is starting")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func sendError(w http.ResponseWriter, err string) {
	log.Printf("ERROR: %s", err)
	http.Error(w, err, http.StatusInternalServerError)
}

type ScanResult struct {
	Status      string
	Description string
}

package main

import (
	"ThyTube/features"
	"ThyTube/utilities"
	"embed"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

//go:embed templates/*.html
var content embed.FS

type WatchData struct {
	Error  string
	Src    string
	Method string
	Title  string
}

type IndexData struct {
	Title string
	Files []string
	Error string
}

func main() {
	slog.Info("starting...")

	if err := os.MkdirAll("./download", 0755); err != nil {
		slog.Error("Failed to create download directory", "error", err)
		return
	}
	http.Handle("/download/", http.StripPrefix("/download/", http.FileServer(http.Dir("./download"))))

	tmpl := template.Must(template.New("").Funcs(template.FuncMap{
		"dict": utilities.Dict,
	}).ParseFS(content, "templates/*"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		indexData := IndexData{Title: "Home"}

		files, err := os.ReadDir("./download")
		if err != nil {
			slog.Error("Error reading download directory", "error", err)
			indexData.Error = "Could not read download directory."
		} else {
			var videoPaths []string
			for _, file := range files {
				if !file.IsDir() && !strings.HasPrefix(file.Name(), ".") {
					videoPaths = append(videoPaths, fmt.Sprintf("/download/%s", file.Name()))
				}
			}
			indexData.Files = videoPaths
		}

		err = tmpl.ExecuteTemplate(w, "index.html", indexData)
		if err != nil {
			slog.Error("Error running template for index.html", "error", err)
			return
		}
	})

	http.HandleFunc("/watch", func(w http.ResponseWriter, r *http.Request) {
		watchData := WatchData{Method: r.Method, Title: "Watch"}

		slog.Info("Handling /watch request", "method", r.Method)

		if r.Method == http.MethodPost {
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					slog.Error("Error closing body", "error", err)
				}
			}(r.Body)

			err := r.ParseForm()
			if err != nil {
				slog.Error("Error parsing form data", "error", err)
				http.Error(w, "Could not parse form data", http.StatusBadRequest)
				return
			}

			searchQuery := r.FormValue("query")

			if searchQuery == "" {
				watchData.Error = "No search query (link) provided!"
			} else if !strings.HasPrefix(searchQuery, "http") {
				watchData.Error = "Invalid search query, it needs protocol! (e.g. http://, https://)"
			}

			if watchData.Error == "" {
				slog.Info("Downloading video", "src", searchQuery)
				err = features.Download(searchQuery)

				if err != nil {
					slog.Error("Error downloading data", "error", err)
					watchData.Error = fmt.Sprintf("Download failed: %s", err.Error())
				} else {
					// todo: janky way to redirect after success, will be replaced with more acceptable solution in future
					http.Redirect(w, r, "/", http.StatusSeeOther)
					return
				}
			}
		} else if r.Method == http.MethodGet {
			videoSrc := r.URL.Query().Get("src")
			if videoSrc != "" {
				watchData.Src = videoSrc
				watchData.Title = "Watch"
			} else {
				watchData.Error = "No video source provided for playback."
			}
		}

		err := tmpl.ExecuteTemplate(w, "watch.html", watchData)
		if err != nil {
			slog.Error("Error executing template", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})

	slog.Info("listening on :9595")
	err := http.ListenAndServe(":9595", nil)
	if err != nil {
		slog.Error("Error running server", "error", err)
		return
	}
}

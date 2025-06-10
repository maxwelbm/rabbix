package web

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

//go:embed static/css/*.css static/js/*.js templates/*.html
var Assets embed.FS

// GetTemplate returns the parsed HTML template
func GetTemplate(name string) (*template.Template, error) {
	return template.ParseFS(Assets, "templates/"+name)
}

// GetStaticHandler returns an HTTP handler for static assets
func GetStaticHandler() http.Handler {
	// Debug: List embedded files
	fs.WalkDir(Assets, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		log.Printf("Embedded file: %s", path)
		return nil
	})

	staticFS, err := fs.Sub(Assets, "static")
	if err != nil {
		log.Printf("Error creating sub filesystem: %v", err)
		panic(err)
	}
	
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Static request: %s", r.URL.Path)
		http.FileServer(http.FS(staticFS)).ServeHTTP(w, r)
	})
}

// GetStaticFile returns the content of a static file
func GetStaticFile(path string) ([]byte, error) {
	return Assets.ReadFile(path)
}

// DebugAssets prints all embedded assets
func DebugAssets() {
	fmt.Println("=== Embedded Assets ===")
	fs.WalkDir(Assets, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			fmt.Printf("File: %s\n", path)
		}
		return nil
	})
	fmt.Println("=====================")
}
package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

func handle(err error) {
	if err != nil {
		log.Panicln(err)
	}
}

func convert(inputFile string, outputFile string) error {
	content, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return err
	}

	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)
	var buf bytes.Buffer
	err = md.Convert(content, &buf)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(outputFile, buf.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if os.Args[1] == "convert" {
		err := convert("community.md", "community.html")
		handle(err)
		fmt.Println("Done")
		return
	}

	err := convert("community.md", "community.html")
	handle(err)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		r, err := os.Open("community.html")
		handle(err)

		// w.Header().Set("content-type", "text/html")
		io.Copy(w, r)
		// fmt.Fprintf(w, "Welcome to the home page!")
	})
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	http.Handle("/", mux)
	log.Println("Listening on :3000...")
	err = http.ListenAndServe(":3000", nil)
	handle(err)
}

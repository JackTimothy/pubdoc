package main

import (
	"flag"
	"log"
	"path/filepath"

	"github.com/JackTimothy/pubdoc/configuration"
	"github.com/JackTimothy/pubdoc/confluence"
	"github.com/joho/godotenv"
)

func main() {
	var markdownFilePath, title string
	flag.StringVar(&markdownFilePath, "markdownFilePath", "", "The path to the Markdown file you want to publish.")
	flag.Parse()

	title = filepath.Base(markdownFilePath)
	title, err := confluence.FormatBase(title)
	if err != nil {
		log.Fatalf("Error formatting title: %v.", err)
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v.", err)
	}

	config := configuration.Init()

	htmlPageContent, err := confluence.ConvertMarkdownToHTML(markdownFilePath)
	if err != nil {
		log.Fatalf("Error reading Markdown file: %v.\n", err)
	}

	newPageId, err := confluence.CreatePage(title, htmlPageContent, "", config)
	if err != nil {
		log.Fatalf("Error creating page: %v\n", err)
	}
	log.Printf("Successfully created new page with ID %s.\n", newPageId)
}

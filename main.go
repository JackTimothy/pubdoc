package main

import (
	"flag"
	"log"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/JackTimothy/pubdoc/confluence"
	"github.com/JackTimothy/pubdoc/configuration"
)

func main() {
	var markdownFilePath, title string
	flag.StringVar(&markdownFilePath, "markdownFilePath", "", "The path to the Markdown file you want to publish.")
	flag.Parse()

	title = filepath.Base(markdownFilePath)

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

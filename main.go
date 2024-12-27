package main

import (
	"flag"
	"log"
	"path/filepath"
	"strings"

	"github.com/JackTimothy/pubdoc/configuration"
	"github.com/JackTimothy/pubdoc/confluence"
	"github.com/JackTimothy/pubdoc/parser"
	"github.com/joho/godotenv"
)

func main() {
	var markdownFilePath, title, htmlPageContent string
	var config configuration.Configuration
	var err error

	// parse flags
	{
		flag.StringVar(&markdownFilePath, "markdownFilePath", "", "The path to the Markdown file you want to publish.")
		flag.Parse()
	}

	// parse title
	{
		title = filepath.Base(markdownFilePath)
		title, err = confluence.FormatBase(title)
		if err != nil {
			log.Fatalf("Error formatting title: %v.", err)
		}
	}

	// load .env and configure confluence api usage
	{
		if err = godotenv.Load(); err != nil {
			log.Fatalf("Error loading .env file: %v.", err)
		}
		config = configuration.Init()
	}

	// generate html from local markdown files
	{
		htmlPageContent, err = parser.ConvertMarkdownToHTML(markdownFilePath)
		if err != nil {
			log.Fatalf("Error reading Markdown file: %v.\n", err)
		}
	}

	// Actually interface with confluence a bit
	// Will create/update the page passed by cmd line args
	{
		newPageId, err := confluence.CreatePage(title, htmlPageContent, "", config)
		if err != nil {
			if strings.Contains(err.Error(), "title already exists") {
				// attempt to update the page to the newly generated htmlPageContent
				updateRespOk, err := confluence.UpdatePage(title, htmlPageContent, config)
				if err != nil {
					log.Fatalf("Error attempting to update page with title: %s. err: %v\n", title, err)
				}
				log.Printf("Successfully updated page with ID %s.\n", updateRespOk.ID)
			} else {
				log.Fatalf("Error creating page: %v\n", err)
			}
		} else {
			log.Printf("Successfully created new page with ID %s.\n", newPageId)
		}
	}
}

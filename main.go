package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/russross/blackfriday/v2"
)

// Send this in a POST to https://{your-domain}/wiki/api/v2/pages to create a new page.
type CreatePageRequestBody struct {
	SpaceID  string `json:"spaceId"`
	Status   string `json:"status"`
	Title    string `json:"title"`
	ParentID string `json:"parentId"`
	Body     struct {
		Representation string `json:"representation"`
		Value          string `json:"value"`
	} `json:"body"`
}

// The body of a '200 OK' response to a Create Page POST for the Confluence v2 API.
// This struct does not actually capture all fields in the full response schema,
// which can be found documented here: https://developer.atlassian.com/cloud/confluence/rest/v2/api-group-page/#api-pages-post-response.
// Only the fields that this program uses are included. Feel free to add more fields as they become used.
type CreatePageResponseBodyStatusOk struct {
	Id    string `json:"id"`    // ID of the page.
	Title string `json:"title"` // Title of the page.
}

type authenticationCredentials struct {
	username string
	apiKey   string
}

func createPage(spaceID, title, htmlPageContent, parentID string, credentials authenticationCredentials) (newPageId string, err error) {
	url := fmt.Sprintf("https://%s/wiki/api/v2/pages", readConfluenceDomainFromEnvironmentVariable())

	payload := CreatePageRequestBody{
		SpaceID:  spaceID,
		Status:   "current",
		Title:    title,
		ParentID: parentID,
	}
	payload.Body.Representation = "storage"
	payload.Body.Value = htmlPageContent

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(credentials.username, credentials.apiKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error response code: %d; body: %s", resp.StatusCode, string(bodyBytes))
	}

	var response CreatePageResponseBodyStatusOk
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response.Id, nil
}

func readAuthenticationCredentialsFromEnvironmentVariables() authenticationCredentials {
	return authenticationCredentials{
		username: os.Getenv("CONFLUENCE_USERNAME"),
		apiKey:   os.Getenv("CONFLUENCE_API_KEY"),
	}
}

func readConfluenceDomainFromEnvironmentVariable() string {
	return os.Getenv("CONFLUENCE_DOMAIN")
}

func convertMarkdownToHTML(markdownFilePath string) (html string, err error) {
	file, err := os.Open(markdownFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to open markdown file: %w", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read markdown file: %w", err)
	}

	htmlContent := blackfriday.Run(content)
	return string(htmlContent), nil
}

func main() {
	var spaceID, markdownFilePath, title string
	flag.StringVar(&spaceID, "spaceID", "", "The ID of the Confluence space where the page will be created.")
	flag.StringVar(&markdownFilePath, "markdownFilePath", "./", "The path to the Markdown file you want to publish.")
	flag.StringVar(&title, "title", "", "What you wish to title the generated Confluence page.")
	flag.Parse()

	htmlPageContent, err := convertMarkdownToHTML(markdownFilePath)
	if err != nil {
		log.Fatalf("Error reading Markdown file: %v.\n", err)
	}

	credentials := readAuthenticationCredentialsFromEnvironmentVariables()

	newPageId, err := createPage(spaceID, title, htmlPageContent, "", credentials)
	if err != nil {
		log.Fatalf("Error creating page: %v\n", err)
	}
	log.Printf("Successfully created new page with ID %s.\n", newPageId)
}

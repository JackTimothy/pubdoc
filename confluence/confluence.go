package confluence

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/JackTimothy/pubdoc/configuration"
	Bodies "github.com/JackTimothy/pubdoc/confluence/types/bodies"
)

// Naive path base formatter. Should a user want something more 
// sophisticated they are free to create their own title.
func FormatBase(base string) (string, error) {
	ext := path.Ext(base)
	fmtBase, ok := strings.CutSuffix(base, ext)
	if !ok {
		return "", fmt.Errorf("could not find extension suffix")
	}
	return fmtBase, nil
}

// Options for the GetPagesInSpace() function
type GetPagesInSpaceOpts struct {
	Title      string
	Limit      int
	BodyFormat string
}

// Will delete a page with given title. 
// Should there be multiple pages with the title nothing occurs.
// Should there not be a page with given title nothing occurs.
func DeletePage(title string, config configuration.Configuration) error {
	gotPages, err := GetPagesInSpace(GetPagesInSpaceOpts{Title: title, Limit: 10}, config)
	if err != nil {
		log.Printf("Warning: Failed to Delete Page due to inablility to get page in space.")
		return nil
	}

	if len(gotPages.Results) > 1 {
		log.Printf("Warning: Failed to Delete Page due to multiple pages with title: %s. Ambiguous deletion avoided.", title)
		return nil
	}

	url := fmt.Sprintf("https://%s/wiki/api/v2/pages/%s", config.Domain, gotPages.Results[0].ID)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(config.Auth.Username, config.Auth.ApiKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("error response code: %d; body: %s", resp.StatusCode, string(bodyBytes))
	}

	// DeletePage returns no response content on success
	// no need to parse response
	return nil
}

// Will create a page with given title containng htmlPageContent.
// If the page already exists an error will be returned.
func CreatePage(title, htmlPageContent, parentID string, config configuration.Configuration) (newPageId string, err error) {
	url := fmt.Sprintf("https://%s/wiki/api/v2/pages", config.Domain)

	payload := Bodies.CreatePageRequestBody{
		SpaceID:  string(config.SpaceID),
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

	req.SetBasicAuth(config.Auth.Username, config.Auth.ApiKey)
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

	var response Bodies.CreatePageResponseBodyStatusOk
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response.Id, nil
}

// Gets pages in a given space based on the passed in opts. 
// See GetPagesInSpaceOpts for query options
func GetPagesInSpace(opts GetPagesInSpaceOpts, config configuration.Configuration) (got *Bodies.GetpagesInSpaceResponseBodyStatusOk, err error) {
	url := fmt.Sprintf("https://%s/wiki/api/v2/spaces/%s/pages", config.Domain, config.SpaceID)

	// Construct query parameters conditionally
	query := make([]string, 0)
	if opts.Title != "" {
		query = append(query, fmt.Sprintf("title=%s", opts.Title))
	}
	if opts.Limit > 0 {
		query = append(query, fmt.Sprintf("limit=%d", opts.Limit))
	}
	if opts.BodyFormat != "" {
		query = append(query, fmt.Sprintf("body-format=%s", opts.BodyFormat))
	}

	// Append query parameters to URL if any
	if len(query) > 0 {
		url = fmt.Sprintf("%s?%s", url, strings.Join(query, "&"))
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(config.Auth.Username, config.Auth.ApiKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response code: %d; body: %s", resp.StatusCode, string(bodyBytes))
	}

	var response Bodies.GetpagesInSpaceResponseBodyStatusOk
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// Will update the page with given pageTitle to contain htmlPageContent
// Does nothing if the page does not exist
// Does nothing if multiple pages exist with given title.
func UpdatePage(pageTitle, htmlPageContent string, config configuration.Configuration) (respOK *Bodies.UpdatePageResponseBodyStatusOk, err error) {
	respBody, err := GetPagesInSpace(GetPagesInSpaceOpts{Title: pageTitle, Limit: 10}, config)
	if err != nil {
		return nil, fmt.Errorf("error attempting to update page: %v", err)
	}

	if len(respBody.Results) == 0 {
		log.Printf("Warning: Attempted to update a page that did not exist. Title was: %s.", pageTitle)
		return nil, nil
	}

	if len(respBody.Results) > 1 {
		log.Printf("Warning: Attempted to update a page that exists in multiple locations. Update prevented. Title was: %s.", pageTitle)
		return nil, nil
	}

	pageID := respBody.Results[0].ID
	versionIncrement := respBody.Results[0].Version.Number + 1
	url := fmt.Sprintf("https://%s/wiki/api/v2/pages/%s", config.Domain, pageID)

	payload := Bodies.UpdatePageRequestBody{
		ID:     pageID,
		Status: "current", // NOTE: tbf perhaps we should always copy the current status. Thus preventing _revealing_ drafts.
		Title:  pageTitle,
		Version: struct {
			Number  string "json:\"number\""
			Message string "json:\"message\""
		}{Number: strconv.Itoa(versionIncrement), Message: "PubDoc automated update"},
	}
	payload.Body.Representation = "storage"
	payload.Body.Value = htmlPageContent

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(config.Auth.Username, config.Auth.ApiKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response code: %d; body: %s", resp.StatusCode, string(bodyBytes))
	}

	if err := json.Unmarshal(bodyBytes, &respOK); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return respOK, nil
}

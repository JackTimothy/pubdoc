package confluence

import (
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/JackTimothy/pubdoc/configuration"
	"github.com/JackTimothy/pubdoc/parser"
	"github.com/joho/godotenv"
)

func TestFormatBase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "valid file with extension",
			input:    "example-file_name.txt",
			expected: "example-file_name",
			wantErr:  false,
		},
		{
			name:     "file without extension",
			input:    "example-file_name",
			expected: "example-file_name",
			wantErr:  false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FormatBase(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FormatBase() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("FormatBase() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestConfluenceAPIGetPagesInSpace(t *testing.T) {
	tests := []struct {
		name           string
		opts           GetPagesInSpaceOpts
		expectedTitles []string
		wantErr        bool
	}{
		{
			name: "page exists",
			opts: GetPagesInSpaceOpts{
				Title: "README",
				Limit: 10,
			},
			expectedTitles: []string{"README"},
			wantErr:        false,
		},
		{
			name: "page does not exist",
			opts: GetPagesInSpaceOpts{
				Title: "READMEEEEEEEE",
				Limit: 10,
			},
			expectedTitles: nil,
			wantErr:        false,
		},
	}

	if err := godotenv.Load("../.env"); err != nil {
		t.Errorf("godotenv could not load .env")
		return
	}

	config := configuration.Init()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetPagesInSpace(tt.opts, config)
			if (err != nil) != tt.wantErr {
				t.Errorf("FormatBase() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got.Results) != len(tt.expectedTitles) && !tt.wantErr {
				t.Errorf("GetPagesInSpace() did not return expectedlen of %d. Len was instead %d.", len(tt.expectedTitles), len(got.Results))
				return
			}
			sort.Strings(tt.expectedTitles)
			var resultTitles []string
			for r := range got.Results {
				resultTitles = append(resultTitles, got.Results[r].Title)
			}
			sort.Strings(resultTitles)
			if !reflect.DeepEqual(tt.expectedTitles, resultTitles) {
				t.Errorf("Slices were not equal. Expected: %s, Response contained: %s", tt.expectedTitles, resultTitles)
				return
			}
		})
	}
}

func TestConfluenceAPIUpdatePage(t *testing.T) {
	if err := godotenv.Load("../.env"); err != nil {
		t.Errorf("godotenv could not load .env")
		return
	}

	config := configuration.Init()
	parsed := parser.MarkdownToHTML([]byte("# 0"))

	_, err := CreatePage("TESTUPDATEPAGE", string(parsed), "", config)
	if err != nil {
		t.Errorf("Failed to create test page. err: %v", err)
		return
	}
	
	time.Sleep(1 * time.Second)

	gotPages, err := GetPagesInSpace(GetPagesInSpaceOpts{Title: "TESTUPDATEPAGE", Limit: 1, BodyFormat: "storage"}, config)
	if err != nil {
		t.Errorf("Error getting pages in space: %v", err)
		return
	}

	tests := []struct {
		name            string
		md string
		expectedHTML string
		wantErr         bool
	}{
		{
			name:            "one",
			md: "# 1",
			expectedHTML: "<h1>1</h1>",
		},
		{
			name:            "two",
			md: "# 2",
			expectedHTML: "<h1>2</h1>",
		},
		{
			name:            "three",
			md: "# 3",
			expectedHTML: "<h1>3</h1>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed = parser.MarkdownToHTML([]byte(tt.md))
			updateResponse, err := UpdatePage(gotPages.Results[0].Title, string(parsed), config)
			if (err != nil) != tt.wantErr {
				t.Errorf("Error getting pages in space: %v", err)
				return
			}
			// does not exactly match parsed. Looks like confluence is removing the ids.
			if updateResponse.Body.Storage.Value != tt.expectedHTML {
				t.Errorf("Updated page does not match expectedContent. Page is: %s. Expected: %s.", updateResponse.Body.Storage.Value, parsed)
			}
		})
	}

	DeletePage("TESTUPDATEPAGE", config)
}

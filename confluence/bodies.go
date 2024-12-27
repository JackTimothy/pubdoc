package confluence

// Send this in a POST to https://{your-domain}/wiki/api/v2/pages to create a new page.
type CreatePageRequestBody struct {
	SpaceID  string `json:"spaceId"`
	Status   string `json:"status,omitempty"`
	Title    string `json:"title,omitempty"`
	ParentID string `json:"parentId,omitempty"`
	Body     struct {
		Representation string `json:"representation"`
		Value          string `json:"value"`
	} `json:"body,omitempty"`
}

// Send this in a PUT to https://{your-domain}/wiki/api/v2/pages/{id} to update the page with that id.
type UpdatePageRequestBody struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	Title    string `json:"title"`
	SpaceID  string `json:"spaceId,omitempty"`
	ParentID string `json:"parentId,omitempty"`
	OwnerID  string `json:"ownerId,omitempty"`
	Body     struct {
		Representation string `json:"representation"`
		Value          string `json:"value"`
	} `json:"body"`
	Version struct {
		Number  string `json:"number"`
		Message string `json:"message"`
	} `json:"version"`
}

// Send this in a GET to https://{your-domain}/wiki/api/v2/spaces/{id}/pages to retreive all pages from
// space with id included in path. NOTE: You can refine the search using the optional (omitempty) fields
// included below.
type GetPagesInSpaceRequestBody struct {
	ID    string `json:"id"`
	Title string `json:"title,omitempty"`
	Limit int    `json:"limit,omitempty"`
}

// The body of a '200 OK' response to a Create Page POST for the Confluence v2 API.
// This struct does not actually capture all fields in the full response schema,
// which can be found documented here: https://developer.atlassian.com/cloud/confluence/rest/v2/api-group-page/#api-pages-post-response.
// Only the fields that this program uses are included. Feel free to add more fields as they become used.
type CreatePageResponseBodyStatusOk struct {
	Id    string `json:"id"`    // ID of the page.
	Title string `json:"title"` // Title of the page.
}

// The body of a '200 OK' response to a GET pages In Space for the Confluence v2 API.
type GetpagesInSpaceResponseBodyStatusOk struct {
	Results []struct {
		ID          string `json:"id"`
		Status      string `json:"status"`
		Title       string `json:"title"`
		SpaceID     string `json:"spaceId"`
		ParentID    string `json:"parentId"`
		ParentType  string `json:"parentType"`
		Position    int    `json:"position"`
		AuthorID    string `json:"authorId"`
		OwnerID     string `json:"ownerId"`
		LastOwnerID string `json:"lastOwnerId"`
		CreatedAt   string `json:"createdAt"`
		Version     struct {
			CreatedAt string `json:"createdAt"`
			Message   string `json:"message"`
			Number    int    `json:"number"`
			MinorEdit bool   `json:"minorEdit"`
			AuthorID  string `json:"authorId"`
		} `json:"version"`
		Body struct {
			Storage struct {
				Representation string `json:"representation"`
				Value          string `json:"value"`
			} `json:"storage"`
			AtlasDocFormat map[string]interface{} `json:"atlas_doc_format"`
		} `json:"body"`
		Links struct {
			WebUI  string `json:"webui"`
			EditUI string `json:"editui"`
			TinyUI string `json:"tinyui"`
		} `json:"_links"`
	} `json:"results"`
	Links struct {
		Next string `json:"next"`
		Base string `json:"base"`
	} `json:"_links"`
}

// The body of a '200 OK' response to a PUT for updating a page for the Confluence v2 API.
type UpdatePageResponseBodyStatusOk struct {
	ID          string `json:"id"`
	Status      string `json:"status"`
	Title       string `json:"title"`
	SpaceID     string `json:"spaceId"`
	ParentID    string `json:"parentId"`
	ParentType  string `json:"parentType"`
	Position    int    `json:"position"`
	AuthorID    string `json:"authorId"`
	OwnerID     string `json:"ownerId"`
	LastOwnerID string `json:"lastOwnerId"`
	CreatedAt   string `json:"createdAt"`
	Version     struct {
		CreatedAt string `json:"createdAt"`
		Message   string `json:"message"`
		Number    int    `json:"number"`
		MinorEdit bool   `json:"minorEdit"`
		AuthorID  string `json:"authorId"`
	} `json:"version"`
	Body struct {
		Storage struct {
			Representation string `json:"representation"`
			Value          string `json:"value"`
		} `json:"storage"`
		AtlasDocFormat map[string]interface{} `json:"atlas_doc_format"`
		View           map[string]interface{} `json:"view"`
	} `json:"body"`
	Labels struct {
		Results []struct {
			ID     string `json:"id"`
			Name   string `json:"name"`
			Prefix string `json:"prefix"`
		} `json:"results"`
		Meta struct {
			HasMore bool   `json:"hasMore"`
			Cursor  string `json:"cursor"`
		} `json:"meta"`
		Links struct {
			Self string `json:"self"`
		} `json:"_links"`
	} `json:"labels"`
	Properties struct {
		Results []struct {
			ID      string `json:"id"`
			Key     string `json:"key"`
			Version struct {
				Number  string `json:"number"`
				Message string `json:"message"`
			} `json:"version"`
		} `json:"results"`
		Meta struct {
			HasMore bool   `json:"hasMore"`
			Cursor  string `json:"cursor"`
		} `json:"meta"`
		Links struct {
			Self string `json:"self"`
		} `json:"_links"`
	} `json:"properties"`
	Operations struct {
		Results []struct {
			Operation  string `json:"operation"`
			TargetType string `json:"targetType"`
		} `json:"results"`
		Meta struct {
			HasMore bool   `json:"hasMore"`
			Cursor  string `json:"cursor"`
		} `json:"meta"`
		Links struct {
			Self string `json:"self"`
		} `json:"_links"`
	} `json:"operations"`
	Likes struct {
		Results []struct {
			AccountID string `json:"accountId"`
		} `json:"results"`
		Meta struct {
			HasMore bool   `json:"hasMore"`
			Cursor  string `json:"cursor"`
		} `json:"meta"`
		Links struct {
			Self string `json:"self"`
		} `json:"_links"`
	} `json:"likes"`
	Versions struct {
		Results []struct {
			CreatedAt string `json:"createdAt"`
			Message   string `json:"message"`
			Number    int    `json:"number"`
			MinorEdit bool   `json:"minorEdit"`
			AuthorID  string `json:"authorId"`
		} `json:"results"`
		Meta struct {
			HasMore bool   `json:"hasMore"`
			Cursor  string `json:"cursor"`
		} `json:"meta"`
		Links struct {
			Self string `json:"self"`
		} `json:"_links"`
	} `json:"versions"`
	IsFavoritedByCurrentUser bool `json:"isFavoritedByCurrentUser"`
	Links                    struct {
		Base string `json:"base,omitempty"`
	} `json:"_links"`
}

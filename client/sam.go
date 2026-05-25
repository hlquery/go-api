package hlquery

import "fmt"

// SAM provides semantic access memory endpoints.
type SAM struct {
	client *Client
}

func (s *SAM) Search(collection, query string, params map[string]interface{}) (*Response, error) {
	if err := requireNonEmpty(query, "query"); err != nil {
		return nil, err
	}
	queryParams := copyMap(params)
	queryParams["q"] = query
	if collection != "" {
		queryParams["collection"] = collection
	}
	return s.client.requestWithQueryValues("GET", "/sam/search", nil, queryParams)
}

func (s *SAM) SearchAll(query string, params map[string]interface{}) (*Response, error) {
	queryParams := copyMap(params)
	queryParams["all"] = true
	return s.Search("", query, queryParams)
}

func (s *SAM) Rebuild(collection string, params map[string]interface{}) (*Response, error) {
	queryParams := copyMap(params)
	if collection != "" {
		queryParams["collection"] = collection
	}
	return s.client.requestWithQueryValues("POST", "/sam/rebuild", nil, queryParams)
}

func (s *SAM) Status(collection string, params ...map[string]interface{}) (*Response, error) {
	queryParams := map[string]interface{}{}
	if len(params) > 0 {
		queryParams = copyMap(params[0])
	}
	if collection != "" {
		queryParams["collection"] = collection
	}
	return s.client.requestWithQueryValues("GET", "/sam/status", nil, queryParams)
}

func (s *SAM) Debug(collection string, params map[string]interface{}) (*Response, error) {
	queryParams := copyMap(params)
	if collection != "" {
		queryParams["collection"] = collection
	}
	return s.client.requestWithQueryValues("GET", "/sam/debug", nil, queryParams)
}

func (s *SAM) History(collection string, limit int, params ...map[string]interface{}) (*Response, error) {
	queryParams := map[string]interface{}{"limit": limit}
	if len(params) > 0 {
		for key, value := range params[0] {
			queryParams[key] = value
		}
	}
	if collection != "" {
		queryParams["collection"] = collection
	}
	return s.client.requestWithQueryValues("GET", "/sam/history", nil, queryParams)
}

func (s *SAM) Pause(pauseUntilMs int64, params map[string]interface{}) (*Response, error) {
	queryParams := copyMap(params)
	queryParams["pause"] = pauseUntilMs
	return s.client.requestWithQueryValues("POST", "/sam/pause", nil, queryParams)
}

func (s *SAM) ClearPause(params map[string]interface{}) (*Response, error) {
	return s.Pause(0, params)
}

func (s *SAM) Improve(params map[string]interface{}) (*Response, error) {
	return s.client.requestWithQueryValues("POST", "/sam/improve", nil, params)
}

func (s *SAM) FlushActorMetadata() (*Response, error) {
	return s.client.request("POST", "/sam/flush_actor_metadata", nil)
}

func (s *SAM) ListDocuments(collection string, offset, limit int, params map[string]interface{}) (*Response, error) {
	if err := requireNonEmpty(collection, "collection name"); err != nil {
		return nil, err
	}
	queryParams := copyMap(params)
	queryParams["collection"] = collection
	queryParams["offset"] = offset
	queryParams["limit"] = limit
	return s.client.requestWithQueryValues("GET", "/sam/documents", nil, queryParams)
}

func (s *SAM) GetDocument(collection, docID string, params map[string]interface{}) (*Response, error) {
	if err := requireNonEmpty(collection, "collection name"); err != nil {
		return nil, err
	}
	if err := requireNonEmpty(docID, "document id"); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/sam/documents/%s/%s", encodePathPart(collection), encodePathPart(docID))
	return s.client.requestWithQueryValues("GET", path, nil, params)
}

func (s *SAM) AddDocumentLabel(collection, docID string, label map[string]interface{}) (*Response, error) {
	if err := requireNonEmpty(collection, "collection name"); err != nil {
		return nil, err
	}
	if err := requireNonEmpty(docID, "document id"); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/sam/label/add/%s/%s", encodePathPart(collection), encodePathPart(docID))
	return s.client.request("POST", path, label)
}

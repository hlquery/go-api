package hlquery

import "fmt"

type Aliases struct {
	client *Client
}

func (a *Aliases) List() (*Response, error) {
	return a.client.request("GET", "/aliases", nil)
}

func (a *Aliases) Get(name string) (*Response, error) {
	if err := requireNonEmpty(name, "alias name"); err != nil {
		return nil, err
	}
	return a.client.request("GET", "/aliases/"+encodePathPart(name), nil)
}

func (a *Aliases) Create(name string, params map[string]interface{}) (*Response, error) {
	if err := requireNonEmpty(name, "alias name"); err != nil {
		return nil, err
	}
	return a.client.request("POST", "/aliases/"+encodePathPart(name), params)
}

func (a *Aliases) Update(name string, params map[string]interface{}) (*Response, error) {
	if err := requireNonEmpty(name, "alias name"); err != nil {
		return nil, err
	}
	return a.client.request("PUT", "/aliases/"+encodePathPart(name), params)
}

func (a *Aliases) Upsert(name string, params map[string]interface{}) (*Response, error) {
	return a.Create(name, params)
}

func (a *Aliases) Delete(name string) (*Response, error) {
	if err := requireNonEmpty(name, "alias name"); err != nil {
		return nil, err
	}
	return a.client.request("DELETE", "/aliases/"+encodePathPart(name), nil)
}

type Synonyms struct {
	client *Client
}

func (s *Synonyms) List(collection string) (*Response, error) {
	if err := requireNonEmpty(collection, "collection name"); err != nil {
		return nil, err
	}
	return s.client.request("GET", fmt.Sprintf("/collections/%s/synonyms", encodePathPart(collection)), nil)
}

func (s *Synonyms) Get(collection, id string) (*Response, error) {
	if err := requireNonEmpty(collection, "collection name"); err != nil {
		return nil, err
	}
	if err := requireNonEmpty(id, "synonym id"); err != nil {
		return nil, err
	}
	return s.client.request("GET", fmt.Sprintf("/collections/%s/synonyms/%s", encodePathPart(collection), encodePathPart(id)), nil)
}

func (s *Synonyms) Create(collection, id string, synonym map[string]interface{}) (*Response, error) {
	if err := requireNonEmpty(collection, "collection name"); err != nil {
		return nil, err
	}
	if err := requireNonEmpty(id, "synonym id"); err != nil {
		return nil, err
	}
	return s.client.request("POST", fmt.Sprintf("/collections/%s/synonyms/%s", encodePathPart(collection), encodePathPart(id)), synonym)
}

func (s *Synonyms) Update(collection, id string, synonym map[string]interface{}) (*Response, error) {
	if err := requireNonEmpty(collection, "collection name"); err != nil {
		return nil, err
	}
	if err := requireNonEmpty(id, "synonym id"); err != nil {
		return nil, err
	}
	return s.client.request("PUT", fmt.Sprintf("/collections/%s/synonyms/%s", encodePathPart(collection), encodePathPart(id)), synonym)
}

func (s *Synonyms) Upsert(collection, id string, synonym map[string]interface{}) (*Response, error) {
	return s.Create(collection, id, synonym)
}

func (s *Synonyms) Delete(collection, id string) (*Response, error) {
	if err := requireNonEmpty(collection, "collection name"); err != nil {
		return nil, err
	}
	if err := requireNonEmpty(id, "synonym id"); err != nil {
		return nil, err
	}
	return s.client.request("DELETE", fmt.Sprintf("/collections/%s/synonyms/%s", encodePathPart(collection), encodePathPart(id)), nil)
}

func (s *Synonyms) ListAll() (*Response, error) {
	return s.client.request("GET", "/synonyms", nil)
}

func (s *Synonyms) ListGlobal() (*Response, error) {
	return s.client.request("GET", "/synonyms/global", nil)
}

func (s *Synonyms) GetGlobal(id string) (*Response, error) {
	if err := requireNonEmpty(id, "synonym id"); err != nil {
		return nil, err
	}
	return s.client.request("GET", "/synonyms/global/"+encodePathPart(id), nil)
}

func (s *Synonyms) CreateGlobal(id string, synonym map[string]interface{}) (*Response, error) {
	if err := requireNonEmpty(id, "synonym id"); err != nil {
		return nil, err
	}
	return s.client.request("POST", "/synonyms/global/"+encodePathPart(id), synonym)
}

func (s *Synonyms) UpdateGlobal(id string, synonym map[string]interface{}) (*Response, error) {
	if err := requireNonEmpty(id, "synonym id"); err != nil {
		return nil, err
	}
	return s.client.request("PUT", "/synonyms/global/"+encodePathPart(id), synonym)
}

func (s *Synonyms) UpsertGlobal(id string, synonym map[string]interface{}) (*Response, error) {
	return s.CreateGlobal(id, synonym)
}

func (s *Synonyms) DeleteGlobal(id string) (*Response, error) {
	if err := requireNonEmpty(id, "synonym id"); err != nil {
		return nil, err
	}
	return s.client.request("DELETE", "/synonyms/global/"+encodePathPart(id), nil)
}

type Stopwords struct {
	client *Client
}

func (s *Stopwords) List(collection string) (*Response, error) {
	if err := requireNonEmpty(collection, "collection name"); err != nil {
		return nil, err
	}
	return s.client.request("GET", fmt.Sprintf("/collections/%s/stopwords", encodePathPart(collection)), nil)
}

func (s *Stopwords) Create(collection string, params map[string]interface{}) (*Response, error) {
	if err := requireNonEmpty(collection, "collection name"); err != nil {
		return nil, err
	}
	return s.client.request("POST", fmt.Sprintf("/collections/%s/stopwords", encodePathPart(collection)), params)
}

func (s *Stopwords) Delete(collection, word string) (*Response, error) {
	if err := requireNonEmpty(collection, "collection name"); err != nil {
		return nil, err
	}
	if err := requireNonEmpty(word, "stopword"); err != nil {
		return nil, err
	}
	return s.client.request("DELETE", fmt.Sprintf("/collections/%s/stopwords/%s", encodePathPart(collection), encodePathPart(word)), nil)
}

func (s *Stopwords) ListAll() (*Response, error) {
	return s.client.request("GET", "/stopwords", nil)
}

func (s *Stopwords) ListGlobal() (*Response, error) {
	return s.client.request("GET", "/stopwords/global", nil)
}

func (s *Stopwords) CreateGlobal(params map[string]interface{}) (*Response, error) {
	return s.client.request("POST", "/stopwords/global", params)
}

func (s *Stopwords) DeleteGlobal(word string) (*Response, error) {
	if err := requireNonEmpty(word, "stopword"); err != nil {
		return nil, err
	}
	return s.client.request("DELETE", "/stopwords/global/"+encodePathPart(word), nil)
}

type Overrides struct {
	client *Client
}

func (o *Overrides) List(collection string) (*Response, error) {
	if err := requireNonEmpty(collection, "collection name"); err != nil {
		return nil, err
	}
	return o.client.request("GET", fmt.Sprintf("/collections/%s/overrides", encodePathPart(collection)), nil)
}

func (o *Overrides) Get(collection, id string) (*Response, error) {
	if err := requireNonEmpty(collection, "collection name"); err != nil {
		return nil, err
	}
	if err := requireNonEmpty(id, "override id"); err != nil {
		return nil, err
	}
	return o.client.request("GET", fmt.Sprintf("/collections/%s/overrides/%s", encodePathPart(collection), encodePathPart(id)), nil)
}

func (o *Overrides) Create(collection, id string, override map[string]interface{}) (*Response, error) {
	if err := requireNonEmpty(collection, "collection name"); err != nil {
		return nil, err
	}
	if err := requireNonEmpty(id, "override id"); err != nil {
		return nil, err
	}
	return o.client.request("POST", fmt.Sprintf("/collections/%s/overrides/%s", encodePathPart(collection), encodePathPart(id)), override)
}

func (o *Overrides) Update(collection, id string, override map[string]interface{}) (*Response, error) {
	if err := requireNonEmpty(collection, "collection name"); err != nil {
		return nil, err
	}
	if err := requireNonEmpty(id, "override id"); err != nil {
		return nil, err
	}
	return o.client.request("PUT", fmt.Sprintf("/collections/%s/overrides/%s", encodePathPart(collection), encodePathPart(id)), override)
}

func (o *Overrides) Upsert(collection, id string, override map[string]interface{}) (*Response, error) {
	return o.Create(collection, id, override)
}

func (o *Overrides) Delete(collection, id string) (*Response, error) {
	if err := requireNonEmpty(collection, "collection name"); err != nil {
		return nil, err
	}
	if err := requireNonEmpty(id, "override id"); err != nil {
		return nil, err
	}
	return o.client.request("DELETE", fmt.Sprintf("/collections/%s/overrides/%s", encodePathPart(collection), encodePathPart(id)), nil)
}

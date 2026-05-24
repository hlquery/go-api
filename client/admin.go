package hlquery

type Keys struct {
	client *Client
}

func (k *Keys) List() (*Response, error) {
	return k.client.request("GET", "/keys", nil)
}

func (k *Keys) Get(id string) (*Response, error) {
	return k.client.request("GET", "/keys/"+encodePathPart(id), nil)
}

func (k *Keys) Create(params map[string]interface{}) (*Response, error) {
	return k.client.request("POST", "/keys", params)
}

func (k *Keys) Update(id string, params map[string]interface{}) (*Response, error) {
	return k.client.request("PUT", "/keys/"+encodePathPart(id), params)
}

func (k *Keys) Delete(id string) (*Response, error) {
	return k.client.request("DELETE", "/keys/"+encodePathPart(id), nil)
}

type Users struct {
	client *Client
}

func (u *Users) List() (*Response, error) {
	return u.client.request("GET", "/users", nil)
}

func (u *Users) Get(id string) (*Response, error) {
	return u.client.request("GET", "/users/"+encodePathPart(id), nil)
}

func (u *Users) Create(params map[string]interface{}) (*Response, error) {
	return u.client.request("POST", "/users", params)
}

func (u *Users) Update(id string, params map[string]interface{}) (*Response, error) {
	return u.client.request("PUT", "/users/"+encodePathPart(id), params)
}

func (u *Users) Delete(id string) (*Response, error) {
	return u.client.request("DELETE", "/users/"+encodePathPart(id), nil)
}

type Modules struct {
	client *Client
}

func (m *Modules) List() (*Response, error) {
	return m.client.request("GET", "/modules", nil)
}

func (m *Modules) Syntax(name string) (*Response, error) {
	return m.client.request("GET", "/modules/"+encodePathPart(name)+"/syntax", nil)
}

func (m *Modules) Load(name string, params map[string]interface{}) (*Response, error) {
	return m.client.request("POST", "/loadmodule/"+encodePathPart(name), params)
}

func (m *Modules) Unload(name string, params map[string]interface{}) (*Response, error) {
	return m.client.request("POST", "/unloadmodule/"+encodePathPart(name), params)
}

func (m *Modules) Execute(method, path string, body interface{}, params map[string]interface{}) (*Response, error) {
	return m.client.requestWithQueryValues(method, path, body, params)
}

// AnalyticsClick records click analytics.
func (c *Client) AnalyticsClick(params map[string]interface{}) (*Response, error) {
	return c.request("POST", "/analytics/click", params)
}

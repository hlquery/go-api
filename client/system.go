package hlquery

// System provides system and operational endpoints.
type System struct {
	client *Client
}

func (s *System) Health() (*Response, error) { return s.client.Health() }
func (s *System) Status() (*Response, error) { return s.client.Status() }
func (s *System) Query() (*Response, error)  { return s.client.request("GET", "/query", nil) }
func (s *System) SearchConfig() (*Response, error) {
	return s.client.request("GET", "/search-config", nil)
}
func (s *System) ConfigFiles() (*Response, error) {
	return s.client.request("GET", "/config-files", nil)
}
func (s *System) Ready() (*Response, error)   { return s.client.Ready() }
func (s *System) Ping() (*Response, error)    { return s.client.Ping() }
func (s *System) Info() (*Response, error)    { return s.client.Info() }
func (s *System) Stats() (*Response, error)   { return s.client.Stats() }
func (s *System) Metrics() (*Response, error) { return s.client.request("GET", "/metrics", nil) }
func (s *System) MetricsJSON() (*Response, error) {
	return s.client.request("GET", "/metrics.json", nil)
}
func (s *System) MetricsHistory() (*Response, error) {
	return s.client.request("GET", "/metrics/history", nil)
}
func (s *System) Connections() (*Response, error) {
	return s.client.request("GET", "/connections", nil)
}
func (s *System) RocksDB() (*Response, error) { return s.client.request("GET", "/rocksdb", nil) }
func (s *System) RocksDBInternal() (*Response, error) {
	return s.client.request("GET", "/_rocksdb", nil)
}
func (s *System) DocTotal() (*Response, error)   { return s.client.request("GET", "/doctotal", nil) }
func (s *System) Etc() (*Response, error)        { return s.client.Etc() }
func (s *System) Startup() (*Response, error)    { return s.client.request("GET", "/startup", nil) }
func (s *System) BootStatus() (*Response, error) { return s.client.request("GET", "/boot-status", nil) }
func (s *System) Integrity() (*Response, error)  { return s.client.request("GET", "/integrity", nil) }
func (s *System) Consistency() (*Response, error) {
	return s.client.request("GET", "/consistency", nil)
}
func (s *System) SelfCheck() (*Response, error) { return s.client.request("GET", "/self-check", nil) }
func (s *System) StorageStatus() (*Response, error) {
	return s.client.request("GET", "/admin/storage_status", nil)
}
func (s *System) DebugCounters() (*Response, error) {
	return s.client.request("GET", "/debug/counters", nil)
}

// UpdateCounters refreshes server counters.
func (s *System) UpdateCounters(params ...map[string]interface{}) (*Response, error) {
	if len(params) > 0 {
		return s.client.request("POST", "/update-counters", params[0])
	}
	return s.client.request("GET", "/update-counters", nil)
}

// Repair runs or inspects repair state.
func (s *System) Repair(params ...map[string]interface{}) (*Response, error) {
	if len(params) > 0 {
		return s.client.request("POST", "/repair", params[0])
	}
	return s.client.request("GET", "/repair", nil)
}

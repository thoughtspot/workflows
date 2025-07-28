package deploykey

import (
	"encoding/json"
	"net/http"
	"time"

	"automation/common"
	"automation/logger"
)

type DeployKey struct {
	Key            string `json:"key"`
	Title          string `json:"title"`
	ReadOnly       bool   `json:"read_only"`
	RepositoryName string `json:"repository_name"`
}

func New(key, title, repositoryName string, readOnly bool) *DeployKey {
	return &DeployKey{
		Key:            key,
		Title:          title,
		ReadOnly:       readOnly,
		RepositoryName: repositoryName,
	}
}

type CreateDeloyKeysResponse struct {
	Key      string `json:"key"`
	Title    string `json:"title"`
	Verified bool   `json:"verified"`
	ReadOnly bool   `json:"read_only"`
	AddedBy  string `json:"added_by"`
	Enabled  bool   `json:"enabled"`
}

func (d *DeployKey) CreateDeployKey() {
	l := logger.New()
	reqBody := struct {
		Key      string `json:"key"`
		Title    string `json:"title"`
		ReadOnly bool   `json:"read_only"`
	}{
		Key:      d.Key,
		Title:    d.Title,
		ReadOnly: d.ReadOnly,
	}

	req, err := http.NewRequest(http.MethodPost, common.CreateDeployKeyEndpointURL(d.RepositoryName), common.RequestBody(reqBody))
	if err != nil {
		l.Fatal(err)
	}

	common.SetHeaders(req)

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		l.Fatal(err)
	}
	defer resp.Body.Close()

	var createDeployKeyResponse CreateDeloyKeysResponse
	if err := json.NewDecoder(resp.Body).Decode(&createDeployKeyResponse); err != nil {
		l.Fatal(err)
	}

	l.Println(resp.StatusCode)
	l.Println(createDeployKeyResponse)
}

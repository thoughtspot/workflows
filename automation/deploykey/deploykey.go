package deploykey

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"automation/common"
)

type DeployKey struct {
	Key            string `json:"key"`
	Title          string `json:"title"`
	ReadOnly       bool   `json:"read_only"`
	RepositoryName string `json:"repository_name"`
}

func (d *DeployKey) CreateDeployKey() {
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
		panic(err)
	}

	common.SetHeaders(req)

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))
}

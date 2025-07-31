package repo

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"automation/common"
	"automation/logger"
)

type OrgPublicKey struct {
	KeyID string `json:"key_id"`
	Key   string `json:"key"`
}

func GetPublicKey(repositoryName string) OrgPublicKey {
	l := logger.New()
	req, err := http.NewRequest(http.MethodGet, common.GetRepoPublicKeyEndpointURL(repositoryName), nil)
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

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Fatal(err)
	}

	var orgPublicKey OrgPublicKey
	if err := json.Unmarshal(bytes, &orgPublicKey); err != nil {
		l.Fatal(err)
	}

	l.Printf("API Response Status Code: %d\n", resp.StatusCode)
	return orgPublicKey
}

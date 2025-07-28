package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"automation/constants"
)

func SetHeaders(req *http.Request) {
	headers := map[string]string{
		"accept":              "application/vnd.github+json",
		"X-GitHub-Api-Versio": "2022-11-28",
		"Authorization":       "Bearer " + os.Getenv("AUTH_TOKEN"),
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}
}

func CreateDeployKeyEndpointURL(repositoryName string) string {
	// https://api.github.com/repos/OWNER/REPO/keys
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/keys", constants.OWNER, repositoryName)
}

func CreateRepositoryEndpointURL() string {
	// https://api.github.com/orgs/ORG/repos
	return fmt.Sprintf("https://api.github.com/orgs/%s/repos", constants.ORGANIZATION)
}

func GetOrgPublicKeyEndpointURL() string {
	// https://api.github.com/orgs/ORG/actions/secrets/public-key
	return fmt.Sprintf("https://api.github.com/orgs/%s/actions/secrets/public-key", constants.ORGANIZATION)
}

func CreateRepositorySecretEndpointURL(repositoryName string) string {
	// https://api.github.com/repos/OWNER/REPO/actions/secrets/SECRET_NAME
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/secrets/%s", constants.OWNER, repositoryName, constants.SECRET_NAME)
}

func RequestBody(data any) io.Reader {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(dataBytes)
}

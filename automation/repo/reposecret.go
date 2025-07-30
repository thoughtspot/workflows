package repo

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"automation/common"
	"automation/logger"
)

type RepositorySecret struct {
	RepositoryName  string `json:"repository_name"`
	SecretName      string `json:"secret_name"`
	EncryptedSecret string `json:"encrypted_value"`
	KeyID           string `json:"key_id"`
}

func NewRepositorySecret(repositoryName, secretName, encryptedSecert, keyID string) *RepositorySecret {
	return &RepositorySecret{
		RepositoryName:  repositoryName,
		SecretName:      secretName,
		EncryptedSecret: encryptedSecert,
		KeyID:           keyID,
	}
}

func (r *RepositorySecret) CreateSecret() {
	l := logger.New()
	bodyData := struct {
		EncryptedSecret string `json:"encrypted_value"`
		KeyID           string `json:"key_id"`
	}{
		EncryptedSecret: r.EncryptedSecret,
		KeyID:           r.KeyID,
	}

	type RequestBody struct {
		EncryptedSecret string `json:"encrypted_value"`
		KeyID           string `json:"key_id"`
	}

	var reqBody RequestBody
	reader := common.RequestBody(bodyData)
	data, _ := io.ReadAll(reader)
	if err := json.Unmarshal(data, &reqBody); err != nil {
		l.Fatal(err)
	}
	l.Println(reqBody)

	req, err := http.NewRequest(http.MethodPut, common.CreateRepositorySecretEndpointURL(r.RepositoryName), common.RequestBody(bodyData))
	if err != nil {
		l.Fatal(err)
	}

	common.SetHeaders(req)

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 201 {
		l.Fatalf("failed to create/update repository secret\n\terr: %v\n\tstatus code: %d\n ", err, resp.StatusCode)
	}
	defer resp.Body.Close()

	l.Println("Repository Secret Successfully")
}

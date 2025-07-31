package repo

import (
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

	l.Printf("API Response Status Code: %d\n", resp.StatusCode)
	l.Printf("Repository Secret Successfully\n")
}

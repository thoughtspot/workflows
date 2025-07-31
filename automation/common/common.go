package common

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"automation/constants"
	"automation/logger"

	"golang.org/x/crypto/nacl/box"
)

func SetHeaders(req *http.Request) {
	fmt.Println(os.Getenv("AUTH_TOKEN"))
	headers := map[string]string{
		"accept":               "application/vnd.github+json",
		"X-GitHub-Api-Version": "2022-11-28",
		"Authorization":        "Bearer " + os.Getenv("AUTH_TOKEN"),
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

func GetRepoPublicKeyEndpointURL(repositoryName string) string {
	// https://api.github.com/repos/OWNER/REPO/actions/secrets/public-key
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/secrets/public-key", constants.ORGANIZATION, repositoryName)
}

func CreateRepositorySecretEndpointURL(repositoryName string) string {
	// https://api.github.com/repos/OWNER/REPO/actions/secrets/SECRET_NAME
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/secrets/%s", constants.OWNER, repositoryName, constants.SECRET_NAME)
}

func RequestBody(data any) io.Reader {
	l := logger.New()
	dataBytes, err := json.Marshal(data)
	if err != nil {
		l.Fatal(err)
	}
	return bytes.NewReader(dataBytes)
}

func EncryptSecret(secretValue []byte, publicKeyB64 string) string {
	l := logger.New()
	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKeyB64)
	if err != nil {
		l.Fatal(err)
	}

	if len(publicKeyBytes) != 32 {
		l.Fatalf("wrong length of public key %d\n", len(publicKeyBytes))
	}

	var publicKey [32]byte
	copy(publicKey[:], publicKeyBytes)

	encryptedSecret, err := box.SealAnonymous(nil, secretValue, &publicKey, rand.Reader)
	if err != nil {
		l.Fatal(err)
	}

	return base64.StdEncoding.EncodeToString(encryptedSecret)
}

package reposecret

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"time"

	"automation/common"

	"golang.org/x/crypto/nacl/box"
	"golang.org/x/crypto/ssh"
)

type RepositorySecret struct {
	RepositoryName  string `json:"repository_name"`
	SecretName      string `json:"secret_name"`
	EncryptedSecret string `json:"encrypted_value"`
	KeyID           string `json:"key_id"`
}

func (r *RepositorySecret) CreateRepositorySecret() {
	bodyData := struct {
		EncryptedSecret string `json:"encrypted_value"`
		KeyID           string `json:"key_id"`
	}{
		EncryptedSecret: r.EncryptedSecret,
		KeyID:           r.KeyID,
	}

	req, err := http.NewRequest(http.MethodPut, common.CreateRepositoryEndpointURL(), common.RequestBody(bodyData))
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

func CreateSSHKeys() {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	privateKeyBlock, err := ssh.MarshalPrivateKey(privateKey, "")
	if err != nil {
		panic(err)
	}
	privateKeyBytes := pem.EncodeToMemory(privateKeyBlock)
	fmt.Println(string(privateKeyBytes))

	sshPublicKey, err := ssh.NewPublicKey(publicKey)
	if err != nil {
		panic(err)
	}

	opensshPublicKey := ssh.MarshalAuthorizedKey(sshPublicKey)
	fmt.Println(string(opensshPublicKey))
}

type OrgPublicKey struct {
	KeyID string `json:"key_id"`
	Key   string `json:"key"`
}

func GetOrgPublicKey() {
	req, err := http.NewRequest(http.MethodGet, common.GetOrgPublicKeyEndpointURL(), nil)
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

	var orgPublicKey OrgPublicKey
	if err := json.Unmarshal(bytes, &orgPublicKey); err != nil {
		panic(err)
	}

	fmt.Println(orgPublicKey)
}

func EncryptSecret(secretValue string, publicKeyB64 string) (string, error) {
	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKeyB64)
	if err != nil {
		return "", err
	}

	if len(publicKeyBytes) != 32 {
		return "", fmt.Errorf("invalid public key length: expected 32 bytes, got %d", len(publicKeyBytes))
	}

	var publicKey [32]byte
	copy(publicKey[:], publicKeyBytes)

	secretBytes := []byte(secretValue)
	encrypted, err := box.SealAnonymous(nil, secretBytes, &publicKey, rand.Reader)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(encrypted), nil
}

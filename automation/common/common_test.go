package common

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"

	"automation/constants"

	"golang.org/x/crypto/nacl/box"
)

func Test_SetHeaders(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "localhost:8080", nil)
	if err != nil {
		t.Errorf("create http request error: %v", err)
	}
	os.Setenv("AUTH_TOKEN", "auth-token-for-testing")

	tests := []struct {
		wantHeaders map[string]string
	}{
		{
			wantHeaders: map[string]string{
				"accept":               "application/vnd.github+json",
				"X-GitHub-Api-Version": "2022-11-28",
				"Authorization":        "Bearer auth-token-for-testing",
			},
		},
	}

	for _, test := range tests {
		SetHeaders(req)
		headers := []string{"accept", "X-GitHub-Api-Version", "Authorization"}
		gotHeaders := make(map[string]string)
		for _, header := range headers {
			gotHeaders[header] = req.Header.Get(header)
		}

		for _, header := range headers {
			if test.wantHeaders[header] != gotHeaders[header] {
				t.Errorf("\ngotHeaders: %v\nwantHeaders: %v\n", gotHeaders[header], test.wantHeaders[header])
			}
		}
	}
}

func Test_CreateDeployKeyEndpointURL(t *testing.T) {
	tests := []struct {
		owner          string
		repositoryName string
		wantURL        string
	}{
		{
			owner:          constants.OWNER,
			repositoryName: "my-repo",
			wantURL:        "https://api.github.com/repos/thoughtspot/my-repo/keys",
		},
		{
			owner:          constants.OWNER,
			repositoryName: "my-repo-private",
			wantURL:        "https://api.github.com/repos/thoughtspot/my-repo-private/keys",
		},
	}

	for _, test := range tests {
		gotURL := CreateDeployKeyEndpointURL(test.repositoryName)
		if gotURL != test.wantURL {
			t.Errorf("\ngotURL: %v\nwantURL: %v\n", gotURL, test.wantURL)
		}
	}
}

func Test_CreateRepositoryEndpointURL(t *testing.T) {
	tests := []struct {
		org     string
		wantURL string
	}{
		{
			org:     constants.ORGANIZATION,
			wantURL: "https://api.github.com/orgs/thoughtspot/repos",
		},
	}

	for _, test := range tests {
		gotURL := CreateRepositoryEndpointURL()
		if gotURL != test.wantURL {
			t.Errorf("\ngotURL: %v\nwantURL: %v", gotURL, test.wantURL)
		}
	}
}

func Test_GetRepoPublicKeyEndpointURL(t *testing.T) {
	tests := []struct {
		repositoryName string
		org            string
		wantURL        string
	}{
		{
			repositoryName: "my-repo-private",
			org:            constants.ORGANIZATION,
			wantURL:        "https://api.github.com/repos/thoughtspot/my-repo-private/actions/secrets/public-key",
		},
	}

	for _, test := range tests {
		gotURL := GetRepoPublicKeyEndpointURL(test.repositoryName)
		if gotURL != test.wantURL {
			t.Errorf("\ngotURL: %v\nwantURL: %v\n", gotURL, test.wantURL)
		}
	}
}

func Test_CreateRepositorySecretEndpointURL(t *testing.T) {
	tests := []struct {
		owner          string
		repositoryName string
		secretName     string
		wantURL        string
	}{
		{
			owner:          constants.OWNER,
			repositoryName: "my-repo-private",
			secretName:     constants.SECRET_NAME,
			wantURL:        "https://api.github.com/repos/thoughtspot/my-repo-private/actions/secrets/SSH_DEPLOY_KEY",
		},
	}

	for _, test := range tests {
		gotURL := CreateRepositorySecretEndpointURL(test.repositoryName)
		if gotURL != test.wantURL {
			t.Errorf("\ngotURL: %v\nwantURL: %v", gotURL, test.wantURL)
		}
	}
}

func Test_RequestBody(t *testing.T) {
	fn := func(data any) io.Reader {
		dataBytes, err := json.Marshal(data)
		if err != nil {
			t.Error(err)
		}
		return bytes.NewReader(dataBytes)
	}
	tests := []struct {
		data any
		want func(data any) io.Reader
	}{
		{
			data: struct {
				a string
				b int
			}{
				a: "string field",
				b: 32,
			},
			want: fn,
		},
		{
			data: struct {
				a string
				b string
			}{
				a: "string-a",
				b: "string-b",
			},
			want: fn,
		},
	}

	for _, test := range tests {
		got := RequestBody(test.data)
		gotBytes, err := io.ReadAll(got)
		if err != nil {
			t.Error(err)
		}

		want := test.want(test.data)
		wantBytes, err := io.ReadAll(want)
		if err != nil {
			t.Error(err)
		}

		if len(wantBytes) != len(gotBytes) {
			t.Errorf("\nwantBytesLen: %v\ngotBytesLen: %v\n", len(wantBytes), len(gotBytes))
		}

		for idx := 0; idx < len(wantBytes); idx++ {
			if wantBytes[idx] != gotBytes[idx] {
				t.Errorf("\nwantBytes: %v\ngotBytes: %v\n", wantBytes, gotBytes)
			}
		}

	}
}

func Test_EncryptSecret(t *testing.T) {
	publicKey, privateKey, err := box.GenerateKey(rand.Reader)
	if err != nil {
		t.Error(err)
	}

	tests := []struct {
		secretMsg    []byte
		publicKeyB64 string
		want         []byte
	}{
		{
			secretMsg:    []byte("super secret message"),
			publicKeyB64: base64.StdEncoding.EncodeToString(publicKey[:]),
			want:         []byte("super secret message"),
		},
	}

	for _, test := range tests {
		encryptedB64Msg := EncryptSecret(test.secretMsg, test.publicKeyB64)
		encryptedMsg, err := base64.StdEncoding.DecodeString(encryptedB64Msg)
		if err != nil {
			t.Error(err)
		}

		got, ok := box.OpenAnonymous(nil, []byte(encryptedMsg), publicKey, privateKey)
		if !ok {
			t.Error("cannot decrypt cypher text")
		}

		if len(got) != len(test.want) {
			t.Errorf("\nlenGotBytes: %v\nlenWantBytes: %v\n", len(got), len(test.want))
		}

		for idx := 0; idx < len(got); idx++ {
			if got[idx] != test.want[idx] {
				t.Errorf("\nGot: %v\nWant: %v\n", got, test.want)
			}
		}

	}
}

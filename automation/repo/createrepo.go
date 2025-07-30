package repo

import (
	"encoding/json"
	"net/http"
	"time"

	"automation/common"
	"automation/logger"
)

type Repository struct {
	Name                string `json:"name"`
	Visibility          string `json:"visibility"`
	DeleteBranchOnMerge bool   `json:"delete_branch_on_merge"`
	Org                 string `json:"org"`
}

func NewRepository(name, visibility, org string, deleteBranchOnMerge bool) *Repository {
	return &Repository{
		Name:                name,
		Visibility:          visibility,
		DeleteBranchOnMerge: deleteBranchOnMerge,
		Org:                 org,
	}
}

type CreateRepoResponse struct {
	Name string `json:"name"`
	URL  string `json:"html_url"`
}

func (r *Repository) Create() {
	l := logger.New()
	req, err := http.NewRequest(http.MethodPost, common.CreateRepositoryEndpointURL(), common.RequestBody(r))
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

	var createRepoResponse CreateRepoResponse
	if err := json.NewDecoder(resp.Body).Decode(&createRepoResponse); err != nil {
		l.Fatal(err)
	}

	l.Println(resp.StatusCode)
	l.Println(createRepoResponse)
}

package repo

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"automation/common"
)

type Repository struct {
	Name                string `json:"name"`
	Visibility          string `json:"visibility"`
	DeleteBranchOnMerge bool   `json:"delete_branch_on_merge"`
	Org                 string `json:"org"`
}

func (r *Repository) CreateRepo() {
	req, err := http.NewRequest(http.MethodPost, common.CreateRepositoryEndpointURL(), common.RequestBody(r))
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

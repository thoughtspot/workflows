package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"automation/common"
	"automation/logger"
	"automation/orgpublickey"
)

type UserInput struct {
	AuthToken             string `json:"auth_token"`
	PrivateRepositoryName string `json:"private_repository_name"`
	PublicRepositoryName  string `json:"public_repository_name"`
	Proceed               string `json:"proceed"`
}

func main() {
	l := logger.New()
	var input UserInput

	qs := []struct {
		Question string
		Message  string
		Fn       func(string)
	}{
		{
			Question: "GitHub Auth Token",
			Message:  "",
			Fn: func(authToken string) {
				input.AuthToken = authToken
				if err := os.Setenv("AUTH_TOKEN", authToken); err != nil {
					l.Fatal(err)
				}
			},
		},
		{
			Question: "Private Repository Name",
			Message:  "",
			Fn: func(repositoryName string) {
				input.PrivateRepositoryName = repositoryName
				input.PublicRepositoryName = strings.TrimSuffix(repositoryName, "-private")
			},
		},
		// {
		// 	Question: "Review the following and confirm to proceed(Y/n)",
		// 	Message:  fmt.Sprintf("\nThe following will be created:\nPrivate Repository: %s\nPublic Repository: %s\n\n", input.RepositoryName, strings.TrimSuffix(input.RepositoryName, "-private")),
		// 	Fn:       func(string) {},
		// },
	}

	for _, q := range qs {
		answer := common.Prompt(q.Question, q.Message)
		q.Fn(answer)
	}

	message := fmt.Sprintf("\nThe following will be created:\nPrivate Repository: %s\nPublic Repository: %s\n\n", input.PrivateRepositoryName, input.PublicRepositoryName)

	input.Proceed = common.Prompt("Review the following and confirm to proceed(Y/n)", message)

	if input.Proceed == "Y" || input.Proceed == "" || input.Proceed == "y" {
		// r := repo.New(input.PrivateRepositoryName, constants.PRIVATE_VISIBILITY, constants.ORGANIZATION, true)
		// fmt.Println(r)
		// r.CreateRepo()
		//
		// r = repo.New(input.PublicRepositoryName, constants.PUBLIC_VISIBILITY, constants.ORGANIZATION, true)
		// fmt.Println(r)
		// r.CreateRepo()
		//
		// opensshKeys := sshkeys.GenerateED25519Keys()
		//
		// d := deploykey.New(string(opensshKeys.PublicKey), input.PrivateRepositoryName, input.PublicRepositoryName, false)
		// fmt.Println(d)
		// d.CreateDeployKey()

		orgPubKey := orgpublickey.GetOrgPublicKey()
		l.Println(orgPubKey)
		// orgPubKeyB64Encoded := base64.StdEncoding.EncodeToString([]byte(orgPubKey.Key))
		// encryptedSecert := common.EncryptSecret(string(opensshKeys.PrivateKey), orgPubKeyB64Encoded)
		//
		// rs := reposecret.New(input.PrivateRepositoryName, constants.SECRET_NAME, encryptedSecert, orgPubKey.KeyID)
		// fmt.Println(rs)
		// rs.CreateRepositorySecret()
	}

	bytes, err := json.MarshalIndent(input, "", "    ")
	if err != nil {
		l.Fatal(err)
	}
	l.Println(string(bytes))
}

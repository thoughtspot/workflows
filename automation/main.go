package main

import (
	"fmt"
	"os"
	"strings"

	"automation/common"
	"automation/constants"
	"automation/deploykey"
	"automation/logger"
	"automation/repo"
	"automation/sshkeys"
)

type UserInput struct {
	/*
		The AuthToken must have the following permissions:
		1. "Administration" repository permissions (write)
		2. "Secrets" organization permissions (read)
		3. "Secrets" repository permissions (write)
	*/
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
	}

	for _, q := range qs {
		answer := common.Prompt(q.Question, q.Message)
		q.Fn(answer)
	}

	message := fmt.Sprintf("\nThe following will be created:\nPrivate Repository: %s\nPublic Repository: %s\n\n", input.PrivateRepositoryName, input.PublicRepositoryName)

	input.Proceed = common.Prompt("Review the following and confirm to proceed(Y/n)", message)

	Run(input)
}

func Run(input UserInput) {
	l := logger.New()
	switch input.Proceed {
	case "Y", "y", "":
		r := repo.NewRepository(input.PrivateRepositoryName, constants.PRIVATE_VISIBILITY, constants.ORGANIZATION, true)
		r.Create()

		r = repo.NewRepository(input.PublicRepositoryName, constants.PUBLIC_VISIBILITY, constants.ORGANIZATION, true)
		r.Create()

		opensshKeys := sshkeys.GenerateED25519Keys()

		d := deploykey.New(string(opensshKeys.PublicKey), input.PrivateRepositoryName, input.PublicRepositoryName, false)
		d.CreateDeployKey()

		orgPubKey := repo.GetPublicKey(input.PrivateRepositoryName)
		encryptedSecert := common.EncryptSecret(opensshKeys.PrivateKey, orgPubKey.Key)

		rs := repo.NewRepositorySecret(input.PrivateRepositoryName, constants.SECRET_NAME, encryptedSecert, orgPubKey.KeyID)
		rs.CreateSecret()
	case "N", "n":
		fmt.Println("Not proceeding...!")
	default:
		l.Fatal("Invalid input! Exiting...")
	}
}

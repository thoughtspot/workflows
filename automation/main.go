package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"automation/reposecret"
)

var questions = []string{
	"GitHub Auth Token",
	"Private Repository Name",
}

type UserInput struct {
	AuthToken      string `json:"auth_token"`
	RepositoryName string `json:"repository_name"`
}

func main() {
	reposecret.CreateSSHKeys()
	var input UserInput

	qs := map[string]func(string){
		questions[0]: func(authToken string) {
			input.AuthToken = authToken
			if err := os.Setenv("AUTH_TOKEN", authToken); err != nil {
				panic(err)
			}
		},
		questions[1]: func(repositoryName string) {
			input.RepositoryName = repositoryName
		},
	}

	for question, fn := range qs {
		answer := prompt(question)
		fn(answer)
	}

	bytes, err := json.MarshalIndent(input, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))
}

func prompt(question string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(question + ": ")
	answer, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(answer)
}

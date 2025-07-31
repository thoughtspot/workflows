# How to

## Prerequisites

### To run this automation you will need a GitHub token with the following permissions:

1. "Administration" repository permissions (write)
2. "Secrets" repository permissions (read)
3. "Secrets" repository permissions (write)

> [!Note] 
> Along with the github token, you will need to provide the name of the repository
you want to create. _Make sure the repository name ends with `-private` suffix.
Failing to provide a repository name with `-private` suffix will cause the cli
to fail._

## Setting up the Environment

1. Install `go` by following the official website to get the instructions for your specific platform https://go.dev/doc/install
2. Clone the repo `git clone git@github.com:thoughtspot/workflows.git`
3. Go to the `automation` directory `cd workflows/automation`
4. Run `go mod tidy` to initialize all dependencies
5. Run `go run main.go` to run the cli.
6. Follow the on-screen prompts

## GitHub Documentations used
- https://docs.github.com/en/rest/repos/repos?apiVersion=2022-11-28#create-an-organization-repository
- https://docs.github.com/en/rest/deploy-keys/deploy-keys?apiVersion=2022-11-28#create-a-deploy-key
- https://docs.github.com/en/rest/actions/secrets?apiVersion=2022-11-28#create-or-update-a-repository-secret
- https://docs.github.com/en/rest/actions/secrets?apiVersion=2022-11-28#get-a-repository-public-key

# How to

## To run this automation you will a github token with the following permissions
1. "Administration" repository permissions (write)
2. "Secrets" repository permissions (read)
3. "Secrets" repository permissions (write)

Note: Along with the github token, you will need to provide the name of the repository
you want to create. Make sure the repository name ends with `-private` suffix.
Failing to provide a repository name with `-private` prefix will cause the cli
to fail.

## Documentations used
- https://docs.github.com/en/rest/repos/repos?apiVersion=2022-11-28#create-an-organization-repository
- https://docs.github.com/en/rest/deploy-keys/deploy-keys?apiVersion=2022-11-28#create-a-deploy-key
- https://docs.github.com/en/rest/actions/secrets?apiVersion=2022-11-28#create-or-update-a-repository-secret
- https://docs.github.com/en/rest/actions/secrets?apiVersion=2022-11-28#get-a-repository-public-key

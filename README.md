# Sync to Public Mirror Workflow

A reusable GitHub Actions workflow that automatically synchronizes private repositories to their public counterparts. This workflow helps maintain public mirrors of private repositories, automatically syncing all content and branch structures while excluding sensitive information.

## Overview

This workflow allows you to maintain a private repository containing sensitive information alongside a public version of the same codebase. When changes are pushed to the private repository, this workflow automatically updates the corresponding public repository.

## Features

- Automatic synchronization on push or PR merge events
- Manual synchronization via workflow dispatch
- Support for custom public repository naming
- Branch-aware synchronization (maintains all branches from the private repo)
- SSH key-based secure authentication

## Prerequisites

1. A private GitHub repository (typically with a name containing "-private")
2. A public GitHub repository to serve as the mirror
3. An SSH deploy key with write access to the public repository

## Setup Instructions

### 1. Generate an SSH Deploy Key

```bash
# Generate a new SSH key pair
ssh-keygen -t ed25519 -C "github-action-sync-key" -f deploy_key -N ""

# deploy_key (private key) - will be used as a secret
# deploy_key.pub (public key) - will be added to the public repository
```

### 2. Configure the Public Repository

1. Go to your public repository on GitHub
2. Navigate to Settings > Deploy keys
3. Click "Add deploy key"
4. Give it a title (e.g., "Sync from Private Repository")
5. Paste the content of `deploy_key.pub`
6. **Check "Allow write access"**
7. Click "Add key"

### 3. Configure the Private Repository

1. Go to your private repository on GitHub
2. Navigate to Settings > Secrets and variables > Actions
3. Click "New repository secret"
4. Set the name to `SSH_DEPLOY_KEY`
5. Paste the content of the private `deploy_key` file
6. Click "Add secret"

### 4. Add the Caller Workflow to Your Private Repository

Create a file in your private repository at `.github/workflows/sync-mirror.yml` with the following content:

```yaml
name: Sync Repository to Public Mirror

on:
  push:
    branches:
      - '**'
  pull_request:
    types: [closed]
  workflow_dispatch:
    inputs:
      force_sync:
        description: 'Force sync all branches'
        required: false
        default: 'true'
        type: boolean
      public_repo_name:
        description: 'Public repository name (leave empty to auto-derive by removing "-private")'
        required: false
        type: string

jobs:
  call-sync-workflow:
    uses: your-org/workflows/.github/workflows/sync-to-public-mirror.yml@main
    with:
      force_sync: ${{ github.event.inputs.force_sync || 'true' }}
      public_repo_name: ${{ github.event.inputs.public_repo_name || '' }}
    secrets:
      SSH_DEPLOY_KEY: ${{ secrets.SSH_DEPLOY_KEY }}
```

Be sure to replace `your-org/workflows` with the actual path to your centralized workflows repository.

## How It Works: Detailed Breakdown

### Workflow Triggers

The workflow runs on:
- **Push events**: Whenever code is pushed to any branch
- **Pull request closures**: When PRs are merged or closed
- **Manual dispatch**: For on-demand synchronization

### Job 1: `check-private-repo`

This job determines if the repository is a private one that should be synced:

1. **Determine repository names**:
   - Extracts the private repository name from `GITHUB_REPOSITORY`
   - Extracts the organization/username
   - Determines the public repository name by either:
     - Using the name provided through manual input
     - Automatically deriving it by removing the "-private" suffix

2. **Check if repository is private**:
   - Verifies if the repository name contains "-private"
   - Sets the `is_private` output flag used by subsequent jobs

### Job 2: `sync-to-public`

This job handles the actual synchronization process:

1. **Check out private repository**:
   - Clones the private repository with all branches (`fetch-depth: 0`)

2. **Set up SSH**:
   - Configures the SSH agent with the deploy key for secure authentication
   - Adds GitHub to known hosts for SSH connections

3. **Set Git identity**:
   - Configures Git with identity information for commits

4. **Get all branch names**:
   - Extracts a list of all branches from the private repository

5. **Create workspace for public mirror**:
   - Sets up a separate workspace for the public repository
   - Initializes Git and configures the remote for the public mirror
   - Verifies the public repository exists and is accessible

6. **Create and update all branches**:
   - For each branch in the private repository:
     - Checks out the branch
     - Creates a temporary directory for the branch content
     - Copies all files except the `.git` directory
     - Switches to the public mirror workspace
     - Creates or updates the branch in the public mirror
     - Commits changes

7. **Push all branches to public mirror**:
   - Pushes each branch to the public repository

## Advanced Configuration

### Custom Public Repository Name

By default, the workflow derives the public repository name by removing the "-private" suffix from the private repository name. For custom naming:

1. **In the GitHub UI**: When manually triggering the workflow, enter the desired name in the "Public repository name" field.
2. **In the caller workflow**: Add a value for the `public_repo_name` parameter.

### Force Sync

The "Force sync all branches" option ensures all branches are synchronized, even if there are no detected changes.

## Troubleshooting

### Common Issues

1. **SSH Key Access Denied**:
   - Verify the SSH deploy key is correctly added to the public repository
   - Ensure "Allow write access" is checked for the deploy key
   - Confirm the private key is correctly stored as a secret

2. **Public Repository Not Found**:
   - Check that the public repository exists
   - Verify the naming convention or custom name is correct

3. **Workflow Not Running**:
   - Ensure your private repository's name contains "-private" or modify the check as needed

## Contributing

Contributions to improve this workflow are welcome! Please submit issues or pull requests to the central workflows repository.

#!/bin/bash

# Path to your workflow file
WORKFLOW_FILE="sync-to-public-mirror.yml"

# Create the workflow content
cat > "$WORKFLOW_FILE" << 'EOF'
name: Sync to Public Mirror

on:
  push:
    branches:
      - '**'  # Matches all branches
  pull_request:
    types: [closed]
    branches:
      - '**'  # Matches all branches
  workflow_dispatch:
    inputs:
      force_sync:
        description: 'Force sync all branches'
        required: false
        default: 'true'
        type: boolean
      public_repo_name:
        description: 'Public repository name'
        required: false
        type: string
        default: ''

jobs:
  debug-info:
    runs-on: ubuntu-latest
    steps:
      - name: Debug Information
        run: |
          echo "Event name: ${{ github.event_name }}"
          echo "Repository: ${{ github.repository }}"
          echo "PR merged?: ${{ github.event_name == 'pull_request' && github.event.pull_request.merged == true }}"
  
  call-sync-workflow:
    needs: debug-info
    # Only run for merged PRs, push events, or manual triggers
    if: |
      github.event_name == 'push' || 
      github.event_name == 'workflow_dispatch' || 
      (github.event_name == 'pull_request' && github.event.pull_request.merged == true)
    uses: thoughtspot/workflows/.github/workflows/sync-to-public-mirror.yml@main
    with:
      force_sync: ${{ github.event.inputs.force_sync == 'true' || github.event_name == 'workflow_dispatch' }}
      public_repo_name: ${{ github.event.inputs.public_repo_name }}
    secrets:
      SSH_DEPLOY_KEY: ${{ secrets.SSH_DEPLOY_KEY }}
EOF

# Base64 encode the workflow file content
ENCODED_CONTENT=$(base64 -w 0 "$WORKFLOW_FILE")

echo "Encoded workflow content:"
echo "$ENCODED_CONTENT"
echo
echo "You can now use this encoded content as input to the deploy-workflow-to-all-branches.yml workflow."

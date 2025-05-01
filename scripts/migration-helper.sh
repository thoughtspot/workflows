#!/bin/bash

# Repository Migration Script
# This script helps migrate from public to private GitHub repositories
# and sets up SSH-based commit signing

# Display help information
show_help() {
    echo "Repository Migration Script"
    echo "============================"
    echo
    echo "This script helps migrate from public to private GitHub repositories"
    echo "and sets up SSH-based commit signing."
    echo
    echo "Usage: $0 [OPTIONS]"
    echo
    echo "Options:"
    echo "  -h, --help     Display this help message and exit"
    echo "  --org NAME     Specify GitHub organization name"
    echo "  --repo NAME    Specify repository name (without -private suffix)"
    echo "  --email EMAIL  Specify your email address for git config and SSH key"
    echo "  --name NAME    Specify your full name for git config"
    echo
    echo "Example:"
    echo "  $0 --org myorg --repo myrepo --email user@example.com --name \"John Doe\""
    echo
    exit 0
}

# Process command line arguments
ORG=""
REPO=""
EMAIL=""
FULLNAME=""

while [[ $# -gt 0 ]]; do
    case "$1" in
        -h|--help)
            show_help
            ;;
        --org)
            ORG="$2"
            shift 2
            ;;
        --repo)
            REPO="$2"
            shift 2
            ;;
        --email)
            EMAIL="$2"
            shift 2
            ;;
        --name)
            FULLNAME="$2"
            shift 2
            ;;
        *)
            echo "Unknown option: $1"
            echo "Use --help to see available options"
            exit 1
            ;;
    esac
done

# Initialize log functions
log_info() {
    echo -e "\033[0;32m[INFO]\033[0m $1"
}

log_warn() {
    echo -e "\033[0;33m[WARNING]\033[0m $1"
}

log_error() {
    echo -e "\033[0;31m[ERROR]\033[0m $1"
}

cleanup() {
    log_info "Cleaning up temporary files..."
    # Add any cleanup tasks here
}

trap cleanup EXIT

# Check git installation
if ! command -v git &> /dev/null; then
    log_error "Git is not installed. Please install Git and try again."
    exit 1
fi

# Default SSH key path
SSH_KEY_PATH="$HOME/.ssh/id_ed25519"

# Functions
check_ssh_key() {
    if [[ ! -f "$SSH_KEY_PATH" ]]; then
        if [[ -f "$HOME/.ssh/id_rsa" ]]; then
            log_info "Found RSA key instead of Ed25519 key."
            SSH_KEY_PATH="$HOME/.ssh/id_rsa"
        else
            return 1
        fi
    fi
    return 0
}

generate_ssh_key() {
    log_info "Generating a new SSH key..."
    read -p "Enter your email address: " email
    ssh-keygen -t ed25519 -C "$email"

    # Check if key was successfully generated
    if [[ ! -f "$SSH_KEY_PATH" ]]; then
        log_error "Failed to generate SSH key. Please try manually."
        exit 1
    fi

    log_info "SSH key generated successfully."
}

add_key_to_agent() {
    log_info "Adding SSH key to ssh-agent..."

    # Start ssh-agent if not already running
    eval "$(ssh-agent -s)"

    # Add the key
    ssh-add "$SSH_KEY_PATH"

    if [ $? -ne 0 ]; then
        log_error "Failed to add key to ssh-agent. Please try manually."
        exit 1
    fi

    log_info "SSH key added to agent successfully."
}

display_public_key() {
    local public_key_path="${SSH_KEY_PATH}.pub"

    log_info "Your public SSH key is:"
    echo "-------------------------------------"
    cat "$public_key_path"
    echo "-------------------------------------"

    # Try to copy to clipboard
    if command -v pbcopy &> /dev/null; then
        cat "$public_key_path" | pbcopy
        log_info "Public key copied to clipboard."
    else
        log_info "Please manually copy the key above."
    fi

    echo ""
    log_info "Please add this key to your GitHub account at: https://github.com/settings/keys"
    log_info "Make sure to enable it for both authentication and signing."
    read -p "Press Enter after adding the key to GitHub... "
}

test_github_connection() {
    log_info "Testing connection to GitHub..."
    ssh -T git@github.com

    if [ $? -eq 1 ]; then
        # Exit code 1 actually means success for ssh -T (it just means no shell)
        log_info "Connection to GitHub successful."
        return 0
    else
        log_error "Failed to connect to GitHub. Please check your SSH configuration."
        return 1
    fi
}

update_remote_url() {
    # Get current directory name as fallback
    local repo_name=$(basename "$(pwd)")

    # Get the current remote URL
    local current_remote=$(git remote get-url origin 2>/dev/null)

    if [[ -z "$current_remote" ]]; then
        log_error "No remote named 'origin' found. Is this a git repository?"
        return 1
    fi

    log_info "Current remote URL: $current_remote"

    # Parse organization and repository name
    local org=""
    local repo=""

    if [[ "$current_remote" =~ github\.com[:/]([^/]+)/([^/\.]+)(\.git)?$ ]]; then
        org="${BASH_REMATCH[1]}"
        repo="${BASH_REMATCH[2]}"
    fi

    # Use org and repo from command line args if provided
    if [[ -n "$ORG" ]]; then
        org="$ORG"
        log_info "Using provided organization: $org"
    fi

    if [[ -n "$REPO" ]]; then
        repo="$REPO"
        log_info "Using provided repository: $repo"
    fi

    # If org or repo still not set, ask user
    if [[ -z "$org" || -z "$repo" ]]; then
        log_warn "Could not automatically detect organization or repository name."

        if [[ -z "$org" ]]; then
            read -p "Enter GitHub organization name: " org
        fi

        if [[ -z "$repo" ]]; then
            read -p "Enter repository name (without -private suffix): " repo
        fi
    else
        log_info "Detected organization: $org"
        log_info "Detected repository: $repo"
    fi

    # If repository already has -private suffix, remove it
    repo="${repo/%-private/}"

    # Construct new remote URL
    local new_remote="git@github.com:$org/$repo-private.git"

    # Update remote URL
    log_info "Updating remote URL to: $new_remote"
    git remote set-url origin "$new_remote"

    if [ $? -ne 0 ]; then
        log_error "Failed to update remote URL. Please try manually."
        return 1
    fi

    log_info "Remote URL updated successfully."
    return 0
}

configure_git_signing() {
    log_info "Configuring Git for SSH commit signing..."

    # Configure SSH for signing
    git config --global gpg.format ssh
    git config --global user.signingkey "$SSH_KEY_PATH"
    git config --global commit.gpgsign true

    # Configure additional git settings for better experience
    git config --global pull.rebase true
    git config --global push.default simple
    git config --global core.editor "nano"  # Use nano as default editor, can be changed

    if [ $? -ne 0 ]; then
        log_error "Failed to configure Git for SSH signing. Please try manually."
        return 1
    fi

    log_info "Git configured for SSH commit signing successfully."
    return 0
}

sync_with_remote() {
    log_info "Syncing with remote repository..."

    # Fetch latest changes
    git fetch origin

    if [ $? -ne 0 ]; then
        log_warn "Failed to fetch from remote. The repository might not exist or you might not have access."
        return 1
    fi

    # Detect current branch
    local current_branch=$(git rev-parse --abbrev-ref HEAD)

    log_info "Current branch: $current_branch"

    # Offer to rebase
    read -p "Do you want to rebase your local branch with the remote? (y/n): " do_rebase

    if [[ "$do_rebase" == "y" || "$do_rebase" == "Y" ]]; then
        git pull --rebase origin "$current_branch"

        if [ $? -ne 0 ]; then
            log_warn "Rebase failed. You might need to resolve conflicts manually."
            return 1
        fi

        log_info "Rebase successful."
    fi

    return 0
}

# Main execution
log_info "Starting repository migration process..."

# First, check if this is a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    log_error "The current directory is not a Git repository. Please navigate to your repository directory."
    exit 1
fi

# Step 1: Always create a new SSH key for this migration
log_info "Creating a new SSH key for repository migration..."

# Use email from command line or prompt
if [[ -z "$EMAIL" ]]; then
    read -p "Enter your email address: " EMAIL
fi

# Set up the user email in git config
log_info "Setting up git config with your email..."
git config --global user.email "$EMAIL"

# Use full name from command line or prompt
if [[ -z "$FULLNAME" ]]; then
    read -p "Enter your full name for git commits: " FULLNAME
fi
git config --global user.name "$FULLNAME"

# Generate a new SSH key with a specific filename for this migration
SSH_KEY_PATH="$HOME/.ssh/github_private_repo_$(date +%Y%m%d)"
log_info "Generating new SSH key at: $SSH_KEY_PATH"
ssh-keygen -t ed25519 -C "$EMAIL" -f "$SSH_KEY_PATH"

# Check if key was successfully generated
if [[ ! -f "$SSH_KEY_PATH" ]]; then
    log_error "Failed to generate SSH key. Please try manually."
    exit 1
fi

log_info "SSH key generated successfully."

# Step 2: Add key to ssh-agent
log_info "Ensuring SSH key is added to ssh-agent..."
add_key_to_agent

# Step 3: Display public key for GitHub
log_info "Your SSH key is set up locally."
read -p "Do you need to add your SSH key to GitHub? (y/n): " add_to_github

if [[ "$add_to_github" == "y" || "$add_to_github" == "Y" ]]; then
    display_public_key
fi

# Step 4: Test GitHub connection
if ! test_github_connection; then
    log_error "Failed to connect to GitHub. Please check your SSH configuration."
    exit 1
fi

# This check was moved to the beginning of the script

# Step 6: Update remote URL
if ! update_remote_url; then
    log_error "Failed to update remote URL. Please update manually."
    exit 1
fi

# Step 7: Configure Git for SSH signing
if ! configure_git_signing; then
    log_error "Failed to configure Git for SSH signing. Please configure manually."
    exit 1
fi

# Step 8: Sync with remote repository
if ! sync_with_remote; then
    log_warn "Failed to sync with remote repository. Please check your repository access."
fi

log_info "Repository migration completed successfully."
log_info "Your repository is now set up to use the private repository with SSH commit signing."
echo ""
log_info "Next steps:"
log_info "1. Make sure you can push to the repository by making a test commit."
log_info "2. Verify that your commits are being signed by checking the 'Verified' badge on GitHub."
echo ""
log_info "For any issues, refer to the README or contact the repository maintainers."

#!/usr/bin/env bash

# migration-helerp.sh
# Repository Migration Script
# This script helps migrate from public to private GitHub repositories
# and sets up SSH-based authentication and commit signing

# Display help information
show_help() {
    echo "Repository Migration Script"
    echo "============================"
    echo
    echo "This script helps migrate from public to private GitHub repositories"
    echo "and sets up SSH-based authentication and/or commit signing."
    echo
    echo "Usage: $0 [OPTIONS]"
    echo
    echo "Options:"
    echo "  -h, --help             Display this help message and exit"
    echo "  --org NAME             Specify GitHub organization name"
    echo "  --repo NAME            Specify repository name (without -private suffix)"
    echo "  --email EMAIL          Specify your email address for git config and SSH key"
    echo "  --name NAME            Specify your full name for git config"
    echo "  --auth-key             Configure key for authentication only"
    echo "  --signing-key          Configure key for signing only"
    echo "  --both-keys            Configure keys for both authentication and signing (default)"
    echo "  --use-existing-key     Use an existing SSH key instead of generating a new one"
    echo "  --auth-key-path PATH   Path to existing authentication SSH key (requires --use-existing-key)"
    echo "  --signing-key-path PATH Path to existing signing SSH key (requires --use-existing-key)"
    echo
    echo "Example:"
    echo "  $0 --org myorg --repo myrepo --email user@example.com --name \"John Doe\""
    echo "  $0 --org myorg --repo myrepo --signing-key --use-existing-key --signing-key-path ~/.ssh/id_ed25519_signing"
    echo
    exit 0
}

# Process command line arguments
ORG=""
REPO=""
EMAIL=""
FULLNAME=""
USE_EXISTING_KEY=false
EXISTING_AUTH_KEY_PATH=""
EXISTING_SIGNING_KEY_PATH=""
SETUP_AUTH_KEY=true
SETUP_SIGNING_KEY=true

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
        --auth-key)
            SETUP_AUTH_KEY=true
            SETUP_SIGNING_KEY=false
            shift
            ;;
        --signing-key)
            SETUP_AUTH_KEY=false
            SETUP_SIGNING_KEY=true
            shift
            ;;
        --both-keys)
            SETUP_AUTH_KEY=true
            SETUP_SIGNING_KEY=true
            shift
            ;;
        --use-existing-key)
            USE_EXISTING_KEY=true
            shift
            ;;
        --auth-key-path)
            EXISTING_AUTH_KEY_PATH="$2"
            shift 2
            ;;
        --signing-key-path)
            EXISTING_SIGNING_KEY_PATH="$2"
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

# Default SSH key paths
AUTH_SSH_KEY_PATH="$HOME/.ssh/id_ed25519_github_auth"
SIGNING_SSH_KEY_PATH="$HOME/.ssh/id_ed25519_github_signing"

# Functions
check_ssh_key() {
    if [[ ! -f "$1" ]]; then
        log_error "SSH key not found at: $1"
        return 1
    fi
    
    # Check if the file is a valid private key
    if ! ssh-keygen -l -f "$1" &>/dev/null; then
        log_error "The file at $1 is not a valid SSH private key."
        return 1
    fi
    
    # Check if corresponding public key exists
    if [[ ! -f "${1}.pub" ]]; then
        log_warn "Public key not found at: ${1}.pub"
        read -p "Do you want to generate the public key from your private key? (y/n): " gen_pub_key
        if [[ "$gen_pub_key" == "y" || "$gen_pub_key" == "Y" ]]; then
            ssh-keygen -y -f "$1" > "${1}.pub"
            if [ $? -ne 0 ]; then
                log_error "Failed to generate public key. The private key might be protected with a passphrase."
                return 1
            fi
            log_info "Public key generated at: ${1}.pub"
        else
            log_error "Public key is required for GitHub authentication and commit signing."
            return 1
        fi
    fi
    
    return 0
}

generate_ssh_key() {
    local key_type="$1"  # "auth" or "signing"
    local email="$2"
    local key_path="$3"
    
    log_info "Generating a new SSH key for $key_type..."
    
    # Use email from command line or prompt
    if [[ -z "$email" ]]; then
        read -p "Enter your email address: " email
    fi
    
    # Generate the key
    ssh-keygen -t ed25519 -C "$email" -f "$key_path"

    # Check if key was successfully generated
    if [[ ! -f "$key_path" ]]; then
        log_error "Failed to generate SSH key for $key_type. Please try manually."
        return 1
    fi

    log_info "SSH key for $key_type generated successfully at: $key_path"
    return 0
}

add_key_to_agent() {
    local key_path="$1"
    local key_type="$2"  # "auth" or "signing"
    
    log_info "Adding $key_type SSH key to ssh-agent..."

    # Start ssh-agent if not already running
    eval "$(ssh-agent -s)"

    # Add the key
    ssh-add "$key_path"

    if [ $? -ne 0 ]; then
        log_error "Failed to add $key_type key to ssh-agent. Please try manually."
        return 1
    fi

    log_info "SSH key for $key_type added to agent successfully."
    return 0
}

display_public_key() {
    local public_key_path="${1}.pub"
    local key_type="$2"  # "auth", "signing", or "both"
    local github_key_type=""

    if [[ "$key_type" == "auth" ]]; then
        github_key_type="authentication only"
        github_instructions="Make sure to NOT enable it for signing."
    elif [[ "$key_type" == "signing" ]]; then
        github_key_type="signing only"
        github_instructions="Make sure to ONLY check 'Signing Key' and NOT enable it for authentication."
    else
        github_key_type="both authentication and signing"
        github_instructions="Make sure to enable it for both authentication and signing."
    fi

    log_info "Your public SSH key for $github_key_type is:"
    echo "-------------------------------------"
    cat "$public_key_path"
    echo "-------------------------------------"

    # Try to copy to clipboard
    if command -v pbcopy &> /dev/null; then
        cat "$public_key_path" | pbcopy
        log_info "Public key copied to clipboard."
    elif command -v xclip &> /dev/null; then
        cat "$public_key_path" | xclip -selection clipboard
        log_info "Public key copied to clipboard."
    elif command -v clip &> /dev/null; then
        cat "$public_key_path" | clip
        log_info "Public key copied to clipboard."
    else
        log_info "Please manually copy the key above."
    fi

    echo ""
    log_info "Please add this key to your GitHub account at: https://github.com/settings/keys"
    log_info "$github_instructions"
    
    if [[ "$key_type" == "signing" || "$key_type" == "both" ]]; then
        log_info "For signing key, click 'New SSH key', add a title, paste your key, and CHECK 'Signing Key'."
    fi
    
    read -p "Press Enter after adding the key to GitHub... "
}

test_github_connection() {
    if [ "$SETUP_AUTH_KEY" = true ]; then
        log_info "Testing connection to GitHub with authentication key..."
        ssh -T git@github.com

        if [ $? -eq 1 ]; then
            # Exit code 1 actually means success for ssh -T (it just means no shell)
            log_info "Connection to GitHub successful."
            return 0
        else
            log_error "Failed to connect to GitHub. Please check your SSH configuration."
            return 1
        fi
    else
        log_info "Skipping GitHub connection test as authentication key is not being set up."
        return 0
    fi
}

update_remote_url() {
    if [ "$SETUP_AUTH_KEY" = true ]; then
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
    else
        log_info "Skipping remote URL update as authentication key is not being set up."
        return 0
    fi
}

configure_git_signing() {
    if [ "$SETUP_SIGNING_KEY" = true ]; then
        local key_path="$1"
        
        log_info "Configuring Git for SSH commit signing..."

        # Configure SSH for signing
        git config --global gpg.format ssh
        git config --global user.signingkey "$key_path.pub"  # Use the public key for signing
        git config --global commit.gpgsign true

        if [ $? -ne 0 ]; then
            log_error "Failed to configure Git for SSH signing. Please try manually."
            return 1
        fi

        log_info "Git configured for SSH commit signing successfully."
        return 0
    else
        log_info "Skipping Git signing configuration as signing key is not being set up."
        return 0
    fi
}

configure_git_settings() {
    # Configure additional git settings for better experience
    log_info "Configuring additional Git settings..."
    git config --global pull.rebase true
    git config --global push.default simple
    git config --global core.editor "nano"  # Use nano as default editor, can be changed

    log_info "Git settings configured successfully."
    return 0
}

sync_with_remote() {
    if [ "$SETUP_AUTH_KEY" = true ]; then
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
    else
        log_info "Skipping remote sync as authentication key is not being set up."
        return 0
    fi
}

prompt_key_setup_type() {
    # Ask user which type of key setup they want if not specified via command line
    if [[ "$SETUP_AUTH_KEY" == true && "$SETUP_SIGNING_KEY" == true ]]; then
        read -p "Do you want to set up keys for: [1] authentication only, [2] signing only, or [3] both (default)? [3] " key_setup_choice
        
        case "$key_setup_choice" in
            1)
                SETUP_AUTH_KEY=true
                SETUP_SIGNING_KEY=false
                log_info "Setting up authentication key only."
                ;;
            2)
                SETUP_AUTH_KEY=false
                SETUP_SIGNING_KEY=true
                log_info "Setting up signing key only."
                ;;
            3|"")
                SETUP_AUTH_KEY=true
                SETUP_SIGNING_KEY=true
                log_info "Setting up both authentication and signing keys."
                ;;
            *)
                log_error "Invalid choice. Defaulting to setting up both keys."
                SETUP_AUTH_KEY=true
                SETUP_SIGNING_KEY=true
                ;;
        esac
    fi
}

# Main execution
log_info "Starting repository migration process..."

# First, check if this is a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    log_error "The current directory is not a Git repository. Please navigate to your repository directory."
    exit 1
fi

# Configure git user.email and user.name
if [[ -z "$EMAIL" ]]; then
    read -p "Enter your email address: " EMAIL
fi
git config --global user.email "$EMAIL"

if [[ -z "$FULLNAME" ]]; then
    read -p "Enter your full name for git commits: " FULLNAME
fi
git config --global user.name "$FULLNAME"

log_info "Git user config set up with email: $EMAIL and name: $FULLNAME"

# Ask user which type of key setup they want (if not specified via command-line)
prompt_key_setup_type

# Handle SSH key setup for authentication
if [ "$SETUP_AUTH_KEY" = true ]; then
    if [ "$USE_EXISTING_KEY" = true ]; then
        # Use existing key path from command line argument
        if [[ -n "$EXISTING_AUTH_KEY_PATH" ]]; then
            AUTH_SSH_KEY_PATH="$EXISTING_AUTH_KEY_PATH"
            log_info "Using existing authentication SSH key at: $AUTH_SSH_KEY_PATH"
        else
            # Prompt for existing key path
            read -p "Enter the full path to your existing authentication SSH key: " AUTH_SSH_KEY_PATH
        fi

        # Validate the existing key
        if ! check_ssh_key "$AUTH_SSH_KEY_PATH"; then
            log_error "The specified authentication SSH key is invalid or doesn't exist."
            exit 1
        fi
    else
        # Generate a new SSH key for authentication
        log_info "Will generate new authentication SSH key at: $AUTH_SSH_KEY_PATH"
        if ! generate_ssh_key "authentication" "$EMAIL" "$AUTH_SSH_KEY_PATH"; then
            log_error "Failed to generate authentication SSH key. Please try manually."
            exit 1
        fi
    fi

    # Add authentication key to ssh-agent
    if ! add_key_to_agent "$AUTH_SSH_KEY_PATH" "authentication"; then
        log_error "Failed to add authentication SSH key to ssh-agent. Please try manually."
        exit 1
    fi

    # Ask if the key needs to be added to GitHub
    read -p "Do you need to add your authentication SSH key to GitHub? (y/n): " add_auth_to_github
    if [[ "$add_auth_to_github" == "y" || "$add_auth_to_github" == "Y" ]]; then
        if [ "$SETUP_SIGNING_KEY" = true ]; then
            display_public_key "$AUTH_SSH_KEY_PATH" "auth"
        else
            display_public_key "$AUTH_SSH_KEY_PATH" "auth"
        fi
    fi
fi

# Handle SSH key setup for signing
if [ "$SETUP_SIGNING_KEY" = true ]; then
    if [ "$USE_EXISTING_KEY" = true ]; then
        # Use existing key path from command line argument
        if [[ -n "$EXISTING_SIGNING_KEY_PATH" ]]; then
            SIGNING_SSH_KEY_PATH="$EXISTING_SIGNING_KEY_PATH"
            log_info "Using existing signing SSH key at: $SIGNING_SSH_KEY_PATH"
        else
            # Prompt for existing key path
            read -p "Enter the full path to your existing signing SSH key: " SIGNING_SSH_KEY_PATH
        fi

        # Validate the existing key
        if ! check_ssh_key "$SIGNING_SSH_KEY_PATH"; then
            log_error "The specified signing SSH key is invalid or doesn't exist."
            exit 1
        fi
    else
        # Generate a new SSH key for signing
        log_info "Will generate new signing SSH key at: $SIGNING_SSH_KEY_PATH"
        if ! generate_ssh_key "signing" "$EMAIL" "$SIGNING_SSH_KEY_PATH"; then
            log_error "Failed to generate signing SSH key. Please try manually."
            exit 1
        fi
    fi

    # Add signing key to ssh-agent
    if ! add_key_to_agent "$SIGNING_SSH_KEY_PATH" "signing"; then
        log_error "Failed to add signing SSH key to ssh-agent. Please try manually."
        exit 1
    fi

    # Ask if the key needs to be added to GitHub
    read -p "Do you need to add your signing SSH key to GitHub? (y/n): " add_signing_to_github
    if [[ "$add_signing_to_github" == "y" || "$add_signing_to_github" == "Y" ]]; then
        display_public_key "$SIGNING_SSH_KEY_PATH" "signing"
    fi
fi

# Test GitHub connection (only for authentication key)
if [ "$SETUP_AUTH_KEY" = true ]; then
    if ! test_github_connection; then
        log_error "Failed to connect to GitHub. Please check your SSH configuration."
        exit 1
    fi
fi

# Update remote URL (only for authentication key)
if [ "$SETUP_AUTH_KEY" = true ]; then
    if ! update_remote_url; then
        log_error "Failed to update remote URL. Please update manually."
        exit 1
    fi
fi

# Configure Git for SSH signing (only for signing key)
if [ "$SETUP_SIGNING_KEY" = true ]; then
    if ! configure_git_signing "$SIGNING_SSH_KEY_PATH"; then
        log_error "Failed to configure Git for SSH signing. Please configure manually."
        exit 1
    fi
fi

# Configure additional Git settings
configure_git_settings

# Sync with remote repository (only for authentication key)
if [ "$SETUP_AUTH_KEY" = true ]; then
    if ! sync_with_remote; then
        log_warn "Failed to sync with remote repository. Please check your repository access."
    fi
fi

log_info "Repository migration completed successfully."

if [ "$SETUP_AUTH_KEY" = true ]; then
    log_info "Your repository is now set up to use the private repository with SSH authentication."
fi

if [ "$SETUP_SIGNING_KEY" = true ]; then
    log_info "Your repository is now set up to use SSH commit signing."
fi

echo ""
log_info "Next steps:"

if [ "$SETUP_AUTH_KEY" = true ]; then
    log_info "1. Make sure you can push to the repository by making a test commit."
fi

if [ "$SETUP_SIGNING_KEY" = true ]; then
    log_info "2. Verify that your commits are being signed by checking the 'Verified' badge on GitHub."
    log_info "   Create a test commit with: git commit --allow-empty -m \"Test signed commit\""
fi

echo ""
log_info "For any issues, refer to the README or contact the repository maintainers."

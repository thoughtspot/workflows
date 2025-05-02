# Repository Migration Script

A robust bash script to streamline the process of migrating GitHub repositories from public to private, while also configuring SSH-based authentication and commit signing.

## Overview

This script automates several important steps when transitioning from a public GitHub repository to a private one:

1. SSH key management for both authentication and signing (use existing or generate new)
2. GitHub SSH authentication setup
3. Git remote URL updates
4. SSH-based commit signing configuration
5. Repository synchronization

## Prerequisites

- Git (installed and available in your PATH)
- SSH client
- Bash shell environment
- GitHub account with appropriate permissions

## Installation

1. Download the script:
   ```bash
   curl -O https://raw.githubusercontent.com/thoughtspot/workflows/refs/heads/main/scripts/migration-helper.sh
   ```

2. Make the script executable:
   ```bash
   chmod +x migration-helper.sh
   ```

## Usage

### Basic Usage

Run the script within your git repository:

```bash
./migration-helper.sh
```

The script will guide you through an interactive process with prompts for all required information, including whether you want to set up authentication keys, signing keys, or both.

### Command Line Arguments

For automated or scripted use, you can provide parameters through command line arguments:

```bash
./migration-helper.sh --org "your-org" --repo "your-repo" --email "your-email@example.com" --name "Your Name"
```

### Available Options

| Option | Description |
|--------|-------------|
| `-h, --help` | Display help information |
| `--org NAME` | Specify GitHub organization name |
| `--repo NAME` | Specify repository name (without -private suffix) |
| `--email EMAIL` | Specify your email address for git config and SSH key |
| `--name NAME` | Specify your full name for git config |
| `--auth-key` | Configure key for authentication only |
| `--signing-key` | Configure key for signing only |
| `--both-keys` | Configure keys for both authentication and signing (default) |
| `--use-existing-key` | Use existing SSH key(s) instead of generating new ones |
| `--auth-key-path PATH` | Path to existing authentication SSH key |
| `--signing-key-path PATH` | Path to existing signing SSH key |

### SSH Key Management

The script offers two key setup options with three configurations:

1. **Key Purpose**:
   * Authentication only: Use SSH key only for connecting to GitHub
   * Signing only: Use SSH key only for signing commits
   * Both (default): Set up keys for both authentication and signing

2. **Key Source**:
   * Use existing key(s): If you already have SSH keys set up
   ```bash
   ./migration-helper.sh --use-existing-key --auth-key-path ~/.ssh/your_auth_key --signing-key-path ~/.ssh/your_signing_key
   ```
   
   * Generate new key(s): The script can generate new Ed25519 SSH keys
   ```bash
   # Generate keys for both authentication and signing
   ./migration-helper.sh --both-keys

   # Generate key for authentication only
   ./migration-helper.sh --auth-key
   
   # Generate key for signing only
   ./migration-helper.sh --signing-key
   ```

When running interactively, the script will prompt you to choose between these options.

## Setting Up SSH Signing Keys in GitHub

After generating or selecting your SSH signing key, you need to configure it in GitHub:

1. **Copy your public signing key**:
   The script will display your public signing key and attempt to copy it to your clipboard.

2. **Add the key to GitHub**:
   - Go to GitHub: [Settings → SSH and GPG keys](https://github.com/settings/keys)
   - Click "New SSH key"
   - Provide a descriptive title (e.g., "Commit Signing Key")
   - Paste your public key
   - **Important**: Select "Signing Key" from the dropdown menu
   - Click "Add SSH key"

   ![SSH Signing Key Selection](https://docs.github.com/assets/images/help/settings/ssh-signing-key-dropdown.png)

3. **Verify setup**:
   - After adding the key and running the script, make a test commit:
   ```bash
   git commit --allow-empty -m "Test signed commit"
   git push
   ```
   - Check on GitHub that your commit shows the "Verified" badge

## How It Works

### Step-by-Step Process

1. **Initialization**:
   - Validates that the current directory is a git repository
   - Sets up git user configuration (email and name)
   - Determines which type of keys to set up (authentication, signing, or both)

2. **SSH Key Setup**:
   - Option to use existing SSH key(s) or generate new one(s)
   - Validates that the key(s) exist and are properly formatted
   - Generates public key(s) if only private key(s) are present
   - Sets up separate keys for authentication and signing if requested

3. **GitHub Integration**:
   - Adds the SSH key(s) to the ssh-agent
   - Displays the public key(s) for adding to GitHub (with clipboard support)
   - Tests the GitHub connection for authentication key
   - Provides specific instructions for adding signing key with proper GitHub settings

4. **Repository Configuration**:
   - Updates the remote URL to point to the new private repository
   - Configures git for SSH-based commit signing with the specified signing key
   - Syncs with the remote repository (optional rebase)

5. **Completion**:
   - Provides next steps and verification guidance

### Repository Naming Convention

The script assumes that your private repository will follow the naming convention of `{original-name}-private`. For example, if your public repository is named `my-project`, the private one should be named `my-project-private`.

## Troubleshooting

### Common Issues

1. **SSH Connection Failure**:
   - Ensure your SSH key is properly added to GitHub
   - Verify your key has been added to the ssh-agent
   - Check GitHub SSH documentation for further assistance

2. **Repository Access Issues**:
   - Confirm you have appropriate permissions to the organization and repository
   - Verify the repository exists with the expected name convention

3. **Key Validation Errors**:
   - Ensure you're providing the path to the private key (not the .pub file)
   - Check that the key file has appropriate permissions (typically 600)
   - Verify the key is in the correct format

4. **Commit Signing Issues**:
   - Ensure your signing key is correctly added to GitHub with the "Signing Key" checkbox checked
   - Verify that your signing key is correctly configured in Git
   - Check that the public key is being used for signing in your Git config

### Logs

The script provides colored log output to help identify information, warnings, and errors:
- <span style="color:green">[INFO]</span> - Informational messages
- <span style="color:yellow">[WARNING]</span> - Potential issues that may need attention
- <span style="color:red">[ERROR]</span> - Critical issues that prevent completion

## Advanced Usage

### CI/CD Integration

For CI/CD pipelines, use the non-interactive mode with all required parameters:

```bash
./migration-helper.sh --org "your-org" --repo "your-repo" --email "ci-bot@example.com" --name "CI Bot" --auth-key --use-existing-key --auth-key-path /path/to/ci/ssh_key
```

### Custom Editor Configuration

The script sets nano as the default git editor. To change this, modify the following line before running:

```bash
# Inside the script, find and modify:
git config --global core.editor "your-preferred-editor"
```

## Security Considerations

- The script generates Ed25519 keys, which offer a good balance of security and performance
- All SSH keys are protected with standard file permissions
- Consider using a passphrase when generating new keys for additional security
- The script adds SSH keys to your ssh-agent for convenience, but these are cleared when you log out
- Using separate keys for authentication and signing provides better security isolation

## FAQ

### What's the difference between authentication and signing keys?

**Authentication keys** are used to securely connect to GitHub when you push, pull, or perform other operations that require remote access. **Signing keys** are used to cryptographically sign your commits, allowing others to verify that commits were actually made by you. While you can use the same key for both purposes, using separate keys provides better security isolation.

### How do I set up commit signing if I already have SSH configured for my repository?

If you already have SSH authentication working with GitHub but want to enable commit signing:

1. **Generate a dedicated signing key** (recommended):
   ```bash
   ssh-keygen -t ed25519 -C "github-signing-key" -f ~/.ssh/github_signing_key
   ```

2. **Add your signing key to GitHub**:
   - Go to GitHub: Settings → SSH and GPG keys
   - Click "New SSH key"
   - Give it a title like "GitHub Commit Signing Key"
   - Paste your public key (~/.ssh/github_signing_key.pub)
   - **Select "Signing Key" from the key type dropdown menu**
   - Save changes

3. **Configure Git to use SSH for signing**:
   ```bash
   # Set the signing format to SSH
   git config --global gpg.format ssh
   
   # Configure your signing key (note: use the public key)
   git config --global user.signingkey ~/.ssh/github_signing_key.pub
   
   # Enable commit signing by default
   git config --global commit.gpgsign true
   ```

4. **Test commit signing**:
   ```bash
   # Make a test commit
   git commit --allow-empty -m "Test signed commit"
   
   # Push to GitHub
   git push
   ```

5. **Verify signature on GitHub**:
   - Go to your repository on GitHub
   - Look at your recent commits
   - You should see a "Verified" badge next to your commit

### Can I use different SSH keys for authentication and signing?

Yes, and the updated script supports this workflow explicitly with the `--auth-key` and `--signing-key` options. Using separate keys is recommended for better security isolation:

```bash
# Configure authentication key in ~/.ssh/config
Host github.com
    IdentityFile ~/.ssh/github_auth_key

# Configure signing key in Git
git config --global gpg.format ssh
git config --global user.signingkey ~/.ssh/github_signing_key.pub
```

### How do I sign a single commit without enabling global signing?

To sign an individual commit without enabling global signing:

```bash
git commit -S -m "Your commit message"
```

The `-S` flag indicates that this commit should be signed.

### How do I disable commit signing for a specific repository?

To disable commit signing for a specific repository:

```bash
# Navigate to your repository
cd /path/to/your/repo

# Disable commit signing for this repository only
git config --local commit.gpgsign false
```

> **⚠️ IMPORTANT NOTE**: This instruction is **NOT RECOMMENDED**. Disabling commit signing will block your commits since signature verification is strictly enforced on our repositories. Only use this if you have been specifically instructed to do so by repository administrators and have an alternative signing method in place.

### Why isn't my commit showing as "Verified" on GitHub?

If your commits aren't showing as verified:

1. Ensure your SSH key is added to GitHub with "Signing Key" selected in the key type dropdown
2. Confirm you're using the correct public key path in your Git config
3. Verify you're using the same email in your Git config as on GitHub
4. Check that Git is configured to use SSH for signing with `gpg.format ssh`
5. Make sure the key exists at the specified path and has correct permissions

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

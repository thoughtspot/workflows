# Repository Migration Script

A robust bash script to streamline the process of migrating GitHub repositories from public to private, while also configuring SSH-based commit signing.

## Overview

This script automates several important steps when transitioning from a public GitHub repository to a private one:

1. SSH key management (use existing or generate new)
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

The script will guide you through an interactive process with prompts for all required information.

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
| `--use-existing-key` | Use an existing SSH key instead of generating a new one |
| `--key-path PATH` | Path to existing SSH key (requires --use-existing-key) |

### SSH Key Management

The script offers two options for SSH key management:

1. **Use an existing key**: If you already have an SSH key set up, you can use it with the script:
   ```bash
   ./migration-helper.sh --use-existing-key --key-path ~/.ssh/your_existing_key
   ```

2. **Generate a new key**: The script can generate a new Ed25519 SSH key:
   ```bash
   # The script will generate a key at ~/.ssh/github_private_repo_YYYYMMDD
   ./migration-helper.sh
   ```

When running interactively, the script will prompt you to choose between these options.

## How It Works

### Step-by-Step Process

1. **Initialization**:
   - Validates that the current directory is a git repository
   - Sets up git user configuration (email and name)

2. **SSH Key Setup**:
   - Option to use an existing SSH key or generate a new one
   - Validates that the key exists and is properly formatted
   - Generates a public key if only a private key is present

3. **GitHub Integration**:
   - Adds the SSH key to the ssh-agent
   - Displays the public key for adding to GitHub (with clipboard support)
   - Tests the GitHub connection

4. **Repository Configuration**:
   - Updates the remote URL to point to the new private repository
   - Configures git for SSH-based commit signing
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

### Logs

The script provides colored log output to help identify information, warnings, and errors:
- <span style="color:green">[INFO]</span> - Informational messages
- <span style="color:yellow">[WARNING]</span> - Potential issues that may need attention
- <span style="color:red">[ERROR]</span> - Critical issues that prevent completion

## Advanced Usage

### CI/CD Integration

For CI/CD pipelines, use the non-interactive mode with all required parameters:

```bash
./migration-helper.sh --org "your-org" --repo "your-repo" --email "ci-bot@example.com" --name "CI Bot" --use-existing-key --key-path /path/to/ci/ssh_key
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

## FAQ

### How do I set up commit signing if I already have SSH configured for my repository?

If you already have SSH authentication working with GitHub but want to enable commit signing:

1. **Identify your existing SSH key**:
   ```bash
   # List your SSH keys
   ls -la ~/.ssh/
   ```
   Look for files like `id_ed25519`, `id_rsa`, or similar (without the .pub extension).

2. **Configure Git to use SSH for signing**:
   ```bash
   # Set the signing format to SSH
   git config --global gpg.format ssh
   
   # Configure your existing key for signing
   git config --global user.signingkey ~/.ssh/your_existing_key
   
   # Enable commit signing by default
   git config --global commit.gpgsign true
   ```

3. **Add your key to GitHub for signing**:
   - Go to GitHub: Settings → SSH and GPG keys
   - Click on your existing SSH key
   - Check the "Enable signing with this key" option
   - Save changes

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

Yes, you can configure Git to use separate SSH keys:

```bash
# Configure authentication key in ~/.ssh/config
Host github.com
    IdentityFile ~/.ssh/github_auth_key

# Configure signing key in Git
git config --global gpg.format ssh
git config --global user.signingkey ~/.ssh/github_signing_key
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

1. Ensure your SSH key is added to GitHub with signing enabled
2. Check that your Git configuration is using the correct key path
3. Confirm you're using the same email in your Git config as on GitHub
4. Verify the SSH key exists at the specified path and has correct permissions

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

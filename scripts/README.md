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
   curl -O https://raw.githubusercontent.com/yourusername/repo-migration/main/repo_migration.sh
   ```

2. Make the script executable:
   ```bash
   chmod +x repo_migration.sh
   ```

## Usage

### Basic Usage

Run the script within your git repository:

```bash
./repo_migration.sh
```

The script will guide you through an interactive process with prompts for all required information.

### Command Line Arguments

For automated or scripted use, you can provide parameters through command line arguments:

```bash
./repo_migration.sh --org "your-org" --repo "your-repo" --email "your-email@example.com" --name "Your Name"
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
   ./repo_migration.sh --use-existing-key --key-path ~/.ssh/your_existing_key
   ```

2. **Generate a new key**: The script can generate a new Ed25519 SSH key:
   ```bash
   # The script will generate a key at ~/.ssh/github_private_repo_YYYYMMDD
   ./repo_migration.sh
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
./repo_migration.sh --org "your-org" --repo "your-repo" --email "ci-bot@example.com" --name "CI Bot" --use-existing-key --key-path /path/to/ci/ssh_key
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

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

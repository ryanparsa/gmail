# Installation Guide

Follow these steps to install and set up the Gmail CLI tool:

## 1. Download the Binary

Using curl

You can download the pre-built binary for your platform directly using curl:

### 1. Linux

```shell
curl -LO https://github.com/your-repo-name/releases/latest/download/app-linux-amd64
chmod +x app-linux-amd64
sudo mv app-linux-amd64 /usr/local/bin/gmail-cli
```

### 2. MacOS (Darwin)

```shell
curl -LO https://github.com/your-repo-name/releases/latest/download/app-darwin-amd64
chmod +x app-darwin-amd64
sudo mv app-darwin-amd64 /usr/local/bin/gmail-cli
```

### 3. Windows

Download using curl and rename the file:

```shell
curl -LO https://github.com/your-repo-name/releases/latest/download/app-windows-amd64.exe
move app-windows-amd64.exe gmail-cli.exe
```

Place the binary (gmail-cli.exe) in a directory included in your PATH.

### Manually Downloading from GitHub

Alternatively, you can download the binary for your platform from the Releases page:

1. Visit the Releases Page.
2. Download the binary matching your operating system and architecture.
3. Extract the binary and move it to a directory included in your $PATH (e.g., /usr/local/bin for Linux/Mac or C:
   \Windows\System32 for Windows).

## 2. Install From Source

If you prefer, you can build the application from source:

Prerequisites
• Go (1.20 or newer)
• Git

Build Steps

1. Clone the repository:

```shell
git clone https://github.com/your-repo-name.git
cd your-repo-name

```
2.	Build the application:

```shell
go build -o gmail-cli

```
3.	
4. Move the binary to a directory in your $PATH:

```shell
mv gmail-cli /usr/local/bin/

```
4. Generate Your Own credentials.json

Why is this required?

For security purposes, the Gmail CLI requires each user to create their own Google OAuth credentials. This ensures that:
• Only your app accesses your Gmail account.
• Your account remains secure and private.

Learn More About Why This is Important

Steps to Generate credentials.json

Follow the steps in the Setting Up Gmail API guide to generate and download the credentials.json file.

1. Save the credentials.json file in the directory where you plan to run the Gmail CLI.
2. Ensure that no one else can access this file for your account’s safety.

4. Run the Application

Check Available Commands

To see all available commands:

gmail-cli --help

Example Commands
• Authenticate with Gmail API:

gmail-cli auth --credentials credentials.json --token token.json

	•	Apply Gmail Filters:

gmail-cli apply --credentials credentials.json --token token.json

	•	Dump Gmail Labels and Filters to a YAML File:

gmail-cli dump --credentials credentials.json --token token.json --scopes "https://www.googleapis.com/auth/gmail.labels"
--config dump.yaml

	•	Load Gmail Labels and Filters from a YAML File:

gmail-cli load --credentials credentials.json --token token.json --config dump.yaml

	•	Wipe All User-Created Labels and Filters:

gmail-cli wipe --credentials credentials.json --token token.json

5. Stay Updated

To ensure you’re using the latest version of the Gmail CLI tool, periodically check the Releases page for updates.

For any issues, feature requests, or contributions, please visit the Issues Page.

Now you’re all set to use the Gmail CLI tool! 🎉
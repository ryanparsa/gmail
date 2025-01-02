# Gmail

A CLI tool to manage Gmail filters and labels using a configuration file. This allows you to avoid Gmail's UI and manage
your filters efficiently in a declarative way.

---

## Features

- Manage Gmail filters using YAML configuration files.
- Create, update, and delete labels programmatically.
- Easily back up and restore Gmail filters.
- Command-line interface for automation.

---

## Installation

You can install the Gmail Filter Manager using one of the following methods:

### **1. Download Binary from GitHub Actions**

1. Go to the [Releases Page](https://github.com/ryanparsa/gmail/releases).
2. Download the binary for your operating system.
3. Make the binary executable:
   ```bash
   chmod +x gmail
   ```
4. Move the binary to a directory in your `$PATH` (e.g., `/usr/local/bin`):
   ```bash
   mv gmail /usr/local/bin
   ```
5. Verify installation:
   ```bash
   gmail --version
   ```

---

### **2. Build from Source**

1. Clone the repository:
   ```bash
   git clone https://github.com/ryanparsa/gmail.git
   cd gmail
   ```
2. Build the project:
   ```bash
   go build -o gmail .
   ```
3. Move the binary to a directory in your `$PATH`:
   ```bash
   mv gmail /usr/local/bin
   ```
4. Verify installation:
   ```bash
   gmail --version
   ```

---

### **3. Install Using `go install`**

If you have Go installed, you can install the tool directly using `go install`:

```bash
go install github.com/ryanparsa/gmail@latest
```

Ensure that your `$GOPATH/bin` is in your `$PATH`. Verify installation:

```bash
gmail --version
```

---

### **4. Using Docker**

1Run the tool:
   ```bash
   docker run --rm -v $(pwd):/data ghcr.io/ryanparsa/gmail backup.yaml
   ```

---

### **Verify Installation**

After installing using any method, verify that the tool is installed correctly:

```bash
gmail --help
```

---

## Post-Installation Steps

After installation, you need to set up credentials to allow the tool to interact with your Gmail account. Follow these
steps:

### 1. Obtain Credentials from Google Cloud Console

- You need a `credentials.json` file to authenticate the tool with Gmail's API.
- Follow the detailed instructions in [docs/credentials.md](docs/credentials.md) to:
    1. Set up a project in the [Google Cloud Console](https://console.cloud.google.com/).
    2. Enable the Gmail API for the project.
    3. Create OAuth 2.0 credentials.
    4. Download the `credentials.json` file and place it in the root directory of this project.

### 2. Authenticate

Run the following command to authenticate the tool with your Gmail account:

```bash
gmail auth
```

This will open a browser where you can grant access to your Gmail account.

### 3. Verify Authentication

Once authentication is complete, the tool will save a `token.json` file. This file is used for subsequent API requests.
To verify that authentication was successful, you can run:

```bash
gmail --help
```

---

## Usage

Here are some common commands you can use with Gmail Filter Manager:

### **1. Back Up Filters**

Back up existing Gmail filters to a YAML file:

```bash
gmail backup
```

### **2. Apply Filters**

Apply filters defined in `backup.yaml` to your Gmail account:

```bash
gmail push
```

### **3. Clean Filters**

Remove filters and labels that match certain criteria:

```bash
gmail clean
```

For more usage examples and details, see the [Usage Guide](docs/usage.md).

---

## Configuration

You can configure the tool using the `backup.yaml` and `config.yaml` files.

### Example `backup.yaml`

```yaml
filters:
    -   query: "from:noreply@example.com"
        label: "Automated"
    -   query: "subject:Important"
        label: "Priority"
```

### Example `config.yaml`

```yaml
gmail:
    client_id: <your_client_id>
    client_secret: <your_client_secret>
    token_file: token.json
```

For more details, see the [Configuration Guide](docs/config.md).

---

## Contributing

Contributions are welcome! Please fork this repository, make your changes, and submit a pull request.

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

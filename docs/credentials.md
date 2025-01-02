# Setting Up the Gmail API

## Why Do You Need to Generate Your Own credentials.json?

To use the Gmail CLI tool securely, you need to create your own credentials.json file. This step is critical for your
privacy and security when accessing the Gmail API. Here’s why:

1. **Protect Your Gmail Account**: By creating your own credentials.json, you’re ensuring that the app accessing your
   Gmail
   account is unique to you. This means that no one else has access to your email or sensitive data.
2. **Avoid Using Shared Configurations**: Many third-party apps rely on shared app configurations, which could
   potentially
   allow other users or developers to access your Gmail account. By generating your own credentials.json, you’re
   ensuring that only your app and your credentials are used for authentication.
3. **Comply with Google Security Guidelines**: Google recommends that each user create their own OAuth credentials to
   avoid
   the risks associated with sharing app configurations.
4. **Full Control Over Access**: With your own credentials, you have complete control over the permissions and access
   granted to the app. This minimizes the risk of unauthorized actions being performed on your Gmail account.

### How Does This Work?

The credentials.json file contains your unique Client ID and Client Secret. These credentials allow Google’s OAuth
servers to authenticate your app and ensure that no one else can use it to access your account.

By requiring users to create their own credentials, this Gmail CLI tool ensures:

- **Security**: Only you have access to the API client credentials linked to your Gmail account.
- **Privacy**: No shared app is used, so your emails and data remain private and secure.

## Step 1: Enable the Gmail API

To use the Gmail API, you need to enable it in a Google Cloud project. You can enable one or more APIs in a single
Google Cloud project.

1. Open the Google Cloud
   Console [Gmail API](https://console.cloud.google.com/flows/enableapi?apiid=gmail.googleapis.com) page.
2. Click on Enable API to activate the Gmail API for your project.

## Step 2: Configure the OAuth Consent Screen

If you’re using a new Google Cloud project, you’ll need to configure the OAuth consent screen and add yourself as a test
user. If this step has already been completed for your project, you can skip to the next section.

1. Go to
   the [OAuth consent screen settings in the Google Cloud Console](https://console.cloud.google.com/apis/credentials/consent).
2. Under User Type, select Internal, and click Create.
3. Fill out the app registration form and click Save and Continue.
4. Skip adding scopes for now by clicking Save and Continue.
   • Note: In the future, if your app will be used outside your Google Workspace organization, you’ll need to set the
   User Type to External and add the necessary authorization scopes.
5. Review the app registration summary.
   • To make changes, click Edit.
   • If everything looks correct, click Back to Dashboard.

## Step 3: Authorize Credentials for a Desktop Application

To authenticate end users and access user data, you’ll need to create one or more OAuth 2.0 Client IDs. Each Client ID
is used to uniquely identify your app with Google’s OAuth servers. If your app runs on multiple platforms, you’ll need
separate Client IDs for each platform.

1. Go to the [Credentials page in the Google Cloud Console](https://console.cloud.google.com/apis/credentials).
2. Click Create Credentials > OAuth client ID.
3. Under Application Type, select Desktop app.
4. In the Name field, provide a name for the credential. This name is only used in the Google Cloud Console for
   identification purposes.
5. Click Create. You’ll see your new Client ID and Client Secret on the confirmation screen.
6. Click OK. The newly created credential will appear under OAuth 2.0 Client IDs.
7. Download the JSON file and save it as credentials.json.
8. Move the credentials.json file to your working directory.

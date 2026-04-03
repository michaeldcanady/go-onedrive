# Getting Started with odc

This tutorial will guide you through installing `odc`, setting up your first profile, and running your first command.

## Prerequisites

- **Go:** Ensure you have Go 1.25 or higher installed.
- **Just:** Recommended for building the project.

## Installation

Currently, the recommended way to install `odc` is to build it from source.

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/michaeldcanady/go-onedrive.git
    cd go-onedrive
    ```

2.  **Build the application:**
    ```bash
    just build
    ```

3.  **Verify the installation:**
    ```bash
    ./odc --help
    ```

## Initial Setup

Before you can interact with your OneDrive, you need to authenticate.

1.  **Create a profile (optional):**
    By default, `odc` uses a profile named `default`. You can create a new profile for a specific account:
    ```bash
    ./odc profile create my-work
    ```

2.  **Login:**
    Run the login command to authenticate your current profile:
    ```bash
    ./odc auth login
    ```
    Follow the instructions in your terminal to complete the authentication in your web browser.

## Your First Commands

Now that you're logged in, try listing the files in your root OneDrive directory:

```bash
./odc ls
```

To create a new directory:

```bash
./odc mkdir MyNewFolder
```

To upload a local file:

```bash
./odc upload local-file.txt MyNewFolder/remote-file.txt
```

Congratulations! You are now set up to use `odc`.

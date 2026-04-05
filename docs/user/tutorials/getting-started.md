# Getting Started with odc

`odc` (OneDrive CLI) is a Unix-style command-line tool designed to interact with
Microsoft OneDrive. This guide covers how to install the tool, set up your
first profile, and run your first commands.

## Prerequisites

Before installing `odc`, ensure your system meets the following requirements:

- **Go:** You must have Go 1.25 or higher installed.
- **Just:** We recommend installing `just` to automate the build process.
- **Git:** You'll need Git to clone the repository.

## Installation

Currently, the recommended way to install `odc` is to build it from the source.

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

## Terminal-Native Documentation (Man Pages)

`odc` can generate traditional Unix man pages for all its commands, allowing
you to access documentation directly from your terminal without a browser.

1.  **Generate man pages:**
    ```bash
    just generate-man
    ```
    This creates a `./man` directory containing man pages for all `odc` 
    commands.

2.  **View a man page:**
    You can view a specific page using the `man` command:
    ```bash
    man -l ./man/odc.1
    man -l ./man/odc-ls.1
    ```

3.  **Install man pages (optional):**
    To make man pages available globally, you can copy them to your system's
    man directory (typically `/usr/local/share/man/man1`):
    ```bash
    sudo cp ./man/*.1 /usr/local/share/man/man1/
    ```
    After installing, you can simply run `man odc` from anywhere.

> **Note:** For easier access to the binary, you can move the `odc` binary to 
> a directory in your `PATH`, such as `/usr/local/bin` or `~/bin`.

## Core Concepts

Understanding these two concepts will help you navigate your OneDrive files
more effectively.

### Profiles
A **profile** represents a single OneDrive account (e.g., your personal account
or your work account). `odc` lets you manage multiple profiles and switch
between them easily.

### Drives
A **drive** is a storage area within a profile. While your personal OneDrive is
the default drive, you can also access shared folders, SharePoint libraries,
and other storage areas as separate drives.

## Initial Setup

Before you can interact with your OneDrive, you must authenticate.

1.  **Create a profile (optional):**
    By default, `odc` uses a profile named `default`. If you want to use a
    different name, you can create a new profile:
    ```bash
    ./odc profile create my-work
    ```

2.  **Log in:**
    Run the login command to authenticate your active profile:
    ```bash
    ./odc auth login
    ```
    Follow the instructions in your terminal to complete the authentication in
    your web browser. You'll be asked to sign in to your Microsoft account and
    grant `odc` permission to access your files.

## Your First Commands

Now that you've logged in, try these commands to explore your OneDrive:

1.  **List files:** View the contents of your root OneDrive directory.
    ```bash
    ./odc ls
    ```

2.  **Create a folder:** Create a new directory for your projects.
    ```bash
    ./odc mkdir /MyNewFolder
    ```

3.  **Upload a file:** Move a local file to your new folder.
    ```bash
    ./odc upload local-file.txt /MyNewFolder/remote-file.txt
    ```

Congratulations! You're now set up and ready to use `odc`.

## Next steps

- **[Basic operations](basic-operations.md)**
- **[Working with drives](../how-to/work-with-drives.md)**
- **[View the command reference](../reference/cli-commands.md)**

# Go Pismo Challenge

## Description

Challenge
Transaction Routine

- Each cardholder (customer) has an account with their data.
- For each operation carried out by the customer, a transaction is created and associated
with this account.
- Each transaction has a type (cash purchase, installment purchase, withdrawal, or
payment), an amount, and a creation date.
- Purchase and withdrawal transactions are recorded with a negative value, while
payment transactions are recorded with a positive value.

## Getting Started

### Prerequisites

- Docker
- Make
- Any IDE

### Installation

1. Clone the repository:

    ```bash
    git clone https://github.com/raxitchauhan/go-pismo-challenge.git
    ```

2. Install dependencies:

### Make
#### On Linux (Debian/Ubuntu-based systems):
```bash
sudo apt update
```
- Install build-essential (includes Make and other necessary tools):
```bash
sudo apt install build-essential
```
- Alternatively, you can specifically install make:
```bash
sudo apt install make
```

Verify installation.
```bash
make -v
```

### Docker

https://www.docker.com/get-started/

### IDE

Suggested: Visual Studio Code

https://code.visualstudio.com/download

#### On Windows:

Using Chocolatey:

Install Chocolatey: (if not already installed) by following instructions on its official website.

Install Make:
```bash
choco install make
```

Restart your terminal or Git Bash session after installation.

Manual Installation (from GnuWin32):

- Download the make executable from a source like GnuWin32.
- Run the installer.
- Add the bin directory of the GnuWin32 installation (e.g., `C:\Program Files (x86)\GnuWin32\bin`) to your system's PATH environment variable.
- Restart any open command prompts or PowerShell windows.
- Verify with make -v.


#### On macOS:

Make is typically included with Xcode Command Line Tools. Install Xcode Command Line Tools.

```bash
xcode-select --install
```

Verify installation.


```bash
make -v
```

## Usage
### Available Makefile Commands

Available commands in the Makefile for managing and working with the project.

---

### `make boot`

**Description:**

Boots up the service and makes it available for testing and development. You can interact with the APIs and access documentation.

**Usage:**

```bash
make boot
```

What it does:

- Starts the application on http://localhost:3000, making the API available for use.
- Provides access to Swagger API Docs at http://localhost:3000/swagger/index.html for easy exploration of available endpoints.
- Useful for development or testing the service in a local environment.

### `make test`

****Description:****  
Runs all the tests in the project, including lint checks and unit tests.

****Usage:****
```bash
make test
```

What it does:

- Runs lint to check for any coding style issues or errors.
- Executes unit tests to verify that the individual components are working as expected.

### `make run-migration`

**Description:**

Sets up a container and runs database migrations within it. Useful for setting up the correct database schema for your application.

**Usage:**

```bash
make run-migration
```

What it does:

- Creates a Docker container for the migration process.
- Runs any pending database migrations within the container.
- Ensures the database is up-to-date with the latest schema changes.

### `make unit-test`

**Description:**

Runs all unit tests within a Docker container. This is a good way to ensure that the application logic works in an isolated environment.

**Usage:**

```bash
make unit-test
```

What it does:

- Spins up a Docker container to run all unit tests.
- Executes all the test cases in your test suite to validate the functionality of your application.

### `make lint`

**Description:**

Runs the linter to check for any syntax issues, formatting errors, or potential bugs in your code.

**Usage:**

```bash
make lint
```

What it does:

- Checks the codebase for adherence to the style guide.
- Useful for maintaining code consistency and quality.

### `make down`

**Description:**

Stops and removes all the Docker containers used by the project.

**Usage:**

```bash
make down
```
What it does:

- Removes all running Docker containers related to the project.

Useful for resetting the environment or clearing up resources.

### `make gen`

**Description:**

Generates all necessary mocks and Swagger API documentation for the project.

**Usage:**

```bash
make gen
```

What it does:

- Generates mocks for interface used for testing and simulation.
- Creates Swagger API documentation that outlines all available endpoints, request/response structures, and more.

Ensures the documentation is up-to-date with the latest changes in the API.

### `make dep`

**Description:**

Vendorise the dependencies according to the package manager configuration.

**Usage:**

```bash
make dep
```

What it does:

- Tidy up the go.mod file and sync the vendor directories
# Contributing to terraform-type-list-helper

Thank you for your interest in contributing to `terraform-type-list-helper`! We welcome contributions from everyone.

## Ways to Contribute

- **Report Bugs:** Submit detailed bug reports on the [Issues](https://github.com/uibm/terraform-type-list-helper/issues) page. Include steps to reproduce the bug, expected behavior, and actual behavior.
- **Suggest Enhancements:** Propose new features or improvements through the [Issues](https://github.com/uibm/terraform-type-list-helper/issues) page.
- **Submit Pull Requests:** Fork the repository, create a new branch for your changes, and submit a pull request. Make sure your code follows the existing style and includes unit tests.

## Development Setup

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/uibm/terraform-type-list-helper.git
    cd terraform-type-list-helper
    ```

2.  **Install dependencies:**
    ```bash
    go mod download
    ```

3.  **Run tests:**
    ```bash
    go test ./...
    ```

## Code Style

-   Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) guidelines.
-   Use `gofmt` to format your code before committing.

## Pull Request Process

1.  Ensure your code builds and all tests pass.
2.  Update the `README.md` with details of changes to the interface or usage, if applicable.
3.  Update the `CHANGELOG.md` with a summary of your changes under the "Unreleased" section.
4.  Submit the pull request, clearly describing the changes and their purpose.

We look forward to your contributions!
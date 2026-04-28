# terraform-type-list-helper

`terraform-type-list-helper` is a Go library designed to simplify the management of `TypeList` attributes within Terraform providers built using the Terraform Plugin SDKv2. It provides helper functions to implement `DiffSuppressFunc` and determine changes (add, update, remove) efficiently and generically for `TypeList` attributes that have a constant field used for identification.

## Features

-   **Generic `DiffSuppressFunc`:** Easily suppress unnecessary diffs for `TypeList` attributes based on a specified constant field.
-   **Change Calculation:**  Accurately determines elements to be added, updated, or removed from a `TypeList`.
-   **`ApplyOnlyOnce` Functionality:** Provides a ready to use function to ensure certain fields are applied only once during resource creation.
-   **Robust and Tested:** Includes comprehensive unit tests to ensure reliability.
-   **Production-Ready:** Designed for use in production-grade Terraform providers.

## Installation

```bash
go get github.com/uibm/terraform-type-list-helper
```
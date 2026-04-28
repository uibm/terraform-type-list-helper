# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2023-11-21

### Added

- Initial release of the `terraform-type-list-helper` library.
- `TypeListHelper` struct to manage TypeList attributes.
- `DiffSuppressFunc` function to suppress unnecessary diffs for TypeLists.
- `CalculateChanges` function to determine changes (add, update, remove) between old and new TypeList values.
- `ApplyOnlyOnce` function to ensure a field is applied only during resource creation.
- Comprehensive unit tests for all core functions.
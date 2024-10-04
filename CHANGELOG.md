# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Support for Vault, Embargo, File-Intel and Redact services
- Plugins capabilities
- Auto update check and `update` command
- Profile capabilities to set different credentials
- File pattern support for File-Intel `/reputation`

### Changed

- Path to `ls`, `migrate`, `create`, `login` and `run`. They all are now inside `plugins` subcommand

### Fixed

- Print helps on commands instead of errors

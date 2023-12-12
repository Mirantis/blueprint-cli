# Boundless Controller

## Overview

This is the private repository for the Boundless Controller. The Boundless Controller is a CLI tool that allows you to manage your Boundless clusters. This repo contains the actual source code for compiling the Boundless Controller (`bctl`) binary.

The release binary, docker images, documentation, and source code in this repo are all private. They are considered dev builds and used for testing and development purposes only. Thought should be given, especially for documentation, as to whether changes made should be located in this repo or the public repo.

The public repo for the Boundless Controller is located at [mirantiscontainers/boundless](https://github.com/mirantiscontainers/boundless). When a release is created in mirantiscontainers/boundless-cli, it will push the generated binary to the public repo. The public repo will then create a release with the binary and documentation. Much of the documentation for working with the Boundless Controller can be found in the public repo.

## Releases

Information on creating a release can be found in the [release documentation](docs/creating-a-release.md).

## CI/CD

Information on the CI/CD pipeline can be found in the [CI/CD documentation](docs/CI.md). This is from a developer perspective to understand what automation will run as you interact with the repo.

If you are working on changes for the CI/CD pipeline, take a look at the [CI/CD development documentation](.github/workflows/README.md).

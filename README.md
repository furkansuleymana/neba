# Neba

Neba is an open-source tool designed to provide a user-friendly graphical user interface (GUI) for managing Axis
Communications devices that support the [VAPIX API](https://developer.axis.com/vapix). This project is currently under
active development. Breaking changes will occur.

## Overview

Neba aims to simplify the management of Axis products, especially in environments with multiple devices. It operates on
a server-client architecture and features a lightweight, cross-platform design.

The backend, written in Go, is compiled for Windows, macOS, and Linux, running as a daemon on a server located at the
user's site. This server can be the same machine hosting Axis Camera Station or a single-board computer like a Raspberry
Pi. It uses a local database to store configuration and operational data.

The frontend is created with web technologies and is accessible through any modern web browser. To simplify installation
and usage, it’s embedded directly into the Go binary. The interface is designed to be responsive, making it easy to
manage devices from smartphones, tablets, or desktop computers.

## Installation

Neba is not yet operational.

## Development

Ensure you have [Go 1.23](https://go.dev/doc/install) or later installed. To build the application, simply execute
`go run ./make.go production` in the main project directory. Additional build commands can be found in the [`make.go`](make.go) file.

## Key Features

- Discover all Axis products using Bonjour and SSDP
- Perform factory resets or restart devices
- Retrieve server reports, system logs, or client logs

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Disclaimer

The name "Neba" stands for "Not Endorsed By Axis", reflecting that while the tool is designed for Axis products, it is
an independent project with no official affiliation or endorsement from Axis Communications. Use it at your own risk.
The maintainer is not responsible for any issues that arise from its use.

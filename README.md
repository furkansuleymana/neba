# Neba

Neba is an open-source tool designed to provide a user-friendly graphical user interface (GUI) for managing Axis Communications devices that support the VAPIX API. The name "Neba" stands for "Not Endorsed By Axis". While the tool is designed for Axis products, it is an independent project with no official affiliation or endorsement from Axis Communications.

This project is currently under active development. Breaking changes will occur.

## Overview

Neba aims to simplify the management of Axis products, especially in environments with multiple devices. It operates on a server-client architecture and features a lightweight, cross-platform design.

Written in Go, the backend is compiled for Windows, macOS, and Linux, running as a daemon on a server at the user's site. That server can be the same machine running Axis Camera Station or a single-board computer like a Raspberry Pi. Acting as a bridge between the Neba frontend and the VAPIX API, the backend ensures the execution of all commands, including recurring events. Configuration and operational data are stored in a SQLite database.

The frontend, developed using web technologies, offers access through any modern web browser. It is designed to be responsive, allowing users to manage devices conveniently from smartphones or tablets.

## Key Features

- Discover all Axis products using Bonjour and SSDP
- Perform factory resets or restart devices
- Retrieve server reports, system logs, or client logs

## Installation

Neba is not yet operational.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Disclaimer

Neba is not endorsed by or affiliated with Axis Communications. Use it at your own risk. The maintainer is not responsible for any issues that arise from its use.

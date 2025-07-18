<img src="assets/kongvisor.png" alt="kongvisor">

# KongVisor

[![Release](https://img.shields.io/github/release/mchlumsky/kongvisor.svg)](https://github.com/mchlumsky/kongvisor/releases/latest)

KongVisor is a TUI for Kong Gateway Admin API.

It supports managing Kong Gateway resources like Workspaces, Services, Routes, and Plugins in a terminal user interface.

## Installation

```bash
go install github.com/mchlumsky/kongvisor@latest
```

## Configuration

KongVisor supports multiple Kong Gateways, each with its own configuration.

It will read this file to find the Kong Gateway Admin API URLs and tokens.

Configuration goes into `$HOME/.config/kongvisor/config.yaml`.

Example configuration:

```yaml
kong-gateway-1: # name of a kong gateway
  url: http://kong-gateway-1.domain:8001 # Kong Admin API URL
  kongAdminToken: secret-token # Kong Admin API token
kong-gateway-2:
  url: http://kong-gateway-2.domain:8001
  kongAdminToken: other-secret-token

```

## Usage

Run KongVisor with the name of the Kong Gateway you want to manage:

```bash
kongvisor kong-gateway-1
```

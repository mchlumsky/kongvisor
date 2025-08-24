# KongVisor

<img src="assets/kongvisor.png" alt="kongvisor">

[![Release](https://img.shields.io/github/release/mchlumsky/kongvisor.svg)](https://github.com/mchlumsky/kongvisor/releases/latest)
[![codecov](https://codecov.io/gh/mchlumsky/kongvisor/graph/badge.svg?token=Nn2v1Xt2MU)](https://codecov.io/gh/mchlumsky/kongvisor)
[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg)](/LICENSE)
[![Build status](https://img.shields.io/github/actions/workflow/status/mchlumsky/kongvisor/build.yml?branch=main)](https://github.com/mchlumsky/kongvisor/actions?workflow=build)
[![Powered By: GoReleaser](https://img.shields.io/badge/powered%20by-goreleaser-green.svg)](https://github.com/goreleaser)

KongVisor is a TUI for Kong Gateway Admin API.

It supports managing Kong Gateway resources like Workspaces, Services, Routes, and Plugins in a terminal user interface.

## Demo

### Listing resources

List resources (workspaces, services, routes, plugins). To see services in a workspace, scroll down to the workspace and press "s". To see routes in a service, scroll down to the service and press "r". To see plugins in a workspace, service or route, scroll down to the workspace, service or route and press "p".

![Listing resources](assets/list.gif)

### Viewing resources

View resource details (in YAML of course!). Scroll down to the resource you want to see and press "y". Press esc to go back.

![Viewing resources](assets/view.gif)

### Deleting resources

Delete resources (**without** confirmation). Scroll down to the resource you want to delete and press "Ctrl-d".

![Deleting resources](assets/delete.gif)

### Editing resources

Edit resources (in your favorite editor defined by the `$EDITOR` environment variable). Scroll down to the resource you want to edit and press "e".

![Editing resources](assets/edit.gif)

## Installation

### YUM/DNF repository (Fedora, CentOS, etc...)

To enable, add the following file /etc/yum.repos.d/mchlumsky.repo:

```
[mchlumsky]
name=mchlumsky repository (kongvisor, etc...)
baseurl=https://yum.fury.io/mchlumsky/
enabled=1
gpgcheck=0
```

and

```shell
dnf install kongvisor
# or
yum install kongvisor
```

### APT/DEB repository (Debian, Ubuntu, etc...)

To enable, add the following file /etc/apt/sources.list.d/mchlumsky.list:

```
deb [trusted=yes] https://apt.fury.io/mchlumsky/ /
```

and

```shell
apt update
apt install kongvisor
```

### Zypper repository (OpenSUSE)

To enable, run the following command:

```shell
zypper addrepo --no-gpgcheck https://yum.fury.io/mchlumsky/ mchlumsky
zypper install kongvisor
```

### Binary Releases

Download and install a prebuilt binary from [github releases](https://github.com/mchlumsky/kongvisor/releases)

### From source

Install from source:

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

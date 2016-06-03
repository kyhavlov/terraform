---
layout: "docker"
page_title: "Provider: Docker"
sidebar_current: "docs-docker-index"
description: |-
  The Docker provider is used to interact with Docker containers and images.
---

# Docker Provider

The Docker provider is used to interact with Docker containers and images.
It uses the Docker API to manage the lifecycle of Docker containers. Because
the Docker provider uses the Docker API, it is immediately compatible not
only with single server Docker but Swarm and any additional Docker-compatible
API hosts.

Use the navigation to the left to read about the available resources.

<div class="alert alert-block alert-info">
<strong>Note:</strong> The Docker provider is new as of Terraform 0.4.
It is ready to be used but many features are still being added. If there
is a Docker feature missing, please report it in the GitHub repo.
</div>

## Example Usage

```
# Configure the Docker provider
provider "docker" {
    host = "tcp://127.0.0.1:1234/"
}

# Create a container
resource "docker_container" "foo" {
    image = "${docker_image.ubuntu.latest}"
    name = "foo"
}

resource "docker_image" "ubuntu" {
    name = "ubuntu:latest"
}
```

## With Registry Credentials

```
# Configure the Docker provider
provider "docker" {
    host = "tcp://127.0.0.1:1234/"

    registry_auth {
        address = "https://registry.hub.docker.com"
        username = "someuser"
        password = "somepass"
    }

    registry_auth {
        address = "https://privateregistry.example.com:5000"
        config_file = "~/.docker/config.json"
    }
}

# Read image digest from the registry
data "docker_registry_image" "ubuntu" {
    name = "someuser/private-ubuntu:precise"
}

resource "docker_image" "ubuntu" {
    name = "${data.docker_registry_image.ubuntu.name}"
    pull_trigger = "${data.docker_registry_image.ubuntu.sha256_digest}"
}

resource "docker_container" "foo" {
    image = "${docker_image.ubuntu.latest}"
    name = "foo"
}
```

## Argument Reference

The following arguments are supported:

* `host` - (Required) This is the address to the Docker host. If this is
  blank, the `DOCKER_HOST` environment variable will also be read.

* `cert_path` - (Optional) Path to a directory with certificate information
  for connecting to the Docker host via TLS. If this is blank, the
  `DOCKER_CERT_PATH` will also be checked.

* `registry_auth` - (Optional, block) See [Registry Auth](#registry_auth) below for details.

<a id="registry_auth"></a>
### Registry Auth

`registry_auth` is a block within the configuration that can be repeated to specify
credentials for a specific Docker registry. If neither username/password nor a config file
path are given, Terraform checks `~/.docker/config.json` for credentials. Each
`registry_auth` block supports the following:

* `address` - (Required, string) Address of the Docker registry.
* `username` - (Optional, string) Username to use for the registry. If this is blank,
  `DOCKER_REGISTRY_USER` will also be checked.
* `password` - (Optional, string) Password to use for the registry. If this is blank,
  `DOCKER_REGISTRY_PASS` will also be checked.
* `config_file` - (Optional, string) Location of a docker config.json file containing
  registry credentials. If this is blank, `DOCKER_CONFIG` will also be checked. Defaults
  to `~/.docker/config.json`.
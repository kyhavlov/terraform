---
layout: "docker"
page_title: "Docker: docker_image"
sidebar_current: "docs-docker-datasource-image"
description: |-
  Finds the latest available version for the image/tag from the registry.
---

# docker\_image

Reads the current image ID from the registry. Used in conjunction with the 
[docker\_image](/docs/providers/docker/r/image.html) resource to to keep up 
to date on the latest available version of the image/tag.

## Example Usage

```
data "docker_image" "ubuntu" {
    name = "ubuntu:precise"
}

resource "docker_image" "ubuntu" {
    name = "${data.docker_image.ubuntu.name}"
    registry_id = "${data.docker_image.ubuntu.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, string) The name of the Docker image, including any tags.

## Attributes Reference

The following attributes are exported in addition to the above configuration:

* `id` (string) - The ID of the image, as stored on the registry.
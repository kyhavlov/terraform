---
layout: "docker"
page_title: "Docker: docker_image"
sidebar_current: "docs-docker-resource-image"
description: |-
  Downloads and exports the ID of a Docker image.
---

# docker\_image

Downloads and exports the ID of a Docker image. This can be used alongside
[docker\_container](/docs/providers/docker/r/container.html)
to programmatically get the latest image IDs without having to hardcode
them.

# Example Usage

### Static image tag
```
resource "docker_image" "ubuntu" {
    name = "ubuntu:precise"
}

# Access it somewhere else with ${docker_image.ubuntu.latest}
```

### Dynamic image tag
```
data "docker_image" "ubuntu" {
    name = "ubuntu:precise"
}

resource "docker_image" "ubuntu" {
    name = "${data.docker_image.ubuntu.name}"
    registry_id = "${data.docker_image.ubuntu.id}"
}

# Access it somewhere else with ${docker_image.ubuntu.latest}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, string) The name of the Docker image, including any tags.
* `keep_updated` - (Optional, boolean) If true, then the Docker image will
  always be updated on the host to the latest. If this is false, as long as an
  image is downloaded with the correct tag, it won't be redownloaded if
  there is a newer image.
* `keep_locally` - (Optional, boolean) If true, then the Docker image won't be
  deleted on destroy operation. If this is false, it will delete the image from
  the docker local storage on destroy operation.
* `registry_id` - (Optional, string) Used to store the image ID from the registry
  and will cause a recreate when changed. Needed when using the docker image
  [data source](/docs/providers/docker/d/image.html) to trigger an update of the 
  image.

## Attributes Reference

The following attributes are exported in addition to the above configuration:

* `latest` (string) - The ID of the image.

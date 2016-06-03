package docker

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

var registryIdRegexp = regexp.MustCompile(`\A([a-zA-Z0-9_-]){64}\z`)

func TestAccDockerImageDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDockerImageDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.docker_image.foo", "id", registryIdRegexp),
				),
			},
		},
	})
}

const testAccDockerImageDataSourceConfig = `
data "docker_image" "foo" {
	name = "alpine:latest"
}
`

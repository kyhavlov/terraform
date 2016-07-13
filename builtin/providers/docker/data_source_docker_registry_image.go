package docker

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceDockerRegistryImage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDockerRegistryImageRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"sha256_digest": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDockerRegistryImageRead(d *schema.ResourceData, meta interface{}) error {
	providerConfig := meta.(*ProviderConfig)
	pullOpts := parseImageOptions(d.Get("name").(string))

	// Use the official Docker Hub if a registry isn't specified
	if pullOpts.Registry == "" {
		pullOpts.Registry = "registry.hub.docker.com"
	} else {
		// Otherwise, filter the registry name out of the repo name
		pullOpts.Repository = strings.Replace(pullOpts.Repository, pullOpts.Registry+"/", "", 1)
	}

	// Docker prefixes 'library' to official images in the path; 'consul' becomes 'library/consul'
	if !strings.Contains(pullOpts.Repository, "/") {
		pullOpts.Repository = "library/" + pullOpts.Repository
	}

	if pullOpts.Tag == "" {
		pullOpts.Tag = "latest"
	}

	registryClient, err := providerConfig.getRegistryClient(pullOpts)

	if err != nil {
		return fmt.Errorf("Error getting registry client: %v", err)
	}

	manifest, err := registryClient.Manifest(pullOpts.Repository, pullOpts.Tag)
	if err != nil {
		return fmt.Errorf("Error getting manifest for image: %v", err)
	}

	m := map[string]interface{}{}
	err = json.Unmarshal([]byte(manifest.History[0].V1Compatibility), &m)
	if err != nil {
		return fmt.Errorf("Error parsing manifest for image %s: %v", d.Get("name").(string), err)
	}

	if _, ok := m["id"]; !ok {
		return fmt.Errorf("Couldn't get image id from manifest for image: %s:%s", pullOpts.Repository, pullOpts.Tag)
	}

	d.SetId(m["id"].(string))
	d.Set("sha256_digest", "sha256:"+m["id"].(string))

	return nil
}

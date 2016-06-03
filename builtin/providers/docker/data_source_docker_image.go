package docker

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/heroku/docker-registry-client/registry"
)

func dataSourceDockerImage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDockerImageRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDockerImageRead(d *schema.ResourceData, meta interface{}) error {
	pullOpts := parseImageOptions(d.Get("name").(string))
	auth := getAuthConfig(pullOpts)

	// Use the official Docker Hub if a registry isn't specified
	if pullOpts.Registry == "" {
		pullOpts.Registry = "registry.hub.docker.com"
	} else {
		pullOpts.Repository = strings.Replace(pullOpts.Repository, pullOpts.Registry+"/", "", 1)
	}

	// The docker registry prefixes 'library' to official images in the path; 'consul' becomes 'library/consul'
	if !strings.Contains(pullOpts.Repository, "/") {
		pullOpts.Repository = "library/" + pullOpts.Repository
	}

	if pullOpts.Tag == "" {
		pullOpts.Tag = "latest"
	}

	url := "https://" + pullOpts.Registry
	username := auth.Username
	password := auth.Password
	hub, err := registry.New(url, username, password)

	if err != nil {
		return fmt.Errorf("Error connecting to registry: %v", err)
	}

	manifest, err := hub.Manifest(pullOpts.Repository, pullOpts.Tag)
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
	d.Set("id", m["id"].(string))

	return nil
}

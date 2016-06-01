package docker

import (
	"encoding/json"
	"fmt"
	"log"
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

	log.Printf("[DEBUG] Got pull and auth config: %v\n%v", pullOpts, auth)

	url := "https://" + pullOpts.Registry //"https://registry.hub.docker.com"
	log.Printf("[DEBUG] Registry url: %s", url)
	username := auth.Username
	password := auth.Password
	hub, _ := registry.New(url, username, password)

	if pullOpts.Registry != "" {
		pullOpts.Repository = strings.Replace(pullOpts.Repository, pullOpts.Registry+"/", "", 1)
	}
	tags, _ := hub.Tags(pullOpts.Repository)
	log.Printf("[DEBUG] Tags: %v", tags)
	manifest, _ := hub.Manifest(pullOpts.Repository, pullOpts.Tag)
	//log.Printf("[DEBUG] Manifest: %v", manifest)

	m := map[string]interface{}{}
	err := json.Unmarshal([]byte(manifest.History[0].V1Compatibility), &m)
	if err != nil {
		return fmt.Errorf("[ERROR] Couldn't parse registry info for image (%s): %v", d.Get("name").(string), err)
	}
	log.Printf("[DEBUG] id: %v", m["id"])
	d.SetId(m["id"].(string))
	d.Set("id", m["id"].(string))

	return nil
}

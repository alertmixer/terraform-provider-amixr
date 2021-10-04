package amixr

import (
	"fmt"
	"github.com/alertmixer/amixr-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
)

func dataSourceAmixrTeam() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAmixrTeamRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"email": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"avatar_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAmixrTeamRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] read team")

	client := m.(*amixr.Client)
	options := &amixr.ListTeamOptions{}
	nameData := d.Get("name").(string)

	options.Name = nameData

	teamsResponse, _, err := client.Teams.ListTeams(options)

	if err != nil {
		return err
	}

	if len(teamsResponse.Teams) == 0 {
		return fmt.Errorf("couldn't find a team matching: %s", options.Name)
	} else if len(teamsResponse.Teams) != 1 {
		return fmt.Errorf("more than one team found matching: %s", options.Name)
	}

	team := teamsResponse.Teams[0]

	d.Set("name", team.Name)
	d.Set("email", team.Email)
	d.Set("avatar_url", team.AvatarUrl)

	d.SetId(team.ID)

	return nil
}

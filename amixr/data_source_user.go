package amixr

import (
	"fmt"
	amixr "github.com/alertmixer/amixr-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func dataSourceAmixrUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAmixrUserRead,
		Schema: map[string]*schema.Schema{
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"team_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"role": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAmixrUserRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] read amixr user")

	client := m.(*amixr.Client)
	options := &amixr.ListUserOptions{}
	emailData := d.Get("email").(string)

	options.Email = emailData

	usersResponse, _, err := client.Users.ListUsers(options)

	if err != nil {
		return err
	}

	if len(usersResponse.Users) == 0 {
		return fmt.Errorf("couldn't find a user matching: %s", options.Email)
	} else if len(usersResponse.Users) != 1 {
		return fmt.Errorf("more than one user found matching: %s", options.Email)
	}

	user := usersResponse.Users[0]

	d.Set("email", user.Email)
	d.Set("name", user.Name)
	d.Set("role", user.Role)
	d.Set("team_id", user.TeamId)

	d.SetId(user.ID)

	return nil
}

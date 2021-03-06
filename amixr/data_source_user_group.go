package amixr

import (
	"fmt"
	amixr "github.com/alertmixer/amixr-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func dataSourceAmixrUserGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAmixrUserGroupRead,
		Schema: map[string]*schema.Schema{
			"slack_handle": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceAmixrUserGroupRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] read amixr user group")

	client := m.(*amixr.Client)
	options := &amixr.ListUserGroupOptions{}
	slackHandleData := d.Get("slack_handle").(string)

	options.SlackHandle = slackHandleData

	userGroupsResponse, _, err := client.UserGroups.ListUserGroups(options)

	if err != nil {
		return err
	}

	if len(userGroupsResponse.UserGroups) == 0 {
		return fmt.Errorf("couldn't find a user group matching: %s", options.SlackHandle)
	} else if len(userGroupsResponse.UserGroups) != 1 {
		return fmt.Errorf("couldn't find a user group matching: %s", options.SlackHandle)
	}

	user_group := userGroupsResponse.UserGroups[0]

	d.SetId(user_group.ID)

	return nil
}

package amixr

//
//import (
//	"fmt"
//	amixr "github.com/alertmixer/amixr-go-client"
//	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
//	"log"
//)
//
//func dataSourceAmixrUserGroup() *schema.Resource {
//	return &schema.Resource{
//		Read: dataSourceAmixrUserGroupRead,
//		Schema: map[string]*schema.Schema{
//			"slack_name": {
//				Type:     schema.TypeString,
//				Required: true,
//			},
//		},
//	}
//}
//
//func dataSourceAmixrUserGroupRead(d *schema.ResourceData, m interface{}) error {
//	log.Printf("[DEBUG] read amixr user group")
//
//	client := m.(*amixr.Client)
//	options := &amixr.ListUserGroupOptions{}
//	slackData := d.Get("slack").([]interface{})
//
//
//	options.SlackName = slackData[0].(map[string]interface{})["name"]
//
//	userGroupsResponse, _, err := client.UserGroups.ListUserGroups(options)
//
//	if err != nil {
//		return err
//	}
//
//	if len(userGroupsResponse.UserGroups) == 0 {
//		return fmt.Errorf("couldn't find a user group matching: %s", options.SlackName)
//	} else if len(userGroupsResponse.UserGroups) != 1 {
//		return fmt.Errorf("couldn't find a user group matching: %s", options.SlackName)
//	}
//
//	user_group := userGroupsResponse.UserGroups[0]
//
//	d.SetId(user_group.ID)
//	d.Set("type", user_group.Type)
//	d.Set("slack", flattenSlackUserGroup(user_group.SlackUserGroup))
//
//	return nil
//}
//
//func flattenSlackUserGroup(in *amixr.SlackUserGroup) []map[string]interface{} {
//	slackUserGroups := make([]map[string]interface{}, 0, 1)
//
//	out := make(map[string]interface{})
//	out["name"] = in.Name
//	out["id"] = in.ID
//
//	slackUserGroups = append(slackUserGroups, out)
//
//	return slackUserGroups
//}

package amixr

import (
	"fmt"
	"github.com/alertmixer/amixr-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
)

var scheduleTypeOptions = []string{
	"ical",
	"calendar",
}

func resourceSchedule() *schema.Resource {
	return &schema.Resource{
		Create: resourceScheduleCreate,
		Read:   resourceScheduleRead,
		Update: resourceScheduleUpdate,
		Delete: resourceScheduleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(scheduleTypeOptions, false),
			},
			"time_zone": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"ical_url_primary": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"ical_url_overrides": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"slack": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"channel_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"user_group_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				MaxItems: 1,
			},
			"shifts": &schema.Schema{
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
		},
	}
}

func resourceScheduleCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] create amixr schedule")

	client := m.(*amixr.Client)

	nameData := d.Get("name").(string)
	typeData := d.Get("type").(string)
	slackData := d.Get("slack").([]interface{})

	createOptions := &amixr.CreateScheduleOptions{
		Name:  nameData,
		Type:  typeData,
		Slack: expandScheduleSlack(slackData),
	}

	iCalUrlPrimaryData, iCalUrlPrimaryOk := d.GetOk("ical_url_primary")
	if iCalUrlPrimaryOk {
		if typeData == "ical" {
			url := iCalUrlPrimaryData.(string)
			createOptions.ICalUrlPrimary = &url
		} else {
			return fmt.Errorf("ical_url_primary can not be set with type: %s", typeData)
		}
	}

	iCalUrlOverridesData, iCalUrlOverridesOk := d.GetOk("ical_url_primary")
	if iCalUrlOverridesOk {
		if typeData == "ical" {
			url := iCalUrlOverridesData.(string)
			createOptions.ICalUrlPrimary = &url
		} else {
			return fmt.Errorf("ical_url_overrides can not be set with type: %s", typeData)
		}
	}

	shiftsData, shiftsOk := d.GetOk("shifts")
	if shiftsOk {
		if typeData == "calendar" {
			createOptions.Shifts = stringSetToStringSlice(shiftsData.(*schema.Set))
		} else {
			return fmt.Errorf("shifts can not be set with type: %s", typeData)
		}
	}

	timeZoneData, timeZoneOk := d.GetOk("time_zone")
	if timeZoneOk {
		if typeData == "calendar" {
			createOptions.TimeZone = timeZoneData.(string)
		} else {
			return fmt.Errorf("time_zone can not be set with type: %s", typeData)
		}
	}

	schedule, _, err := client.Schedules.CreateSchedule(createOptions)
	if err != nil {
		return err
	}

	d.SetId(schedule.ID)

	return resourceScheduleRead(d, m)
}

func resourceScheduleUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] update amixr schedule")

	client := m.(*amixr.Client)

	nameData := d.Get("name").(string)
	slackData := d.Get("slack").([]interface{})
	typeData := d.Get("type").(string)

	updateOptions := &amixr.UpdateScheduleOptions{
		Name:  nameData,
		Slack: expandScheduleSlack(slackData),
	}

	iCalUrlPrimaryData, iCalUrlPrimaryOk := d.GetOk("ical_url_primary")
	if iCalUrlPrimaryOk {
		if typeData == "ical" {
			url := iCalUrlPrimaryData.(string)
			updateOptions.ICalUrlPrimary = &url
		} else {
			return fmt.Errorf("ical_url_primary can not be set with type: %s", typeData)
		}
	}

	iCalUrlOverridesData, iCalUrlOverridesOk := d.GetOk("ical_url_primary")
	if iCalUrlOverridesOk {
		if typeData == "ical" {
			url := iCalUrlOverridesData.(string)
			updateOptions.ICalUrlPrimary = &url
		} else {
			return fmt.Errorf("ical_url_overrides can not be set with type: %s", typeData)
		}
	}

	timeZoneData, timeZoneOk := d.GetOk("time_zone")
	if timeZoneOk {
		if typeData == "calendar" {
			updateOptions.TimeZone = timeZoneData.(string)
		} else {
			return fmt.Errorf("time_zone can not be set with type: %s", typeData)
		}
	}

	shiftsData, shiftsOk := d.GetOk("shifts")
	if shiftsOk {
		if typeData == "calendar" {
			updateOptions.Shifts = stringSetToStringSlice(shiftsData.(*schema.Set))
		} else {
			return fmt.Errorf("shifts can not be set with type: %s", typeData)
		}
	}

	schedule, _, err := client.Schedules.UpdateSchedule(d.Id(), updateOptions)
	if err != nil {
		return err
	}

	d.SetId(schedule.ID)

	return resourceScheduleRead(d, m)
}

func resourceScheduleRead(d *schema.ResourceData, m interface{}) error {

	client := m.(*amixr.Client)
	options := &amixr.GetScheduleOptions{}
	schedule, _, err := client.Schedules.GetSchedule(d.Id(), options)

	if err != nil {
		return err
	}

	d.Set("name", schedule.Name)
	d.Set("type", schedule.Type)
	d.Set("ical_url_primary", schedule.ICalUrlPrimary)
	d.Set("ical_url_overrides", schedule.ICalUrlOverrides)
	d.Set("time_zone", schedule.TimeZone)
	d.Set("slack", flattenScheduleSlack(schedule.Slack))
	d.Set("shifts", schedule.Shifts)

	return nil
}

func resourceScheduleDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] delete amixr schedule")

	client := m.(*amixr.Client)
	options := &amixr.DeleteScheduleOptions{}
	_, err := client.Schedules.DeleteSchedule(d.Id(), options)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func flattenScheduleSlack(in *amixr.SlackSchedule) []map[string]interface{} {
	slack := make([]map[string]interface{}, 0, 1)

	out := make(map[string]interface{})

	if in.ChannelId != nil {
		out["channel_id"] = in.ChannelId
	}

	if in.UserGroupId != nil {
		out["user_group_id"] = in.UserGroupId
	}

	slack = append(slack, out)
	return slack
}

func expandScheduleSlack(in []interface{}) *amixr.SlackSchedule {
	slackSchedule := amixr.SlackSchedule{}

	for _, r := range in {
		inputMap := r.(map[string]interface{})
		if inputMap["channel_id"] != "" {
			channelId := inputMap["channel_id"].(string)
			slackSchedule.ChannelId = &channelId
		}
		if inputMap["user_group_id"] != "" {
			userGroupId := inputMap["user_group_id"].(string)
			slackSchedule.UserGroupId = &userGroupId
		}
	}

	return &slackSchedule
}

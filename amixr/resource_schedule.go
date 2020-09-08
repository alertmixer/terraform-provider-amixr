package amixr

import (
	amixr "github.com/alertmixer/amixr-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
)

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
			"time_zone": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "UTC",
			},
			"ical_url": &schema.Schema{
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
					},
				},
				MaxItems: 1,
			},
		},
	}
}

func resourceScheduleCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] create amixr schedule")

	client := m.(*amixr.Client)

	nameData := d.Get("name").(string)
	slackData := d.Get("slack").([]interface{})

	createOptions := &amixr.CreateScheduleOptions{
		Name:  nameData,
		Slack: expandScheduleSlack(slackData),
	}

	iCalUrlData, iCalUrlOk := d.GetOk("ical_url")
	if iCalUrlOk {
		url := iCalUrlData.(string)
		createOptions.ICalUrl = &url
	}

	timeZoneData, timeZoneOk := d.GetOk("time_zone")
	if timeZoneOk {
		createOptions.TimeZone = timeZoneData.(string)
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

	updateOptions := &amixr.UpdateScheduleOptions{
		Name:  nameData,
		Slack: expandScheduleSlack(slackData),
	}

	iCalUrlData, iCalUrlOk := d.GetOk("ical_url")
	if iCalUrlOk {
		url := iCalUrlData.(string)
		updateOptions.ICalUrl = &url
	}

	timeZoneData, timeZoneOk := d.GetOk("time_zone")
	if timeZoneOk {
		updateOptions.TimeZone = timeZoneData.(string)
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
	d.Set("ical_url", schedule.ICalUrl)
	d.Set("time_zone", schedule.TimeZone)
	d.Set("slack", flattenScheduleSlack(schedule.Slack))

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
		slack = append(slack, out)
	}
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
	}

	return &slackSchedule
}

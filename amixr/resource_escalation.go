package amixr

import (
	"fmt"
	amixr "github.com/alertmixer/amixr-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
)

var escalationOptions = []string{
	"wait",
	"notify_persons",
	"notify_person_next_each_time",
	"notify_on_call_from_schedule",
	"trigger_action",
	"notify_user_group",
	"resolve",
	"notify_whole_channel",
	"notify_if_time_from_to",
}

var stepsWithImportant = []string{
	"notify_persons",
	"notify_user_group",
	"notify_on_call_from_schedule",
}

var durationOptions = []int{
	60,
	300,
	900,
	1800,
	3600,
}

func resourceEscalation() *schema.Resource {
	return &schema.Resource{
		Create: resourceEscalationCreate,
		Read:   handleNonExistentResource(resourceEscalationRead),
		Update: resourceEscalationUpdate,
		Delete: resourceEscalationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"route_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"position": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(escalationOptions, false),
			},
			"important": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"duration": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ConflictsWith: []string{
					"notify_on_call_from_schedule",
					"persons_to_notify",
					"persons_to_notify_next_each_time",
					"action_to_trigger",
					"group_to_notify",
					"notify_if_time_from",
					"notify_if_time_to",
				},
				ValidateFunc: validation.IntInSlice(durationOptions),
			},
			"notify_on_call_from_schedule": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ConflictsWith: []string{
					"duration",
					"persons_to_notify",
					"persons_to_notify_next_each_time",
					"action_to_trigger",
					"group_to_notify",
					"notify_if_time_from",
					"notify_if_time_to",
				},
			},
			"persons_to_notify": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				ConflictsWith: []string{
					"duration",
					"notify_on_call_from_schedule",
					"persons_to_notify_next_each_time",
					"action_to_trigger",
					"group_to_notify",
					"notify_if_time_from",
					"notify_if_time_to",
				},
			},
			"persons_to_notify_next_each_time": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				ConflictsWith: []string{
					"duration",
					"notify_on_call_from_schedule",
					"persons_to_notify",
					"action_to_trigger",
					"group_to_notify",
					"notify_if_time_from",
					"notify_if_time_to",
				},
			},
			"action_to_trigger": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ConflictsWith: []string{
					"duration",
					"notify_on_call_from_schedule",
					"persons_to_notify",
					"persons_to_notify_next_each_time",
					"group_to_notify",
					"notify_if_time_from",
					"notify_if_time_to",
				},
			},
			"group_to_notify": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ConflictsWith: []string{
					"duration",
					"notify_on_call_from_schedule",
					"persons_to_notify",
					"persons_to_notify_next_each_time",
					"action_to_trigger",
					"notify_if_time_from",
					"notify_if_time_to",
				},
			},
			"notify_if_time_from": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ConflictsWith: []string{
					"duration",
					"notify_on_call_from_schedule",
					"persons_to_notify",
					"persons_to_notify_next_each_time",
					"action_to_trigger",
				},
				RequiredWith: []string{
					"notify_if_time_to",
				},
			},
			"notify_if_time_to": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ConflictsWith: []string{
					"duration",
					"notify_on_call_from_schedule",
					"persons_to_notify",
					"persons_to_notify_next_each_time",
					"action_to_trigger",
				},
				RequiredWith: []string{
					"notify_if_time_from",
				},
			},
		},
	}
}

func resourceEscalationCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] create amixr escalation")

	client := m.(*amixr.Client)

	routeIdData := d.Get("route_id").(string)
	typeData := d.Get("type").(string)

	createOptions := &amixr.CreateEscalationOptions{
		RouteId:     routeIdData,
		Type:        typeData,
		ManualOrder: true,
	}

	durationData, durationOk := d.GetOk("duration")
	if durationOk {
		if typeData == "wait" {
			createOptions.Duration = durationData.(int)
		} else {
			return fmt.Errorf("duration can not be set with type: %s", typeData)
		}
	}

	personsToNotifyData, personsToNotifyDataOk := d.GetOk("persons_to_notify")
	if personsToNotifyDataOk {
		if typeData == "notify_persons" {
			createOptions.PersonsToNotify = stringSetToStringSlice(personsToNotifyData.(*schema.Set))
		} else {
			return fmt.Errorf("persons_to_notify can not be set with type: %s", typeData)
		}
	}

	notifyOnCallFromScheduleData, notifyOnCallFromScheduleDataOk := d.GetOk("notify_on_call_from_schedule")
	if notifyOnCallFromScheduleDataOk {
		if typeData == "notify_on_call_from_schedule" {
			createOptions.NotifyOnCallFromSchedule = notifyOnCallFromScheduleData.(string)
		} else {
			return fmt.Errorf("notify_on_call_from_schedule can not be set with type: %s", typeData)
		}
	}

	personsToNotifyNextEachTimeData, personsToNotifyNextEachTimeDataOk := d.GetOk("persons_to_notify_next_each_time")
	if personsToNotifyNextEachTimeDataOk {
		if typeData == "notify_person_next_each_time" {
			createOptions.PersonsToNotify = stringSetToStringSlice(personsToNotifyNextEachTimeData.(*schema.Set))
		} else {
			return fmt.Errorf("persons_to_notify_next_each_time can not be set with type: %s", typeData)
		}
	}

	notifyToGroupData, notifyToGroupDataOk := d.GetOk("group_to_notify")
	if notifyToGroupDataOk {
		if typeData == "notify_user_group" {
			createOptions.GroupToNotify = notifyToGroupData.(string)
		} else {
			return fmt.Errorf("notify_to_group can not be set with type: %s", typeData)
		}
	}

	actionToTriggerData, actionToTriggerDataOk := d.GetOk("action_to_trigger")
	if actionToTriggerDataOk {
		if typeData == "trigger_action" {
			createOptions.ActionToTrigger = actionToTriggerData.(string)
		} else {
			return fmt.Errorf("action to trigger can not be set with type: %s", typeData)
		}
	}

	notifyIfTimeFromData, notifyIfTimeFromDataOk := d.GetOk("notify_if_time_from")
	if notifyIfTimeFromDataOk {
		if typeData == "notify_if_time_from_to" {
			createOptions.NotifyIfTimeFrom = notifyIfTimeFromData.(string)
		} else {
			return fmt.Errorf("notify_if_time_from can not be set with type: %s", typeData)
		}
	}

	notifyIfTimeToData, notifyIfTimeToDataOk := d.GetOk("notify_if_time_to")
	if notifyIfTimeToDataOk {
		if typeData == "notify_if_time_from_to" {
			createOptions.NotifyIfTimeTo = notifyIfTimeToData.(string)
		} else {
			return fmt.Errorf("notify_if_time_to can not be set with type: %s", typeData)
		}
	}

	importanceData := d.Get("important").(bool)
	createOptions.Important = &importanceData

	positionData := d.Get("position").(int)
	createOptions.Position = &positionData

	escalation, _, err := client.Escalations.CreateEscalation(createOptions)
	if err != nil {
		return err
	}

	d.SetId(escalation.ID)

	return resourceEscalationRead(d, m)
}

func resourceEscalationRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] read amixr escalation")

	client := m.(*amixr.Client)

	escalation, _, err := client.Escalations.GetEscalation(d.Id(), &amixr.GetEscalationOptions{})
	if err != nil {
		return err
	}

	d.Set("route_id", escalation.RouteId)
	d.Set("position", escalation.Position)
	d.Set("type", escalation.Type)
	d.Set("duration", escalation.Duration)
	d.Set("notify_on_call_from_schedule", escalation.NotifyOnCallFromSchedule)
	d.Set("persons_to_notify", escalation.PersonsToNotify)
	d.Set("persons_to_notify_next_each_time", escalation.PersonsToNotifyEachTime)
	d.Set("group_to_notify", escalation.GroupToNotify)
	d.Set("action_to_trigger", escalation.ActionToTrigger)
	d.Set("important", escalation.Important)
	d.Set("notify_if_time_from", escalation.NotifyIfTimeFrom)
	d.Set("notify_if_time_to", escalation.NotifyIfTimeTo)

	return nil
}

func resourceEscalationUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] update amixr escalation")
	client := m.(*amixr.Client)

	typeData := d.Get("type").(string)

	updateOptions := &amixr.UpdateEscalationOptions{
		Type:        typeData,
		ManualOrder: true,
	}

	durationData, durationOk := d.GetOk("duration")
	if durationOk {
		if typeData == "wait" {
			updateOptions.Duration = durationData.(int)
		}
	}

	personsToNotifyData, personsToNotifyDataOk := d.GetOk("persons_to_notify")
	if personsToNotifyDataOk {
		if typeData == "notify_persons" {
			updateOptions.PersonsToNotify = stringSetToStringSlice(personsToNotifyData.(*schema.Set))
		}
	}

	notifyOnCallFromScheduleData, notifyOnCallFromScheduleDataOk := d.GetOk("notify_on_call_from_schedule")
	if notifyOnCallFromScheduleDataOk {
		if typeData == "notify_on_call_from_schedule" {
			updateOptions.NotifyOnCallFromSchedule = notifyOnCallFromScheduleData.(string)
		}
	}

	personsToNotifyNextEachTimeData, personsToNotifyNextEachTimeDataOk := d.GetOk("persons_to_notify_next_each_time")
	if personsToNotifyNextEachTimeDataOk {
		if typeData == "notify_person_next_each_time" {
			updateOptions.PersonsToNotify = stringSetToStringSlice(personsToNotifyNextEachTimeData.(*schema.Set))
		}
	}

	notifyToGroupData, notifyToGroupDataOk := d.GetOk("group_to_notify")
	if notifyToGroupDataOk {
		if typeData == "notify_user_group" {
			updateOptions.GroupToNotify = notifyToGroupData.(string)
		}
	}

	actionToTriggerData, actionToTriggerDataOk := d.GetOk("action_to_trigger")
	if actionToTriggerDataOk {
		if typeData == "trigger_action" {
			updateOptions.ActionToTrigger = actionToTriggerData.(string)
		}
	}

	notifyIfTimeFromData, notifyIfTimeFromDataOk := d.GetOk("notify_if_time_from")
	if notifyIfTimeFromDataOk {
		if typeData == "notify_if_time_from_to" {
			updateOptions.NotifyIfTimeFrom = notifyIfTimeFromData.(string)
		}
	}

	notifyIfTimeToData, notifyIfTimeToDataOk := d.GetOk("notify_if_time_to")
	if notifyIfTimeToDataOk {
		if typeData == "notify_if_time_from_to" {
			updateOptions.NotifyIfTimeTo = notifyIfTimeToData.(string)
		}
	}

	positionData := d.Get("position").(int)
	updateOptions.Position = &positionData

	importanceData := d.Get("important").(bool)
	updateOptions.Important = &importanceData

	escalation, _, err := client.Escalations.UpdateEscalation(d.Id(), updateOptions)
	if err != nil {
		return err
	}

	d.SetId(escalation.ID)
	return resourceEscalationRead(d, m)
}

func resourceEscalationDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] delete amixr escalation")

	client := m.(*amixr.Client)

	_, err := client.Escalations.DeleteEscalation(d.Id(), &amixr.DeleteEscalationOptions{})
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

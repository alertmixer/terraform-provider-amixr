package amixr

import (
	"fmt"
	"github.com/alertmixer/amixr-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
)

var onCallShiftTypeOptions = []string{
	"rolling_users",
	"recurrent_event",
	"single_event",
}

var onCallShiftFrequencyOptions = []string{
	"daily",
	"weekly",
	"monthly",
}

var onCallShiftWeekDayOptions = []string{
	"MO",
	"TU",
	"WE",
	"TH",
	"FR",
	"SA",
	"SU",
}

var sourceTerraform = 3

func resourceOnCallShift() *schema.Resource {
	return &schema.Resource{
		Create: resourceOnCallShiftCreate,
		Read:   resourceOnCallShiftRead,
		Update: resourceOnCallShiftUpdate,
		Delete: resourceOnCallShiftDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
            "team_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(onCallShiftTypeOptions, false),
			},
			"level": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"start": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"duration": &schema.Schema{
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"frequency": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(onCallShiftFrequencyOptions, false),
			},
			"users": &schema.Schema{
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"rolling_users": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeSet,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				Optional: true,
			},
			"interval": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"week_start": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(onCallShiftWeekDayOptions, false),
			},
			"by_day": &schema.Schema{
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(onCallShiftWeekDayOptions, false),
				},
				Optional: true,
			},
			"by_month": &schema.Schema{
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeInt,
					ValidateFunc: validation.IntBetween(1, 12),
				},
				Optional: true,
			},
			"by_monthday": &schema.Schema{
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeInt,
					ValidateFunc: validation.IntBetween(-31, 31),
				},
				Optional: true,
			},
			"time_zone": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"start_rotation_from_user_index": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},
		},
	}
}

func resourceOnCallShiftCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] create amixr on-call shift")

	client := m.(*amixr.Client)

	teamIdData := d.Get("team_id").(string)
	typeData := d.Get("type").(string)
	nameData := d.Get("name").(string)
	startData := d.Get("start").(string)
	durationData := d.Get("duration").(int)

	createOptions := &amixr.CreateOnCallShiftOptions{
	    TeamId:   teamIdData,
		Type:     typeData,
		Name:     nameData,
		Start:    startData,
		Duration: durationData,
		Source:   sourceTerraform,
	}

	levelData := d.Get("level").(int)
	createOptions.Level = &levelData

	frequencyData, frequencyOk := d.GetOk("frequency")
	if frequencyOk {
		if typeData != "single_event" {
			f := frequencyData.(string)
			createOptions.Frequency = &f
		} else {
			return fmt.Errorf("frequency can not be set with type: %s", typeData)
		}
	}

	usersData, usersDataOk := d.GetOk("users")
	if usersDataOk {
		if typeData != "rolling_users" {
			createOptions.Users = stringSetToStringSlice(usersData.(*schema.Set))
		} else {
			return fmt.Errorf("`users` can not be set with type: %s, use `rolling_users` field instead", typeData)
		}
	}

	intervalData, intervalOk := d.GetOk("interval")
	if intervalOk {
		if typeData != "single_event" {
			i := intervalData.(int)
			createOptions.Interval = &i
		} else {
			return fmt.Errorf("interval can not be set with type: %s", typeData)
		}
	}

	weekStartData, weekStartOk := d.GetOk("week_start")
	if weekStartOk {
		if typeData != "single_event" {
			w := weekStartData.(string)
			createOptions.WeekStart = &w
		} else {
			return fmt.Errorf("week_start can not be set with type: %s", typeData)
		}
	}

	byDayData, byDayOk := d.GetOk("by_day")
	if byDayOk {
		if typeData != "single_event" {
			createOptions.ByDay = stringSetToStringSlice(byDayData.(*schema.Set))
		} else {
			return fmt.Errorf("by_day can not be set with type: %s", typeData)
		}
	}

	byMonthData, byMonthOk := d.GetOk("by_month")
	if byMonthOk {
		if typeData != "single_event" {
			createOptions.ByMonth = intSetToIntSlice(byMonthData.(*schema.Set))
		} else {
			return fmt.Errorf("by_month can not be set with type: %s", typeData)
		}
	}

	byMonthdayData, byMonthdayOk := d.GetOk("by_monthday")
	if byMonthdayOk {
		if typeData != "single_event" {
			createOptions.ByMonthday = intSetToIntSlice(byMonthdayData.(*schema.Set))
		} else {
			return fmt.Errorf("by_monthday can not be set with type: %s", typeData)
		}
	}

	rollingUsersData, rollingUsersOk := d.GetOk("rolling_users")
	if rollingUsersOk {
		if typeData == "rolling_users" {
			createOptions.RollingUsers = listOfSetsToStringSlice(rollingUsersData.([]interface{}))
		} else {
			return fmt.Errorf("`rolling_users` can not be set with type: %s, use `users` field instead", typeData)
		}
	}

	timeZoneData, timeZoneOk := d.GetOk("time_zone")
	if timeZoneOk {
		tz := timeZoneData.(string)
		createOptions.TimeZone = &tz
	}

	if typeData == "rolling_users" {
		startRotationFromUserIndexData := d.Get("start_rotation_from_user_index")
		i := startRotationFromUserIndexData.(int)
		createOptions.StartRotationFromUserIndex = &i
	} // todo: add validation for start_rotation_from_user_index

	onCallShift, _, err := client.OnCallShifts.CreateOnCallShift(createOptions)
	if err != nil {
		return err
	}

	d.SetId(onCallShift.ID)

	return resourceOnCallShiftRead(d, m)
}

func resourceOnCallShiftUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] update amixr on-call shift")

	client := m.(*amixr.Client)

	typeData := d.Get("type").(string)
	nameData := d.Get("name").(string)
	startData := d.Get("start").(string)
	durationData := d.Get("duration").(int)

	updateOptions := &amixr.UpdateOnCallShiftOptions{
		Type:     typeData,
		Name:     nameData,
		Start:    startData,
		Duration: durationData,
		Source:   sourceTerraform,
	}

	levelData := d.Get("level").(int)
	updateOptions.Level = &levelData

	frequencyData, frequencyOk := d.GetOk("frequency")
	if frequencyOk {
		if typeData != "single_event" {
			f := frequencyData.(string)
			updateOptions.Frequency = &f
		} else {
			return fmt.Errorf("frequency can not be set with type: %s", typeData)
		}
	}

	usersData, usersDataOk := d.GetOk("users")
	if usersDataOk {
		if typeData != "rolling_users" {
			updateOptions.Users = stringSetToStringSlice(usersData.(*schema.Set))
		} else {
			return fmt.Errorf("`users` can not be set with type: %s, use `rolling_users` field instead", typeData)
		}
	}

	intervalData, intervalOk := d.GetOk("interval")
	if intervalOk {
		if typeData != "single_event" {
			i := intervalData.(int)
			updateOptions.Interval = &i
		} else {
			return fmt.Errorf("interval can not be set with type: %s", typeData)
		}
	}

	weekStartData, weekStartOk := d.GetOk("week_start")
	if weekStartOk {
		if typeData != "single_event" {
			w := weekStartData.(string)
			updateOptions.WeekStart = &w
		} else {
			return fmt.Errorf("week_start can not be set with type: %s", typeData)
		}
	}

	byDayData, byDayOk := d.GetOk("by_day")
	if byDayOk {
		if typeData != "single_event" {
			updateOptions.ByDay = stringSetToStringSlice(byDayData.(*schema.Set))
		} else {
			return fmt.Errorf("by_day can not be set with type: %s", typeData)
		}
	}

	byMonthData, byMonthOk := d.GetOk("by_month")
	if byMonthOk {
		if typeData != "single_event" {
			updateOptions.ByMonth = intSetToIntSlice(byMonthData.(*schema.Set))
		} else {
			return fmt.Errorf("by_month can not be set with type: %s", typeData)
		}
	}

	byMonthDayData, byMonthDayOk := d.GetOk("by_monthday")
	if byMonthDayOk {
		if typeData != "single_event" {
			updateOptions.ByMonthday = intSetToIntSlice(byMonthDayData.(*schema.Set))
		} else {
			return fmt.Errorf("by_monthday can not be set with type: %s", typeData)
		}
	}

	timeZoneData, timeZoneOk := d.GetOk("time_zone")
	if timeZoneOk {
		tz := timeZoneData.(string)
		updateOptions.TimeZone = &tz
	}

	rollingUsersData, rollingUsersOk := d.GetOk("rolling_users")
	if rollingUsersOk {
		if typeData == "rolling_users" {
			updateOptions.RollingUsers = listOfSetsToStringSlice(rollingUsersData.([]interface{}))
		} else {
			return fmt.Errorf("`rolling_users` can not be set with type: %s, use `users` field instead", typeData)
		}
	}

	if typeData == "rolling_users" {
		startRotationFromUserIndexData := d.Get("start_rotation_from_user_index")
		i := startRotationFromUserIndexData.(int)
		updateOptions.StartRotationFromUserIndex = &i
	} // todo: add validation for start_rotation_from_user_index

	onCallShift, _, err := client.OnCallShifts.UpdateOnCallShift(d.Id(), updateOptions)
	if err != nil {
		return err
	}

	d.SetId(onCallShift.ID)

	return resourceOnCallShiftRead(d, m)
}

func resourceOnCallShiftRead(d *schema.ResourceData, m interface{}) error {

	client := m.(*amixr.Client)
	options := &amixr.GetOnCallShiftOptions{}
	onCallShift, _, err := client.OnCallShifts.GetOnCallShift(d.Id(), options)

	if err != nil {
		return err
	}

	d.Set("team_id", onCallShift.TeamId)
	d.Set("name", onCallShift.Name)
	d.Set("type", onCallShift.Type)
	d.Set("level", onCallShift.Level)
	d.Set("start", onCallShift.Start)
	d.Set("duration", onCallShift.Duration)
	d.Set("frequency", onCallShift.Frequency)
	d.Set("week_start", onCallShift.WeekStart)
	d.Set("interval", onCallShift.Interval)
	d.Set("users", onCallShift.Users)
	d.Set("rolling_users", onCallShift.RollingUsers)
	d.Set("by_day", onCallShift.ByDay)
	d.Set("by_month", onCallShift.ByMonth)
	d.Set("by_monthday", onCallShift.ByMonthday)
	d.Set("time_zone", onCallShift.TimeZone)
	d.Set("start_rotation_from_user_index", onCallShift.StartRotationFromUserIndex)

	return nil
}

func resourceOnCallShiftDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] delete amixr on-call shift")

	client := m.(*amixr.Client)
	options := &amixr.DeleteOnCallShiftOptions{}
	_, err := client.OnCallShifts.DeleteOnCallShift(d.Id(), options)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

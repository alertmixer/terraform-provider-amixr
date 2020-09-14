---
layout: "amixr"
page_title: "Amixr: amixr_schedule"
sidebar_current: "docs-amixr-resource-schedule"
description: |-
  Creates and manages schedules in Amixr.
---

# amixr\_schedule

Manage a [schedule](https://api-docs.amixr.io/#schedules).

## Example Usage

```hcl
data "amixr_slack_channel" "example_slack_channel" {
  name = "example_slack_channel"
}

resource "amixr_schedule" "example_schedule" {
  name      = "Example Schadule"
  time_zone = "America/New_York"
  ical_url  = "https://example.com/example_ical.ics"
  slack {
    channel_id = data.amixr_slack_channel.example_slack_channel.slack_id
  }
}

```

## Argument Reference

The following arguments are supported:

  * `name` - (Required) The schedule's name
  * `time_zone` - (Optional) The schedule's time zone. More information about [time zones](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones).
  * `ical_url` - (Optional) The URL of the external calendar iCal file
  * `slack` - (Optional) Dictionary with slack-specific settings for a schedule. Includes:
    - `channel_id` - Slack channel id. Reminder about schedule shifts will be directed to this channel in Slack


## Attributes Reference

The following attributes are exported:

  * `id` - The ID of the schedule.
  

## Import

Existing schedule can be imported using the schedule ID:

```sh
$ terraform import amixr_schedule.example_schedule {schedule ID}
```


---
layout: "amixr"
page_title: "Amixr: amixr_escalation"
sidebar_current: "docs-amixr-resource-escalation"
description: |-
  Configures escalation policies in Amixr.
---

# amixr\_escalation

[Escalation policy](https://api-docs.amixr.io/#escalation-policies) configures what happened after incident is triggered: who will be notified first, second, etc., and delay before notifications. 

## Example Usage

```hcl
data "amixr_user" "alex" {
  username = "alex"
}

data "amixr_escalation_chain" "default" {
    name = "default"
}

resource "amixr_escalation" "example_notify_step" {
  escalation_chain_id = data.amixr_escalation_chain.default.id
  type = "notify_persons"
  persons_to_notify = [
    data.amixr_user.alex.id
  ]
  position = 0
}

resource "amixr_escalation" "example_wait_step" {
  escalation_chain_id = data.amixr_escalation_chain.default.id
  type = "wait"
  duration = 60
  position = 1
}

resource "amixr_escalation" "example_notify_step_important" {
  escalation_chain_id = data.amixr_escalation_chain.default.id
  type = "notify_persons"
  important = true
  persons_to_notify = [
    data.amixr_user.alex.id
  ]
  position = 2
}
```

## Argument Reference

The following arguments are supported:

  * `escalation_chain_id` - (Required) The ID of the escalation chain.
  * `position` - (Required) The position of the escalation step (starts from 0)
  * `type` - (Optional) The type of escalation policy. Can be:
    - `wait` - just wait
    - `notify_persons` - notify multiple users at the same time
    - `notify_person_next_each_time` - notify one user from queue
    - `notify_on_call_from_schedule` - notify by schedule
    - `notify_user_group` - notify User Group (available for teams with Slack integration)
    - `trigger_action` - trigger action (outgoing webhook)
    - `notify_if_time_from_to` - continue escalation only if the time is within the selected interval
    - `notify_whole_channel` - notify a channel in Slack (available for teams with Slack integration)
    - `resolve` - resolve incident
    - `null` - do nothing
  * `important` - (Optional) Will activate "important" personal notification rules. Can be `true` or `false`. Actual for steps: `notify_persons`, `notify_on_call_from_schedule` and `notify_user_group`.
  * `duration` - (Optional) The duration of delay for `wait` type step.
  * `persons_to_notify` - (Optional) The list of ID's of users for `notify_persons` type step.
  * `persons_to_notify_next_each_time` - (Optional) The list of ID's of users for `notify_person_next_each_time` type step.
  * `notify_on_call_from_schedule` - (Optional) ID of a Schedule for `notify_on_call_from_schedule` type step.
  * `action_to_trigger` - (Optional) ID of an Action for `trigger_action` type step.
  * `group_to_notify` - (Optional) ID of a User Group for `notify_user_group` type step.
  * `notify_if_time_from` - (Optional) The beginning of the time interval for `notify_if_time_from_to` type step in UTC (for example 08:00:00Z).
  * `notify_if_time_to` - (Optional) The end of the time interval for `notify_if_time_from_to` type step in UTC (for example 18:00:00Z).

## Attributes Reference

The following attributes are exported:

  * `id` - The ID of the escalation policy.


## Import

Existing escalations can be imported using the escalation ID:

```sh
$ terraform import amixr_escalation.example_escalation {escalation ID}
```

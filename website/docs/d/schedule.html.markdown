---
layout: "amixr"
page_title: "Amixr: amixr_schedule"
sidebar_current: "docs-amixr-datasource-schedule"
description: |-
  Get information about a schedule.
---

# amixr\_schedule

Get information about a specific [schedule](https://api-docs.amixr.io/#schedules).

## Example Usage

```hcl
data "amixr_schedule" "schedule" {
  name = "example_schedule"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The schedule's name.

## Attributes Reference

* `id` - The ID of the found schedule.
* `type` - The type of the found schedule.

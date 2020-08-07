---
layout: "amixr"
page_title: "Amixr: amixr_user_group"
sidebar_current: "docs-amixr-datasource-user_group"
description: |-
  Get information about a User Group.
---

# amixr\_user_group

Get information about a specific [User Group](https://api-docs.amixr.io/#user-groups).

## Example Usage

```hcl
data "amixr_user_group" "example_user_group" {
  slack_handle = "example_slack_handle"
}
```

## Argument Reference

The following arguments are supported:

* `slack_handle` - (Required) The User Group Slack handle.

## Attributes Reference

* `id` - The ID of the found User Group.

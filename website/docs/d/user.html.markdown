---
layout: "amixr"
page_title: "Amixr: amixr_user"
sidebar_current: "docs-amixr-datasource-user"
description: |-
  Get information about a user.
---

# amixr\_user

Get information about a specific [user](https://api-docs.amixr.io/#users).

## Example Usage

```hcl
data "amixr_user" "alex" {
  email = "alex@example.com"
}
```
## Argument Reference

The following arguments are supported:

* `email` - (Required) User's email.


## Attributes Reference

* `id` - The ID of the found user.
* `name` - The name of the found user.
* `role` - User's role in team.
* `team_id` - The ID of current team.
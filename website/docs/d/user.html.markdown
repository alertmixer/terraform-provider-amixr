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
  username = "alex"
}
```
## Argument Reference

The following arguments are supported:

* `username` - (Required) User's username.


## Attributes Reference

* `id` - The ID of the found user.
* `username` - The username of the found user.
* `role` - User's role in organization.
* `email` - The email of the found user.
---
layout: "amixr"
page_title: "Amixr: amixr_escalation_chain"
sidebar_current: "docs-amixr-resource-escalation-chain"
description: |-
  Configures escalation chains in Amixr.
---

# amixr\_escalation_chain

[Escalation chain](https://api-docs.amixr.io/#escalation-chains) is a reusable sequence of escalation policies.

## Example Usage

```hcl
resource "amixr_escalation_chain" "sre_east" {
  name = "sre-east"
}

// or use datasource to get the default escalation chain
data "amixr_escalation_chain" "default" {
  name = "default"
}

resource "amixr_escalation" "resolve" {
  escalation_chain_id = amixr_escalation_chain.sre_east.id
  type = "resolve"
  position = 0
}
```

## Argument Reference

The following arguments are supported:

  * `name` - (Required) Name of the escalation chain.

## Attributes Reference

The following attributes are exported:

  * `id` - The ID of the escalation chain.


## Import

Existing escalation chains can be imported using the escalation chain ID:

```sh
$ terraform import amixr_escalation_chain.example_escalation_chain {escalation chain ID}
```

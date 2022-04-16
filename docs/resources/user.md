---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "drone_user Resource - terraform-provider-drone"
subcategory: ""
description: |-
  
---

# drone_user (Resource)

Manage a user.

~> In order to use the `drone_user` resource you must have admin privileges within your Drone environment.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `active` (Boolean)
- `admin` (Boolean)
- `login` (String)

### Optional

- `id` (String) The ID of this resource.
- `machine` (Boolean)

### Read-Only

- `token` (String)


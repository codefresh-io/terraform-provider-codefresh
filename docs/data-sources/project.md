---
page_title: "codefresh_project Data Source - terraform-provider-codefresh"
subcategory: ""
description: |-
  This data source retrieves a project by its ID or name.
---

# codefresh_project (Data Source)

This data source retrieves a project by its ID or name.

## Example Usage

```hcl
data "codefresh_project" "myapp" {
  name = "myapp"
}


resource "codefresh_pipeline" "myapp-deploy" {

  name    = "${data.codefresh_project.myapp.projectName}/myapp-deploy"

  ...
}

```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `_id` (String)
- `name` (String)
- `tags` (List of String)

### Read-Only

- `id` (String) The ID of this resource.
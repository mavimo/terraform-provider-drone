# drone_repo Resource

Activate and configure a repository.

## Example Usage

```hcl
resource "drone_repo" "hello_world" {
  repository = "octocat/hello-world"
  visibility = "public"
}
```

## Argument Reference

* `repository` - (Required) Repository name (e.g. `octocat/hello-world`).
* `trusted` - (Optional) Repository is trusted (default: `false`).
* `protected` - (Optional) Repository is protected (default: `false`).
* `timeout` - (Optional) Repository timeout (default: `60`).
* `visibility` - (Optional) Repository visibility (default: `private`).
* `configuration` - (Optional) Drone Configuration file (default `.drone.yml`).

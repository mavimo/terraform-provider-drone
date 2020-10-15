# drone_orgsecret Resource

Manage an org secret.

## Example Usage

```hcl
resource "drone_orgsecret" "octocat" {
  namespace  = "myorg"
  name       = "master_password"
  value      = "correct horse battery staple"
}
```

## Argument Reference

* `namespace` - (Required) Organization name (e.g. `octocat`).
* `name` - (Required) Secret name.
* `value` - (Required) Secret value.
* `allow_on_pull_request` - (Optional) Allow retrieving the secret on pull requests.  (Default `false`)
* `allow_push_on_pull_request` - (Optional) Allow pushing on pull requests.  (Default `false`)

~> In order to use the `drone_orgsecret` resource you must have admin privileges within your Drone environment.

# drone_secret Resource

Manage a repository secret.

## Example Usage

```hcl
resource "drone_secret" "master_password" {
  repository = "octocat/hello-world"
  name       = "master_password"
  value      = "correct horse battery staple"
}
```

## Argument Reference

* `repository` - (Required) Repository name (e.g. `octocat/hello-world`).
* `name` - (Required) Secret name.
* `value` - (Required) Secret value.
* `allow_on_pull_request` - (Optional) Allow retrieving the secret on pull requests.  (Default `false`)

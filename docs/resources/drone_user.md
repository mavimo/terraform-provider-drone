# drone_user Resource

Manage a repository cron job.

## Example Usage

```hcl
resource "drone_user" "octocat" {
  login = "octocat"
}
```

## Argument Reference

* `login` - (Required) Login name.
* `admin` - (Required) Is the user an admin?
* `active` - (Required) Is the user active?
* `machine` - (Optional) Is the user a machine?  (Default `false`)

## Attributes Reference

* `token` - The user's access token

~> In order to use the `drone_user` resource you must have admin privileges within your Drone environment.

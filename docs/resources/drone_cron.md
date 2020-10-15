# drone_cron Resource

Manage a repository cron job.

## Example Usage

```hcl
resource "drone_cron" "cron_job" {
  repository = "octocat/hello-world"
  name       = "the_cron_job"
  expr       = "@monthly"
  event      = "push"
}
```

## Argument Reference

* `repository` - (Required) Repository name (e.g. `octocat/hello-world`).
* `name` - (Required) Cron job name.
* `event` - (Required) The event for this cron job.  Only allowed value is `push`.
* `expr` - (Required) The cron interval.  Allowed values are `@daily`, `@weekly`, `@monthly`, and `@yearly`.
* `branch` - (Optional) The branch to run this cron job on.  (Default `master`)
* `disabled` - (Optional) Is this cron job disabled?

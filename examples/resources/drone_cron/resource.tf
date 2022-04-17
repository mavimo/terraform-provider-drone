resource "drone_cron" "cron_job" {
  repository = "octocat/hello-world"
  name       = "the_cron_job"
  expr       = "@monthly"
  event      = "push"
}
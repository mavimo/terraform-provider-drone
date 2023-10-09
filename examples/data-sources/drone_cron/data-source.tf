#retrieve a specific cron from drone instance
data "drone_cron" "cron_job" {
  repository = "octocat/hello-world"
  name       = "the_cron_job"
}

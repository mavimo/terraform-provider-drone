#retrieve all crons in a repositorie from drone instance
data "drone_crons" "all" {
  repository = "octolab/hello-world"
}

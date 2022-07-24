#retrieve a specific repository from drone instance
data "drone_repo" "hello_world" {
  repository = "octolab/hello-world"
}

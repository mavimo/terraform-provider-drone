resource "drone_repo" "hello_world" {
  repository = "octocat/hello-world"
  visibility = "public"
}
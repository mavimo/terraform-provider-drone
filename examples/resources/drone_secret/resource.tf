resource "drone_secret" "master_password" {
  repository = "octocat/hello-world"
  name       = "master_password"
  value      = "correct horse battery staple"
}
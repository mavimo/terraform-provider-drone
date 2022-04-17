resource "drone_orgsecret" "octocat" {
  namespace = "myorg"
  name      = "master_password"
  value     = "correct horse battery staple"
}
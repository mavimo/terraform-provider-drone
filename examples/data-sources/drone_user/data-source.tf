#retrieve a specific user from drone instance
data "drone_user" "octocat" {
  login = "octocat"
}

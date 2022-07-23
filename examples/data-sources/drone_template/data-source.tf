#retrieve a specific te,èòate from drone instance
data "drone_template" "octocat" {
  name      = "common.yaml"
  namespace = "octolab"
}

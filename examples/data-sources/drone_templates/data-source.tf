#retrieve all templates from drone instance
data "drone_templates" "all" {
  namespace = "octolab"
}

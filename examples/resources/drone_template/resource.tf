resource "drone_template" "my_pipeline" {
  namespace = "myorg"
  name      = "pipeline1.yaml"
  data      = file("${path.module}/templates/pipeline1.yaml")
}

resource "drone_template" "other_pipeline" {
  namespace = "myorg"
  name      = "pipeline2.json"
  data      = file("${path.module}/templates/pipeline2.json")
}

resource "drone_template" "other_pipeline" {
  namespace = "myorg"
  name      = "pipeline3.jsonnet"
  data      = file("${path.module}/templates/pipeline3.jsonnet")
}

resource "drone_template" "other_pipeline" {
  namespace = "myorg"
  name      = "pipeline4.starlark"
  data      = file("${path.module}/templates/pipeline4.starlark")
}
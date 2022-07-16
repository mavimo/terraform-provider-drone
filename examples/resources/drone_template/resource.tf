resource "drone_template" "my_pipeline" {
  namespace = "myorg"
  name      = "pipeline1.yaml"
  data      = "lorem: ipsum"
}

resource "drone_template" "other_pipeline" {
  namespace = "myorg"
  name      = "pipeline2.jsonnet"
  data      = file("${path.module}/templates/pipeline2.jsonnet")
}

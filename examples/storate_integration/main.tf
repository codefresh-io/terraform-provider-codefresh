resource "codefresh_context" "storage_integration" {
  for_each = toset(["create"])
  name = "pasha-test-t2"
  spec {
    storagegc {
      data {
        auth {
          type = "basic"
          json_config = tomap({
            "fd": "fd"
          })
        }
      }
    }
  }
}

## example
```terraform
resource "headless_chrome_request" "hoge" {
  url     = ""
  headers = {

  }
}
output "response" {
  value = data.headless_chrome_request.hoge.body
}
```


## note

This repository is based on [this template](https://github.com/hashicorp/terraform-provider-template)

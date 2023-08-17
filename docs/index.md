---
page_title: "Provider: Chrome Headless"
description: "Chrome Headless Provider is a provider for working with Chrome on Terraform"
---

# Chrome Headless Provider

Summary of what the provider is for, including use cases and links to
app/service documentation.

## Example

```hcl
# Basic Usage
data "headless_chrome_request" "example" {
  url    = "https://yasudacloud.github.io"
  width  = 1024
  height = 768

  screenshot = {
    dist_path = "/var/app/dist"     # local directory
    file_name = "example.png"       # PNG File Name
  }
}

# You can retrieve the response for a given url
output "example-output" {
  value = headless_chrome_request.example.body
}
```


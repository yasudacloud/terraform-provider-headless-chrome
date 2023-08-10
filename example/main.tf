terraform {
  required_providers {
    headless = {
      source  = "hashicorp.com/edu/headless-chrome"
      version = "1.0.0"
    }
  }
}

data "headless_chrome_request" "example" {
  url        = "https://yasudacloud.github.io"
  screenshot = {
    dist_path = var.dist_path
    file_name = var.file_name
  }
}

output "response" {
  value = {
    "file_name" : data.headless_chrome_request.example.screenshot.file_name
    "status_code" : data.headless_chrome_request.example.status_code
    "response_headers" : data.headless_chrome_request.example.response_headers
  }
}

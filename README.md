# Chrome Headless Provider

This repository is a provider for headless use of Chrome from Terraform.

## Example

```terraform
data "headless_chrome_request" "example" {
  url    = "https://yasudacloud.github.io"
  width  = 1024
  height = 768

  screenshot = {
    dist_path = "/var/app/dist"     # local directory
    file_name = "example.png"       # PNG File Name
  }
}

output "response" {
  value = data.headless_chrome_request.example.body
}
```

## Reference

It is listed [here](https://github.com/yasudacloud/terraform-provider-headless-chrome/blob/main/docs/data-sources/request.md)

## Feedback

Please write [Issue](https://github.com/yasudacloud/terraform-provider-headless-chrome/issues)

## Thanks

- [terraform-provider-template](https://github.com/hashicorp/terraform-provider-template)

- [chromedp](https://github.com/chromedp/chromedp)

Thanks for the great OSS!

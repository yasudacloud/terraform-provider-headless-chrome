# headless_chrome_request Data Source

headless_chrome_request makes Http requests through Chrome.

You can set up a screenshot to save an image locally, but it is optional.

Note that the request is still made at the point of plan.

## Example Usage

```hcl
data "headless_chrome_request" "example" {
  url = "https://yasudacloud.github.io"
}
output "example-output" {
  value = headless_chrome_request.example.body
}
```

## Schema

### Required

- `url` (String) Web URL

### Optional

- `screenshot` (Object) Setting up to take a screenshot
- `useragent` (String) Browser User Agent
- `width` Screen Width
- `height` Screen Height

### Read-Only

- `body` (String) Http Response Body
- `status_code` (String) Http Response Status Code
- `response_headers` (String) Http Response Headers

### Nested Schema for `screenshot`

### Read-Only

- `dist_path` (String) absolute path of the destination of the screenshot
- `file_name` (String) File name of the image to be saved. The extension should be png

Only one of dist_path and file_name cannot be left blank.
package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"regexp"
	"testing"
)

func TestAccDataRequest(t *testing.T) {
	name := "data.headless_chrome_request.example"
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testCase1,
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testCase2,
				ExpectError: regexp.MustCompile("Incorrect attribute value type"),
				//Check: resource.ComposeTestCheckFunc(
				//	resource.TestCheckResourceAttr(name, "url", "https://yasudacloud.github.io"),
				//),
			},
			{
				Config:      testCase3,
				ExpectError: regexp.MustCompile("key 'dist_path' is required"),
			},
			{
				Config:      testCase4,
				ExpectError: regexp.MustCompile("key 'file_name' is required"),
			},
			{
				Config:      testCase5,
				ExpectError: regexp.MustCompile("key 'dist_path' is a directory that does not exist"),
			},
			{
				Config: testCase6,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "url", "https://yasudacloud.github.io"),
					resource.TestCheckResourceAttr(name, "status_code", "200"),
				),
			},
		},
	})
}

// reference: docs/test/data_headless_chrome_request.md
var (
	testCase1 = `
data "headless_chrome_request" "example" {
 provider = "headless-chrome"
}`
	testCase2 = `
data "headless_chrome_request" "example" {
 provider = "headless-chrome"
 url      = "https://yasudacloud.github.io"
 screenshot = ""
}
`
	testCase3 = `
data "headless_chrome_request" "example" {
 provider = "headless-chrome"
 url      = "https://yasudacloud.github.io"
 screenshot = {
   file_name = "test.png"
 }
}
`
	testCase4 = `
data "headless_chrome_request" "example" {
 provider = "headless-chrome"
 url      = "https://yasudacloud.github.io"
 screenshot = {
   dist_path = "/var/www"
 }
}
`
	testCase5 = `
data "headless_chrome_request" "example" {
 provider = "headless-chrome"
 url      = "https://yasudacloud.github.io"
 screenshot = {
   dist_path = "/var/testing"
   file_name = "test.png"
 }
}
`
	testCase6 = `
data "headless_chrome_request" "example" {
 provider = "headless-chrome"
 url      = "https://yasudacloud.github.io"
}
`
)

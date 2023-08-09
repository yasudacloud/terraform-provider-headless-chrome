package provider

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/hashicorp/terraform/helper/schema"
	"io/ioutil"
	"log"
	"os"
)

func dataRequest() *schema.Resource {
	return &schema.Resource{
		Read: dataRequestRead,
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Description: "HTTP Request URL",
				Required:    true,
			},
			"screenshot": {
				Type:         schema.TypeMap,
				Optional:     true,
				Default:      make(map[string]interface{}),
				Description:  "ScreenShot Setting",
				ValidateFunc: validateScreenShotAttribute,
			},
			"body": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "HTTP Response Body",
				Computed:    true,
			},
			"status_code": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "HTTP Response Status Code",
				Computed:    true,
			},
			"response_headers": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "HTTP Response Header",
				Computed:    true,
			},
		},
	}
}

func dataRequestRead(d *schema.ResourceData, _ interface{}) error {
	ctx, cancel := chromedp.NewContext(context.Background(), chromedp.WithDebugf(log.Printf))
	defer cancel()
	var responseReceived *network.EventResponseReceived
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch ev := ev.(type) {
		case *network.EventResponseReceived:
			responseReceived = ev
		}
	})
	var htmlContent string
	var imageBuf []byte
	var screenshotAttr, isScreenShot = d.GetOk("screenshot")
	url := d.Get("url").(string)

	actions := []chromedp.Action{
		network.Enable(),
		chromedp.Navigate(url),
		chromedp.OuterHTML(`html`, &htmlContent),
	}
	if isScreenShot {
		actions = append(actions, chromedp.Screenshot(`html`, &imageBuf, chromedp.NodeVisible, chromedp.ByQuery))
	}
	err := chromedp.Run(ctx, actions...)
	if err != nil {
		log.Fatal(err)
	} else {
		d.SetId(url)
		d.Set("body", htmlContent)
		if responseReceived != nil {
			d.Set("status_code", responseReceived.Response.Status)
			d.Set("response_headers", responseReceived.Response.Headers)
		}
	}
	if isScreenShot {
		m := screenshotAttr.(map[string]interface{})
		var distPath = m["dist_path"].(string)
		var fileName = m["file_name"].(string)
		var path = fmt.Sprintf("%s/%s", distPath, fileName)
		if err := ioutil.WriteFile(path, imageBuf, 0644); err != nil {
			log.Fatal(err)
		}
	}
	return nil
}
func validateScreenShotAttribute(v interface{}, key string) (warnings []string, errors []error) {
	m, ok := v.(map[string]interface{})
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of <attr_name> to be map[string]interface{}"))
		return warnings, errors
	}
	if _, ok := m["file_name"]; !ok {
		errors = append(errors, fmt.Errorf("key 'file_name' is required"))
	}

	if _, ok := m["dist_path"]; !ok {
		errors = append(errors, fmt.Errorf("key 'dist_path' is required"))
	} else {
		var distPath = m["dist_path"].(string)
		info, err := os.Stat(distPath)
		if err != nil {
			if os.IsNotExist(err) {
				errors = append(errors, fmt.Errorf("key 'dist_path' is a directory that does not exist"))
			} else {
				errors = append(errors, fmt.Errorf(fmt.Sprintf("Error checking path %s", distPath)))
			}
		} else {
			if !info.IsDir() {
				errors = append(errors, fmt.Errorf("key 'dist_path' is not a directory"))
			}
		}
	}
	return warnings, errors
}

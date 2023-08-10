package provider

import (
	"bytes"
	"context"
	"fmt"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"os"
)

func dataRequest() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataRequestRead,
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
			"useragent": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "HTTP Request Header User Agent",
				ValidateFunc: func(i interface{}, s string) (warnings []string, errors []error) {
					if value := i.(string); len(value) > 512 {
						errors = append(errors, fmt.Errorf("upper limit is 512 characters"))
					}
					return warnings, errors
				},
			},
			"width": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Window Width",
			},
			"height": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Window Height",
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

func dataRequestRead(c context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	var iWidth int
	var iHeight int

	ctx, cancel := chromedp.NewContext(context.Background(), chromedp.WithDebugf(log.Printf))
	defer cancel()
	var httpResponse *network.EventResponseReceived
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch ev := ev.(type) {
		case *network.EventResponseReceived:
			httpResponse = ev
		}
	})
	var htmlContent string
	var imageBuf []byte
	var screenshotAttr, isScreenShot = d.GetOk("screenshot")
	url := d.Get("url").(string)

	actions := []chromedp.Action{
		network.Enable(),
	}

	if useragent := d.Get("useragent").(string); useragent != "" {
		actions = append(
			actions,
			emulation.SetUserAgentOverride(useragent),
		)
	}

	width, wOk := d.GetOk("width")
	height, hOk := d.GetOk("height")
	if wOk && hOk {
		iWidth = width.(int)
		iHeight = height.(int)
	} else if (wOk && !hOk) || (!wOk && hOk) {
		return diag.Errorf("'width' and 'height' must both be set or both be unset")
	} else {
		iWidth = 1280
		iHeight = 768
	}
	actions = append(actions,
		chromedp.EmulateViewport(int64(iWidth), int64(iHeight)),
		chromedp.Navigate(url),
		chromedp.OuterHTML(`html`, &htmlContent),
	)

	if isScreenShot {
		actions = append(actions, chromedp.Screenshot("html", &imageBuf, chromedp.NodeVisible, chromedp.ByQuery))
	}
	err := chromedp.Run(ctx, actions...)
	if err != nil {
		log.Fatal(err)
	} else {
		d.SetId(url)
		d.Set("body", htmlContent)
		if httpResponse != nil {
			d.Set("status_code", httpResponse.Response.Status)
			d.Set("response_headers", httpResponse.Response.Headers)
		}
	}
	if isScreenShot {
		img, _, err := image.Decode(bytes.NewReader(imageBuf))
		if err != nil {
			return diag.Errorf(err.Error())
		}
		rect := image.Rect(0, 0, iWidth, iHeight)
		trimmed := img.(interface {
			SubImage(r image.Rectangle) image.Image
		}).SubImage(rect)

		var buf bytes.Buffer
		if err := png.Encode(&buf, trimmed); err != nil {
			print("@error image")
			return diag.Errorf("resize image error")
		}
		imageBuf = buf.Bytes()

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

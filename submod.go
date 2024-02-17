package xpostblue

import "github.com/playwright-community/playwright-go"

func context(device *playwright.DeviceDescriptor, browser playwright.Browser) (playwright.BrowserContext, error) {
	context, err := browser.NewContext(playwright.BrowserNewContextOptions{
		// sample location in Tokyo
		Geolocation: &playwright.Geolocation{
			Longitude: 139.749281,
			Latitude:  35.6959983,
		},
		Viewport:          device.Viewport,
		JavaScriptEnabled: playwright.Bool(true),
		UserAgent:         playwright.String(device.UserAgent),
		DeviceScaleFactor: playwright.Float(device.DeviceScaleFactor),

		HasTouch: playwright.Bool(device.HasTouch),

		// button value="Next" to "次へ"
		Locale: playwright.String("ja-JP"),
		// Timezone: "Asia/Tokyo"
		TimezoneId: playwright.String("Asia/Tokyo"),
	})
	if err != nil {
		return nil, err
	}

	return context, err
}

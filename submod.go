package xpostblue

import "github.com/playwright-community/playwright-go"

func context(device *playwright.DeviceDescriptor, browser playwright.Browser) (playwright.BrowserContext, error) {
	context, err := browser.NewContext(playwright.BrowserNewContextOptions{
		Geolocation: &playwright.Geolocation{
			Longitude: 139.749281,
			Latitude:  35.6959983,
		},
		Permissions:       []string{"geolocation"},
		Viewport:          device.Viewport,
		JavaScriptEnabled: playwright.Bool(true),
		UserAgent:         playwright.String(device.UserAgent),
		DeviceScaleFactor: playwright.Float(device.DeviceScaleFactor),
		IsMobile:          playwright.Bool(device.IsMobile),
		HasTouch:          playwright.Bool(device.HasTouch),
	})
	if err != nil {
		return nil, err
	}

	return context, err
}

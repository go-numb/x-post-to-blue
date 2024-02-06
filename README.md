# X-POST-TO-BLUE

This Golang package provides long-form posting on Twitter Blue using the Playwright library.

### Installation
To use this package, you need to have Go installed. You can install it using the following command:
```bash
go get github.com/username/xpostblue

# playwright
go run github.com/playwright-community/playwright-go/cmd/playwright@latest install --with-deps
# Or
go install github.com/playwright-community/playwright-go/cmd/playwright@latest
playwright install --with-deps
```

### Usage
```go
package main

import (
	"github.com/username/x-post-to-blue"
)

func main() {
	// Create a new client
	client := xpostblue.New(true)
	// Close the client
	defer client.Close()

	// Login to Twitter
	err := client.Login("yourusername", "yourpassword")
	if err != nil {
		panic(err)
	}

	// Post a message on Twitter
    files := []string{"./images/1.jpg", "./images/2.jpg"}
	err = client.Post(true, 5, "Hello, world! and long long-text", nil)
	if err != nil {
		panic(err)
	}
}
```

See [module test file](./mod_test.go) for details.

### Structs
- `ClientBody`: Main client struct containing Playwright instances and page methods
- `PostLocator`: Locator struct containing the elements for login and post sections

### Methods
- `New(isHeadless bool) *ClientBody`: Initialize a new client
- `Close()`: Close the client and browser
- `Login(username, password string) error`: Login to Twitter
- `Post(isPost bool, sleepSecForUpload int, msg string, files []string) error`: Post a message on Twitter

### Dependencies
- Playwright Go
- Zerolog

### License
This package is released under the MIT License.

For more information on Playwright Go, visit [https://github.com/mxschmitt/playwright-go](https://github.com/mxschmitt/playwright-go).

Feel free to contribute to this package by submitting issues or pull requests on GitHub.
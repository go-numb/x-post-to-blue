package xpostblue

import (
	"fmt"
	"net/url"
	"time"

	"github.com/playwright-community/playwright-go"
	"github.com/rs/zerolog/log"
)

const (
	// /i/flow/login
	TWITTER    = "https://twitter.com"
	TWITTERPRO = "https://pro.twitter.com"
	PATHLOGIN  = "/i/flow/login"

	// Login section
	INPUTID   = "input[type='text']"
	BTNID     = "xpath=//span[text()='次へ']"
	INPUTPASS = "input[type='password']"
	BTNPASS   = "[data-testid='LoginForm_Login_Button']"

	// Post section
	CONFIRMAREA = "[data-testid='tweetTextarea_0']"
	TOPOST      = "div[aria-label='ポストを作成']"
	INPUTMSG    = "[data-testid='tweetTextarea_0']"
	SELECTFILE  = "input[data-testid='fileInput']"
	BTNPOST     = "xpath=//span[text()='ポストする']"
)

type ClientBody struct {
	Pw      *playwright.Playwright
	Browser playwright.Browser
	Context playwright.BrowserContext
	Page    playwright.Page

	URL *url.URL

	PostLocator *PostLocator
}

type PostLocator struct {
	LoginURL string
	ProURL   string

	InputID   string
	BtnID     string
	InputPass string
	BtnPass   string

	ConfirmArea string
	ToPost      string
	InputMsg    string
	SelectFile  string
	BtnPost     string
}

func New(isHeadless bool) *ClientBody {
	pw, err := playwright.Run()
	if err != nil {
		log.Fatal().Msgf("could not run playwright: %v", err)
		return nil
	}

	// is_post = If false, display and operate GUI browser
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(isHeadless),
	})
	if err != nil {
		log.Fatal().Msgf("could not launch browser: %v", err)
		return nil
	}

	// Mobile device settings
	// Use older Pixel 5 model to avoid known bugs
	// Specify latitude and longitude of Tokyo, Japan
	device := pw.Devices["Pixel 7"]
	context, err := context(device, browser)
	if err != nil {
		log.Fatal().Msgf("could not create context: %v", err)
		return nil
	}

	page, err := context.NewPage()
	if err != nil {
		log.Fatal().Msgf("could not create page: %v", err)
		return nil
	}

	u, err := url.Parse(TWITTER)
	if err != nil {
		log.Fatal().Msgf("could not parse url: %v", err)
		return nil
	}

	return &ClientBody{
		Pw:      pw,
		Browser: browser,
		Context: context,
		Page:    page,

		URL: u,
		PostLocator: &PostLocator{
			LoginURL: TWITTER + PATHLOGIN,
			ProURL:   TWITTERPRO,

			InputID:   INPUTID,
			BtnID:     BTNID,
			InputPass: INPUTPASS,
			BtnPass:   BTNPASS,

			ConfirmArea: CONFIRMAREA,
			ToPost:      TOPOST,
			InputMsg:    INPUTMSG,
			SelectFile:  SELECTFILE,
			BtnPost:     BTNPOST,
		},
	}
}

func (p *ClientBody) Close() {
	p.Page.Close()
	p.Browser.Close()
	p.Pw.Stop()
}

func (p *ClientBody) Login(username, password string) error {
	// to twitter.com/i/flow/login
	u, _ := url.Parse(p.PostLocator.LoginURL)
	if _, err := p.Page.Goto(u.String()); err != nil {
		return fmt.Errorf("%v > could not goto", err)
	}
	log.Debug().Msgf("target url: %s", u.String())

	// input Username/Email
	if err := p.Page.Locator(p.PostLocator.InputID).Fill(username); err != nil {
		return fmt.Errorf("%v > could not fill to account input", err)
	}

	if err := p.Page.Locator(p.PostLocator.BtnID).Tap(); err != nil {
		return fmt.Errorf("%v > could not click to next button", err)
	}

	// input Password
	if err := p.Page.Locator(p.PostLocator.InputPass).Fill(password); err != nil {
		return fmt.Errorf("%v > could not fill to password input", err)
	}

	if err := p.Page.Locator(p.PostLocator.BtnPass).Nth(0).Tap(); err != nil {
		return fmt.Errorf("%v > could not click to login button", err)
	}

	return nil
}

func (p *ClientBody) Post(isPost bool, sleepSecForUpload int, msg string, files []string) error {
	// to pro.twitter.com
	u, _ := url.Parse(p.PostLocator.ProURL)
	if _, err := p.Page.Goto(u.String()); err != nil {
		return fmt.Errorf("%v > could not goto", err)
	}
	log.Debug().Msgf("target url: %s", u.String())

	// Enter text into elements with contententeditable attribute
	isVisible, err := p.Page.Locator(p.PostLocator.ConfirmArea).IsVisible()
	if err != nil {
		return fmt.Errorf("%v > could not check the element is visible", err)
	}
	if !isVisible {
		// If there is no input screen, tap the Display Input Screen button. If there are notifications, etc., the notification screen will be prioritized for display.
		if err := p.Page.Locator(p.PostLocator.ToPost).Tap(); err != nil {
			log.Debug().Msgf("%v", fmt.Errorf("%v > ok or could not tap to %s element", err, p.PostLocator.ToPost))
		}
	}
	if err := p.Page.Locator(p.PostLocator.InputMsg).Fill(msg); err != nil {
		return fmt.Errorf("%v > could not fill to tweet input", err)
	}

	// upload files
	if files != nil {
		if err := p.Page.Locator(p.PostLocator.SelectFile).SetInputFiles(files); err != nil {
			return fmt.Errorf("%v > could not upload files", err)
		}
		time.Sleep(time.Second * time.Duration(sleepSecForUpload))
	}

	// click to post button
	if !isPost {
		return fmt.Errorf("is_post is false")
	}
	if err := p.Page.Locator(p.PostLocator.BtnPost).Tap(); err != nil {
		return fmt.Errorf("%v > could not click to post button", err)
	}

	return nil
}

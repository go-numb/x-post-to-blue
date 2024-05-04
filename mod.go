package xpostblue

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/playwright-community/playwright-go"
	"github.com/rs/zerolog/log"
)

const (
	// /i/flow/login
	TWITTER    = "https://twitter.com"
	TWITTERPRO = "https://pro.twitter.com"
	PATHHOME   = "/home"
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

	MaxWaitSecForRequest int
	MaxWaitSecForInput   int
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
	// Why Firefox required?
	// Because it is the only browser that can upload files in headless mode for video files.
	browser, err := pw.Firefox.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(isHeadless),
	})
	if err != nil {
		log.Fatal().Msgf("could not launch browser: %v", err)
		return nil
	}

	// Mobile device settings, etc.
	// Use older Pixel 5 model to avoid known bugs
	// Specify latitude and longitude of Tokyo, Japan
	device := pw.Devices["iPad Pro 11 landscape"]
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

	page.SetDefaultTimeout(*playwright.Float(120 * 1000))

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

		MaxWaitSecForRequest: 120,
		MaxWaitSecForInput:   5,
	}
}

func (p *ClientBody) Close() {
	p.Page.Close()
	p.Browser.Close()
	p.Pw.Stop()
}

func (p *ClientBody) SetTimeout(sec int) *ClientBody {
	p.Page.SetDefaultTimeout(*playwright.Float(float64(sec * 1000)))
	p.MaxWaitSecForRequest = sec
	return p
}

func (p *ClientBody) Login(username, password string) error {
	maxWaitSec := p.MaxWaitSecForInput * 1000

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

	wait(maxWaitSec)

	if err := p.Page.Locator(p.PostLocator.BtnID).Tap(); err != nil {
		return fmt.Errorf("%v > could not click to next button", err)
	}

	wait(maxWaitSec)

	// input Password
	if err := p.Page.Locator(p.PostLocator.InputPass).Fill(password); err != nil {
		return fmt.Errorf("%v > could not fill to password input", err)
	}

	wait(maxWaitSec)

	if err := p.Page.Locator(p.PostLocator.BtnPass).Nth(0).Tap(); err != nil {
		return fmt.Errorf("%v > could not click to login button", err)
	}

	wait(maxWaitSec)

	return nil
}

func (p *ClientBody) Post(isPost bool, sleepSecForUpload int, msg string, files ...string) error {
	maxWaitSec := p.MaxWaitSecForInput * 1000

	// to pro.twitter.com
	u, _ := url.Parse(p.PostLocator.ProURL)
	if _, err := p.Page.Goto(u.String()); err != nil {
		return fmt.Errorf("%v > could not goto", err)
	}
	log.Debug().Msgf("target url: %s", u.String())

	wait(maxWaitSec)

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

		wait(maxWaitSec)
	}
	if err := p.Page.Locator(p.PostLocator.InputMsg).Fill(msg); err != nil {
		return fmt.Errorf("%v > could not fill to tweet input", err)
	}

	wait(maxWaitSec)

	// upload files
	if err := p.uploadFiles(isPost, files...); err != nil {
		return fmt.Errorf("%v > could not upload files", err)
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

// uploadFiles ファイルをアップロードする
func (p *ClientBody) uploadFiles(with_files bool, files ...string) error {
	if len(files) == 0 {
		log.Debug().Msgf("no files to upload, files: %v", files)
		return nil
	}

	// GUIが求める型式に変更する
	inputFiles, err := filesToInputFiles(files)
	if err != nil {
		return SetError(err, "could not convert files to input files")
	}

	// ファイルをアップロード
	if err := p.Page.Locator("input[data-testid='fileInput']").SetInputFiles(inputFiles); err != nil {
		if with_files { // ファイルの選択ができない場合、エラーを返す
			return SetError(err, "could not upload file")
		} else { // ファイルの投稿がなくても続行する
			log.Debug().Msgf("ok or could not upload file: %v", err)
		}
	}

	// ファイルの表示を確認する
	var isOK bool
	for i := 0; i < p.MaxWaitSecForRequest; i++ {
		isThere, err := p.Page.Locator("div[data-testid='attachments']").IsVisible()
		if err != nil {
			return SetError(err, "could not check the element is visible")
		}

		// 投稿画像及び動画が表示された
		if isThere {
			isOK = true
			log.Debug().Int("gui upload wait sec", p.MaxWaitSecForRequest-i).Msg("ok or could not upload file")
			break
		}

		time.Sleep(time.Second)
	}
	if with_files { // ファイル必須ならば、ファイルの表示を確認してから判断する
		if !isOK {
			return SetError(fmt.Errorf("could not upload file, timeout: past %ds", p.MaxWaitSecForRequest), "could not upload file")
		}
	}

	return nil

}

// filesToInputFiles ファイルをアップロードするための目的の型式に変換する
func filesToInputFiles(files []string) ([]playwright.InputFile, error) {
	var inputFiles []playwright.InputFile
	for _, file := range files {
		name, buffer, err := readFile(file)
		if err != nil {
			return nil, err
		}

		fileType := http.DetectContentType(buffer)

		log.Debug().Str("function", "filesToInputFiles").Msgf("file: %s, type: %s, byte size: %d", name, fileType, len(buffer))

		inputFiles = append(inputFiles, playwright.InputFile{
			Name:     name,
			MimeType: fileType,
			Buffer:   buffer,
		})
	}

	return inputFiles, nil
}

// readFile 小さいインスタンスでも大きいファイルを扱うため、チャンクで読み込む
func readFile(file string) (string, []byte, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", nil, SetError(err, "could not open file")
	}
	defer f.Close()

	// チャンクで読み込みを行う
	var buffer []byte
	buf := make([]byte, 1024*1024) // 1MBのバッファ
	for {
		// ファイルからデータを読み込む
		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			return "", nil, err // 読み込み中にエラーが発生した場合
		}
		if n == 0 {
			break // ファイルの終端に達した場合、読み込みを終了
		}

		// 読み込んだデータをバッファに追加
		buffer = append(buffer, buf[:n]...)
	}

	return f.Name(), buffer, nil
}

func SetError(err error, msg any) error {
	var s string
	switch v := msg.(type) {
	case string:
		s = v
	case error:
		s = v.Error()
	default:
		s = fmt.Sprintf("%v", msg)
	}

	return fmt.Errorf("%v > %v", err, errors.New(s))
}

func wait(ms int) {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	millisec := time.Duration(r.Intn(ms)) * time.Millisecond
	time.Sleep(time.Second + millisec)
}

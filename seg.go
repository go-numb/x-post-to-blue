package xpostblue

import (
	"fmt"
	"net/url"

	"github.com/rs/zerolog/log"
)

// 各処理を細分化

// ログインページへ遷移
func (p *ClientBody) ToLogin() error {
	u, _ := url.Parse(p.PostLocator.LoginURL)
	if _, err := p.Page.Goto(u.String()); err != nil {
		return fmt.Errorf("%v > could not goto", err)
	}
	log.Debug().Msgf("target url: %s", u.String())

	return nil
}

// ログインページで、InputID ユーザー名/メールアドレスを入力
// ID/Email/TEL入力欄となっている
func (p *ClientBody) InputID(usernameOrEmail string) error {
	if isThere, err := p.IsThere(p.PostLocator.InputID); err != nil || !isThere {
		return err
	}

	// input Username/Email
	if err := p.Page.Locator(p.PostLocator.InputID).Fill(usernameOrEmail); err != nil {
		return fmt.Errorf("%v > could not fill to account input", err)
	}

	return nil
}

// ログインページで、ClickBtnLogin IDを入力後ボタンをクリック
func (p *ClientBody) ClickBtnLogin() error {
	if isThere, err := p.IsThere(p.PostLocator.BtnID); err != nil || !isThere {
		return err
	}

	if err := p.Page.Locator(p.PostLocator.BtnID).Tap(); err != nil {
		return fmt.Errorf("%v > could not click to next button", err)
	}

	return nil
}

// ログインページで、InputTel 電話番号を入力
func (p *ClientBody) InputTel(tel string) error {
	if isThere, err := p.IsThere(p.PostLocator.InputTel); err != nil || !isThere {
		return err
	}

	if err := p.Page.Locator(p.PostLocator.InputTel).Fill(tel); err != nil {
		return fmt.Errorf("%v > could not fill to tel input", err)
	}

	return nil
}

// ログインページで、ClickTelBtn 電話番号入力後ボタンをクリック
func (p *ClientBody) ClickTelBtn() error {
	if isThere, err := p.IsThere(p.PostLocator.BtnTel); err != nil || !isThere {
		return err
	}

	if err := p.Page.Locator(p.PostLocator.BtnTel).Nth(0).Tap(); err != nil {
		return fmt.Errorf("%v > could not click to next button", err)
	}

	return nil
}

// ログインページで、InputPassword パスワードを入力
func (p *ClientBody) InputPassword(password string) error {
	if isThere, err := p.IsThere(p.PostLocator.InputPass); err != nil || !isThere {
		return err
	}

	if err := p.Page.Locator(p.PostLocator.InputPass).Fill(password); err != nil {
		return fmt.Errorf("%v > could not fill to password input", err)
	}

	return nil
}

// ログインページで、ClickPassBtn パスワード入力後ボタンをクリック
func (p *ClientBody) ClickPassBtn() error {
	if isThere, err := p.IsThere(p.PostLocator.BtnPass); err != nil || !isThere {
		return err
	}

	if err := p.Page.Locator(p.PostLocator.BtnPass).Nth(0).Tap(); err != nil {
		return fmt.Errorf("%v > could not click to next button", err)
	}

	return nil
}

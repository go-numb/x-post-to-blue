package xpostblue

import (
	"fmt"
	"math/rand"
	"net/url"
	"time"

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

	switch p.ClickType {
	case ClickTypeClick:
		if err := p.Page.Locator(p.PostLocator.BtnID).Click(); err != nil {
			return fmt.Errorf("%v > could not click to next button", err)
		}

	case ClickTypeTap:
		if err := p.Page.Locator(p.PostLocator.BtnID).Tap(); err != nil {
			return fmt.Errorf("%v > could not click to next button", err)
		}

	default:
		return fmt.Errorf("click type is not defined")
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

// ログインページで、ClickBtnTel 電話番号入力後ボタンをクリック
func (p *ClientBody) ClickBtnTel() error {
	if isThere, err := p.IsThere(p.PostLocator.BtnTel); err != nil || !isThere {
		return err
	}

	switch p.ClickType {
	case ClickTypeClick:
		if err := p.Page.Locator(p.PostLocator.BtnTel).Nth(0).Click(); err != nil {
			return fmt.Errorf("%v > could not click to next button", err)
		}

	case ClickTypeTap:
		if err := p.Page.Locator(p.PostLocator.BtnTel).Nth(0).Tap(); err != nil {
			return fmt.Errorf("%v > could not click to next button", err)
		}

	default:
		return fmt.Errorf("click type is not defined")
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
func (p *ClientBody) ClickBtnPass() error {
	if isThere, err := p.IsThere(p.PostLocator.BtnPass); err != nil || !isThere {
		return err
	}

	switch p.ClickType {
	case ClickTypeClick:
		if err := p.Page.Locator(p.PostLocator.BtnPass).Nth(0).Click(); err != nil {
			return fmt.Errorf("%v > could not click to next button", err)
		}

	case ClickTypeTap:
		if err := p.Page.Locator(p.PostLocator.BtnPass).Nth(0).Tap(); err != nil {
			return fmt.Errorf("%v > could not click to next button", err)
		}

	default:
		return fmt.Errorf("click type is not defined")
	}

	return nil
}

// 投稿ページへ遷移
func (p *ClientBody) ToPost() error {
	u, _ := url.Parse(p.PostLocator.ProURL)
	if _, err := p.Page.Goto(u.String()); err != nil {
		return fmt.Errorf("%v > could not goto post page", err)
	}
	log.Debug().Msgf("target url: %s", u.String())

	return nil
}

// 投稿ページで、InputText 投稿内容を入力
func (p *ClientBody) InputText(msg string) error {
	if isThere, err := p.IsThere(p.PostLocator.InputMsg); err != nil || !isThere {
		return err
	}

	if err := p.Page.Locator(p.PostLocator.InputMsg).Fill(msg); err != nil {
		return fmt.Errorf("%v > could not fill to post", err)
	}

	return nil
}

// 投稿ページで、InputFiles 投稿ファイルをアップロード
func (p *ClientBody) InputFiles(with_file bool, files ...string) error {
	// upload files
	if err := p.uploadFiles(with_file, files...); err != nil {
		return fmt.Errorf("%v > could not upload files", err)
	}

	return nil
}

// 投稿ページで、ClickBtnPost 投稿ボタンをクリック
func (p *ClientBody) ClickBtnPost(isPost bool) error {
	if !isPost {
		return fmt.Errorf("post enabled, is_post is false")
	}

	switch p.ClickType {
	case ClickTypeClick:
		if err := p.Page.Locator(p.PostLocator.BtnPost).Click(); err != nil {
			return fmt.Errorf("%v > could not click to post button", err)
		}

	case ClickTypeTap:
		if err := p.Page.Locator(p.PostLocator.BtnPost).Tap(); err != nil {
			return fmt.Errorf("%v > could not click to post button", err)
		}

	default:
		return fmt.Errorf("click type is not defined")
	}

	return nil
}

// Wait 指定時間（ミリ秒）待機
// minwaitmillsec引数を最低待機時間、ms引数を上限とした乱数msを追加し、待機
func (p *ClientBody) Wait(minWaitMillisec, ms int) {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	millisec := time.Duration(r.Intn(ms)) * time.Millisecond

	// 最低待機時間
	time.Sleep(time.Duration(minWaitMillisec * int(time.Millisecond)))
	// 乱数待機時間
	time.Sleep(millisec)
}

package xpostblue

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	client := New(false, nil)
	// to not blue accont for post
	client.PostLocator.ProURL = "https://twitter.com/compose/post"
	client.SetDefaultTimeout(60)
	defer client.Close()

	var (
		id   = ""
		tel  = ""
		pass = ""

		minWaitMs  = 1000
		randWaitMs = 2000
	)

	assert.NoError(t, client.ToLogin(), "could not goto")
	client.Wait(minWaitMs, randWaitMs)

	assert.NoError(t, client.InputID(id), "could not fill to account input")
	client.Wait(minWaitMs, randWaitMs)
	assert.NoError(t, client.ClickBtnLogin(), "could not click to next button")
	client.Wait(minWaitMs, randWaitMs)

	assert.NoError(t, client.InputTel(tel), "could not fill to tel input")
	client.Wait(minWaitMs, randWaitMs)
	assert.NoError(t, client.ClickBtnTel(), "could not click to next button")
	client.Wait(minWaitMs, randWaitMs)

	assert.NoError(t, client.InputPassword(pass), "could not fill to password input")
	client.Wait(minWaitMs, randWaitMs)
	assert.NoError(t, client.ClickBtnPass(), "could not click to next button")
	client.Wait(minWaitMs, randWaitMs)

	assert.NoError(t, client.Post(true, 60, "test tweet"), "could not post")
}

func TestWait(t *testing.T) {
	client := New(false, nil)
	defer client.Close()

	minWait := 1000
	randWait := 2000

	start := time.Now()
	defer func() {
		log.Printf("elapsed: %v", time.Since(start))
	}()

	client.Wait(minWait, randWait)

	// 待機時間がminWait以上であることを確認する
	assert.True(t, time.Since(start) > time.Duration(minWait)*time.Millisecond)

	// 待機時間がminWait + randWait未満であることを確認する
	assert.True(t, time.Since(start) < time.Duration(minWait+randWait)*time.Millisecond)
}

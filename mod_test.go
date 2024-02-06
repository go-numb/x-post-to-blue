package xpostblue

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	USERNAME = ""
	PASSWORD = ""

	MSG = `むすめのむかし、とある村に小さなおうじ様が住んでいました。おうじ様は美しい王国の王子であり、優しい心を持っていました。

	ある日、おうじ様は王国を旅している途中で、嵐に遭遇し道に迷ってしまいました。そこへやってきたのは、小さな光を身にまとった不思議な少女でした。
	
	少女は「私は光の妖精です。おうじ様、道を教えましょうか？」と声をかけてくれたので、おうじ様は驚きながらも嬉しさを感じました。
	
	光の妖精はおうじ様に、嵐が過ぎ去るまで一緒に待っていてくれると言ってくれました。おうじ様は小さな妖精に感謝しながら、一緒に庇を作って嵐をしのぎました。
	
	嵐が過ぎ去り、おうじ様は妖精と共に王国に帰還しました。おうじ様は妖精にお礼を言おうとしましたが、妖精は微笑みながら「おうじ様の心が輝いていることが私の報酬です」と言って姿を消してしまいました。
	
	おうじ様はその後も妖精のことを忘れずに、王国での生活を続けました。ある日、おうじ様は王国の中庭で美しい光の粒を見つけ、それが妖精からの贈り物だと気付きました。
	
	おうじ様はその光を宝物にし、王国の人々に誇らしげに見せました。それ以降、王国はますます繁栄し、おうじ様と妖精のふれあいは語り継がれるようになりました。`
)

func TestDo(t *testing.T) {
	isHeadless := false
	isPost := false
	sleepSec := 5
	files := []string{"./images/1.jpg", "./images/2.jpg"}

	client := New(isHeadless)
	defer client.Close()

	err := client.Login(USERNAME, PASSWORD)
	assert.NoError(t, err)

	err = client.Post(isPost, sleepSec, MSG, files)
	assert.NoError(t, err)

	time.Sleep(10 * time.Second)

	fmt.Println("login success!")
}

package sensitive_words

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
)

func TestHit(t *testing.T) {
	st := New(
		buildWordsCall,
		WithMode(ModePinyin),
		WithMaskWord('*'),
		WithRebuildWordsInterval(time.Second*10),
	)
	ctx := context.Background()
	for word, hit := range map[string]bool{
		"":      false,
		"shazi": true,
		"傻子":    true,
		"傻逼":    true,
		"大傻逼":   true,
	} {
		isHit, hitWord, err := st.Hit(ctx, word)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(fmt.Sprintf("hit: word: %s, hit_word: %s, want: %t, result: %t", word, hitWord, hit, isHit))
		assert.Equal(t, isHit, hit)
	}
}

func TestHitMust(t *testing.T) {
	st := New(
		buildWordsCall,
		WithMode(ModePinyin),
		WithMaskWord('*'),
		WithRebuildWordsInterval(time.Second*10),
	)
	ctx := context.Background()
	for word, hit := range map[string]bool{
		"你这个傻子": true,
		"你这个傻瓜": false,
		"shazi": true,
		"傻子":    true,
		"大傻逼":   true,
	} {
		isHit, _, err := st.HitMust(ctx, word, 1)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(fmt.Sprintf("hit: word: %s, want: %t, result: %t", word, hit, isHit))
		assert.Equal(t, isHit, hit)
	}

	isHit, _, err := st.HitMust(ctx, "傻子", 1)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, isHit, true)

	isHit, _, err = st.HitMust(ctx, "傻瓜", 1)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, isHit, false)
}

func TestMatchReplace(t *testing.T) {
	st := New(
		buildWordsCall,
		WithMode(ModePinyin),
		WithMaskWord('*'),
		WithRebuildWordsInterval(time.Second*10),
	)
	ctx := context.Background()
	for text, data := range map[string]struct {
		IsHit   bool
		NewText string
	}{
		"你这个小美女":      {false, "你这个小美女"},
		"你这个丑八怪":      {true, "你这个***"},
		"丑东西":         {true, "***"},
		"丑（）东西":       {true, "*（）**"},
		"你也太丑了吧":      {false, "你也太丑了吧"},
		"色情直播":        {true, "**直播"},
		"色--。。。//情直播": {true, "*--。。。//*直播"},
		"方舟子我问候你全家":   {false, "方舟子我问候你全家"},
		"方舟子傻逼我问候你全家": {true, "方舟子**我问候你全家"},
		"方舟子傻逼早就该死了":  {true, "*****早就该**"},
		"司马南在美国买房子":   {true, "***在**买房子"},
		"司马南在中国买房子":   {false, "司马南在中国买房子"},
	} {
		isHit, newText, err := st.MatchReplace(ctx, text)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(fmt.Sprintf(
			"hit: text: %s, want_hit: %t, result_hit:%t, want_text: %s, result_text: %s",
			text, data.IsHit, isHit, data.NewText, newText))
		assert.Equal(t, isHit, data.IsHit)
		assert.Equal(t, newText, data.NewText)
	}
}

func TestInfos(t *testing.T) {
	st := New(
		buildWordsCall,
		WithMode(ModePinyin),
		WithMaskWord('*'),
		WithRebuildWordsInterval(time.Second*10),
	)
	ctx := context.Background()
	t.Log(st.DebugInfos(ctx))
}

func buildWordsCall(ctx context.Context) (words []string, err error) {
	return []string{
		"丑八怪",
		"丑东西",
		"丑女",
		"色情",
		"色魔",
		"傻子",
		"傻逼",
		"司马南|美国",
		"方舟子|死了",
	}, nil
}

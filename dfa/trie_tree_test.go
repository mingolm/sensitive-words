package dfa

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestFilterChar(t *testing.T) {
	tree := NewTrieTree(nil)
	runeMap := map[rune]bool{
		'a': false,
		'0': false,
		'1': false,
		'你': false,
		'-': true,
		')': true,
		']': true,
		'💗': true,
	}

	for ch, want := range runeMap {
		assert.Equal(t, tree.isFilterChar(ch), want)
	}
}

func TestHit(t *testing.T) {
	tree := NewTrieTree(nil)
	tree.AddWords([]string{
		"傻逼", "煞笔", "垃圾", "小啦",
	}...)

	wordMap := map[string]struct {
		isHit bool
		word  string
	}{
		"我觉得你是傻逼":     {true, "傻逼"},
		"我觉得你是垃圾":     {true, "垃圾"},
		"我觉得你是小可爱":    {false, ""},
		"我觉得你是，垃、！！圾": {true, "，垃、！！圾"},
		"我觉得你是，-- 垃":  {false, ""},
		"我觉得你是小可爱啦":   {false, ""},
	}
	for word, want := range wordMap {
		isHit, hitWord := tree.Detect(word, false)
		assert.Equal(t, isHit, want.isHit)
		assert.Equal(t, hitWord, want.word)
	}
}

func TestReplace(t *testing.T) {
	tree := NewTrieTree(nil)
	tree.AddWords([]string{
		"傻逼", "煞笔", "垃圾", "小啦",
	}...)

	wordMap := map[string]struct {
		isHit bool
		word  string
	}{
		"我觉得你是傻逼":      {true, "我觉得你是**"},
		"我觉得你是垃圾":      {true, "我觉得你是**"},
		"我觉得你是垃00圾":    {false, "我觉得你是垃00圾"},
		"我觉得你是-=-垃=-圾": {true, "我觉得你是-=-*=-*"},
		"我觉得你是小可爱":     {false, "我觉得你是小可爱"},
		"我觉得你是--小可爱":   {false, "我觉得你是--小可爱"},
	}
	for word, want := range wordMap {
		isHit, hitWord := tree.Replace(word, '*')
		assert.Equal(t, isHit, want.isHit)
		assert.Equal(t, hitWord, want.word)
	}
}

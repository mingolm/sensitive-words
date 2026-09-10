package sensitive_words

import (
	"errors"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestModeContain(t *testing.T) {
	mode := ModePinyin | ModeStats
	assert.Equal(t, mode.Contain(ModePinyin), true)
	assert.Equal(t, mode.Contain(ModeStats), true)

	pinyin := ModePinyin
	assert.Equal(t, pinyin.Contain(ModePinyin), true)
	assert.Equal(t, pinyin.Contain(ModeStats), false)

	stats := ModeStats
	assert.Equal(t, stats.Contain(ModeStats), true)
}

func TestModeRange(t *testing.T) {
	var visited []Mode
	err := (ModePinyin | ModeStats).Range(func(value Mode) error {
		visited = append(visited, value)
		return nil
	})
	assert.Equal(t, err, nil)
	assert.Equal(t, visited, []Mode{ModePinyin, ModeStats})

	visited = nil
	err = ModePinyin.Range(func(value Mode) error {
		visited = append(visited, value)
		return nil
	})
	assert.Equal(t, err, nil)
	assert.Equal(t, visited, []Mode{ModePinyin})

	stopErr := errors.New("stop")
	callCount := 0
	err = (ModePinyin | ModeStats).Range(func(value Mode) error {
		callCount++
		return stopErr
	})
	assert.Equal(t, err, stopErr)
	assert.Equal(t, callCount, 1)
}

func TestPinyinWordReg(t *testing.T) {
	assert.Equal(t, pinyinWordReg.MatchString("傻子"), true)
	assert.Equal(t, pinyinWordReg.MatchString("司马南|美国"), true)
	assert.Equal(t, pinyinWordReg.MatchString("罗永浩|直播|翻车"), true)
	assert.Equal(t, pinyinWordReg.MatchString("shazi"), false)
	assert.Equal(t, pinyinWordReg.MatchString("傻子abc"), false)
	assert.Equal(t, pinyinWordReg.MatchString(""), false)
	assert.Equal(t, pinyinWordReg.MatchString("abc"), false)
	assert.Equal(t, pinyinWordReg.MatchString("123"), false)
}

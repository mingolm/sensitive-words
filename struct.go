package sensitive_words

import (
	"context"
	"regexp"
)

type Mode int

const (
	ModeNormal Mode = iota
	ModePinyin
	ModeEn
	ModeStrict
)

func (t *Mode) Contain(m Mode) bool {
	return *t&m == m
}

type BuildWordsFn func(ctx context.Context) ([]string, error)

// 所有非中英文数字
var defaultFilterSpecialReg = regexp.MustCompile("[^\u4e00-\u9fa5a-zA-Z0-9]")

// 中文
var chineseReg = regexp.MustCompile("^\\p{Han}+([\u00B7\u2022\u2027\u30FB\u002E\u0387\u16EB\u2219\u22C5\uFF65\u05BC]\\p{Han}+)*?$")

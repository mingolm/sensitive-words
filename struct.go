package sensitive_words

import (
	"context"
	"regexp"
)

type Mode int

const (
	ModePinyin Mode = iota
	ModeEn
)

func (t *Mode) Contain(m Mode) bool {
	return *t&m == m
}

type BuildWordsFn func(ctx context.Context) ([]string, error)

// 中文
var chineseReg = regexp.MustCompile("^\\p{Han}+([\u00B7\u2022\u2027\u30FB\u002E\u0387\u16EB\u2219\u22C5\uFF65\u05BC]\\p{Han}+)*?$")

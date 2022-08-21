package sensitive_words

import (
	"time"

	"go.uber.org/zap"
)

type options struct {
	// 掩码字符，默认使用 *
	maskWord rune
	// 查找模式
	mode Mode
	// 过滤特殊字符，默认过滤除中英文数字之外的所有字符
	filterChars []rune
	// 定时触发的 callback
	rebuildWordsInterval time.Duration
	// 创建敏感词回调方法
	buildWordsCall BuildWordsFn
	// 日志
	logger *zap.SugaredLogger
}

type Option func(*options)

func WithMaskWord(word rune) Option {
	return func(o *options) {
		o.maskWord = word
	}
}

func WithMode(modes ...Mode) Option {
	return func(o *options) {
		for _, m := range modes {
			o.mode |= m
		}
	}
}

func WithFilterChars(filterChars ...rune) Option {
	return func(o *options) {
		o.filterChars = filterChars
	}
}

func WithRebuildWordsInterval(interval time.Duration) Option {
	return func(o *options) {
		o.rebuildWordsInterval = interval
	}
}

func WithLogger(logger *zap.SugaredLogger) Option {
	return func(o *options) {
		o.logger = logger
	}
}

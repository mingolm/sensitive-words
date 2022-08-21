package sensitive_words

import (
	"context"
	"strings"
	"sync/atomic"
	"time"

	"github.com/mingolm/sensitive-words/dfa"
	"github.com/mozillazg/go-pinyin"
	"go.uber.org/zap"
)

type SensitiveWorder interface {
	// Hit 判断是否命中敏感词
	Hit(ctx context.Context, word string) (isHit bool, hitWord string, err error)
	// MatchReplace 敏感词替换
	MatchReplace(ctx context.Context, text string) (isHit bool, lastText string, err error)
	// DebugInfos 输出当前所有敏感词
	DebugInfos(ctx context.Context) (words []string)
}

var _ SensitiveWorder = (*sensitiveWord)(nil)

func New(buildWords BuildWordsFn, opts ...Option) SensitiveWorder {
	o := options{
		maskWord:       '*',
		buildWordsCall: buildWords,
		logger:         zap.S().Named("sensitive"),
	}
	for _, fn := range opts {
		fn(&o)
	}

	st := &sensitiveWord{
		options: o,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	_ = cancel
	if err := st.buildWords(ctx); err != nil {
		st.logger.Panicw("build words failed",
			"err", err,
		)
	}

	st.logger.Debug("init success")

	if st.rebuildWordsInterval > 0 {
		go func() {
			ticker := time.NewTicker(st.rebuildWordsInterval)
			for {
				select {
				case <-ticker.C:
					if err := st.buildWords(ctx); err != nil {
						st.logger.Errorw("rebuild words failed",
							"err", err,
						)
					}
				}
			}
		}()
	}

	return st
}

type sensitiveWord struct {
	options
	trieTree atomic.Value
}

func (st *sensitiveWord) buildWords(ctx context.Context) error {
	st.logger.Debugw("rebuild words",
		"start_time", time.Now().Format("2006-01-02 15:04:05"),
	)
	words, err := st.buildWordsCall(ctx)
	if err != nil {
		return err
	}

	tree := dfa.NewTrieTree(st.filterChars)
	// 开启拼音模式
	if st.mode.Contain(ModePinyin) {
		for _, word := range words {
			if !chineseReg.MatchString(word) {
				continue
			}
			words = append(words, strings.Join(pinyin.LazyConvert(word, nil), ""))
		}
	}

	tree.AddWords(words...)

	st.trieTree.Store(tree)

	st.logger.Debugw("rebuild words success",
		"end_time", time.Now().Format("2006-01-02 15:04:05"),
	)

	return nil
}

func (st *sensitiveWord) Hit(ctx context.Context, word string) (isHit bool, hitWord string, err error) {
	tree := st.trieTree.Load().(*dfa.TrieTree)
	isHit, hitWord = tree.Detect(word, st.mode.Contain(ModeStrict))
	return isHit, hitWord, nil
}

func (st *sensitiveWord) MatchReplace(ctx context.Context, text string) (isHit bool, lastText string, err error) {
	tree := st.trieTree.Load().(*dfa.TrieTree)
	isHit, lastText = tree.Replace(text, st.maskWord)
	return isHit, lastText, nil
}

func (st *sensitiveWord) DebugInfos(ctx context.Context) (words []string) {
	tree := st.trieTree.Load().(*dfa.TrieTree)
	return tree.DebugInfos()
}

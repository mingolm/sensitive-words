package sensitive_words

import (
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	"go.uber.org/zap"
)

func TestWithMaskWord(t *testing.T) {
	o := options{maskWord: '*'}
	WithMaskWord('#')(&o)
	assert.Equal(t, o.maskWord, '#')
}

func TestWithMode(t *testing.T) {
	o := options{}
	WithMode(ModePinyin, ModeStats)(&o)
	assert.Equal(t, o.mode, ModePinyin|ModeStats)
	assert.Equal(t, o.mode.Contain(ModeStats), true)
}

func TestWithFilterChars(t *testing.T) {
	o := options{}
	WithFilterChars('-', '*')(&o)
	assert.Equal(t, o.filterChars, []rune{'-', '*'})
}

func TestWithRebuildWordsInterval(t *testing.T) {
	o := options{}
	WithRebuildWordsInterval(time.Second * 10)(&o)
	assert.Equal(t, o.rebuildWordsInterval, time.Second*10)
}

func TestWithLogger(t *testing.T) {
	logger := zap.S().Named("test")
	o := options{}
	WithLogger(logger)(&o)
	assert.Equal(t, o.logger, logger)
}

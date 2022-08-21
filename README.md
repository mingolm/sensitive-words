# Sensitive-Words

```shell
1. 支持敏感词的查找
2. 支持敏感词替换
3. 支持组合词的查找
4. 支持组合词的替换
```

### 用法

```go

package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"go.uber.org/zap"

	sw "github.com/mingolm/sensitive-words"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	st := sw.New(
		buildWordsCall,
		sw.WithMode(sw.ModePinyin, sw.ModeStats),
		sw.WithMaskWord('*'),
		sw.WithRebuildWordsInterval(time.Second*5),
		sw.WithLogger(logger.Sugar()),
	)
	ctx := context.Background()

	// 判断敏感词是否命中
	isHit, hitWord, err := st.Hit(ctx, "你这个丑八怪")
	if err != nil {
		panic(err)
	}
	// 输出：is_hit: true, hit_word: 丑八怪
	fmt.Printf("is_hit: %t, hit_word: %s\n", isHit, hitWord)

	// 敏感词替换
	isHit, newText, err := st.MatchReplace(ctx, "你这个丑逼")
	if err != nil {
		panic(err)
	}
	// 输出：is_hit: true, new_text: 你这个**
	fmt.Printf("is_hit: %t, new_text: %s\n", isHit, newText)

	// 组合词匹配
	isHit, hitWord, err = st.Hit(ctx, "听说司马南在美国买房子")
	if err != nil {
		panic(err)
	}
	// 输出：is_hit: true, hit_word: 司马南|美国
	fmt.Printf("is_hit: %t, hit_word: %s\n", isHit, hitWord)

	// 组合词替换
	isHit, newText, err = st.MatchReplace(ctx, "听说司马南在美国买房子")
	if err != nil {
		panic(err)
	}
	// 输出：is_hit: true, new_text: 听说***在**买房子
	fmt.Printf("is_hit: %t, new_text: %s\n", isHit, newText)

	// debug info
	infos := st.DebugInfos(ctx)
	for _, info := range infos {
		// 输出：word: 丑八怪, hit_count: 1
		// 输出：word: 丑逼, hit_count: 1
		// 输出：word: 司马南|美国, hit_count: 2
		// 输出：word: 方舟子|死了, hit_count: 0
		// 输出：word: choubaguai, hit_count: 0
		// 输出：word: choubi, hit_count: 0
		// 输出：word: simananmeiguo, hit_count: 0
		// 输出：word: fangzhouzisile, hit_count: 0
		fmt.Printf("word: %s, hit_count: %d\n", info.Word, info.HitCount)
	}

	time.Sleep(time.Minute)
}

func buildWordsCall(ctx context.Context) (words []string, err error) {
	time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
	return []string{
		"丑八怪",
		"丑逼",
		"司马南|美国", // 组合词
		"方舟子|死了",
	}, nil
}


```
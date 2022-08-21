# Sensitive-Words

```shell
1. 支持敏感词的查找
2. 支持敏感词替换
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
	st := sw.New(
		buildWordsCall,
		sw.WithMode(sw.ModePinyin),
		sw.WithRebuildWordsInterval(time.Second*5),
	)
	ctx := context.Background()
	for word, hit := range map[string]bool{
		"你这个傻瓜": false,
		"shazi":     true,
		"傻子":    true,
		"傻逼":    true,
		"大傻逼":   true,
	} {
		isHit, err := st.Hit(ctx, word)
		if err != nil {
			panic(err)
		}

		fmt.Printf("hit: word: %s, want: %t, result:%t\n", word, hit, isHit)
	}

	time.Sleep(time.Minute)
}

func buildWordsCall(ctx context.Context) (words []string, err error) {
	time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
	return []string{
		"傻逼",
		"傻子",
	}, nil
}
```
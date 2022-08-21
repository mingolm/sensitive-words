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
		sw.WithMode(sw.ModePinyin),
		sw.WithMaskWord('*'),
		sw.WithRebuildWordsInterval(time.Second*5),
		sw.WithLogger(logger.Sugar()),
	)
	ctx := context.Background()
	for word, hit := range map[string]bool{
		"你这个傻瓜":    false,
		"shazi":    true,
		"傻子":       true,
		"傻逼":       true,
		"大傻逼":      true,
		"方舟子早就该死了": true,
		"方舟子还活着":   false,
	} {
		isHit, hitWord, err := st.Hit(ctx, word)
		if err != nil {
			panic(err)
		}

		fmt.Printf("hit: word: %s, hit_word: %s, want: %t, result: %t\n", word, hitWord, hit, isHit)
	}

	time.Sleep(time.Minute)
}

func buildWordsCall(ctx context.Context) (words []string, err error) {
	time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
	return []string{
		"傻逼",
		"傻子",
		"司马南|美国", // 组合词
		"方舟子|死了",
	}, nil
}

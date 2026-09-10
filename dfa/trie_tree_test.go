package dfa

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestFilterChar(t *testing.T) {
	tree := NewTrieTree()
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
	tree := NewTrieTree()
	tree.AddWords([]string{
		"傻逼", "煞笔", "垃圾", "小啦", "傻瓜｜笨猪", "司马南|美国",
	}...)

	wordMap := map[string]struct {
		isHit bool
		word  string
	}{
		"我觉得你是傻逼":     {true, "傻逼"},
		"我觉得你是垃圾":     {true, "垃圾"},
		"我觉得你是小可爱":    {false, ""},
		"我觉得你是，垃、！！圾": {true, "垃圾"},
		"我觉得你是，-- 垃":  {false, ""},
		"我觉得你是小可爱啦":   {false, ""},
		"司马南是两面派":     {false, "司马南"},
	}
	for word, want := range wordMap {
		isHit, hitWords := tree.Detect(word, 1)
		assert.Equal(t, isHit, want.isHit)
		if isHit {
			assert.Equal(t, hitWords[0], want.word)
		}
	}

	// 命中多次
	isHit, hitWords := tree.Detect("我觉得你是个垃圾傻逼", 2)
	assert.Equal(t, isHit, true)
	assert.Equal(t, hitWords, []string{"垃圾", "傻逼"})

	isHit, hitWords = tree.Detect("我觉得你是个垃圾傻瓜", 4)
	assert.Equal(t, isHit, false)
	assert.Equal(t, hitWords, []string{"垃圾"})

	// 组合词
	isHit, hitWords = tree.Detect("我觉得司马南是傻逼", 1)
	assert.Equal(t, isHit, true)
	assert.Equal(t, hitWords, []string{"傻逼"})

	isHit, hitWords = tree.Detect("我觉得司马南是人才", 1)
	assert.Equal(t, isHit, false)

	isHit, hitWords = tree.Detect("司马南否认在美国买房子", 1)
	assert.Equal(t, isHit, true)
	assert.Equal(t, hitWords, []string{"司马南|美国"})
}

func TestReplace(t *testing.T) {
	tree := NewTrieTree()
	tree.AddWords([]string{
		"傻逼", "煞笔", "垃圾", "小啦", "司马南|美国", "方舟子|死了",
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
		"司马南在美国买房子":    {true, "***在**买房子"},
		"司马南在中国买房子":    {false, "司马南在中国买房子"},
		"方舟子我问候你全家":    {false, "方舟子我问候你全家"},
		"方舟子傻逼我问候你全家":  {true, "方舟子**我问候你全家"},
		"方舟子傻逼早就该死了":   {true, "*****早就该**"},
	}
	for word, want := range wordMap {
		isHit, hitWord := tree.Replace(word, '*')
		assert.Equal(t, isHit, want.isHit)
		assert.Equal(t, hitWord, want.word)
	}
}

func TestNewTrieTree(t *testing.T) {
	tree := NewTrieTree()
	assert.NotEqual(t, tree, nil)
	assert.NotEqual(t, tree.root, nil)
	assert.Equal(t, tree.root.isRoot, true)
	assert.Equal(t, tree.root.character, '0')
	assert.NotEqual(t, tree.root.children, nil)
	assert.NotEqual(t, tree.comboRoot, nil)
	assert.Equal(t, tree.comboRoot.isRoot, true)
	assert.Equal(t, tree.comboRoot.character, '0')
	assert.NotEqual(t, tree.comboRoot.children, nil)
	assert.Equal(t, tree.openStats, false)
	assert.NotEqual(t, tree.filterRuneMap, nil)
}

func TestWithFilterChars(t *testing.T) {
	tree := NewTrieTree().WithFilterChars([]rune{'-', '*'})
	assert.Equal(t, tree.isFilterChar('-'), true)
	assert.Equal(t, tree.isFilterChar('*'), true)
	assert.Equal(t, tree.isFilterChar('a'), false)
	assert.Equal(t, tree.isFilterChar('你'), false)
	assert.Equal(t, tree.isFilterChar(')'), false)
	assert.Equal(t, tree.isFilterChar('💗'), false)
}

func TestWithStats(t *testing.T) {
	tree := NewTrieTree()
	tree.AddWords("傻逼")
	tree.Detect("傻逼", 1)
	for _, stats := range tree.DebugInfos() {
		if stats.Word == "傻逼" {
			assert.Equal(t, stats.HitCount, uint64(0))
		}
	}

	tree = NewTrieTree().WithStats()
	tree.AddWords("傻逼")
	tree.Detect("傻逼", 1)
	found := false
	for _, stats := range tree.DebugInfos() {
		if stats.Word == "傻逼" {
			found = true
			assert.Equal(t, stats.HitCount, uint64(1))
		}
	}
	assert.Equal(t, found, true)

	tree.Detect("你是个傻逼", 1)
	for _, stats := range tree.DebugInfos() {
		if stats.Word == "傻逼" {
			assert.Equal(t, stats.HitCount, uint64(2))
		}
	}
}

func TestAddWordsEdgeCases(t *testing.T) {
	tree := NewTrieTree()
	tree.AddWords("", "傻", "傻逼", "傻逼")

	isHit, hitWords := tree.Detect("", 1)
	assert.Equal(t, isHit, false)
	assert.Equal(t, len(hitWords), 0)

	isHit, hitWords = tree.Detect("傻", 1)
	assert.Equal(t, isHit, true)
	assert.Equal(t, hitWords[0], "傻")

	isHit, _ = tree.Detect("傻逼", 1)
	assert.Equal(t, isHit, true)

	emptyTree := NewTrieTree()
	emptyTree.AddWords("")
	assert.Equal(t, len(emptyTree.DebugInfos()), 0)
}

func TestDebugInfos(t *testing.T) {
	tree := NewTrieTree().WithStats()
	tree.AddWords("傻逼", "垃圾", "司马南|美国")

	infos := tree.DebugInfos()
	wordSet := map[string]uint64{}
	for _, stats := range infos {
		wordSet[stats.Word] = stats.HitCount
	}
	assert.Equal(t, wordSet["傻逼"], uint64(0))
	assert.Equal(t, wordSet["垃圾"], uint64(0))
	assert.Equal(t, wordSet["司马南|美国"], uint64(0))

	tree.Detect("傻逼", 1)
	for _, stats := range tree.DebugInfos() {
		if stats.Word == "傻逼" {
			assert.Equal(t, stats.HitCount, uint64(1))
		}
	}

	empty := &TrieTree{}
	assert.Equal(t, empty.DebugInfos(), nil)
}

func TestNodeIsEnd(t *testing.T) {
	n := NewNode('a')
	assert.Equal(t, n.IsEnd(), false)
	n.isEnd = true
	assert.Equal(t, n.IsEnd(), true)
}

func TestIncrStats(t *testing.T) {
	n := NewNode('a')
	n.incrStats(false)
	assert.Equal(t, n.hitCount.Load(), uint64(0))
	n.incrStats(true)
	assert.Equal(t, n.hitCount.Load(), uint64(1))
	n.incrStats(true)
	assert.Equal(t, n.hitCount.Load(), uint64(2))
}

func TestDetectComboWords(t *testing.T) {
	tree := NewTrieTree()
	tree.AddWords("司马南|美国", "罗永浩|直播|翻车")

	isHit, hitWords := tree.Detect("司马南在美国买房子", 1)
	assert.Equal(t, isHit, true)
	assert.Equal(t, hitWords, []string{"司马南|美国"})

	isHit, _ = tree.Detect("司马南在中国买房子", 1)
	assert.Equal(t, isHit, false)

	isHit, hitWords = tree.Detect("罗永浩在第一场直播的时候肯定翻车", 1)
	assert.Equal(t, isHit, true)
	assert.Equal(t, hitWords, []string{"罗永浩|直播|翻车"})

	isHit, _ = tree.Detect("罗永浩在第一场直播的时候很成功", 1)
	assert.Equal(t, isHit, false)

	isHit, _ = tree.Detect("司马南在美--国买房子", 1)
	assert.Equal(t, isHit, true)

	isHit, lastText := tree.Replace("司马南在美国买房子", '*')
	assert.Equal(t, isHit, true)
	assert.Equal(t, lastText, "***在**买房子")
}

package dfa

import (
	"go.uber.org/atomic"
	"strings"
	"unicode"
)

type TrieTree struct {
	root          *Node
	comboRoot     *Node
	openStats     bool
	filterRuneMap map[rune]struct{}
}

type Node struct {
	isRoot    bool
	isEnd     bool
	character rune
	words     []string
	children  map[rune]*Node
	hitCount  atomic.Uint64
}

// Stats 敏感词统计
type Stats struct {
	Word     string
	HitCount uint64
}

func NewTrieTree() *TrieTree {
	return &TrieTree{
		root: &Node{
			isRoot:    true,
			character: '0',
			children:  make(map[rune]*Node, 0),
		},
		comboRoot: &Node{
			isRoot:    true,
			character: '0',
			children:  make(map[rune]*Node, 0),
		},
		filterRuneMap: map[rune]struct{}{},
	}
}

func (tree *TrieTree) WithFilterChars(filterChars []rune) *TrieTree {
	filterCharMap := make(map[rune]struct{}, len(filterChars))
	for _, c := range filterChars {
		filterCharMap[c] = struct{}{}
	}
	tree.filterRuneMap = filterCharMap
	return tree
}

func (tree *TrieTree) WithStats() *TrieTree {
	tree.openStats = true
	return tree
}

func (tree *TrieTree) AddWords(words ...string) {
	for _, word := range words {
		tree.addWord(false, word)
	}
}

func (tree *TrieTree) addWord(isCombo bool, word string) {
	if word == "" {
		return
	}

	var cur *Node
	if isCombo {
		cur = tree.comboRoot
	} else {
		cur = tree.root
	}
	words := strings.Split(word, "|")
	characters := []rune(words[0])
	for position := 0; position < len(characters); position++ {
		ch := characters[position]
		if tree.isFilterChar(ch) {
			continue
		}
		if next, ok := cur.children[ch]; ok {
			cur = next
		} else {
			newNode := NewNode(ch)
			cur.children[ch] = newNode
			cur = newNode
		}
	}

	cur.isEnd = true
	// 新增组合词
	if len(words) > 1 {
		cur.words = words[1:]
		for _, word = range words[1:] {
			tree.addWord(true, word)
		}
	}
}

func (tree *TrieTree) detectInCombo(text string, words ...string) ([]int, bool) {
	var (
		parent         = tree.comboRoot
		cur            *Node
		found          bool
		runes          = []rune(text)
		length         = len(runes)
		left           = 0
		filterIndexMap = map[int]struct{}{}
		wordMap        = make(map[string]struct{}, len(words))
		indexes        []int
	)
	for _, word := range words {
		wordMap[word] = struct{}{}
	}
	for position := 0; position < length; position++ {
		ch := runes[position]
		if tree.isFilterChar(ch) {
			filterIndexMap[position] = struct{}{}
			continue
		}
		cur, found = parent.children[ch]

		if !found || (!cur.IsEnd() && position == length-1) {
			parent = tree.comboRoot
			position = left
			left++
			continue
		}

		if cur.IsEnd() && left <= position {
			var word []rune
			var indexCh []int
			for i := left; i <= position; i++ {
				// 特殊字符不替换
				if _, ok := filterIndexMap[i]; ok {
					continue
				}
				word = append(word, runes[i])
				indexCh = append(indexCh, i)
			}
			wordStr := string(word)
			if _, ok := wordMap[wordStr]; ok {
				delete(wordMap, wordStr)
				indexes = append(indexes, indexCh...)
				if len(wordMap) == 0 {
					return indexes, true
				}
			}
		}

		parent = cur
	}

	return nil, false
}

func (tree *TrieTree) Detect(text string, times int) (bool, []string) {
	var (
		parent         = tree.root
		cur            *Node
		found          bool
		runes          = []rune(text)
		length         = len(runes)
		left           = 0
		filterIndexMap = map[int]struct{}{}
		hitWords       []string
		isHit          bool
	)

	for position := 0; position < length; position++ {
		ch := runes[position]
		if tree.isFilterChar(ch) {
			filterIndexMap[position] = struct{}{}
			continue
		}
		cur, found = parent.children[ch]

		if !found || (!cur.IsEnd() && position == length-1) {
			parent = tree.root
			position = left
			left++
			continue
		}

		if cur.IsEnd() && left <= position {
			var word []rune
			for i := left; i <= position; i++ {
				// 特殊字符不替换
				if _, ok := filterIndexMap[i]; ok {
					continue
				}
				word = append(word, runes[i])
			}
			// 组合词的情况下，需要另外处理
			if len(cur.words) == 0 {
				isHit = true
				hitWords = append(hitWords, string(word))
				times--
			} else if _, comboHit := tree.detectInCombo(text, cur.words...); comboHit {
				isHit = true
				times -= len(cur.words) + 1
				hitWords = append(hitWords, string(word)+"|"+strings.Join(cur.words, "|"))
			}
		}

		if isHit {
			cur.incrStats(tree.openStats)
		}

		if times <= 0 {
			return isHit, hitWords
		}

		parent = cur
	}

	return times <= 0, hitWords
}

func (tree *TrieTree) Replace(text string, replace rune) (bool, string) {
	var (
		parent         = tree.root
		cur            *Node
		runes          = []rune(text)
		length         = len(runes)
		left           = 0
		found          bool
		isHit          bool
		filterIndexMap = map[int]struct{}{}
	)

	for position := 0; position < len(runes); position++ {
		ch := runes[position]
		if tree.isFilterChar(ch) {
			filterIndexMap[position] = struct{}{}
			continue
		}
		cur, found = parent.children[ch]

		if !found || (!cur.IsEnd() && position == length-1) {
			parent = tree.root
			position = left
			left++
			continue
		}

		if cur.IsEnd() && left <= position {
			// 组合词的情况下，需要另外处理
			if len(cur.words) == 0 {
				isHit = true
			} else {
				replaceIndexes, comboHit := tree.detectInCombo(text, cur.words...)
				if comboHit {
					isHit = true
				}
				for _, i := range replaceIndexes {
					runes[i] = replace
				}
			}
			if isHit {
				cur.incrStats(tree.openStats)
				for i := left; i <= position; i++ {
					// 特殊字符不替换
					if _, ok := filterIndexMap[i]; ok {
						continue
					}
					runes[i] = replace
				}
			}
		}

		parent = cur
	}

	return isHit, string(runes)
}

func (tree *TrieTree) DebugInfos() []*Stats {
	node := tree.root
	if node == nil {
		return nil
	}

	return mapDeepRange([]*Stats{}, "", node.children)
}

func (tree *TrieTree) isFilterChar(ch rune) bool {
	// 过滤指定字符
	if len(tree.filterRuneMap) > 0 {
		_, ok := tree.filterRuneMap[ch]
		return ok
	}

	// 默认过滤非中英文数字
	switch {
	case unicode.Is(unicode.Han, ch): // 汉字
		return false
	case unicode.IsLetter(ch): // 字母
		return false
	case unicode.IsDigit(ch): // 数字
		return false
	}
	return true
}

func NewNode(character rune) *Node {
	return &Node{
		character: character,
		children:  make(map[rune]*Node, 0),
	}
}

func (node *Node) IsEnd() bool {
	return node.isEnd
}

func (node *Node) incrStats(openStats bool) {
	if !openStats {
		return
	}
	node.hitCount.Inc()
}

func mapDeepRange(results []*Stats, word string, maps map[rune]*Node) []*Stats {
	for ch, cur := range maps {
		currentWord := word
		currentWord += string(ch)
		if cur.children != nil {
			results = mapDeepRange(results, currentWord, cur.children)
		}
		if cur.IsEnd() {
			if len(cur.words) > 0 {
				currentWord += "|" + strings.Join(cur.words, "|")
			}
			results = append(results, &Stats{
				Word:     currentWord,
				HitCount: cur.hitCount.Load(),
			})
		}
	}
	return results
}

package dfa

import (
	"unicode"
)

type TrieTree struct {
	root          *Node
	filterRuneMap map[rune]struct{}
}

type Node struct {
	isRoot    bool
	isEnd     bool
	Character rune
	Children  map[rune]*Node
}

func NewTrieTree(filterChars []rune) *TrieTree {
	filterCharMap := make(map[rune]struct{}, len(filterChars))
	for _, c := range filterChars {
		filterCharMap[c] = struct{}{}
	}
	return &TrieTree{
		root: &Node{
			isRoot:    true,
			Character: '0',
			Children:  make(map[rune]*Node, 0),
		},
		filterRuneMap: filterCharMap,
	}
}

func (tree *TrieTree) AddWords(words ...string) {
	for _, word := range words {
		tree.addWord(word)
	}
}

func (tree *TrieTree) addWord(word string) {
	if word == "" {
		return
	}

	cur := tree.root
	characters := []rune(word)
	for position := 0; position < len(characters); position++ {
		ch := characters[position]
		if tree.isFilterChar(ch) {
			continue
		}

		if next, ok := cur.Children[ch]; ok {
			cur = next
		} else {
			newNode := NewNode(ch)
			cur.Children[ch] = newNode
			cur = newNode
		}

		if position == len(characters)-1 {
			cur.isEnd = true
		}
	}
}

func (tree *TrieTree) Detect(text string, times int) (bool, []string) {
	var (
		parent         = tree.root
		current        *Node
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
		current, found = parent.Children[ch]

		if !found || (!current.IsEnd() && position == length-1) {
			parent = tree.root
			position = left
			left++
			continue
		}

		if current.IsEnd() && left <= position {
			isHit = true
			var word []rune
			for i := left; i <= position; i++ {
				// 特殊字符不替换
				if _, ok := filterIndexMap[i]; ok {
					continue
				}
				word = append(word, runes[i])
			}
			hitWords = append(hitWords, string(word))
			times--
		}

		if times == 0 {
			return isHit, hitWords
		}

		parent = current
	}

	return times == 0, hitWords
}

func (tree *TrieTree) Replace(text string, replace rune) (bool, string) {
	var (
		parent         = tree.root
		current        *Node
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
		current, found = parent.Children[ch]

		if !found || (!current.IsEnd() && position == length-1) {
			parent = tree.root
			position = left
			left++
			continue
		}

		if current.IsEnd() && left <= position {
			isHit = true
			for i := left; i <= position; i++ {
				// 特殊字符不替换
				if _, ok := filterIndexMap[i]; ok {
					continue
				}
				runes[i] = replace
			}
		}

		parent = current
	}

	return isHit, string(runes)
}

func (tree *TrieTree) DebugInfos() []string {
	node := tree.root
	if node == nil {
		return nil
	}

	return mapDeepRange([]string{}, "", node.Children)
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
		Character: character,
		Children:  make(map[rune]*Node, 0),
	}
}

func (node *Node) IsEnd() bool {
	return node.isEnd
}

func mapDeepRange(results []string, word string, maps map[rune]*Node) []string {
	for ch, node := range maps {
		currentWord := word
		currentWord += string(ch)
		if node.Children != nil {
			results = mapDeepRange(results, currentWord, node.Children)
		}
		if node.IsEnd() {
			results = append(results, currentWord)
		}
	}
	return results
}

package unicodeid

import "unicode/utf8"

// Check ascii and unicode for id_start
func IsIDStart(ch rune) bool {
	if ch < utf8.RuneSelf {
		return ASCII_START[ch]
	}
	return IsIDStartUnicode(ch)
}

// Check unicode only for id_start
func IsIDStartUnicode(ch rune) bool {
	chunkIdx := uint(ch) / 8 / CHUNK
	var chunk uint8
	if chunkIdx < uint(len(TRIE_START)) {
		chunk = TRIE_START[chunkIdx]
	} else {
		chunk = 0
	}
	offset := uint(chunk)*CHUNK/2 + uint(ch)/8%CHUNK
	return (LEAF[offset]>>uint(ch%8))&1 != 0
}

// Check ascii and unicode for id_continue
func IsIDContinue(ch rune) bool {
	if ch < utf8.RuneSelf {
		return ASCII_CONTINUE[ch]
	}
	return IsIDContinueUnicode(ch)
}

// Check unicode only for id_continue
func IsIDContinueUnicode(ch rune) bool {
	chunkIdx := uint(ch) / 8 / CHUNK
	var chunk uint8
	if chunkIdx < uint(len(TRIE_CONTINUE)) {
		chunk = TRIE_CONTINUE[chunkIdx]
	} else {
		chunk = 0
	}
	offset := uint(chunk)*CHUNK/2 + uint(ch)/8%CHUNK
	return (LEAF[offset]>>uint(ch%8))&1 != 0
}

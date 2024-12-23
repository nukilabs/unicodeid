package main

import (
	"bytes"
	"fmt"
)

const Head = `const T = true
const F = false
`

func output(properties *Properties, indexStart, indexContinue, halfDense []uint8) string {
	var out bytes.Buffer

	out.WriteString("package unicodeid\n\n")
	out.WriteString("const T = true\n")
	out.WriteString("const F = false\n\n")

	// ASCII_START array
	out.WriteString("var ASCII_START = [128]bool{\n")
	for i := 0; i < 4; i++ {
		out.WriteString("    ")
		for j := 0; j < 32; j++ {
			ch := rune(i*32 + j)
			isIDStart := properties.IsIDStart(ch)
			if isIDStart {
				out.WriteString("T, ")
			} else {
				out.WriteString("F, ")
			}
		}
		out.WriteString("\n")
	}
	out.WriteString("}\n\n")

	// ASCII_CONTINUE array
	out.WriteString("var ASCII_CONTINUE = [128]bool{\n")
	for i := 0; i < 4; i++ {
		out.WriteString("    ")
		for j := 0; j < 32; j++ {
			ch := rune(i*32 + j)
			isIDContinue := properties.IsIDContinue(ch)
			if isIDContinue {
				out.WriteString("T, ")
			} else {
				out.WriteString("F, ")
			}
		}
		out.WriteString("\n")
	}
	out.WriteString("}\n\n")

	// CHUNK constant
	out.WriteString(fmt.Sprintf("const CHUNK = %d\n\n", CHUNK))

	// TRIE_START array
	out.WriteString(fmt.Sprintf("var TRIE_START = [%d]uint8{\n", len(indexStart)))
	for i := 0; i < len(indexStart); i += 16 {
		out.WriteString("    ")
		for j := i; j < i+16 && j < len(indexStart); j++ {
			out.WriteString(fmt.Sprintf("0x%02X, ", indexStart[j]))
		}
		out.WriteString("\n")
	}
	out.WriteString("}\n\n")

	// TRIE_CONTINUE array
	out.WriteString(fmt.Sprintf("var TRIE_CONTINUE = [%d]uint8{\n", len(indexContinue)))
	for i := 0; i < len(indexContinue); i += 16 {
		out.WriteString("    ")
		for j := i; j < i+16 && j < len(indexContinue); j++ {
			out.WriteString(fmt.Sprintf("0x%02X, ", indexContinue[j]))
		}
		out.WriteString("\n")
	}
	out.WriteString("}\n\n")

	// LEAF array
	out.WriteString(fmt.Sprintf("var LEAF = [%d]uint8{\n", len(halfDense)))
	for i := 0; i < len(halfDense); i += 16 {
		out.WriteString("    ")
		for j := i; j < i+16 && j < len(halfDense); j++ {
			out.WriteString(fmt.Sprintf("0x%02X, ", halfDense[j]))
		}
		out.WriteString("\n")
	}
	out.WriteString("}\n")

	return out.String()
}

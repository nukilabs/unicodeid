package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const filename = "DerivedCoreProperties.txt"

type Properties struct {
	ID_Start    map[uint32]struct{}
	ID_Continue map[uint32]struct{}
}

func (p *Properties) IsIDStart(r rune) bool {
	_, ok := p.ID_Start[uint32(r)]
	return ok
}

func (p *Properties) IsIDContinue(r rune) bool {
	_, ok := p.ID_Continue[uint32(r)]
	return ok
}

func parse() (*Properties, error) {
	props := &Properties{
		ID_Start:    make(map[uint32]struct{}),
		ID_Continue: make(map[uint32]struct{}),
	}

	path := filepath.Join(unicodeDir, filename)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}
		lo, hi, name := parseLine(line)
		switch name {
		case "ID_Start":
			for i := lo; i <= hi; i++ {
				props.ID_Start[i] = struct{}{}
			}
		case "ID_Continue":
			for i := lo; i <= hi; i++ {
				props.ID_Continue[i] = struct{}{}
			}
		default:
			continue
		}
	}

	return props, nil
}

func parseLine(line string) (uint32, uint32, string) {
	parts := strings.SplitN(line, ";", 2)
	codepoint, rest := parts[0], parts[1]

	var lo, hi uint32
	if parts2 := strings.SplitN(codepoint, "..", 2); len(parts2) == 2 {
		lo = parseCodepoint(strings.TrimSpace(parts2[0]))
		hi = parseCodepoint(strings.TrimSpace(parts2[1]))
	} else {
		lo = parseCodepoint(strings.TrimSpace(parts2[0]))
		hi = lo
	}

	name := strings.TrimSpace(strings.SplitN(rest, "#", 2)[0])
	return lo, hi, name
}

func parseCodepoint(codepoint string) uint32 {
	n, _ := strconv.ParseUint(codepoint, 16, 32)
	return uint32(n)
}

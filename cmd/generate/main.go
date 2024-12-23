package main

import (
	"log"
	"os"
)

const CHUNK uint = 64

func main() {
	if err := download(); err != nil {
		log.Fatal(err)
	}
	if err := unzip(); err != nil {
		log.Fatal(err)
	}
	props, err := parse()
	if err != nil {
		log.Fatal(err)
	}

	chunkmap := make(map[[CHUNK]uint8]uint8)
	dense := make([][CHUNK]uint8, 0)
	newChunk := func(chunk [CHUNK]uint8) uint8 {
		if prev, ok := chunkmap[chunk]; ok {
			return prev
		}
		dense = append(dense, chunk)
		if n := len(chunkmap); n >= 256 {
			panic("exceeded 256 unique chunks")
		} else {
			chunkmap[chunk] = uint8(n)
			return uint8(n)
		}
	}

	emptychunk := [CHUNK]uint8{}
	newChunk(emptychunk)

	idxStart := make([]uint8, 0)
	idxContinue := make([]uint8, 0)
	for i := uint32(0); i < (uint32(0x10FFFF)+1)/uint32(CHUNK)/8; i++ {
		var startBits, continueBits [CHUNK]uint8
		for j := uint32(0); j < uint32(CHUNK); j++ {
			for k := uint32(0); k < 8; k++ {
				code := (i*uint32(CHUNK)+j)*8 + k
				if code >= 0x80 {
					if props.IsIDStart(rune(code)) {
						startBits[j] |= 1 << k
					}
					if props.IsIDContinue(rune(code)) {
						continueBits[j] |= 1 << k
					}
				}
			}
		}
		idxStart = append(idxStart, newChunk(startBits))
		idxContinue = append(idxContinue, newChunk(continueBits))
	}

	for len(idxStart) > 0 && idxStart[len(idxStart)-1] == 0 {
		idxStart = idxStart[:len(idxStart)-1]
	}
	for len(idxContinue) > 0 && idxContinue[len(idxContinue)-1] == 0 {
		idxContinue = idxContinue[:len(idxContinue)-1]
	}

	halfchunkmap := make(map[[CHUNK / 2]uint8][][CHUNK / 2]uint8)
	for _, chunk := range dense {
		var front, back [CHUNK / 2]uint8
		copy(front[:], chunk[:CHUNK/2])
		copy(back[:], chunk[CHUNK/2:])
		halfchunkmap[front] = append(halfchunkmap[front], back)
	}

	halfdense := make([]uint8, 0)
	denseToHalfdense := make(map[uint8]uint8)
	for _, chunk := range dense {
		originalPos := chunkmap[chunk]
		if _, ok := denseToHalfdense[originalPos]; ok {
			continue
		}
		var front, back [CHUNK / 2]uint8
		copy(front[:], chunk[:CHUNK/2])
		copy(back[:], chunk[CHUNK/2:])

		// Map original dense index to half-dense index
		halfIndex := uint8(uint(len(halfdense)) / (CHUNK / 2))
		denseToHalfdense[originalPos] = halfIndex

		// Add front and back to halfDense
		halfdense = append(halfdense, front[:]...)
		halfdense = append(halfdense, back[:]...)

		// Handle sequential halves
		for {
			nexts, ok := halfchunkmap[back]
			if !ok || len(nexts) == 0 {
				break
			}
			next := nexts[0]
			halfchunkmap[back] = nexts[1:]
			var combined [CHUNK]uint8
			copy(combined[:CHUNK/2], back[:])
			copy(combined[CHUNK/2:], next[:])
			originalPos := chunkmap[combined]
			if _, exists := denseToHalfdense[originalPos]; exists {
				continue
			}
			denseToHalfdense[originalPos] = uint8(uint(len(halfdense))/(CHUNK/2) - 1)
			halfdense = append(halfdense, next[:]...)
			back = next
		}
	}

	for i, idx := range idxStart {
		idxStart[i] = denseToHalfdense[idx]
	}
	for i, idx := range idxContinue {
		idxContinue[i] = denseToHalfdense[idx]
	}

	out := output(props, idxStart, idxContinue, halfdense)
	if err := os.WriteFile("tables.go", []byte(out), 0644); err != nil {
		log.Fatal(err)
	}
}

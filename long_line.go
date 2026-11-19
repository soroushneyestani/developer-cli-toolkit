package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"hash"
	"io"
	"unicode/utf8"
)

type longLineResult struct {
	valid    bool
	matched  bool
	size     int64
	consumed int64
	digest   [sha256.Size]byte
}

// scanLongLine probes a line in fixed-size fragments without buffering its contents.
func scanLongLine(reader *bufio.Reader, first, literal []byte) (longLineResult, error) {
	probe := longLineProbe{valid: true, literal: literal, digest: sha256.New()}
	part, err := first, bufio.ErrBufferFull
	var consumed int64
	pendingCR := false
	carriageReturn := []byte{'\r'}
	for {
		if err != nil && err != io.EOF && err != bufio.ErrBufferFull {
			return longLineResult{}, err
		}
		consumed += int64(len(part))
		if err == nil {
			part = part[:len(part)-1]
			if len(part) == 0 {
				// The previous fragment's trailing CR is part of this CRLF.
				pendingCR = false
			} else {
				part = bytes.TrimSuffix(part, carriageReturn)
			}
		}
		if pendingCR {
			probe.add(carriageReturn)
		}
		// A CR at the end of a fragment might precede a LF in the next read.
		pendingCR = err == bufio.ErrBufferFull && len(part) > 0 && part[len(part)-1] == '\r'
		if pendingCR {
			part = part[:len(part)-1]
		}
		probe.add(part)
		if err != bufio.ErrBufferFull {
			break
		}
		if !probe.valid {
			return longLineResult{}, nil
		}
		part, err = reader.ReadSlice('\n')
	}
	result := longLineResult{
		valid:    probe.valid && probe.utf8.pendingLen == 0,
		matched:  probe.matched,
		size:     probe.size,
		consumed: consumed,
	}
	if result.valid && result.matched {
		probe.digest.Sum(result.digest[:0])
	}
	return result, err
}

type longLineProbe struct {
	valid   bool
	matched bool
	literal []byte
	overlap []byte
	utf8    utf8Boundary
	digest  hash.Hash
	size    int64
}

func (p *longLineProbe) add(part []byte) {
	if !p.valid || len(part) == 0 {
		return
	}
	p.valid = bytes.IndexByte(part, 0) < 0 && p.utf8.accept(part)
	if !p.valid {
		return
	}
	p.size += int64(len(part))
	p.digest.Write(part)
	if p.matched || bytes.Contains(part, p.literal) {
		p.matched = true
		return
	}
	keep := len(p.literal) - 1
	if keep <= 0 {
		return
	}
	p.overlap = append(p.overlap, part[:min(keep, len(part))]...)
	if bytes.Contains(p.overlap, p.literal) {
		p.matched = true
		return
	}
	switch {
	case len(part) >= keep:
		p.overlap = append(p.overlap[:0], part[len(part)-keep:]...)
	case len(p.overlap) > keep:
		copy(p.overlap, p.overlap[len(p.overlap)-keep:])
		p.overlap = p.overlap[:keep]
	}
}

// utf8Boundary retains an incomplete rune between otherwise validated fragments.
type utf8Boundary struct {
	pending    [utf8.UTFMax]byte
	pendingLen int
}

func (v *utf8Boundary) accept(part []byte) bool {
	if v.pendingLen > 0 {
		copied := copy(v.pending[v.pendingLen:], part)
		prefix := v.pending[:v.pendingLen+copied]
		if !utf8.FullRune(prefix) {
			v.pendingLen += copied
			return true
		}
		_, width := utf8.DecodeRune(prefix)
		if width == 1 {
			return false
		}
		part = part[width-v.pendingLen:]
		v.pendingLen = 0
	}
	if len(part) == 0 {
		return true
	}
	tail := len(part) - 1
	for tail > 0 && len(part)-tail < utf8.UTFMax && !utf8.RuneStart(part[tail]) {
		tail--
	}
	if !utf8.FullRune(part[tail:]) {
		v.pendingLen = copy(v.pending[:], part[tail:])
		part = part[:tail]
	}
	return utf8.Valid(part)
}

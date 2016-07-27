package bitset

import (
	"errors"
)

var (
	ones     []byte
	zeros    []byte
	ones32   []uint32
	setbits  []byte
	zerobits []byte
	leftmz   []byte
	rightmz  []byte
	leftm1   []byte
	rightm1  []byte
	ErrRange = errors.New("Index out of range")
	ErrMaxR  = errors.New("Maximum bit range allowed for GetVal and SetVal is 32")
)

const (
	and uint32 = iota
	or
	xor
)

func init() {
	ones = make([]byte, 8)
	zeros = make([]byte, 8)
	setbits = make([]byte, 256)
	zerobits = make([]byte, 256)
	rightm1 = make([]byte, 256)
	leftm1 = make([]byte, 256)
	leftmz = make([]byte, 256)
	rightmz = make([]byte, 256)

	ones[0] = 1
	zeros[0] = 254
	i := 1
	for i < 8 {
		ones[i] = ones[i-1] << 1
		zeros[i] = ^ones[i]
		i++
	}
	ones32 = make([]uint32, 32)
	ones32[0] = 1
	i = 1
	for i < 32 {
		ones32[i] = (ones32[i-1] << 1) | 1
		i++
	}
	i = 0
	for i < 256 {
		setbits[i] = byte(i&1) + setbits[i>>uint32(1)]
		zerobits[i] = 8 - setbits[i]
		i++
	}
	i = 0
	j := 0
	k := 0
	for i < 256 {
		j = 0
		rightm1[i] = 8
		rightmz[i] = 8
		k = i
		for j < 8 {
			if rightmz[i] == 8 && (k&1 == 0) {
				rightmz[i] = byte(j)
			}
			if rightm1[i] == 8 && (k&1 == 1) {
				rightm1[i] = byte(j)
			}
			if rightmz[i] != 8 && rightm1[i] != 8 {
				break
			}
			j++
			k >>= 1
		}
		i++
	}
	i = 0
	j = 0
	k = 0
	for i < 256 {
		j = 7
		leftm1[i] = 8
		leftmz[i] = 8
		k = i
		for {
			if leftmz[i] == 8 && (k&128 == 0) {
				leftmz[i] = byte(j)
			}
			if leftm1[i] == 8 && (k&128 == 128) {
				leftm1[i] = byte(j)
			}
			if leftmz[i] != 8 && leftm1[i] != 8 {
				break
			}
			if j == 0 {
				break
			}
			j--
			k <<= 1
		}
		i++
	}
}

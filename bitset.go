// package bitset provides facilities to manipulate bits in a bitset.
// It is thread-safe.
package bitset

import (
	"sync"
)

type Bitset struct {
	size  uint32
	buf   []byte
	mutex *sync.RWMutex
}

// Get a new instace of Bitset with at least specified size in bytes
func NewBitset(size uint32) *Bitset {
	buf := make([]byte, size)
	return &Bitset{size: size, buf: buf, mutex: &sync.RWMutex{}}
}

// Resize expands or contracts a bitset keeping the content intact for
// the copied bytes
func (bs *Bitset) Resize(newsize uint32) {
	newbf := make([]byte, newsize)
	bs.mutex.Lock()
	defer bs.mutex.Unlock()
	copy(newbf, bs.buf)
	bs.size = newsize
	bs.buf = newbf
}

// Clone makes a copy of the current Bitset
func (bs *Bitset) Clone() *Bitset {
	mutex := &sync.RWMutex{}
	bs.mutex.Lock()
	buf := make([]byte, bs.size)
	copy(buf, bs.buf)
	bs.mutex.Unlock()
	return &Bitset{size: uint32(len(buf)), buf: buf, mutex: mutex}
}

// SetBit sets the bit at some position. It returns false if the position exceeds the size of the
// bitset, true otherwise
func (bs *Bitset) SetBit(position uint32) bool {
	bytepos := position >> 3
	bitpos := 7 - (position & 7)
	bs.mutex.Lock()
	defer bs.mutex.Unlock()
	if bytepos >= bs.size {
		return false
	}
	bs.buf[bytepos] |= ones[bitpos]
	return true
}

// ResetBit reset the bit at some position. It returns false if the position exceeds the size of the
// bitset, true otherwise
func (bs *Bitset) ResetBit(position uint32) bool {
	bytepos := position >> 3
	bitpos := 7 - (position & 7)
	bs.mutex.Lock()
	defer bs.mutex.Unlock()
	if bytepos >= bs.size {
		return false
	}
	bs.buf[bytepos] &= zeros[bitpos]
	return true
}

// IsSet returns true if the bit is set at position, false otherwise. error retuned will be
// non-nil if the position exceeds the bitset capacity, nil otherwise
func (bs *Bitset) IsSet(position uint32) (bool, error) {
	bs.mutex.RLock()
	defer bs.mutex.RUnlock()
	bytepos, bitpos, err := bs.getBitBytePosition(position)
	if err != nil {
		return false, err
	}
	b := bs.buf[bytepos] & ones[bitpos]
	return b != 0, nil
}

// GetByte returns byte that contains bit corresponding to the position.
// Non-nil error is returned in case the position is out of range
func (bs *Bitset) GetByte(position uint32) (byte, error) {
	bs.mutex.RLock()
	defer bs.mutex.RUnlock()
	bytepos, _, err := bs.getBitBytePosition(position)
	if err != nil {
		return byte(0), err
	}
	return bs.buf[bytepos], nil
}

// Getsize returns the size of the bitset in bytes
func (bs *Bitset) GetSize() uint32 {
	bs.mutex.RLock()
	defer bs.mutex.RUnlock()
	return bs.size
}

// SetVal assigns the value from fromval on the bits on the bitset, it returns nil on success.
// It returns out of range error if end exceeds size of the bitset or start - end >  31
// Also it will take the lowest (end - start + 1) from fromval
func (bs *Bitset) SetVal(start uint32, end uint32, fromval uint32) error {
	if end < start {
		start, end = end, start
	}
	numbit_to_set := end - start + 1
	if numbit_to_set >= 31 {
		return ErrRange
	}
	startbyte := start >> 3
	startbitpos := start & 7
	endbyte := end >> 3
	endbitpos := end & 7
	fromval &= ones32[numbit_to_set-1]
	if endbitpos != 7 {
		fromval <<= 7 - endbitpos
	}
	bs.mutex.Lock()
	defer bs.mutex.Unlock()
	if endbyte >= bs.size {
		return ErrRange
	}
	tmp := byte(0)
	for i := endbyte; i >= startbyte; i-- {
		cur_byte := byte(0xff & fromval)
		if i == startbyte && startbitpos != 0 {
			tmp |= ^byte(ones32[7-startbitpos])
		}
		if i == endbyte && endbitpos != 7 {
			tmp |= byte(ones32[7-endbitpos-1])
		}
		if tmp != 0 {
			bs.buf[i] = (cur_byte & ^tmp) | (bs.buf[i] & tmp)
		} else {
			bs.buf[i] = cur_byte
		}
		if i == startbyte {
			break
		}
		fromval >>= 8
		tmp = 0
	}
	return nil
}

// ClearAll sets all the bits in the bitset to zero
func (bs *Bitset) ClearAll() {
	bs.mutex.Lock()
	defer bs.mutex.Unlock()
	for i := uint32(0); i < bs.size; i++ {
		bs.buf[i] = 0
	}
}

// ClearRange clears the bits in positions start <= position <= end. It returns non-nil errors if
// any of the position passed is out of range
func (bs *Bitset) ClearRange(start uint32, end uint32) error {
	if start > end {
		start, end = end, start
	}
	startbyte := start >> 3
	startbitpos := start & 7
	endbyte := end >> 3
	endbitpos := end & 7
	bs.mutex.Lock()
	defer bs.mutex.Unlock()
	if endbyte >= bs.size {
		return ErrRange
	}
	var i uint32 = startbyte
	var andwith byte = 0
	for {
		if i == startbyte && startbitpos != 0 {
			andwith = ^(0xff >> startbitpos)
		}
		if i == endbyte && endbitpos != 7 {
			andwith |= ^(0xff << (7 - endbitpos))
		}
		bs.buf[i] &= andwith
		if i >= endbyte {
			break
		}
		i++
		andwith = 0
	}
	return nil
}

// SetAll sets all the bits in the bitset to 1
func (bs *Bitset) SetAll() {
	bs.mutex.Lock()
	defer bs.mutex.Unlock()
	for i := uint32(0); i < bs.size; i++ {
		bs.buf[i] = 255
	}
}

// SetRange sets the bits in positions start <= position <= end. It returns non-nil error if
// any of the position passed is out of range
func (bs *Bitset) SetRange(start uint32, end uint32) error {
	if start > end {
		start, end = end, start
	}
	startbyte := start >> 3
	startbitpos := start & 7
	endbyte := end >> 3
	endbitpos := end & 7
	bs.mutex.Lock()
	defer bs.mutex.Unlock()
	if endbyte >= bs.size {
		return ErrRange
	}
	var i uint32 = startbyte
	var orwith byte = 0xff
	for {
		if i == startbyte && startbitpos != 0 {
			orwith = 0xff >> startbitpos
		}
		if i == endbyte && endbitpos != 7 {
			orwith &= 0xff << (7 - endbitpos)
		}
		bs.buf[i] |= orwith
		if i >= endbyte {
			break
		}
		i++
		orwith = 0xff
	}
	return nil
}

// GetBytes returns a clone of unnderlying byte array of the Bitset
func (bs *Bitset) GetBytes() []byte {
	bs.mutex.RLock()
	defer bs.mutex.RUnlock()
	ret := make([]byte, bs.size)
	copy(ret, bs.buf)
	return ret
}

// Flip flips the bit at some position, returns non-nil error if position is out of range
func (bs *Bitset) Flip(position uint32) error {
	bs.mutex.Lock()
	defer bs.mutex.Unlock()
	bytepos, bitpos, err := bs.getBitBytePosition(position)
	if err != nil {
		return err
	}
	bs.buf[bytepos] ^= ones[bitpos]
	return nil
}

// Flip flips the bits for positions start <= position <= end, return non-nil errors if any of
// positions is out of range
func (bs *Bitset) FlipRange(start uint32, end uint32) error {
	if start > end {
		start, end = end, start
	}
	startbyte := start >> 3
	startbitpos := start & 7
	endbyte := end >> 3
	endbitpos := end & 7
	bs.mutex.Lock()
	defer bs.mutex.Unlock()
	if endbyte >= bs.size {
		return ErrRange
	}
	var i uint32 = startbyte
	var xorwith byte = 255
	for {
		if i == startbyte && startbitpos != 0 {
			xorwith >>= startbitpos
		}
		if i == endbyte && endbitpos != 7 {
			xorwith &= 0xff << (7 - endbitpos)
		}
		bs.buf[i] ^= xorwith
		if i >= endbyte {
			break
		}
		i++
		xorwith = 255
	}
	return nil
}

// IsAllZero returns true if all the bits in the set are zero
func (bs *Bitset) IsAllZero() bool {
	bs.mutex.RLock()
	defer bs.mutex.RUnlock()
	for _, buf := range bs.buf {
		if buf != 0 {
			return false
		}
	}
	return true
}

// IsAllSet returns true if all the bits in the set are 1
func (bs *Bitset) IsAllSet() bool {
	bs.mutex.RLock()
	defer bs.mutex.RUnlock()
	for _, buf := range bs.buf {
		if buf != 0xff {
			return false
		}
	}
	return true
}

// And bitwise ands the bits with the other bitset
func (bs *Bitset) And(other *Bitset) {
	bs.op(other, and)
}

// Or bitwise ors the bits with the other bitset
func (bs *Bitset) Or(other *Bitset) {
	bs.op(other, or)
}

// Xor bitwise xors the bits with the other bitset
func (bs *Bitset) Xor(other *Bitset) {
	bs.op(other, xor)
}

// op performs and, or, xor operation on two bitsets
func (bs *Bitset) op(other *Bitset, opcode uint32) {
	bs.mutex.Lock()
	defer bs.mutex.Unlock()
	other.mutex.RLock()
	defer other.mutex.RUnlock()
	var i uint32 = bs.size
	if i > other.size {
		i = other.size
	}
	var j uint32 = 0
	for j < i {
		switch opcode {
		case and:
			bs.buf[j] &= other.buf[j]
		case or:
			bs.buf[j] |= other.buf[j]
		case xor:
			bs.buf[j] ^= other.buf[j]
		}
		j++
	}
}

// GetSetbitCount returns the number of set or 1 bits in the bitset
func (bs *Bitset) GetSetbitCount() uint64 {
	return bs.getSetbitc()
}

// GetZerobitCount returns the number of 0 bits in the bitset
func (bs *Bitset) GetZerobitCount() uint64 {
	return (uint64(bs.size) << 3) - bs.getSetbitc()
}

func (bs *Bitset) getSetbitc() uint64 {
	bs.mutex.RLock()
	defer bs.mutex.RUnlock()
	var ret uint64 = 0
	var i uint32 = 0
	for i < bs.size {
		ret += uint64(setbits[bs.buf[i]])
		i++
	}
	return ret
}

// getBitBytePosition returns the corresponding byte position and bit position within the byte
// for the absolute bit position passed
func (bs *Bitset) getBitBytePosition(position uint32) (uint32, uint32, error) {
	bytepos := position >> 3
	bitpos := 7 - (position & 7)
	if bytepos >= bs.size {
		return 0, 0, ErrRange
	}
	return bytepos, bitpos, nil
}

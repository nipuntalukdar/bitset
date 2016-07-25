package bitset

import (
	"testing"
)

func TestSetReset(t *testing.T) {
	bs := NewBitset(80)
	bs.SetBit(9)
	ret, err := bs.IsSet(9)
	if err != nil {
		t.Fatalf("set bit failed %s", err)
	}
	if !ret {
		t.Fatal("set bit failed")
	}

	var expect byte = 0x40
	recvd, err := bs.GetByte(9)
	if err != nil {
		t.Fatalf("set bit failed %s", err)
	}
	if recvd != expect {
		t.Fatal("Set bit failed")
	}
	bs.ResetBit(9)
	ret, err = bs.IsSet(9)
	if err != nil {
		t.Fatalf("test failed %s", err)
	}
	if ret {
		t.Fatal("Reset bit failed")
	}
	bs.SetVal(1, 10, 511)

	expect = 63
	recvd, err = bs.GetByte(0)
	if err != nil {
		t.Fatalf("set val failed %s", err)
	}
	if recvd != expect {
		t.Fatalf("Set val failed expected: %d, found: %d", expect, recvd)
	}
	expect = 224
	recvd, err = bs.GetByte(9)
	if err != nil {
		t.Fatalf("set val failed %s", err)
	}
	if recvd != expect {
		t.Fatalf("Set val failed expected: %d, found: %d", expect, recvd)
	}
	bs.SetVal(6, 8, 0)
	expect = 60
	recvd, err = bs.GetByte(0)
	if err != nil {
		t.Fatalf("set val failed %s", err)
	}
	if recvd != expect {
		t.Fatalf("Set val failed expected: %d, found: %d", expect, recvd)
	}
	expect = 224 - 128
	recvd, err = bs.GetByte(9)
	if err != nil {
		t.Fatalf("set val failed %s", err)
	}
	if recvd != expect {
		t.Fatalf("Set val failed expected: %d, found: %d", expect, recvd)
	}

	//clearall test
	bs.ClearAll()
	bts := bs.GetBytes()
	for start := 0; start < len(bts); start++ {
		if bts[start] != 0 {
			t.Fatal("ClearAll failed")
		}
	}

	//setall test
	bs.SetAll()
	bts = bs.GetBytes()
	for start := 0; start < len(bts); start++ {
		if bts[start] != 255 {
			t.Fatal("SetAll failed")
		}
	}

	// Flip test
	err = bs.Flip(10)
	if err != nil {
		t.Fatalf("Flip test failed")
	}
	recvd, err = bs.GetByte(10)
	if err != nil || recvd != 223 {
		t.Fatalf("Flip test failed")
	}
	err = bs.Flip(10)
	if err != nil {
		t.Fatalf("Flip test failed")
	}
	recvd, err = bs.GetByte(10)
	if err != nil || recvd != 255 {
		t.Fatalf("Flip test failed")
	}
	bs.SetAll()
	setcount := bs.GetSetbitCount()
	if setcount != 80*8 {
		t.Fatal("GetSetbitCount failed")
	}
	zerocount := bs.GetZerobitCount()
	if zerocount != 0 {
		t.Fatal("GetZerobitCount failed")
	}
	bs.ClearAll()
	setcount = bs.GetSetbitCount()
	if setcount != 0 {
		t.Fatal("GetSetbitCount failed")
	}
	zerocount = bs.GetZerobitCount()
	if zerocount != 80*8 {
		t.Fatal("GetZerobitCount failed")
	}

	bs.SetAll()
	if bs.Flip(40) != nil || bs.Flip(60) != nil || bs.Flip(55) != nil {
		t.Fatalf("Flip failed")
	}
	setcount = bs.GetSetbitCount()
	if setcount != (80*8 - 3) {
		t.Fatal("GetSetbitCount failed")
	}
	zerocount = bs.GetZerobitCount()
	if zerocount != 3 {
		t.Fatal("GetZerobitCount failed")
	}
	bs.SetAll()
	bytes := bs.GetBytes()
	for _, b := range bytes {
		if b != 255 {
			t.Fatalf("Set all failed, expected value 255, got %d", b)
		}
	}
	err = bs.FlipRange(1, 640)
	if err == nil {
		t.Fatalf("FlipRange failed to detect invalid range")
	}
	err = bs.FlipRange(0, 639)
	if err != nil {
		t.Fatalf("FlipRange failed")
	}
	bytes = bs.GetBytes()
	for _, b := range bytes {
		if b != 0 {
			t.Fatalf("Flip range failed, expected value 0, got %d", b)
		}
	}
	err = bs.FlipRange(10, 630)
	if err != nil {
		t.Fatalf("FlipRange failed")
	}
	bytes = bs.GetBytes()
	expected_vals := make([]byte, 80)
	i := 0
	for i < 80 {
		if i == 0 {
			expected_vals[i] = 0
		} else if i == 1 {
			expected_vals[i] = 63
		} else if i == 78 {
			expected_vals[i] = 254
		} else if i == 79 {
			expected_vals[i] = 0
		} else {
			expected_vals[i] = 255
		}
		i++
	}
	for i, b := range bytes {
		if b != expected_vals[i] {
			t.Fatalf("FlipRange failed")
		}
	}

	bs.ClearAll()
	if err = bs.SetRange(20, 40); err != nil {
		t.Fatal("SetRange failed")
	}
	if recvd, err = bs.GetByte(20); err != nil || recvd != 15 {
		t.Fatalf("SetRange failed, expected 7, got %d", recvd)
	}
	if recvd, err = bs.GetByte(24); err != nil || recvd != 255 {
		t.Fatal("SetRange failed, expected 255, got %d", recvd)
	}
	if recvd, err = bs.GetByte(32); err != nil || recvd != 255 {
		t.Fatal("SetRange failed, expected 255, got %d", recvd)
	}
	if recvd, err = bs.GetByte(40); err != nil || recvd != 128 {
		t.Fatal("SetRange failed, expected 128, got %d", recvd)
	}
	bs.SetAll()
	if err = bs.ClearRange(20, 40); err != nil {
		t.Fatal("ClearRange failed")
	}
	if recvd, err = bs.GetByte(20); err != nil || recvd != 240 {
		t.Fatalf("ClearRange failed, expected 7, got %d", recvd)
	}
	if recvd, err = bs.GetByte(24); err != nil || recvd != 0 {
		t.Fatal("ClearRange failed, expected 255, got %d", recvd)
	}
	if recvd, err = bs.GetByte(32); err != nil || recvd != 0 {
		t.Fatal("ClearRange failed, expected 255, got %d", recvd)
	}
	if recvd, err = bs.GetByte(40); err != nil || recvd != 127 {
		t.Fatal("ClearRange failed, expected 128, got %d", recvd)
	}

	if bs.IsAllSet() || bs.IsAllZero() {
		t.Fatal("IsAllSet or IsAllZero failed")
	}
	bs.ClearAll()
	if bs.IsAllSet() || !bs.IsAllZero() {
		t.Fatal("IsAllSet or IsAllZero failed")
	}
	bs.SetAll()
	if !bs.IsAllSet() || bs.IsAllZero() {
		t.Fatal("IsAllSet or IsAllZero failed")
	}

	bs.ClearAll()
	other := NewBitset(2)
	other.SetAll()
	bs.Or(other)

	bytes = bs.GetBytes()
	if bytes[0] != 0xff || bytes[1] != 0xff {
		t.Fatal("Or test failed")
	}
	bytes[0] = 0
	bytes[1] = 0
	for _, b := range bytes {
		if b != 0 {
			t.Fatal("Or test failed")
		}
	}

	bs.SetAll()
	bs.Xor(other)
	bytes = bs.GetBytes()
	if bytes[0] != 0 || bytes[1] != 0 {
		t.Fatal("Xor test failed")
	}
	bytes[0] = 0xff
	bytes[1] = 0xff
	for _, b := range bytes {
		if b != 0xff {
			t.Fatal("Xor test failed")
		}
	}

	bs.ClearRange(0, 15)
	bs.And(other)
	bytes = bs.GetBytes()
	if bytes[0] != 0 || bytes[1] != 0 {
		t.Fatal("And test failed")
	}
	bytes[0] = 0xff
	bytes[1] = 0xff
	for _, b := range bytes {
		if b != 0xff {
			t.Fatal("And test failed")
		}
	}
}

func TestBitSetClone(t *testing.T) {
	bs1 := NewBitset(100)
	bs1.SetAll()
	bs2 := bs1.Clone()
	if bs1.GetSize() != bs2.GetSize() {
		t.Fatal("Clone failed")
	}
	if bs2.GetSetbitCount() != 100*8 {
		t.Fatal("Clone failed")
	}
	if bs1.GetSetbitCount() != bs2.GetSetbitCount() {
		t.Fatal("Clone failed")
	}
}

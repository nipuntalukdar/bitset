Bitset
=========

**A bitset library for Go**

-----------------------
Bitset is a library useful for manipulation of bits in a set of bits. It provides operations like,
Set, Reset, Flip, Clear, XOR, AND, OR. Also, the bitset can be resized. It is thread-safe or rather
go-routine safe :) hence can be used concurrently from multiple goroutines.

Godoc for this library is available [here](https://godoc.org/github.com/nipuntalukdar/bitset)

**Below is an example regarding how to use the library**

---
```go
package main

import (
    "fmt"
    "github.com/nipuntalukdar/bitset"
)

func main() {
    bs := bitset.NewBitset(80)
    if bs.IsAllZero() {
        fmt.Printf("All bits are set to zero\n")
    }
    bs.Flip(60)
    fmt.Printf("Number of ones %d\n", bs.GetSetbitCount())
    fmt.Printf("Number of zeros %d\n", bs.GetZerobitCount())
    bs.SetBit(9)
    ret, err := bs.IsSet(9)
    if err != nil || !ret {
        fmt.Printf("Some issue\n")
    }

    bs.SetVal(1, 10, 511)
    recvd, err := bs.GetByte(0)
    if err == nil {
        fmt.Printf("First byte in the underlying array %d\n", recvd)
    }

    // set bytes 6,7,8 to zero and don't touch other bits
    bs.SetVal(6, 8, 0)

    // Clear all the bits, or make them all zero
    bs.ClearAll()

    // Set all the bits to 1
    bs.SetAll()

    // Flip 10th bit
    bs.Flip(10)
    recvd, _ = bs.GetByte(10)
    fmt.Printf("The 10 th byte %d\n", recvd)

    // Get a copy of underlying bytes of the bitset
    bytes := bs.GetBytes()
    fmt.Printf("The first four bytes in the bitset %v\n", bytes[0:4])

    // Flip all the bits
    bs.FlipRange(1, 639)

    // Set the bits 20,21,....,39,40 to 1
    bs.SetRange(20, 40)

    // Set the bits from 20,21,22,....,79,80 to zero
    bs.ClearRange(20, 80)

    other1 := bitset.NewBitset(4)
    other1.SetAll()
    bs.Xor(other1)

    other2 := bitset.NewBitset(10)
    other2.SetRange(40, 60)
    bs.And(other2)
    bs.ResetBit(80)

    // Get next zero bit position after index 0
    indx, _ := bs.GetNextZeroBit(0)
    fmt.Printf("Index of next zero bit %d\n", indx)

    // Get the zero bit before bit 80
    indx, _ = bs.GetPrevSetBit(80)

    // Get the bits from 10 to 28 packed in an uint32
    retuint32, _ := bs.GetVal(10, 28)
    fmt.Printf("The packed uint32 returned is %d\n", retuint32)
}
```

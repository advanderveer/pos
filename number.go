package pos

import (
	"fmt"
	"math/big"
	"strconv"
)

// Num describes arbitrary large number with an implied domain
type Num struct {
	*big.Int
	z uint
}

// NewNum inits a num with implied domain of 'z' bits
func NewNum(n *big.Int, z uint) (num Num) {
	if z == 0 {
		panic("pos: number with implied domain '0' not allowed")
	}
	num = Num{n, z}
	num.checkOverflow()
	return
}

// Num64 creates a number with a implied domain of 'z' bits from 'n'
func Num64(n uint64, z uint) (num Num) {
	return NewNum(new(big.Int).SetUint64(n), z)
}

// checkOverflow will assert the num itself and panic if the implied domain is too small to represent
// the number this num holds.
func (num Num) checkOverflow() {
	if num.BitLen() <= int(num.z) {
		return
	}
	panic("pos: number too large for implied domain")
}

// Format is called by stringer
func (num Num) Format(f fmt.State, verb rune) {
	if num.BitLen() <= int(num.z) {
		fmt.Fprintf(f, "%03db%0"+strconv.Itoa(int(num.z))+"b", num.z, num.Int)
		return
	}
	fmt.Fprintf(f, "errb%0"+strconv.Itoa(int(num.z))+"b", num.Int)
}

// Domain returns the implied domain of this number
func (num Num) Domain() uint { return num.z }

// Uint64 returns the num as a uint64 but panics if the implied domain is larger then 64 bits
func (num Num) Uint64() uint64 {
	if num.Domain() > 64 {
		panic("pos: number's implied domain to large for uint64, got: " + strconv.Itoa(int(num.Domain())))
	}

	return num.Int.Uint64()
}

// ToBlake bytes implements the perculiar conversion to bytes before the number is fed to the
// blake hasher. It reads all the bits in the implied domain and sets them on 'b' in reverse order.
func (num Num) ToBlakeBytes(b []byte) {
	for i, j := int(num.z)-1, 0; i >= 0; i-- {
		cell, pos := j/8, j%8
		b[cell] |= byte(num.Bit(i) << (7 - pos))
		j++
	}
}

package pos

import (
	"math/big"
)

// F1 performs a ChaCha8 cipher of the provided x value
func F1(params *Params, x Num) Num {
	ciphered := readChaCha(params, x)                      // ChaCha8(...)
	return Concat(ciphered, Slice(x, 0, uint(params.ext))) // || x[:param_ext]
}

// readChaCha peforms the chacha byte reading
func readChaCha(params *Params, x Num) Num {
	q, r := divmod(x.Uint64()*uint64(params.k), 512)
	ciphertext0, end := ChaCha8(q, params.pseed), r+uint64(params.k)
	if end < 512 { // the bytes we need to read can be found in the first round of chacha8
		return Slice(NewNum(new(big.Int).SetBytes(ciphertext0[:]), 512), uint(r), uint(end))
	}

	// else, append extra bytes of ciphertext, and slice that instead
	ciphertext1 := ChaCha8(q+1, params.pseed)
	comb := new(big.Int).SetBytes(append(ciphertext0[:], ciphertext1[:]...))
	return Slice(NewNum(comb, 1024), uint(r), uint(end))
}

// divmod as taken from: https://stackoverflow.com/questions/43945675/division-with-returning-quotient-and-remainder
// and tested again python implementation https://www.programiz.com/python-programming/methods/built-in/divmod
func divmod(x, m uint64) (quo, rem uint64) {
	quo = x / m
	rem = x % m
	return
}

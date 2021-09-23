package pos

import "math/big"

// Concat implements zero-padding concatenation for domain z. NOTE: i'm assuming that the zero-padding only
// applies to all elements except the first, where extra zeros are always ignored
func Concat(z uint, xs ...*big.Int) (res *big.Int) {
	res = new(big.Int)
	for _, x := range xs {
		res.Lsh(res, uint(z)).Add(x, res)
	}
	return res
}

// Concat64 takes 64bit x-values as argument
func Concat64(z uint, xs ...uint64) *big.Int {
	bxs := make([]*big.Int, len(xs))
	for i, x := range xs {
		bxs[i] = new(big.Int).SetUint64(x)
	}
	return Concat(z, bxs...)
}

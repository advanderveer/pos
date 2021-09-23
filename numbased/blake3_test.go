package pos

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBlake3(t *testing.T) {
	s1 := Blake3([]byte{0x01})
	require.Equal(t, byte(0x48), s1[0])
	require.Equal(t, byte(0xfc), s1[1])
}

func TestA(t *testing.T) {
	require.Equal(t, "256b0000000110010110110111101000111110100111011011010110100010110010110111100111010100101011101001001111101111100100111110100110011001111100001111110000101101100111111001111011011001110111111100111001000111101101111111011000011111110100010101100010010111111101",
		fmt.Sprint(A(Num64(1, 25), Num64(2, 4), Num64(1, 5))))
}

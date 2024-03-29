package pos

import (
	"bufio"
	"embed"
	"math/big"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProofFromBytes(t *testing.T) {
	proof, err := ProofFromHex("228f532336a70179e3a96fef5d43cfc7753a527e876cfe328d7a169b4632bf95c62863df453c2d36e6f49a6967e7d58a57249a02c36638676117a73ca0db52c12a118e359346115a75ca5c454a67f8a3de32832801d33dab42246890142e247237f77dfae81c108cd1e01d9e195a9d4cee6491abf509acb301cc00b9bd2dab5a18aa6c07ee3583afd0b24937077557eb52797161b25ba308a440fbd4d35365d08d56d58d74028355ba33a44bef583f1af1801f995d32f4b228002d93c79a7555c87cdb00d7d11670", 25)
	require.NoError(t, err)
	require.Equal(t, uint64(30479984), proof[63].Uint64())

	// x1, x2 double checked with reference implementation
	require.Equal(t, uint64(0b0010001010001111010100110), proof[0].Uint64())
	require.Equal(t, uint64(0b0100011001101101010011100), proof[1].Uint64())

	// x3, x4 double check with reference
	require.Equal(t, uint64(0b0000010111100111100011101), proof[2].Uint64())
	require.Equal(t, uint64(0b0100101101111111011110101), proof[3].Uint64())
}

func TestVerify(t *testing.T) {
	params := NewParams(25, [32]byte{0x01, 0x02, 0x2f, 0xb4, 0x2c, 0x08, 0xc1, 0x2d, 0xe3, 0xa6, 0xaf, 0x05, 0x38, 0x80, 0x19, 0x98, 0x06, 0x53, 0x2e, 0x79, 0x51, 0x5f, 0x94, 0xe8, 0x34, 0x61, 0x61, 0x21, 0x01, 0xf9, 0x41, 0x2f})
	proof, err := ProofFromHex("228f532336a70179e3a96fef5d43cfc7753a527e876cfe328d7a169b4632bf95c62863df453c2d36e6f49a6967e7d58a57249a02c36638676117a73ca0db52c12a118e359346115a75ca5c454a67f8a3de32832801d33dab42246890142e247237f77dfae81c108cd1e01d9e195a9d4cee6491abf509acb301cc00b9bd2dab5a18aa6c07ee3583afd0b24937077557eb52797161b25ba308a440fbd4d35365d08d56d58d74028355ba33a44bef583f1af1801f995d32f4b228002d93c79a7555c87cdb00d7d11670", 25)
	require.NoError(t, err)
	chall, err := ChallengeFromHex("1000000000000000000000000000000000000000000000000000000000000000")
	require.NoError(t, err)

	require.True(t, Verify(params, proof, chall))
}

func TestFailVerify(t *testing.T) {
	params := NewParams(25, [32]byte{0x01, 0x02, 0x2f, 0xb4, 0x2c, 0x08, 0xc1, 0x2d, 0xe3, 0xa6, 0xaf, 0x05, 0x38, 0x80, 0x19, 0x98, 0x06, 0x53, 0x2e, 0x79, 0x51, 0x5f, 0x94, 0xe8, 0x34, 0x61, 0x61, 0x21, 0x01, 0xf9, 0x41, 0x2f})
	proof, err := ProofFromHex("228f532336a70179e3a96fef5d43cfc7753a527e876cfe328d7a169b4632bf95c62863df453c2d36e6f49a6967e7d58a57249a02c36638676117a73ca0db52c12a118e359346115a75ca5c454a67f8a3de32832801d33dab42246890142e247237f77dfae81c108cd1e01d9e195a9d4cee6491abf509acb301cc00b9bd2dab5a18aa6c07ee3583afd0b24937077557eb52797161b25ba308a440fbd4d35365d08d56d58d74028355ba33a44bef583f1af1801f995d32f4b228002d93c79a7555c87cdb00d7d11670", 25)
	require.NoError(t, err)
	chall, err := ChallengeFromHex("1000000000000000000000000000000000000000000000000000000000000000")
	require.NoError(t, err)

	proof[5].Int = proof[5].Add(proof[5].Int, big.NewInt(1)) // single change should invalidate the proof

	require.False(t, Verify(params, proof, chall))
}

func TestFailVerifyChallenge(t *testing.T) {
	params := NewParams(25, [32]byte{0x01, 0x02, 0x2f, 0xb4, 0x2c, 0x08, 0xc1, 0x2d, 0xe3, 0xa6, 0xaf, 0x05, 0x38, 0x80, 0x19, 0x98, 0x06, 0x53, 0x2e, 0x79, 0x51, 0x5f, 0x94, 0xe8, 0x34, 0x61, 0x61, 0x21, 0x01, 0xf9, 0x41, 0x2f})
	proof, err := ProofFromHex("228f532336a70179e3a96fef5d43cfc7753a527e876cfe328d7a169b4632bf95c62863df453c2d36e6f49a6967e7d58a57249a02c36638676117a73ca0db52c12a118e359346115a75ca5c454a67f8a3de32832801d33dab42246890142e247237f77dfae81c108cd1e01d9e195a9d4cee6491abf509acb301cc00b9bd2dab5a18aa6c07ee3583afd0b24937077557eb52797161b25ba308a440fbd4d35365d08d56d58d74028355ba33a44bef583f1af1801f995d32f4b228002d93c79a7555c87cdb00d7d11670", 25)
	require.NoError(t, err)
	chall, err := ChallengeFromHex("0000000000000000000000000000000000000000000000000000000000000000")
	require.NoError(t, err)

	require.False(t, Verify(params, proof, chall))
}

//go:embed test_proofs.txt
var testProofs embed.FS

func TestBatchProofValidation(t *testing.T) {
	params := NewParams(25, [32]byte{0x01, 0x02, 0x2f, 0xb4, 0x2c, 0x08, 0xc1, 0x2d, 0xe3, 0xa6, 0xaf, 0x05, 0x38, 0x80, 0x19, 0x98, 0x06, 0x53, 0x2e, 0x79, 0x51, 0x5f, 0x94, 0xe8, 0x34, 0x61, 0x61, 0x21, 0x01, 0xf9, 0x41, 0x2f})

	f, err := testProofs.Open("test_proofs.txt")
	require.NoError(t, err)

	scan := bufio.NewScanner(f)
	for scan.Scan() {
		parts := strings.SplitN(scan.Text(), ":", 2)

		proof, err := ProofFromHex(parts[1], params.k)
		require.NoError(t, err)
		chall, err := ChallengeFromHex(parts[0])
		require.NoError(t, err)

		require.True(t, Verify(params, proof, chall))
	}
}

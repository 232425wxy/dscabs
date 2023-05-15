package chaincode

import (
	"fmt"
	"testing"

	"github.com/232425wxy/dscabs/algorithm"
	"github.com/232425wxy/dscabs/ecdsa/bigint"
	"github.com/stretchr/testify/assert"
)

func TestProcess(t *testing.T) {
	params := algorithm.Setup(256)
	fmt.Println(params.Curve.Params().N)

	ak := algorithm.ExtractAK(params, []string{"a", "b", "c", "d"})
	fmt.Println(ak.SecretKey)

	pk := algorithm.GenPK(params, "{a,b,c,d,[4,4]}")

	sig := algorithm.Sign(params, []byte("hello"), ak.SecretKey)

	ok := algorithm.Verify(params, ak.PublicKey, pk, []byte("hello"), sig)

	assert.Equal(t, ok, true)
}

func TestSign(t *testing.T) {
	params := algorithm.Setup(256)
	fmt.Println(params.Curve.Params().N)
	ori := "80848746029104293953772004333084988230793558049441846281290986991604797601500"
	sk, _ := new(bigint.BigInt).SetString(ori, 10)

	sig := algorithm.Sign(params, []byte("DSCABS.FuncXXX1"), sk)

	fmt.Println(sig.S, sig.R)

}
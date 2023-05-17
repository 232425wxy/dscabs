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
	ori := "34454378694695527068725690557774082674590862555097716080739531125196898673466"
	sk, _ := new(bigint.BigInt).SetString(ori, 10)

	sig := algorithm.Sign(params, []byte("20:tom:DogContract:GetDog"), sk)

	fmt.Println(sig.S, sig.R)

}
package algorithm

import (
	"crypto/elliptic"

	"github.com/232425wxy/dscabs/ecdsa"
	"github.com/232425wxy/dscabs/ecdsa/bigint"
)

func Setup(securityLevel int) *SystemParams {
	params := &SystemParams{}
	switch securityLevel {
	case 224:
		params.Curve = elliptic.P224()
	case 256:
		params.Curve = elliptic.P256()
	case 384:
		params.Curve = elliptic.P384()
	default:
		params.Curve = elliptic.P256()
	}

	params.MSK = ecdsa.RandNumOnCurve(params.Curve)

	ecdsa.Bottom.X = new(bigint.BigInt).SetInt64(-1)
	ecdsa.Bottom.Y = new(bigint.BigInt).SetInt64(-1)

	return params
}

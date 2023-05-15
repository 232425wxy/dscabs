package algorithm

import (
	"crypto/elliptic"

	"github.com/232425wxy/dscabs/ecdsa"
	"github.com/232425wxy/dscabs/ecdsa/bigint"
)

type SystemParams struct {
	MSK   *bigint.BigInt `json:"msk"`
	Curve elliptic.Curve `json:"curve"`
}

type AttributeKey struct {
	SecretKey  *bigint.BigInt                       `json:"secret-key"`
	PublicKey  map[string]*ecdsa.EllipticCurvePoint `json:"public-key"`
	Attributes []string                             `json:"attributes"`
}

type Polynomial struct {
	Coefficients map[int]*bigint.BigInt
	Curve        elliptic.Curve
}

type AccessTreeNode struct {
	Policy      string
	Index, N, T int
	attribute   *attribute
	parent      *AccessTreeNode
	Children    map[string]*AccessTreeNode
	polynomial  *Polynomial
}

type Key struct {
	HashVal     string                    `json:"hash-val"`
	Index       int                       `json:"index"`
	N           int                       `json:"n"`
	T           int                       `json:"t"`
	Du          *bigint.BigInt            `json:"du"`
	Parent      *Key                      `json:"parent"`
	Children    map[string]*Key           `json:"children"`
	UseToVerify *ecdsa.EllipticCurvePoint `json:"use-to-verify"`
}

type attribute struct {
	value   string
	hashVal string
	x       *bigint.BigInt
	y       *ecdsa.EllipticCurvePoint
}

type track struct {
	m map[*Key][]struct {
		key   *Key
		point *ecdsa.EllipticCurvePoint
	}
}

var universe map[string]*attribute

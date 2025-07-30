package sdk

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

type Signature struct {
	R hexutil.Bytes `json:"r"`
	S hexutil.Bytes `json:"s"`
	V uint8         `json:"v"`
}

func NewSignature(sig []byte) (*Signature, error) {
	if len(sig) != 65 {
		return nil, fmt.Errorf("signature length is not 65 bytes, instead %d", len(sig))
	}

	r := make([]byte, 32)
	s := make([]byte, 32)
	copy(r, sig[:32])
	copy(s, sig[32:64])

	v := sig[64]
	if v < 27 {
		v = v + 27
	}

	return &Signature{
		R: r,
		S: s,
		V: v,
	}, nil
}

func (sig *Signature) Encode() ([]byte, error) {
	if sig.V < 27 {
		return nil, fmt.Errorf("invalid V value: %d, expected at least 27", sig.V)
	}
	result := make([]byte, 65)
	copy(result[:32], sig.R)
	copy(result[32:64], sig.S)
	result[64] = sig.V - 27
	return result, nil
}

type Signer interface {
	Address() common.Address
	Sign(msg []byte) ([]byte, error)
}

func SignInner(signer Signer, msg apitypes.TypedData) (*Signature, error) {
	bytes, _, err := apitypes.TypedDataAndHash(msg)
	if err != nil {
		return nil, fmt.Errorf("SignInner hash on data: %v, err: %#v", msg, err)
	}
	sig, err := signer.Sign(bytes)
	if err != nil {
		return nil, fmt.Errorf("SignInner sign err on data: %v, err: %#v", msg, err)
	}
	return NewSignature(sig)
}

type LocalSigner struct {
	key *ecdsa.PrivateKey
}

func NewLocalSignerFromHex(hexKey string) (*LocalSigner, error) {
	key, err := crypto.HexToECDSA(hexKey)
	if err != nil {
		return nil, err
	}
	return &LocalSigner{key: key}, nil
}

func (signer *LocalSigner) Address() common.Address {
	return crypto.PubkeyToAddress(signer.key.PublicKey)
}

func (signer *LocalSigner) Sign(msg []byte) ([]byte, error) {
	sig, err := crypto.Sign(msg, signer.key)
	if err != nil {
		return nil, err
	}
	return sig, nil
}

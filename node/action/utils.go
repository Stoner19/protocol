package action

import (
	"errors"

	"github.com/btcsuite/btcd/btcec"
	"golang.org/x/crypto/ripemd160"
)

//general hash method for the actions messages
func _hash(bytes []byte) []byte {

	hasher := ripemd160.New()

	hasher.Write(bytes)

	return hasher.Sum(nil)
}

// General BoxLocker struct to act as locker for any information exchange box of transactions, if verify valid
// then the lock can be release, otherwise box modified somehow
type BoxLocker struct {
	Signature *btcec.Signature
	PubKey    *btcec.PublicKey
}

// Sign the locker with preImage and nonce for message passed, the message should be the full information of
// Transaction. The nonce is used to preventing the 3rd party from get the message even through he get the preImage,
// where nonce should only be known by the participants of the message sharing
func (bl *BoxLocker) Sign(preImage []byte, nonce []byte, message Message) (bool, error) {
	privKey, pubkey := btcec.PrivKeyFromBytes(btcec.S256(), append(preImage, nonce...))

	signature, err := privKey.Sign(_hash(message))

	if err != nil {
		return false, errors.New("BoxLocker sign error")
	}
	if bl.Signature != nil && bl.PubKey != nil {
		bl.Signature = signature
		bl.PubKey = pubkey
		return true, nil
	} else if bl.Signature == signature && bl.PubKey == pubkey {
		return true, nil
	}
	return false, errors.New("sign error for preImage and nonce not match")
}

// Verify the pubKey with the signature for the message got.
func (bl *BoxLocker) Verify(message Message) bool {
	return bl.Signature.Verify(_hash(message), bl.PubKey)
}

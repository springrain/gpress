// Copyright (c) 2023 gpress Authors.
//
// This file is part of gpress.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"
)

func TestRecoverPublicKey(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	msg := "hello"
	hash := keccak256Hash([]byte(msg))

	// 签名（强制添加恢复ID）
	r, s, _ := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	// 恢复公钥
	v := new(big.Int).Mod(privateKey.PublicKey.Y, big.NewInt(2)) // 奇偶性
	publicKey, err := recoverPublicKey(hash[:], r, s, uint(v.Int64()))
	if err != nil {
		panic(err)
	}

	fmt.Println(hex.EncodeToString(publicKey.X.Bytes()))
	fmt.Println(hex.EncodeToString(privateKey.PublicKey.X.Bytes()))

	// 比较恢复的公钥与原始公钥
	if publicKey.X.Cmp(privateKey.PublicKey.X) != 0 || publicKey.Y.Cmp(privateKey.PublicKey.Y) != 0 {
		t.Error("恢复失败")
	} else {
		fmt.Println("恢复成功")
	}

}
func TestVerifySecp256r1Signature(t *testing.T) {
	// Example message to sign
	msg := "123"

	// Generate a key pair
	privateKey, err := GenerateKeyPair()
	if err != nil {
		fmt.Println("Error generating key pair:", err)
		return
	}

	prefix := fmt.Sprintf("\x86XuperChain Signed Message:\n%d%s", len(msg), msg)
	message := keccak256Hash([]byte(prefix))

	// Sign the message
	signature, err := SignMessage(privateKey, message)
	if err != nil {
		fmt.Println("Error signing message:", err)
		return
	}

	// Output the message and signature
	fmt.Println("Signature:", signature)

	ok, _, err := verifySecp256r1Signature(msg, signature)
	fmt.Println(ok)
	if err != nil {
		t.Error(err)
	}

}
func TestVerifySecp256k1Signature(t *testing.T) {
	address := "0xbe153AE90F5f114EF48A0e4279c565Be726302F6"
	sign := "0x812a04f34f988692682412010dee232f7b09e4ce96a6a3a4c5a37373db008312213f882c2248cbfbdf16b75ec595aeb75c4f7fd743e5b061bcdac1cd6e1e64931b"
	msg := "123"

	ok, err := verifySecp256k1Signature(address, msg, sign)
	fmt.Println(ok)
	if err != nil {
		t.Error(err)
	}

}

// SignMessage signs a given message using the private key and returns the signature as a base64 string.
func SignMessage(privateKey *ecdsa.PrivateKey, hash []byte) (string, error) {

	// Sign the hash
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash)
	if err != nil {
		return "", fmt.Errorf("failed to sign message: %v", err)
	}

	// Ensure s is in the lower half of the curve order to prevent malleability
	curveOrder := privateKey.Curve.Params().N
	halfCurveOrder := new(big.Int).Rsh(curveOrder, 1)
	if s.Cmp(halfCurveOrder) > 0 {
		s.Sub(curveOrder, s)
	}

	// Determine v (recovery identifier) based on the parity of the y-coordinate
	isOdd := privateKey.PublicKey.Y.Bit(0) == 1
	v := byte(27) // Ethereum uses 27 or 28 for the recovery id
	if isOdd {
		v = 28
	}

	// Prepare signature data: r (32 bytes) + s (32 bytes) + v (1 byte)
	rBytes := r.Bytes()
	sBytes := s.Bytes()

	fmt.Println("SignMessage-hash:", hex.EncodeToString(hash))
	fmt.Println("SignMessage-r:", hex.EncodeToString(r.Bytes()))
	fmt.Println("SignMessage-s:", hex.EncodeToString(s.Bytes()))
	fmt.Println("SignMessage-x:", hex.EncodeToString(privateKey.PublicKey.X.Bytes()))
	fmt.Println("SignMessage-y:", hex.EncodeToString(privateKey.PublicKey.Y.Bytes()))

	// Ensure r and s are padded to 32 bytes
	rPadded := make([]byte, 32)
	sPadded := make([]byte, 32)
	copy(rPadded[32-len(rBytes):], rBytes)
	copy(sPadded[32-len(sBytes):], sBytes)

	// Concatenate r, s, and v into a single byte slice
	signature := append(rPadded, sPadded...)
	signature = append(signature, v)

	// Encode the signature to base64 for output
	return hex.EncodeToString(signature), nil
}

// GenerateKeyPair generates an ECDSA key pair using the Secp256r1 curve.
func GenerateKeyPair() (*ecdsa.PrivateKey, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate key pair: %v", err)
	}
	return privateKey, nil
}

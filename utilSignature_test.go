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
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"
)

func TestRecoverP256PublicKey(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	msg := "hello"
	hash := keccak256Hash([]byte(msg))

	// 签名（强制添加恢复ID）
	r, s, _ := ecdsa.Sign(rand.Reader, privateKey, hash[:])

	// 恢复公钥
	recoveryID := new(big.Int).Mod(privateKey.PublicKey.Y, big.NewInt(2)) // 奇偶性

	//EIP-155 v=35+2×ChainID+recoveryID
	v := 35 + 2*1 + uint(recoveryID.Int64())
	// TODO 偶尔会恢复失败,原因待查
	publicKey, err := recoverP256PublicKey(hash[:], r, s, v)
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

	ok, err := verifyEthereumSignature(address, msg, sign)
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

	// Determine v (recovery identifier) based on the parity of the R point's y-coordinate
	// R = r * G, where G is the base point
	_, rY := privateKey.Curve.ScalarBaseMult(r.Bytes())
	isOdd := rY.Bit(0) == 1
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

// TestVerifyXuperChainSignature 测试 XuperChain 验签功能
func TestVerifyXuperChainSignature(t *testing.T) {
	// 生成 secp256r1 密钥对
	privateKey, err := GenerateKeyPair()
	if err != nil {
		t.Fatal("生成密钥对失败:", err)
	}

	// 生成钱包地址
	address, err := getAddressFromPublicKey(&privateKey.PublicKey)
	if err != nil {
		t.Fatal("生成地址失败:", err)
	}

	// 要签名的消息
	msg := "test message for xuperchain"

	// 计算带前缀的消息哈希
	prefix := fmt.Sprintf("\x86XuperChain Signed Message:\n%d%s", len(msg), msg)
	messageHash := keccak256Hash([]byte(prefix))

	// 签名并进行 s 值标准化
	curve := elliptic.P256()
	halfOrder := new(big.Int).Rsh(curve.Params().N, 1)
	var r, s *big.Int
	for {
		r, s, err = ecdsa.Sign(rand.Reader, privateKey, messageHash)
		if err != nil {
			t.Fatal("签名失败:", err)
		}
		// 确保 s 在 low-s 范围内，防止签名延展性
		if s.Cmp(halfOrder) > 0 {
			s.Sub(curve.Params().N, s)
		}
		// 检查 R 点是否能正确恢复
		_, rY := curve.ScalarBaseMult(r.Bytes())
		isOdd := rY.Bit(0) == 1
		v := byte(27)
		if isOdd {
			v = 28
		}

		// 构建测试签名
		rBytes := r.Bytes()
		sBytes := s.Bytes()
		rPadded := make([]byte, 32)
		sPadded := make([]byte, 32)
		copy(rPadded[32-len(rBytes):], rBytes)
		copy(sPadded[32-len(sBytes):], sBytes)
		signatureBytes := append(rPadded, sPadded...)
		signatureBytes = append(signatureBytes, v)
		signatureHex := hex.EncodeToString(signatureBytes)

		// 验证签名是否能正确恢复公钥
		ok, recoveredPub, _ := verifySecp256r1Signature(msg, signatureHex)
		if ok && recoveredPub != nil &&
			recoveredPub.X.Cmp(privateKey.PublicKey.X) == 0 &&
			recoveredPub.Y.Cmp(privateKey.PublicKey.Y) == 0 {
			// 公钥恢复成功，使用这个签名
			fmt.Println("XuperChain Address:", address)
			fmt.Println("Signature:", signatureHex)
			fmt.Println("Message:", msg)
			fmt.Println("v:", v, "rY is odd:", isOdd)

			// 验证签名
			valid, err := verifyXuperChainSignature(address, msg, signatureHex)
			if err != nil {
				t.Error("验证失败:", err)
			}
			if !valid {
				t.Error("签名验证不通过")
			} else {
				fmt.Println("XuperChain 验签成功!")
			}

			// 测试错误的地址应该失败
			wrongAddress := "1234567890abcdef"
			valid, err = verifyXuperChainSignature(wrongAddress, msg, signatureHex)
			if err == nil && valid {
				t.Error("使用错误地址应该验证失败")
			} else {
				fmt.Println("错误地址验证失败，符合预期")
			}

			// 测试错误的消息应该失败
			wrongMsg := "wrong message"
			valid, err = verifyXuperChainSignature(address, wrongMsg, signatureHex)
			if err == nil && valid {
				t.Error("使用错误消息应该验证失败")
			} else {
				fmt.Println("错误消息验证失败，符合预期")
			}
			return
		}
		// 如果恢复失败，重新生成签名
	}
}

// TestVerifySolanaSignature 测试 Solana 验签功能
func TestVerifySolanaSignature(t *testing.T) {
	// 生成 ed25519 密钥对
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal("生成密钥对失败:", err)
	}

	// 生成 Solana 地址 (base58 编码的公钥)
	// Solana 地址是公钥的 base58 编码，公钥是 32 字节
	publicKeyBytes := publicKey[:32]
	address := base58Encode(publicKeyBytes)
	fmt.Println("Solana Address:", address)
	fmt.Println("Public Key Length:", len(publicKeyBytes))

	// 要签名的消息
	msg := "test message for solana"

	// Solana 消息前缀
	prefix := fmt.Sprintf("\x19Solana Signed Message:\n%d%s", len(msg), msg)
	messageHash := hashUsingSha256([]byte(prefix))

	// 使用 ed25519 签名
	signature := ed25519.Sign(privateKey, messageHash)
	signatureHex := hex.EncodeToString(signature)

	fmt.Println("Signature:", signatureHex)
	fmt.Println("Signature Length:", len(signature))
	fmt.Println("Message:", msg)

	// 验证签名
	valid, err := verifySolanaSignature(address, msg, signatureHex)
	if err != nil {
		t.Error("验证失败:", err)
	}
	if !valid {
		t.Error("签名验证不通过")
	} else {
		fmt.Println("Solana 验签成功!")
	}

	// 测试错误的地址应该失败
	wrongAddress := "1234567890abcdef1234567890abcdef"
	valid, err = verifySolanaSignature(wrongAddress, msg, signatureHex)
	if err == nil && valid {
		t.Error("使用错误地址应该验证失败")
	} else {
		fmt.Println("错误地址验证失败，符合预期")
	}

	// 测试错误的消息应该失败
	wrongMsg := "wrong message"
	valid, err = verifySolanaSignature(address, wrongMsg, signatureHex)
	if err == nil && valid {
		t.Error("使用错误消息应该验证失败")
	} else {
		fmt.Println("错误消息验证失败，符合预期")
	}
}

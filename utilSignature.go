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
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"

	dcrdEcdsa "github.com/decred/dcrd/dcrec/secp256k1/v4/ecdsa"
	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"
)

// verifySecp256k1Signature 验证secp256k1的签名
func verifySecp256k1Signature(chainAddress string, msg string, signature string) (bool, error) {
	signatureBytes, err := fromHex(signature)
	if err != nil {
		return false, err
	}
	if len(signatureBytes) < 65 {
		return false, errors.New("invalid signature")
	}
	// 计算消息的哈希,包括 MetaMask 的消息前缀
	prefix := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(msg), msg)
	messageHash := keccak256Hash([]byte(prefix))
	r, s, v := signatureBytes[:32], signatureBytes[32:64], signatureBytes[64]
	sign, err := hex.DecodeString(fmt.Sprintf("%x%x%x", v, r, s))
	if err != nil {
		return false, err
	}
	//fmt.Println("verifySecp256k1Signature-r:", hex.EncodeToString(r))
	//fmt.Println("verifySecp256k1Signature-s:", hex.EncodeToString(s))
	dcrdPublicKey, _, err := dcrdEcdsa.RecoverCompact(sign, messageHash)
	if err != nil {
		return false, err
	}
	//fmt.Println("verifySecp256k1Signature-x:", hex.EncodeToString(dcrdPublicKey.X().Bytes()))
	//fmt.Println("verifySecp256k1Signature-y:", hex.EncodeToString(dcrdPublicKey.Y().Bytes()))
	pubKeyBytes := dcrdPublicKey.SerializeUncompressed()[1:]
	addressHash := keccak256Hash(pubKeyBytes)
	address := ""
	if len(addressHash) > 20 {
		address = fmt.Sprintf("0x%x", addressHash[len(addressHash)-20:])
	}
	return strings.EqualFold(address, chainAddress), nil
}

// verifySecp256r1Signature XuperChain使用NIST标准的公钥,验证签名
func verifySecp256r1Signature(chainAddress string, msg string, signature string) (valid bool, err error) {
	signatureBytes, err := fromHex(signature)
	if err != nil {
		return false, err
	}
	if len(signatureBytes) < 65 {
		return false, errors.New("invalid signature")
	}
	// 计算消息的哈希,包括消息前缀
	prefix := fmt.Sprintf("\x86XuperChain Signed Message:\n%d%s", len(msg), msg)
	messageHash := keccak256Hash([]byte(prefix))
	r := new(big.Int).SetBytes(signatureBytes[:32])
	s := new(big.Int).SetBytes(signatureBytes[32:64])
	v := signatureBytes[64]
	publicKey, err := recoverPublicKey(messageHash, r, s, v)
	if err != nil {
		return false, err
	}
	verifyAddress, _ := verifyAddressUsingPublicKey(chainAddress, publicKey)
	if !verifyAddress {
		return false, errors.New(funcT("The public key in the signature does not match the address"))
	}

	return true, nil
}

// fromHex 将16进制字符串解码为字节数组
func fromHex(s string) ([]byte, error) {
	if len(s) >= 2 && s[0] == '0' && (s[1] == 'x' || s[1] == 'X') {
		s = s[2:]
	}
	if len(s)%2 == 1 {
		s = "0" + s
	}
	return hex.DecodeString(s)
}

// keccak256Hash 对字节数字进行hash
func keccak256Hash(data []byte) []byte {
	d := sha3.NewLegacyKeccak256()
	d.Write(data)
	return d.Sum(nil)
}

/*
// checkKeyCurve 判断是否是NIST标准的公钥

	func checkKeyCurve(k *ecdsa.PublicKey) bool {
		if k.X == nil || k.Y == nil {
			return false
		}
		switch k.Params().Name {
		case "P-256": // NIST
			return true
		default: // 不支持的密码学类型
			return false
		}
	}

type ECDSASignature struct {
	R, S *big.Int
}
*/
// verifyAddressUsingPublicKey 验证钱包地址是否和指定的公钥匹配. 如果成功,返回true和对应的密码学标记位;如果失败,返回false和默认的密码学标记位0
func verifyAddressUsingPublicKey(address string, pub *ecdsa.PublicKey) (bool, uint8) {
	//base58反解回byte[]数组
	slice := base58Decode(address)
	//检查是否是合法的base58编码
	if len(slice) < 1 {
		return false, 0
	}
	//拿到密码学标记位
	byteVersion := slice[:1]
	nVersion := uint8(byteVersion[0])

	realAddress, error := getAddressFromPublicKey(pub)
	if error != nil {
		return false, 0
	}

	if realAddress == address {
		return true, nVersion
	}

	return false, 0
}

// getAddressFromPublicKey 返回33位长度的Address
func getAddressFromPublicKey(pub *ecdsa.PublicKey) (string, error) {
	// 将ECDSA公钥转换为ECDH公钥
	ecdhPublicKey, err := pub.ECDH()
	if err != nil {
		return "", err
	}
	// 替换废弃的 elliptic.Marshal 函数
	data := ecdhPublicKey.Bytes()

	outputSha256 := hashUsingSha256(data)
	OutputRipemd160 := hashUsingRipemd160(outputSha256)

	//暂时只支持一个字节长度,也就是uint8的密码学标志位
	// 判断是否是nist标准的私钥
	nVersion := 1

	switch pub.Params().Name {
	case "P-256": // NIST
	case "SM2-P-256": // 国密
		nVersion = 2
	default: // 不支持的密码学类型
		return "", fmt.Errorf("this cryptography[%v] has not been supported yet", pub.Params().Name)
	}

	bufVersion := []byte{byte(nVersion)}

	strSlice := make([]byte, len(bufVersion)+len(OutputRipemd160))
	copy(strSlice, bufVersion)
	copy(strSlice[len(bufVersion):], OutputRipemd160)

	//using double SHA256 for future risks
	checkCode := doubleSha256(strSlice)
	simpleCheckCode := checkCode[:4]

	slice := make([]byte, len(strSlice)+len(simpleCheckCode))
	copy(slice, strSlice)
	copy(slice[len(strSlice):], simpleCheckCode)

	//使用base58编码,手写不容易出错.
	//相比Base64,Base58不使用数字"0",字母大写"O",字母大写"I",和字母小写"l",以及"+"和"/"符号
	strEnc := base58Encode(slice)

	return strEnc, nil
}

// hashUsingSha256 使用sha256 Hash
func hashUsingSha256(data []byte) []byte {
	h := sha256.New()
	h.Write(data)
	out := h.Sum(nil)
	return out
}

// doubleSha256 执行2次SHA256,这是为了防止SHA256算法被攻破
func doubleSha256(data []byte) []byte {
	return hashUsingSha256(hashUsingSha256(data))
}

// hashUsingRipemd160 Ripemd160 hash算法可以缩短长度
func hashUsingRipemd160(data []byte) []byte {
	h := ripemd160.New()
	h.Write(data)
	out := h.Sum(nil)
	return out
}

// recoverPublicKey recovers the public key from r, s, v (all in []byte format) and the hash.
func recoverPublicKey(hash []byte, r *big.Int, s *big.Int, v byte) (*ecdsa.PublicKey, error) {
	curve := elliptic.P256()

	// The curve order (N)
	N := curve.Params().N

	// Ensure r and s are valid
	if r.Sign() <= 0 || r.Cmp(N) >= 0 || s.Sign() <= 0 || s.Cmp(N) >= 0 {
		return nil, fmt.Errorf("invalid r or s values")
	}

	// Determine the x-coordinate and y's parity
	x, _ := recoverX(curve, r, s, hash)
	// 奇数
	isOdd := v&1 == 1 // Determine y's parity based on v

	// Calculate the y-coordinate
	y := calculateY(curve, x, isOdd)
	if y == nil {
		return nil, fmt.Errorf("failed to calculate y coordinate")
	}

	// Create the public key
	pubKey := ecdsa.PublicKey{Curve: curve, X: x, Y: y}

	// Verify if the recovered public key is correct
	if ecdsa.Verify(&pubKey, hash, r, s) {
		return &pubKey, nil
	}
	fmt.Println("recoverPublicKey-hash:", hex.EncodeToString(hash))
	fmt.Println("recoverPublicKey-r:", hex.EncodeToString(r.Bytes()))
	fmt.Println("recoverPublicKey-s:", hex.EncodeToString(s.Bytes()))
	fmt.Println("recoverPublicKey-x:", hex.EncodeToString(pubKey.X.Bytes()))
	fmt.Println("recoverPublicKey-y:", hex.EncodeToString(pubKey.Y.Bytes()))
	return nil, fmt.Errorf("failed to verify the signature with the recovered public key")
}

// calculateY calculates the y-coordinate of the curve given x and isOdd.
func calculateY(curve elliptic.Curve, x *big.Int, isOdd bool) *big.Int {
	params := curve.Params()

	// y² = x³ - 3x + b (mod p)
	x3 := new(big.Int).Exp(x, big.NewInt(3), params.P)
	x3.Sub(x3, new(big.Int).Mul(big.NewInt(3), x))
	x3.Add(x3, params.B)
	x3.Mod(x3, params.P)

	// Calculate square root of x³ - 3x + b mod p
	y := new(big.Int).ModSqrt(x3, params.P)
	if y == nil {
		return nil // No valid y found
	}

	// Ensure the correct parity based on isOdd
	if y.Bit(0) != 0 == isOdd {
		y.Sub(params.P, y)
	}
	return y
}

// 计算x坐标的恢复算法
func recoverX(curve elliptic.Curve, r *big.Int, s *big.Int, hashByte []byte) (*big.Int, error) {
	// N是椭圆曲线的阶
	N := curve.Params().N
	hash := new(big.Int).SetBytes(hashByte)
	// 计算k的逆
	kInv := new(big.Int).ModInverse(s, N)

	// 计算k = s^(-1) * (hash + r * d) mod N
	k := new(big.Int).Mul(kInv, new(big.Int).Add(hash, new(big.Int).Mul(r, kInv)))
	k.Mod(k, N)

	// 计算恢复的x坐标
	return r, nil
}

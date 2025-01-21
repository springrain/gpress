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
	dcrdPublicKey, _, err := dcrdEcdsa.RecoverCompact(sign, messageHash)
	if err != nil {
		return false, err
	}
	pubKeyBytes := dcrdPublicKey.SerializeUncompressed()[1:]
	addressHash := keccak256Hash(pubKeyBytes)
	address := ""
	if len(addressHash) > 20 {
		address = fmt.Sprintf("0x%x", addressHash[len(addressHash)-20:])
	}
	return strings.EqualFold(address, chainAddress), nil
}

// verifyXuperChainSignature XuperChain使用NIST标准的公钥,验证签名
func verifyXuperChainSignature(chainAddress string, msg string, signature string) (valid bool, err error) {

	verify, publicKey, err := verifySecp256r1Signature(msg, signature)
	if verify == false || err != nil {
		return false, err
	}
	// 验证XuperChain的address
	verifyAddress, _ := verifyAddressUsingPublicKey(chainAddress, publicKey)
	if !verifyAddress {
		return false, errors.New(funcT("The public key in the signature does not match the address"))
	}

	return true, nil
}

func verifySecp256r1Signature(msg string, signature string) (bool, *ecdsa.PublicKey, error) {
	signatureBytes, err := fromHex(signature)
	if err != nil {
		return false, nil, err
	}
	if len(signatureBytes) < 65 {
		return false, nil, errors.New("invalid signature")
	}
	// 计算消息的哈希,包括消息前缀
	prefix := fmt.Sprintf("\x86XuperChain Signed Message:\n%d%s", len(msg), msg)
	messageHash := keccak256Hash([]byte(prefix))
	r := new(big.Int).SetBytes(signatureBytes[:32])
	s := new(big.Int).SetBytes(signatureBytes[32:64])
	v := signatureBytes[64]
	publicKey, err := recoverPublicKey(messageHash, r, s, uint(v))
	if err != nil {
		return false, publicKey, err
	}

	return true, publicKey, nil
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

// TODO 不稳定,恢复的公钥有可能验签失败
// recoverPublicKey recovers the public key from r, s, v (all in []byte format) and the hash.
func recoverPublicKey(hash []byte, r *big.Int, s *big.Int, recoveryID uint) (*ecdsa.PublicKey, error) {
	curve := elliptic.P256()
	params := curve.Params()
	recoveryID = recoveryID % 2
	//recoveryID := (uint(v) + 1) % 2

	// 检查r和s范围
	if r.Sign() <= 0 || s.Sign() <= 0 || r.Cmp(params.N) >= 0 || s.Cmp(params.N) >= 0 {
		return nil, errors.New("invalid r/s value")
	}

	// 计算R点x坐标
	x := new(big.Int).Set(r)
	if x.Cmp(params.P) >= 0 {
		return nil, errors.New("r >= P")
	}

	// 计算y² = x³ - 3x + b mod P
	x3 := new(big.Int).Exp(x, big.NewInt(3), params.P)
	threeX := new(big.Int).Mul(x, big.NewInt(3))
	threeX.Mod(threeX, params.P)
	ySquared := new(big.Int).Sub(x3, threeX)
	ySquared.Add(ySquared, params.B)
	ySquared.Mod(ySquared, params.P)

	// 计算y坐标
	y := new(big.Int).ModSqrt(ySquared, params.P)
	if y == nil {
		return nil, errors.New("invalid R point")
	}

	// 根据恢复ID调整y奇偶性
	if (y.Bit(0) == 0 && recoveryID == 1) || (y.Bit(0) == 1 && recoveryID == 0) {
		y.Sub(params.P, y)
	}

	// 计算r的模逆元
	rInv := new(big.Int).ModInverse(r, params.N)
	if rInv == nil {
		return nil, errors.New("r is not invertible")
	}

	// 计算sR点
	sRx, sRy := curve.ScalarMult(x, y, s.Bytes())

	// 计算e = hash mod N
	e := new(big.Int).SetBytes(hash)
	e.Mod(e, params.N)

	// 计算eG点
	eGx, eGy := curve.ScalarBaseMult(e.Bytes())

	// 计算sR - eG
	minusEGy := new(big.Int).Neg(eGy)
	minusEGy.Mod(minusEGy, params.P)
	sumX, sumY := curve.Add(sRx, sRy, eGx, minusEGy)

	// 乘以r逆元得到公钥Q
	qX, qY := curve.ScalarMult(sumX, sumY, rInv.Bytes())

	// 验证点有效性
	if !curve.IsOnCurve(qX, qY) {
		return nil, errors.New("recovered point is invalid")
	}

	return &ecdsa.PublicKey{Curve: curve, X: qX, Y: qY}, nil
}

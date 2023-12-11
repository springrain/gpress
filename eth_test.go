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
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
)

func TestEth(t *testing.T) {

	// 生成新的以太坊密钥对
	privateKey, err := generatePrivateKey()
	if err != nil {
		log.Fatal("生成以太坊私钥失败:", err)
	}
	privateKeyHex := hex.EncodeToString(privateKey.D.Bytes())
	fmt.Println("以太坊私钥:", privateKeyHex)
	publicKey := privateKey.Public().(*ecdsa.PublicKey)

	/*
		publicKeyHex := "0434d53f7efc6d3c63d0a3eb010c6bea6b153dca5f5a0a1e747a3d5a64a62959494af78b111d5025b866d3c90e2a100c7e77a7cc6ed0a3e15c7a9e65c3f7b5f0a3"
		privateKeyHex := "f75c9ff1b8c2e0c1f78a39bc67a92dbb26cc727bac9e0a9f4c9c2d0a337d0ef5"

		// 解析公钥
		publicKey, err := parsePublicKey(publicKeyHex)
		if err != nil {
			log.Fatal("解析公钥失败:", err)
		}

		// 解析私钥
		privateKey, err := parsePrivateKey(privateKeyHex)
		if err != nil {
			log.Fatal("解析私钥失败:", err)
		}
	*/

	// 生成以太坊地址
	address := generateAddress(publicKey)
	fmt.Println("以太坊地址:", address)

	// 要加密的数据
	plaintext := []byte("Hello, World!")

	// 使用以太坊公钥加密数据
	ciphertext, err := encryptWithPublicKey(publicKey, plaintext)
	if err != nil {
		log.Fatal("加密失败:", err)
	}

	fmt.Println("加密后的数据:", hex.EncodeToString(ciphertext))

	// 使用以太坊私钥解密数据
	decrypted, err := decryptWithPrivateKey(privateKey, ciphertext)
	if err != nil {
		log.Fatal("解密失败:", err)
	}

	fmt.Println("解密后的数据:", string(decrypted))
}

// 生成以太坊私钥
func generatePrivateKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(secp256k1(), rand.Reader)
}

// 使用以太坊公钥加密数据
func encryptWithPublicKey(publicKey *ecdsa.PublicKey, plaintext []byte) ([]byte, error) {
	// 生成临时私钥
	tempPrivateKey, err := generatePrivateKey()
	if err != nil {
		return nil, err
	}

	// 计算临时公钥的坐标
	tempPublicKeyX, tempPublicKeyY := secp256k1().ScalarBaseMult(tempPrivateKey.D.Bytes())

	// 计算共享密钥的坐标
	sharedKeyX, _ := secp256k1().ScalarMult(publicKey.X, publicKey.Y, tempPrivateKey.D.Bytes())

	// 使用共享密钥的坐标作为初始化向量
	iv := sharedKeyX.Bytes()

	// 加密数据
	ciphertext := xor(plaintext, iv)

	// 组合加密后的临时公钥和密文
	encryptedData := append(tempPublicKeyX.Bytes(), tempPublicKeyY.Bytes()...)
	encryptedData = append(encryptedData, ciphertext...)

	return encryptedData, nil
}

// 使用以太坊私钥解密数据
func decryptWithPrivateKey(privateKey *ecdsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	// 解析临时公钥的坐标
	tempPublicKeyX := new(big.Int).SetBytes(ciphertext[:32])
	tempPublicKeyY := new(big.Int).SetBytes(ciphertext[32:64])

	// 计算共享密钥的坐标
	sharedKeyX, _ := secp256k1().ScalarMult(tempPublicKeyX, tempPublicKeyY, privateKey.D.Bytes())

	// 使用共享密钥的坐标作为初始化向量
	iv := sharedKeyX.Bytes()

	// 解密数据
	plaintext := xor(ciphertext[64:], iv)

	return plaintext, nil
}

// 执行异或操作
func xor(a, b []byte) []byte {
	result := make([]byte, len(a))
	for i := range a {
		result[i] = a[i] ^ b[i]
	}
	return result
}

// 返回 secp256k1 椭圆曲线
func secp256k1() elliptic.Curve {
	return elliptic.P256()
}

// 解析公钥
func parsePublicKey(publicKeyHex string) (*ecdsa.PublicKey, error) {
	publicKeyBytes, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		return nil, err
	}

	// 由于以太坊的公钥前面有一个字节的标志位，需要去掉
	publicKey := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     new(big.Int).SetBytes(publicKeyBytes[1:33]),
		Y:     new(big.Int).SetBytes(publicKeyBytes[33:65]),
	}

	return publicKey, nil
}

// 解析私钥
func parsePrivateKey(privateKeyHex string) (*ecdsa.PrivateKey, error) {
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return nil, err
	}

	privateKey := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: elliptic.P256(),
		},
		D: new(big.Int).SetBytes(privateKeyBytes),
	}

	privateKey.PublicKey.X, privateKey.PublicKey.Y = privateKey.PublicKey.Curve.ScalarBaseMult(privateKey.D.Bytes())

	return privateKey, nil
}

// 生成以太坊地址
func generateAddress(publicKey *ecdsa.PublicKey) string {
	publicKeyBytes := elliptic.Marshal(publicKey.Curve, publicKey.X, publicKey.Y)
	hash := sha256.Sum256(publicKeyBytes)
	address := hash[12:]

	return fmt.Sprintf("0x%x", address)
}

// https://www.jianshu.com/p/71a4454c74da
// https://github.com/wenweih/bitcoin_address_protocol
// https://studygolang.com/articles/25128

// https://segmentfault.com/a/1190000018359512
func TestEthBTCPrivateKey(t *testing.T) {
	// 生成私钥
	//privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	privateKey, err := crypto.HexToECDSA("3e5adb36d7975cfb4af7fd7adcba310d915ec529400806516a2a92df6a6010c5")
	//privateKey, err := crypto.GenerateKey()

	if err != nil {
		fmt.Println("生成私钥出错：", err)
		return
	}

	// 根据私钥生成公钥
	publicKey := privateKey.PublicKey

	// 计算公钥的keccak256哈希
	keccak := sha3.NewLegacyKeccak256()
	keccak.Write(publicKey.X.Bytes())
	keccak.Write(publicKey.Y.Bytes())
	ethHash := keccak.Sum(nil)

	// 从哈希中提取最后的20个字节（40个十六进制字符）
	ethAddress := hex.EncodeToString(ethHash[12:])

	// 将私钥和公钥转换为十六进制字符串
	privateKeyHex := hex.EncodeToString(privateKey.D.Bytes())
	publicKeyHex := hex.EncodeToString(append(publicKey.X.Bytes(), publicKey.Y.Bytes()...))

	// 打印结果
	fmt.Println("私钥：", privateKeyHex)
	fmt.Println("公钥：", publicKeyHex)
	fmt.Println("以太坊地址：", "0x"+ethAddress)

}
func TestEth2(t *testing.T) {
	// 生成以太坊私钥
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	// 将私钥转换为字节切片
	privateKeyBytes := crypto.FromECDSA(privateKey)

	// 生成以太坊公钥
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("无法获取公钥")
	}

	// 将公钥转换为字节切片
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)

	// 生成以太坊地址
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	// 将私钥、公钥和地址打印出来
	fmt.Println("私钥:", hex.EncodeToString(privateKeyBytes))
	fmt.Println("公钥:", hex.EncodeToString(publicKeyBytes))
	fmt.Println("地址:", address)
}
func TestEth3(t *testing.T) {
	// MetaMask 签名数据
	signature := "0x4acafcdd5ee478e14453a36c074dee7d142dca7ead7a2029a0c6b7a3e547ee46379018cd8fc661450c84fd90bcca34e2e0008a02b27eab7944a97947c0d8bfa71b"
	// 消息字符串
	message := "20230522151922392009508861"
	// 发送者地址
	senderAddress := "0xD530eC9517C20DE518345A7210338dFB6279f454"
	verify, err := verifySecp256k1Signature(senderAddress, message, signature)
	fmt.Println(verify)
	fmt.Println(err)
}

func TestECDH(t *testing.T) {
	// 生成ECDH私钥
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Println("私钥生成失败:", err)
		return
	}

	// 生成对应的ECDH公钥
	publicKey := privateKey.PublicKey

	// 生成对方的ECDH公钥
	otherPublicKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Println("对方公钥生成失败:", err)
		return
	}

	// ECDH密钥交换
	x, _ := publicKey.Curve.ScalarMult(otherPublicKey.X, otherPublicKey.Y, privateKey.D.Bytes())

	// 计算共享密钥
	sharedKey := sha256.Sum256(x.Bytes())

	// 加密数据
	plaintext := []byte("Hello, world!")
	ciphertext, err := encryptAES(sharedKey[:], plaintext)
	if err != nil {
		fmt.Println("加密失败:", err)
		return
	}

	// 解密数据
	decryptedText, err := decryptAES(sharedKey[:], ciphertext)
	if err != nil {
		fmt.Println("解密失败:", err)
		return
	}

	fmt.Printf("加密前的数据: %s\n", plaintext)
	fmt.Printf("解密后的数据: %s\n", decryptedText)
}

// 使用AES-CBC模式加密数据
func encryptAES(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext, nil
}

// 使用AES-CBC模式解密数据
func decryptAES(key, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	plaintext := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, ciphertext)

	return plaintext, nil
}

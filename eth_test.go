package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"testing"

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

// https://segmentfault.com/a/1190000018359512
func TestEthBTCPrivateKey(t *testing.T) {
	// 生成私钥
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
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

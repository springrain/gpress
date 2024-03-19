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

	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/ripemd160"
)

func verifySecp256k1Signature(senderAddress string, signatureData string, signature string) (bool, error) {
	// 将签名数据解码为字节数组
	signatureBytes := common.FromHex(signature)

	// 将发送者地址解码为以太坊地址类型
	sender := common.HexToAddress(senderAddress)

	// 计算消息的哈希,包括 MetaMask 的消息前缀
	prefix := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(signatureData), signatureData)
	messageBytes := []byte(prefix)
	messageHash := ethcrypto.Keccak256Hash(messageBytes)

	// 提取恢复 ID
	recoveryID := signatureBytes[64]
	if recoveryID != 27 && recoveryID != 28 {
		return false, errors.New("invalid recovery ID")
	}

	// 修复恢复 ID 的值
	if recoveryID == 27 {
		signatureBytes[64] = 0
	} else {
		signatureBytes[64] = 1
	}

	// 使用签名数据验证消息哈希
	signaturePublicKey, err := ethcrypto.SigToPub(messageHash.Bytes(), signatureBytes)
	if err != nil {
		return false, err
	}

	signerAddress := ethcrypto.PubkeyToAddress(*signaturePublicKey)
	if signerAddress != sender {
		return false, errors.New("signature verification failed")
	}
	return true, nil

}

// XuperChain使用NIST标准的公钥
func verifyXuperSignature(chainAddress string, sig, msg []byte) (valid bool, err error) {
	/*k := &ecdsa.PublicKey{}
	err = json.Unmarshal([]byte(chainAddress), k)
	if err != nil {
		return false, err //json有问题
	}

	k.Curve = elliptic.P256()

	// 判断是否是NIST标准的公钥
	isNistCurve := checkKeyCurve(k)
	if isNistCurve == false {
		return false, fmt.Errorf("this cryptography curve[%s] has not been supported yet.", k.Params().Name)
	}

	r, s, err := unmarshalECDSASignature(sig)
	if err != nil {
		return false, fmt.Errorf("failed to unmarshal the ecdsa signature [%s]", err)
	}

	return ecdsa.Verify(k, msg, r, s), nil*/

	s2 := string(sig)
	signature, e := hex.DecodeString(s2)
	if e != nil {
		return false, err
	}

	r := new(big.Int).SetBytes(signature[:32])
	s := new(big.Int).SetBytes(signature[32:64])
	publicKeyX := new(big.Int).SetBytes(signature[64:96])
	publicKeyY := new(big.Int).SetBytes(signature[96:128])
	data := signature[128:]
	if string(msg) != string(data) { //数据不一致
		return false, errors.New("原始数据不一致")
	}

	/*x := new(big.Int)
	y := new(big.Int)
	e = x.UnmarshalText(publicKeyX)
	if e != nil {
		return false, e
	}
	e = y.UnmarshalText(publicKeyY)
	if e != nil {
		return false, e
	}*/
	pub := ecdsa.PublicKey{Curve: elliptic.P256(), X: publicKeyX, Y: publicKeyY}
	pubKey := &pub
	verifyAddress, _ := verifyAddressUsingPublicKey(chainAddress, pubKey)
	if !verifyAddress {
		return false, errors.New("签名中的公钥和address不匹配")
	}

	return ecdsa.Verify(pubKey, data, r, s), nil
}

// 判断是否是NIST标准的公钥
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

/*
// hashSha256 使用sha256计算hash值

	func hashSha256(str string) string {
		hashByte := sha256.Sum256([]byte(str))
		hashStr := hex.EncodeToString(hashByte[:])
		return hashStr
	}

// use DER-encoded ASN.1 octet standard to represent the signature
// 与比特币算法一样,基于DER-encoded ASN.1 octet标准,来表达使用椭圆曲线签名算法返回的结果

	func MarshalECDSASignature(r, s *big.Int) ([]byte, error) {
		return asn1.Marshal(ECDSASignature{r, s})
	}

// 将公钥序列化成byte数组

	func MarshalPublicKey(publicKey *ecdsa.PublicKey) []byte {
		return elliptic.Marshal(publicKey.Curve, publicKey.X, publicKey.Y)
	}

func unmarshalECDSASignature(rawSig []byte) (*big.Int, *big.Int, error) {
	sig := new(ECDSASignature)
	_, err := asn1.Unmarshal(rawSig, sig)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to unmashal the signature [%v] to R & S, and the error is [%s]", rawSig, err)
	}

	if sig.R == nil {
		return nil, nil, errors.New("invalid signature, R is nil")
	}
	if sig.S == nil {
		return nil, nil, errors.New("invalid signature, S is nil")
	}

	if sig.R.Sign() != 1 {
		return nil, nil, errors.New("invalid signature, R must be larger than zero")
	}
	if sig.S.Sign() != 1 {
		return nil, nil, errors.New("invalid signature, S must be larger than zero")
	}

	return sig.R, sig.S, nil
}
*/
// 验证钱包地址是否和指定的公钥match
// 如果成功,返回true和对应的密码学标记位；如果失败,返回false和默认的密码学标记位0
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

// 返回33位长度的地址
func getAddressFromPublicKey(pub *ecdsa.PublicKey) (string, error) {
	//using SHA256 and Ripemd160 for hash summary
	//data := elliptic.Marshal(pub.Curve, pub.X, pub.Y)
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

	//使用base58编码,手写不容易出错。
	//相比Base64,Base58不使用数字"0",字母大写"O",字母大写"I",和字母小写"l",以及"+"和"/"符号。
	strEnc := base58Encode(slice)

	return strEnc, nil
}
func hashUsingSha256(data []byte) []byte {
	h := sha256.New()
	h.Write(data)
	out := h.Sum(nil)

	return out
}

// 执行2次SHA256,这是为了防止SHA256算法被攻破。
func doubleSha256(data []byte) []byte {
	return hashUsingSha256(hashUsingSha256(data))
}

// Ripemd160,这种hash算法可以缩短长度
func hashUsingRipemd160(data []byte) []byte {
	h := ripemd160.New()
	h.Write(data)
	out := h.Sum(nil)

	return out
}

/*
	type Signature struct {
		KeyID     string `json:"keyId"`
		Algorithm string `json:"algorithm"`
		Headers   string `json:"headers"`
		Value     string `json:"signature"`
	}

	func parseSignature(signatureString string) (*Signature, error) {
		//逗号分割签名的字符串
		s1 := strings.Split(signatureString, ",")
		sigMap := make(map[string]string, 0)
		for _, s2 := range s1 {
			dIndex := strings.Index(s2, "=")
			if dIndex < 0 {
				continue
			}
			sigMap[s2[:dIndex]] = strings.Trim(strings.TrimSpace(s2[dIndex+1:]), `"`)
		}
		sigByte, _ := json.Marshal(sigMap)
		// 解析签名字符串,提取相关信息
		signature := &Signature{}
		err := json.Unmarshal(sigByte, signature)
		if err != nil {
			return nil, err
		}
		return signature, nil
	}

	func verifySignature(signature *Signature, data string) (bool, error) {
		switch signature.Algorithm {

			case "rsa-sha256":
				publicKey, err := getRSAPublicKeyPem(signature.KeyID)
				if err != nil {
					return false, err
				}
				// 验证签名
				hashed := sha256.Sum256([]byte(data))
				signatureBytes, err := base64.StdEncoding.DecodeString(signature.Value)
				if err != nil {
					return false, err
				}

				err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], signatureBytes)
				if err != nil {
					return false, err
				}

		case "secp256k1": //以太坊账号的签名算法
			// KeyID应为 chain://域名[address],合约地址,链ID    域名下的 publicKey 值
			// 主要就是要解析IP地址,备选  #域名[address]#合约地址#链ID,通过#后缀跟上信息,前端点击时,使用js ajax获取需要处理的数据
			// KeyID应为 address,用于和签名数据里获取的address进行比较
			// 这里KeyID暂时定为address,实际应该为区块链域名,从域名反查合约获取address.这里比较简单
			return verifySecp256k1Signature(signature.KeyID, data, signature.Value)
		}
		return true, nil
	}

	func getRSAPublicKeyPem(publicKeyID string) (*rsa.PublicKey, error) {
		// 根据公钥 ID 获取对应的公钥
		// 这里使用假数据,实际使用时需要替换为真实的公钥获取逻辑
		//publicKeyPEM := "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAr9HicDyHYlpGVYVHrm7j\nU7Nq4z9SeynK8UUi+JoBWuotChg2oSDQtWuj+zdQSKM3g27+sqNNw/BuZp85BVT6\n8PRyamTHjVrZPj6JIC+A/EGeJTqycODoMTDTTdz3evxBUbPAH7By91VrMNE5i8zl\nJ40IqAYYNLjmUdvQliGmGpX/xmPAfIeJ/mMQ3kCq/2uSICrL1ORicAB/qqXgyPsB\nWZCTYOOdJsV9bbbhAQUqRjevZrRIdaVcrIObxTDY0VgtBJgsElGNxbnb/g4vfPgy\nWdi/E0qLSRyayml8lGZhPccgY3PnqGO765X/j0tra/I4JIjLC0AOV0nLs0fLmH72\nEwIDAQAB\n-----END PUBLIC KEY-----\n"
		publicKeyPEM, err := responseJsonValue(publicKeyID, "publicKey.publicKeyPem", publicKeyID)
		if err != nil {
			return nil, err
		}

		if publicKeyPEM == nil {
			return nil, errors.New("获取公钥值为nil")
		}

		// 解析公钥 PEM 格式
		block, _ := pem.Decode([]byte(publicKeyPEM.(string)))
		if block == nil || block.Type != "PUBLIC KEY" {
			return nil, errors.New("公钥解析失败")
		}
		// 解析公钥 DER 格式
		publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}

		// 转换为 RSA 公钥
		publicKey, ok := publicKeyInterface.(*rsa.PublicKey)
		if !ok {
			return nil, errors.New("公钥类型错误")
		}

		return publicKey, nil
	}

	func buildSignatureData(c *app.RequestContext, headers string) string {
		// 构建签名字符串
		var comparisonStrings []string
		signedHeaders := strings.Split(headers, " ")
		for _, header := range signedHeaders {
			value := ""
			header = strings.TrimSpace(header)
			if header == "(request-target)" {
				method := string(c.Method())
				method = strings.ToLower(method)
				uri := string(c.Request.URI().Path())
				value = fmt.Sprintf("%s %s", method, uri)
			} else {
				value = string(c.GetHeader(header))
			}
			comparisonStrings = append(comparisonStrings, header+": "+value)
		}
		return strings.Join(comparisonStrings, "\n")
	}

// generateRSASignature 对字符串签名,并将签名结果进行 Base64 编码

	func generateRSASignature(signingString string) (string, error) {
		// 读取私钥文件
		privateKeyFile := datadir + "pem/private.pem"
		privateKeyPEM, err := os.ReadFile(privateKeyFile)
		if err != nil {
			return "", fmt.Errorf("读取私钥文件失败:%w", err)
		}

		// 解析私钥
		privateKeyBlock, _ := pem.Decode(privateKeyPEM)
		privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
		if err != nil {
			return "", fmt.Errorf("解析私钥失败:%w", err)
		}

		// 对签名前的字符串进行签名
		hashed := sha256.Sum256([]byte(signingString))
		signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
		if err != nil {
			return "", fmt.Errorf("签名失败:%w", err)
		}

		// 将签名结果进行 Base64 编码
		signatureBase64 := base64.StdEncoding.EncodeToString(signature)
		return signatureBase64, nil
	}
*/

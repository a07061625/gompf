// Package mpencrypt rsa
// User: 姜伟
// Time: 2020-02-19 06:29:32
package mpencrypt

import (
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "encoding/pem"
    "errors"
)

// 公钥文件生成
// openssl genrsa -out rsa_private_key.pem 1024
// 私钥文件生成
// openssl rsa -in rsa_private_key.pem -pubout -out rsa_public_key.pem

// RsaEncrypt 加密
// originData []byte 原始数据
// publicKey []byte 公钥,公钥文件内容,包含-----BEGIN PUBLIC KEY-----和-----END PUBLIC KEY-----
func RsaEncrypt(originData []byte, publicKey []byte) ([]byte, error) {
    // 解密pem格式的公钥
    block, _ := pem.Decode(publicKey)
    if block == nil {
        return nil, errors.New("public key error")
    }
    // 解析公钥
    pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
    if err != nil {
        return nil, err
    }
    // 类型断言
    pub := pubKey.(*rsa.PublicKey)
    return rsa.EncryptPKCS1v15(rand.Reader, pub, originData)
}

// RsaDecrypt 解密
// encryptData []byte 加密数据
// privateKey []byte 私钥,私钥文件内容,包含-----BEGIN RSA PRIVATE KEY-----和-----END RSA PRIVATE KEY-----
func RsaDecrypt(encryptData []byte, privateKey []byte) ([]byte, error) {
    // 解密pem格式的私钥
    block, _ := pem.Decode(privateKey)
    if block == nil {
        return nil, errors.New("private key error! ")
    }
    // 解析PKCS1格式的私钥
    priKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
    if err != nil {
        return nil, err
    }

    return rsa.DecryptPKCS1v15(rand.Reader, priKey, encryptData)
}

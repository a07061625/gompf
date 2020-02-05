/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2019/12/26 0026
 * Time: 9:09
 */
package mpencrypt

import (
    "bytes"
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "errors"
    "io"
)

func AesPaddingPKCS7(cipherText []byte, blockSize int) []byte {
    padding := blockSize - len(cipherText)%blockSize
    padText := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(cipherText, padText...)
}

func AesUnPaddingPKCS7(origData []byte) []byte {
    length := len(origData)
    unPadding := int(origData[length-1])
    return origData[:(length - unPadding)]
}

func AesEncryptCBC(origData, key []byte) []byte {
    // 分组秘钥
    // NewCipher该函数限制了输入k的长度必须为16, 24或者32
    block, _ := aes.NewCipher(key)
    // 获取秘钥块的长度
    blockSize := block.BlockSize()
    // 补全码
    origData = AesPaddingPKCS7(origData, blockSize)
    // 加密模式
    blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
    // 创建数组
    encrypted := make([]byte, len(origData))
    // 加密
    blockMode.CryptBlocks(encrypted, origData)
    return encrypted
}

func AesDecryptCBC(encrypted, key []byte) []byte {
    // 分组秘钥
    block, _ := aes.NewCipher(key)
    // 获取秘钥块的长度
    blockSize := block.BlockSize()
    // 加密模式
    blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
    // 创建数组
    decrypted := make([]byte, len(encrypted))
    // 解密
    blockMode.CryptBlocks(decrypted, encrypted)
    // 去除补全码
    decrypted = AesUnPaddingPKCS7(decrypted)
    return decrypted
}

func AesEncryptCBCPKCS7(plainData, aesKey []byte) ([]byte, error) {
    k := len(aesKey)
    if len(plainData)%k != 0 {
        plainData = AesPaddingPKCS7(plainData, k)
    }

    block, err := aes.NewCipher(aesKey)
    if err != nil {
        return nil, err
    }

    iv := make([]byte, aes.BlockSize)
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return nil, err
    }

    cipherData := make([]byte, len(plainData))
    blockMode := cipher.NewCBCEncrypter(block, iv)
    blockMode.CryptBlocks(cipherData, plainData)
    return cipherData, nil
}

func AesDecryptCBCPKCS7(cipherData, aesKey []byte) ([]byte, error) {
    k := len(aesKey)
    if len(cipherData)%k != 0 {
        return nil, errors.New("密文长度不是密钥长度的整数倍")
    }

    block, err := aes.NewCipher(aesKey)
    if err != nil {
        return nil, err
    }

    iv := make([]byte, aes.BlockSize)
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return nil, err
    }

    blockMode := cipher.NewCBCDecrypter(block, iv)
    plainData := make([]byte, len(cipherData))
    blockMode.CryptBlocks(plainData, cipherData)
    return plainData, nil
}

func ecbGenerateKey(key []byte) []byte {
    genKey := make([]byte, 16)
    copy(genKey, key)
    for i := 16; i < len(key); {
        for j := 0; j < 16 && i < len(key); j, i = j+1, i+1 {
            genKey[j] ^= key[i]
        }
    }
    return genKey
}

func AesEncryptECB(origData []byte, key []byte) []byte {
    cipherText, _ := aes.NewCipher(ecbGenerateKey(key))
    length := (len(origData) + aes.BlockSize) / aes.BlockSize
    plain := make([]byte, length*aes.BlockSize)
    copy(plain, origData)
    pad := byte(len(plain) - len(origData))
    for i := len(origData); i < len(plain); i++ {
        plain[i] = pad
    }
    encrypted := make([]byte, len(plain))
    // 分组分块加密
    for bs, be := 0, cipherText.BlockSize(); bs <= len(origData); bs, be = bs+cipherText.BlockSize(), be+cipherText.BlockSize() {
        cipherText.Encrypt(encrypted[bs:be], plain[bs:be])
    }

    return encrypted
}

func AesDecryptECB(encrypted []byte, key []byte) []byte {
    decrypted := make([]byte, len(encrypted))
    cipherText, _ := aes.NewCipher(ecbGenerateKey(key))
    for bs, be := 0, cipherText.BlockSize(); bs < len(encrypted); bs, be = bs+cipherText.BlockSize(), be+cipherText.BlockSize() {
        cipherText.Decrypt(decrypted[bs:be], encrypted[bs:be])
    }

    trim := 0
    if len(decrypted) > 0 {
        trim = len(decrypted) - int(decrypted[len(decrypted)-1])
    }

    return decrypted[:trim]
}

func AesEncryptCFB(origData []byte, key []byte) []byte {
    block, err := aes.NewCipher(key)
    if err != nil {
        panic(err)
    }
    encrypted := make([]byte, aes.BlockSize+len(origData))
    iv := encrypted[:aes.BlockSize]
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        panic(err)
    }
    stream := cipher.NewCFBEncrypter(block, iv)
    stream.XORKeyStream(encrypted[aes.BlockSize:], origData)
    return encrypted
}

func AesDecryptCFB(encrypted []byte, key []byte) []byte {
    block, _ := aes.NewCipher(key)
    if len(encrypted) < aes.BlockSize {
        panic("ciphertext too short")
    }
    iv := encrypted[:aes.BlockSize]
    decrypted := encrypted[aes.BlockSize:]

    stream := cipher.NewCFBDecrypter(block, iv)
    stream.XORKeyStream(decrypted, encrypted)
    return decrypted
}

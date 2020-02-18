// Package mpf hash
// User: 姜伟
// Time: 2020-02-19 05:30:04
package mpf

import (
    "crypto/md5"
    "crypto/sha1"
    "crypto/sha256"
    "crypto/sha512"
    "encoding/hex"
    "hash/crc32"
)

// HashCrc32 Crc32
func HashCrc32(str string, key string) string {
    c := crc32.NewIEEE()
    c.Write([]byte(str))
    extend := c.Sum(nil)
    if len(key) > 0 {
        extend = c.Sum([]byte(key))
    }
    return hex.EncodeToString(extend)
}

// HashMd5 Md5
func HashMd5(str string, key string) string {
    c := md5.New()
    c.Write([]byte(str))
    extend := c.Sum(nil)
    if len(key) > 0 {
        extend = c.Sum([]byte(key))
    }
    return hex.EncodeToString(extend)
}

// HashMd5Sign md5签名
func HashMd5Sign(data string, secret string) string {
    h := md5.New()
    h.Write([]byte(data + secret))
    return hex.EncodeToString(h.Sum(nil))
}

// HashMd5Verify 验证md5签名
func HashMd5Verify(data string, secret string, sign string) bool {
    nowSign := HashMd5Sign(data, secret)
    return sign == nowSign
}

// HashSha1 Sha1
func HashSha1(str string, key string) string {
    c := sha1.New()
    c.Write([]byte(str))
    extend := c.Sum(nil)
    if len(key) > 0 {
        extend = c.Sum([]byte(key))
    }
    return hex.EncodeToString(extend)
}

// HashSha256 Sha256
func HashSha256(str string, key string) string {
    c := sha256.New()
    c.Write([]byte(str))
    extend := c.Sum(nil)
    if len(key) > 0 {
        extend = c.Sum([]byte(key))
    }
    return hex.EncodeToString(extend)
}

// HashSha512 Sha512
func HashSha512(str string, key string) string {
    c := sha512.New()
    c.Write([]byte(str))
    extend := c.Sum(nil)
    if len(key) > 0 {
        extend = c.Sum([]byte(key))
    }
    return hex.EncodeToString(extend)
}

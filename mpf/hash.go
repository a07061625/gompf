/**
 * hash签名
 * User: 姜伟
 * Date: 2019/12/19 0019
 * Time: 17:30
 */
package mpf

import (
    "crypto/md5"
    "crypto/sha1"
    "crypto/sha256"
    "crypto/sha512"
    "encoding/hex"
    "hash/crc32"
)

func HashCrc32(str string, key string) string {
    c := crc32.NewIEEE()
    c.Write([]byte(str))
    extend := c.Sum(nil)
    if len(key) > 0 {
        extend = c.Sum([]byte(key))
    }
    return hex.EncodeToString(extend)
}

func HashMd5(str string, key string) string {
    c := md5.New()
    c.Write([]byte(str))
    extend := c.Sum(nil)
    if len(key) > 0 {
        extend = c.Sum([]byte(key))
    }
    return hex.EncodeToString(extend)
}

// md5签名
func HashMd5Sign(data string, secret string) string {
    h := md5.New()
    h.Write([]byte(data + secret))
    return hex.EncodeToString(h.Sum(nil))
}

// 验证md5签名
func HashMd5Verify(data string, secret string, sign string) bool {
    nowSign := HashMd5Sign(data, secret)
    return sign == nowSign
}

func HashSha1(str string, key string) string {
    c := sha1.New()
    c.Write([]byte(str))
    extend := c.Sum(nil)
    if len(key) > 0 {
        extend = c.Sum([]byte(key))
    }
    return hex.EncodeToString(extend)
}

func HashSha256(str string, key string) string {
    c := sha256.New()
    c.Write([]byte(str))
    extend := c.Sum(nil)
    if len(key) > 0 {
        extend = c.Sum([]byte(key))
    }
    return hex.EncodeToString(extend)
}

func HashSha512(str string, key string) string {
    c := sha512.New()
    c.Write([]byte(str))
    extend := c.Sum(nil)
    if len(key) > 0 {
        extend = c.Sum([]byte(key))
    }
    return hex.EncodeToString(extend)
}

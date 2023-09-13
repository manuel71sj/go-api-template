package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"manuel71sj/go-api-template/pkg/str"
)

func MD5(s string) string {
	sum := md5.Sum(str.S(s).Bytes())
	return hex.EncodeToString(sum[:])
}

func SHA1(s string) string {
	sum := sha1.Sum(str.S(s).Bytes())
	return hex.EncodeToString(sum[:])
}

func SHA256(s string) string {
	sum := sha256.Sum256(str.S(s).Bytes())
	return hex.EncodeToString(sum[:])
}

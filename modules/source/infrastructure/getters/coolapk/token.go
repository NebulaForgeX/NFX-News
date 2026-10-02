package coolapk

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func Token() string {
	deviceID := deviceID()
	now := time.Now().Unix()
	hexNow := fmt.Sprintf("0x%x", now)
	md5Now := md5Hex(strconv.FormatInt(now, 10))
	s := "token://com.coolapk.market/c67ef5943784d09750dcfbb31020f0ab?" + md5Now + "$" + deviceID + "&com.coolapk.market"
	return md5Hex(base64.StdEncoding.EncodeToString([]byte(s))) + deviceID + hexNow
}

func deviceID() string {
	lens := []int{8, 4, 4, 4, 12}
	parts := make([]string, len(lens))
	for i, n := range lens {
		buf := make([]byte, n)
		_, _ = rand.Read(buf)
		const alphabet = "0123456789abcdefghijklmnopqrstuvwxyz"
		for j := range buf {
			buf[j] = alphabet[int(buf[j])%len(alphabet)]
		}
		parts[i] = string(buf)
	}
	return strings.Join(parts, "-")
}

func md5Hex(s string) string {
	sum := md5.Sum([]byte(s))
	return hex.EncodeToString(sum[:])
}

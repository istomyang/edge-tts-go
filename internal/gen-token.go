package internal

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

// GenSecMsGec generate Sec-MS-GEC token.
//
// Use algo from: https://github.com/rany2/edge-tts/issues/290#issuecomment-2464956570
func GenSecMsGec() string {
	var sec uint64 = uint64(time.Now().Unix())
	sec += 11644473600
	sec -= (sec % 300)
	nsec := sec * 1e9 / 100
	str := fmt.Sprintf("%d6A5AA1D4EAFF4E9FB37E23D68491D6F4", nsec)
	hasher := sha256.New()
	hasher.Write([]byte(str))
	hash := hasher.Sum(nil)
	encoding := hex.EncodeToString(hash)
	return strings.ToUpper(encoding)
}

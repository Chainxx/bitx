package signed

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func Sign(salt string, s string) (string, error) {
	mac := hmac.New(sha256.New, []byte(salt))
	_, err := mac.Write([]byte(s))
	
	if err != nil {
		return "", err
	}
	r := hex.EncodeToString(mac.Sum(nil))
	
	return r, nil
}

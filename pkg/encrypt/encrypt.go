package encrypt

import (
	"encoding/base64"
	"github.com/zeromicro/go-zero/core/codec"
)

const (
	emailAesKey = "5A2E746B08D846502F37A6E2D85D583B"
)

func EncEmail(mobile string) (string, error) {
	data, err := codec.EcbEncrypt([]byte(emailAesKey), []byte(mobile))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(data), nil
}

func DecEmail(mobile string) (string, error) {
	originalData, err := base64.StdEncoding.DecodeString(mobile)
	if err != nil {
		return "", err
	}
	data, err := codec.EcbDecrypt([]byte(emailAesKey), originalData)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

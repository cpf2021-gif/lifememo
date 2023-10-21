package encrypt

import (
	"testing"
)

func TestEncryptMobile(t *testing.T) {
	mobile := "1139841@qq.com"
	encryptedMobile, err := EncEmail(mobile)
	if err != nil {
		t.Fatal(err)
	}
	decryptedMobile, err := DecEmail(encryptedMobile)
	if err != nil {
		t.Fatal(err)
	}
	if mobile != decryptedMobile {
		t.Fatalf("expected %s, but got %s", mobile, decryptedMobile)
	}
}

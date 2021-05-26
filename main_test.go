package hmsLicenter

import (
	"testing"
	"time"
)

func TestBll(t *testing.T) {
	repo := NewRepo("127.0.0.1", "xky.seckiller", time.Now().UnixNano(), "aabbcc11", time.Now(),
		time.Now().Format("20060102150405")).WithEncryptKey([]byte("xxxxxxxx"))
	repo.Start()
	time.Sleep(time.Second * 30)
}

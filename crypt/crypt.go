package crypt

import (
	"fmt"
	"github.com/xukgo/gsaber/encrypt/aes"
)

func EncryptClientData(data []byte, ckey []byte) ([]byte, error) {
	if len(data) == 0 {
		return data, nil
	}
	return aes.Encrypt(data, ckey), nil
}

func DecryptServerData(data []byte, ckey []byte) (res []byte, err error) {
	if len(data) == 0 {
		return data, nil
	}

	res = nil
	err = fmt.Errorf("aes.Decrypt error")

	defer func() {
		//recover() //可以打印panic的错误信息
		//fmt.Println(recover())
		if err := recover(); err != nil { //产生了panic异常
			fmt.Println(err)
		}

	}() //别忘了(), 调用此匿名函数

	res = aes.Decrypt(data, ckey)
	err = nil
	return
}

package rsa

import (
	"fmt"
	"testing"
)

func TestRsa(t *testing.T) {
	//RsaGenKey(2048)
	var data = []byte("hello world")
	encrypt := RSAEncrypt(data, []byte("public.pem"))
	fmt.Println("加密后的数据:", string(encrypt))
	decrypt := RSADecrypt(encrypt, []byte("private.pem"))
	fmt.Println("解密后的数据:", string(decrypt))
}

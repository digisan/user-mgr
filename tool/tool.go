package tool

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
	"os"
	"sync"

	fd "github.com/digisan/gotk/filedir"
	lk "github.com/digisan/logkit"
)

var (
	commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
	mu       sync.Mutex
)

func Maxlen(s string, length int) string {
	if len(s) < length {
		return s
	}
	return s[:length]
}

func Encrypt(plain string, key []byte) []byte {

	mu.Lock()
	defer mu.Unlock()

	c, err := aes.NewCipher(key)
	if err != nil {
		fmt.Printf("Error: NewCipher(%d bytes) = %s", len(key), err)
		os.Exit(-1)
	}

	cfb := cipher.NewCFBEncrypter(c, commonIV)
	cipherBuf := make([]byte, len(plain))
	cfb.XORKeyStream(cipherBuf, []byte(plain))
	// fmt.Printf("%s => %x\n", []byte(plain), cipherBuf)

	return cipherBuf // fmt.Sprintf("%x", cipherBuf)
}

func Decrypt(cipherBuf, key []byte) string {

	mu.Lock()
	defer mu.Unlock()

	c, err := aes.NewCipher(key)
	if err != nil {
		fmt.Printf("Error: NewCipher(%d bytes) = %s", len(key), err)
		os.Exit(-1)
	}

	cfbdec := cipher.NewCFBDecrypter(c, commonIV)
	plainBuf := make([]byte, 1024)
	cfbdec.XORKeyStream(plainBuf, cipherBuf)
	plainBuf = bytes.TrimRight(plainBuf, "\x00")
	// fmt.Printf("%x => %s\n", cipherBuf, plainBuf)
	return string(plainBuf)
}

// h : [md5.New() / sha1.New() / sha256.New()]
func FileHash(file string, h hash.Hash) string {
	if !fd.FileExists(file) {
		return ""
	}
	f, err := os.Open(file)
	lk.FailOnErr("%v", err)
	defer f.Close()
	_, err = io.Copy(h, f)
	lk.FailOnErr("%v", err)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func SelfMD5() string {
	return FileHash(os.Args[0], md5.New())
}

func SelfSHA1() string {
	return FileHash(os.Args[0], sha1.New())
}

func SelfSHA256() string {
	return FileHash(os.Args[0], sha256.New())
}

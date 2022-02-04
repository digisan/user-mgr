package user

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
	"os"

	fd "github.com/digisan/gotk/filedir"
	lk "github.com/digisan/logkit"
)

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

package tool

import (
	"fmt"
	"testing"
	"time"
)

func TestEnDe(t *testing.T) {

	original := "AA"
	key := []byte(fmt.Sprintf("%d", time.Now().UnixNano())[3:19])

	secret := Encrypt(original, key)
	// fmt.Println(secret)
	// fmt.Println(string(secret))
	// fmt.Printf("%x\n", secret)

	plain := Decrypt(secret, key)
	// fmt.Println(plain)

	if plain == original {
		fmt.Println("OK")
	} else {
		fmt.Println("ERROR")
	}
}

func TestMaxlen(t *testing.T) {
	fmt.Println(Maxlen("abcdef", 6))
	fmt.Println(Maxlen("abcdefg", 6))
	fmt.Println(Maxlen("abcde", 6))
}

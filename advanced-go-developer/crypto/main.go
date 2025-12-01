package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func main() {
	fmt.Println(randomBase64String(8))
}

func randomBase64String(l int) string {
	b := make([]byte, l)
	_, _ = rand.Read(b)
	return base64.RawStdEncoding.EncodeToString(b)
}

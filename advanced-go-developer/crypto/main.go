package main

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log/slog"
	"os"
)

func main() {
	var (
		data  = make([]byte, 512) // слайс случайных байт
		hash1 []byte              // хеш с использованием интерфейса hash.Hash
		hash2 [md5.Size]byte      // хеш, возвращаемый функцией md5.Sum
	)

	_, _ = rand.Read(data)

	h := md5.New()
	_, err := h.Write(data)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	hash1 = h.Sum(hash1)
	hash2 = md5.Sum(data)

	// hash2[:] приводит массив байт к слайсу
	if bytes.Equal(hash1, hash2[:]) {
		fmt.Println("Всё правильно! Хеши равны")
	} else {
		fmt.Println("Что-то пошло не так")
	}
}

func randomBase64String(l int) string {
	b := make([]byte, l)
	_, _ = rand.Read(b)
	return base64.RawStdEncoding.EncodeToString(b)
}

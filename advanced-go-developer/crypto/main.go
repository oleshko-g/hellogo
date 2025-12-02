package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log/slog"
)

var secretkey = []byte("secret key")

func main() {
	var (
		data = make([]byte, 0) // декодированное сообщение с подписью
		id   uint32            // значение идентификатора
		err  error
		sign []byte // HMAC-подпись от идентификатора
	)

	msg := "048ff4ea240a9fdeac8f1422733e9f3b8b0291c969652225e25c5f0f9f8da654139c9e21"
	// 1) декодируйте msg в data
	data = make([]byte, hex.DecodedLen(len(msg)))
	_, err = hex.Decode(data, []byte(msg))
	if err != nil {
		slog.Error(err.Error())
	}
	id = binary.BigEndian.Uint32(data[:4])
	// 3) вычислите HMAC-подпись sign для этих четырёх байт
	h := hmac.New(sha256.New, secretkey)
	h.Write(data[:4])
	sign = h.Sum(nil)

	if hmac.Equal(sign, data[4:]) {
		fmt.Println("Подпись подлинная. ID:", id)
	} else {
		fmt.Println("Подпись неверна. Где-то ошибка")
	}
}

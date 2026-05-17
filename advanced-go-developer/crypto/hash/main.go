package main

import (
    "bytes"
    "crypto/md5"
    "crypto/rand"
    "fmt"
)

func main() {
    var (
        data  []byte         // слайс случайных байт
        hash1 []byte         // хеш с использованием интерфейса hash.Hash
        hash2 [md5.Size]byte // хеш, возвращаемый функцией md5.Sum
    )
    // допишите код
    // 1) сгенерируйте data длиной 512 байт
    data = make([]byte, 512)
    rand.Read(data)

    // 2) вычислите hash1 с использованием md5.New
    h := md5.New()
    h.Write(data)
    hash1 = h.Sum(nil)

    // 3) вычислите hash2 функцией md5.Sum
    hash2 = md5.Sum(data)

    // hash2[:] приводит массив байт к слайсу
    if bytes.Equal(hash1, hash2[:]) {
    		fmt.Println(hash1)
    		fmt.Println(hash2)
        fmt.Println("Всё правильно! Хеши равны")
    } else {
        fmt.Println("Что-то пошло не так")
    }
}

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	// Подключение к сетевой службе.
	conn, err := net.Dial("tcp4", "localhost:12345")
	if err != nil {
		log.Fatal(err)
	}
	// Не забываем закрыть ресурс.
	defer conn.Close()

	// Прервет соединение через 15 секунд
	go func() {
		time.Sleep(time.Second * 15)
		_, err = conn.Write([]byte("exit\n"))
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Буфер для чтения данных из соединения.
	reader := bufio.NewReader(conn)
	// Считывание массива байт до перевода строки.
	for {
		b, err := reader.ReadBytes('\n')
		if err != nil {
			log.Fatal(err)
		}
		// Обработка ответа.
		fmt.Println("Ответ от сервера:", string(b))
	}
}

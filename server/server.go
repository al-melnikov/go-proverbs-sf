package main

import (
	"bufio"
	"log"
	"math/rand"
	"net"
	"strings"
	"time"
)

const addr = "0.0.0.0:12345"

// Протокол сетевой службы.
const proto = "tcp4"

func main() {
	// Запуск сетевой службы по протоколу TCP
	// на порту 12345.
	listener, err := net.Listen(proto, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	// Фразы берутся из общего канала.
	// Это нужно чтобы при разных подключениях они не совпадали
	chProverbs := make(chan string)
	go fillCh(getProverbs(), chProverbs)

	// Подключения обрабатываются в бесконечном цикле.
	// Иначе после обслуживания первого подключения сервер
	//завершит работу.
	for {
		// Принимаем подключение.
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		// Вызов обработчика подключения.
		go handleConn(conn, chProverbs)
	}
}

// Обработчик. Вызывается для каждого соединения.
// Фразы берутся из канала, он один общий поэтому передается как аргумент
func handleConn(conn net.Conn, ch <-chan string) {

	// Печатает фразу из канала раз в 3 секунды
	go func() {
		for {
			conn.Write([]byte(<-ch + "\n"))
			time.Sleep(time.Second * 3)
		}
	}()

	// Закрытие соединения.
	// Чтение сообщения от клиента.
	reader := bufio.NewReader(conn)
	for {
		b, err := reader.ReadBytes('\n')
		if err != nil {
			log.Println(err)
			return
		}
		// Удаление символов конца строки.
		msg := strings.TrimSuffix(string(b), "\n")
		msg = strings.TrimSuffix(msg, "\r")

		// соединение закроется если поступят такие сообщения
		if msg == "quit" || msg == "exit" || msg == "q" {
			conn.Close()
		}
	}
}

// массив строк нельзя объявить константой, поэтому так
func getProverbs() []string {
	return []string{
		"Don't communicate by sharing memory, share memory by communicating.",
		"Concurrency is not parallelism.",
		"Channels orchestrate; mutexes serialize.",
		"The bigger the interface, the weaker the abstraction.",
		"Make the zero value useful.",
		"interface{} says nothing.",
		"Gofmt's style is no one's favorite, yet gofmt is everyone's favorite.",
		"A little copying is better than a little dependency.",
		"Syscall must always be guarded with build tags.",
		"Cgo must always be guarded with build tags.",
		"Cgo is not Go.",
		"With the unsafe package there are no guarantees.",
		"Clear is better than clever.",
		"Reflection is never clear.",
		"Errors are values.",
		"Don't just check errors, handle them gracefully.",
		"Design the architecture, name the components, document the details.",
		"Documentation is for users.",
		"Don't panic.",
	}
}

func fillCh(arr []string, ch chan<- string) {
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s) // initialize local pseudorandom generator
	for {
		ch <- arr[r.Intn(len(arr))]
	}
}

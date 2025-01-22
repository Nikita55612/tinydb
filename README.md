# tinydb
Мой первый проект на Go

SQLite база данных работающая через TCP протокол

Usage of tinydb.exe:
  -b uint
        Max conn buff size (default 8388608)
  -p uint
        Port (default 4004)

Сообщение: [hash][data]

```go
package main

import (
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"testing"
)

const Addr = ":4000"

func TestW(t *testing.T) {
	// Первые 32 байта пустые для записи данных
	slice := make([]byte, 32)

	conn, err := net.Dial("tcp", Addr)
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Сообщение для записи данных
	mess := append(slice, []byte("Типо данные")...)
	_, err = conn.Write(mess)
	if err != nil {
		t.Fatalf("Failed to send data: %v", err)
	}

	// Получение hash данных для чтения
	buff := make([]byte, 1024)
	n, err := conn.Read(buff)
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}

	fmt.Printf("Received %d bytes: %v\n", n, buff[:n])

	// hash в hex encode
	fmt.Printf("Hex: %s\n", hex.EncodeToString(buff[:n]))

	t.Logf("Received %d bytes: %v", n, buff[:n])
}

func TestR(t *testing.T) {
	// hash данных
	hash, err := hex.DecodeString("a7cf70c3b3e47ca7278f623f50ab59b02447ae4c69c25700a26db7b647d03fd6")
	if err != nil {
		t.Fatalf("%v", err)
	}

	conn, err := net.Dial("tcp", Addr)
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Читаем данные по hash
	_, err = conn.Write(hash)
	if err != nil {
		t.Fatalf("Failed to send data: %v", err)
	}

	buff := make([]byte, 1024)
	n, err := conn.Read(buff)
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}

	// Полученные данные
	fmt.Printf("Received %d bytes from data: %s\n", n, string(buff[32:n]))

	t.Logf("Received %d bytes: %v", n, buff[:n])
}

func TestReplase(t *testing.T) {
	// hash данных
	hash, err := hex.DecodeString("a7cf70c3b3e47ca7278f623f50ab59b02447ae4c69c25700a26db7b647d03fd6")
	if err != nil {
		t.Fatalf("%v", err)
	}

	conn, err := net.Dial("tcp", Addr)
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Перезапись данных
	_, err = conn.Write(append(hash, []byte("Новые данные")...))
	if err != nil {
		t.Fatalf("Failed to send data: %v", err)
	}

	buff := make([]byte, 1024)
	n, err := conn.Read(buff)
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}

	// Новый hash
	fmt.Printf("New hash: %s\n", string(buff[:32]))

	t.Logf("Received %d bytes: %v", n, buff[:n])
}

func TestWImg(t *testing.T) {
	// Запись картинки
	data, err := os.ReadFile("img.jpg")
	if err != nil {
		fmt.Println(err)
		return
	}
	slice := make([]byte, 32)

	conn, err := net.Dial("tcp", Addr)
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	mess := append(slice, data...)
	_, err = conn.Write(mess)
	if err != nil {
		t.Fatalf("Failed to send data: %v", err)
	}

	buff := make([]byte, 1024)
	n, err := conn.Read(buff)
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}

	fmt.Printf("Received %d bytes: %v\n", n, buff[:n])
	fmt.Printf("Hex: %s\n", hex.EncodeToString(buff[:n]))

	t.Logf("Received %d bytes: %v", n, buff[:n])
}

```

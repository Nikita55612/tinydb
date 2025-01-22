package main

import (
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

const Addr = ":4000"

func TestW(t *testing.T) {
	slice := make([]byte, 32)

	conn, err := net.Dial("tcp", Addr)
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	mess := append(slice, []byte("0000000000123")...)
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

func TestR(t *testing.T) {
	hash, err := hex.DecodeString("a7cf70c3b3e47ca7278f623f50ab59b02447ae4c69c25700a26db7b647d03fd6")
	if err != nil {
		t.Fatalf("%v", err)
	}

	conn, err := net.Dial("tcp", Addr)
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	_, err = conn.Write(hash)
	if err != nil {
		t.Fatalf("Failed to send data: %v", err)
	}

	buff := make([]byte, 1024)
	n, err := conn.Read(buff)
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}

	fmt.Printf("Received %d bytes from data: %s\n", n, string(buff[32:n]))

	t.Logf("Received %d bytes: %v", n, buff[:n])
}

func TestWImg(t *testing.T) {
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

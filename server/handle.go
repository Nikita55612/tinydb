package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"lukechampine.com/blake3"
)

var (
	MaxBuffSize uint64
	EmptyHash   = [32]byte{}
)

func handleAcceptConn(listener net.Listener, ctx context.Context) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				log.Println("Server shutdown complete")
				return
			default:
				log.Println("Error accepting connection:", err)
				continue
			}
		}

		go handleConn(conn)
	}
}

func hashFromScan(scan []byte) [32]byte {
	var hash [32]byte
	copy(hash[:], scan[:32])
	return hash
}

func handleConn(conn net.Conn) {
	remAddr := conn.RemoteAddr().String()
	log.Printf("[INFO] Connection established from %s", remAddr)
	defer conn.Close()

	buff := make([]byte, MaxBuffSize+1)
	for {
		n, err := conn.Read(buff)
		if err != nil {
			fmt.Printf("[ERROR] Failed to read from connection %s: %v\n", remAddr, err)
			break
		}
		if uint64(n) > MaxBuffSize {
			fmt.Println("[ERROR] Message exceeded maximum allowed size")
			break
		}

		scan := buff[:n]
		if len(scan) < 32 {
			fmt.Printf("[WARN] Received undersized message from %s: %d bytes\n", remAddr, len(scan))
			continue
		}

		log.Printf("[INFO] Processing message from %s: %d bytes", remAddr, n)

		hash := hashFromScan(scan)
		response := hash[:]

		switch {
		case hash != EmptyHash && len(scan) > 32:
			data := scan[32:]
			new_hash := blake3.Sum256(data)
			if err := Replace(hash[:], new_hash[:], data); err != nil {
				fmt.Printf("[ERROR] Failed to replace data for hash %x: %v\n", hash, err)
				continue
			}
			fmt.Printf("[INFO] Successfully updated data for hash %x\n", hash)
		case hash != EmptyHash:
			data, err := Select(hash[:])
			if err != nil {
				fmt.Printf("[ERROR] Failed to select data for hash %x: %v\n", hash, err)
				continue
			}
			fmt.Printf("[INFO] Successfully retrieved data for hash %x\n", hash)
			response = append(response, data...)
		case len(scan) > 32:
			data := scan[32:]
			dataHash := blake3.Sum256(data)
			if err := Insert(dataHash[:], data); err != nil {
				fmt.Printf("[ERROR] Failed to insert data with hash %x: %v\n", dataHash, err)
				continue
			}
			fmt.Printf("[INFO] Successfully inserted data with hash %x\n", dataHash)
			copy(response, dataHash[:])
		}

		if n, err = conn.Write(response); err != nil {
			fmt.Printf("[ERROR] Failed to write response to %s: %v\n", remAddr, err)
			break
		}
		if n != len(response) {
			fmt.Printf("[ERROR] Incomplete write to %s: sent %d of %d bytes\n", remAddr, n, len(response))
			break
		}
	}
}

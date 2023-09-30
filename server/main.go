package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	// hash "incremental-backup-system/cmd"
)

// type message struct {
// 	ID      int    `json:"id"`
// 	Type    string `json:"type"`
// 	Message string `json:"message"`
// }

type Request struct {
	DirPath       string
	IsFirstBackup string
	SaveHistory   string
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Deserializa o JSON de volta para a estrutura Request
		var request Request
		err = json.Unmarshal(buf[:n], &request)
		if err != nil {
			fmt.Println("Erro ao deserializar a requisição:", err)
			continue
		}

		fmt.Printf("Diretório: %s\nPrimeiro Backup: %v\nSave History: %s", request.DirPath, request.IsFirstBackup, request.SaveHistory)

	}
}

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide port number")
		return
	}

	PORT := ":" + arguments[1]
	ln, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		go handleClient(conn)
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/fatih/color"
)

type Request struct {
	DirPath       string
	IsFirstBackup string
	SaveHistory   string
}

func handleClient(conn net.Conn) {
	color.Green("\nConexão estabelecida!\n\n")
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

		fmt.Printf("Diretório: %s\nPrimeiro Backup: %v\nSave History: %s\n", request.DirPath, request.IsFirstBackup, request.SaveHistory)

	}
}

func main() {
	if len(os.Args) == 1 {
		color.Yellow("Port not provided, for default port use 6677")
		os.Args = append(os.Args, "6677")
	}

	PORT := ":" + os.Args[1]
	ln, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Server is running on port ")
	color.Green(PORT)

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

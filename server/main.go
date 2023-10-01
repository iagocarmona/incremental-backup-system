package main

import (
	"encoding/json"
	"fmt"
	hash "incremental-backup-system/cmd"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

type Request struct {
	DirPath       string
	IsFirstBackup string
	SaveHistory   string
}

var hashServer map[string]hash.DirEntry

func getFileSize(conn net.Conn) (int64, error) {
	//receber o tamanho do arquivo em 8 bytes
	buffer := make([]byte, 8)
	n, err := conn.Read(buffer)
	if err != nil {
		color.Red("Erro ao receber o tamanho do arquivo: %v\n", err)
	}

	// Converte o tamanho do arquivo para int64
	fileSize := int64(0)
	for i := 0; i < n; i++ {
		fileSize += int64(buffer[i]) << uint(8*(7-i))
	}

	return fileSize, err
}

func getRelativePath(conn net.Conn, backupPath string) (string, *os.File, error) {
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		color.Red("Erro ao receber o path do arquivo: %v\n", err)
		return "", nil, err
	}

	fileName := string(buffer[:n])

	// Cria o diretório pai, se necessário
	err = os.MkdirAll(filepath.Dir(backupPath+fileName), os.ModePerm)
	if err != nil {
		color.Red("Erro ao criar o diretório pai:", err)
		return "", nil, err
	}

	// Cria o arquivo no diretório de backup
	file, err := os.Create(backupPath + fileName)
	if err != nil {
		color.Red("Erro ao criar o arquivo: %v\n", err)
		return "", nil, err
	}

	return fileName, file, err
}

func getFileContent(conn net.Conn, fileSize int64, file *os.File) ([]byte, error) {
	fileContent := make([]byte, fileSize)

	bufferSize := 1024
	buffer := make([]byte, bufferSize)

	err := error(nil)

	for bytesRead := int64(0); bytesRead < fileSize; {
		// Verifique se o tamanho restante é menor que o tamanho do buffer
		if remaining := fileSize - bytesRead; remaining < int64(bufferSize) {
			bufferSize = int(remaining)
			buffer = make([]byte, bufferSize)
		}

		// Leia um cluster de dados da conexão
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break // Chegou ao fim da transferência
			}
			fmt.Println("Erro ao ler da conexão:", err)
		}

		// Copie os dados lidos para o slice 'fileContent'
		copy(fileContent[bytesRead:], buffer[:n])

		// Atualize o número de bytes lidos
		bytesRead += int64(n)
	}

	return fileContent, err
}

func getHashClient(conn net.Conn) (map[string]hash.DirEntry, error) {
	// Pega o tamanho do buffer
	bufferSize, err := getFileSize(conn)
	if err != nil {
		return nil, err
	}

	// Envia uma confirmação de 1 byte para o cliente
	_, err = conn.Write([]byte{1})
	if err != nil {
		color.Red("Erro ao enviar confirmação: %v\n", err)
		return nil, err
	}

	// Cria um buffer para armazenar os dados recebidos
	buffer := make([]byte, bufferSize)

	// Lê os dados do buffer
	_, err = io.ReadFull(conn, buffer)
	if err != nil {
		color.Red("Erro ao ler os dados do buffer: %v\n", err)
		return nil, err
	}

	// Agora você pode desserializar os dados para hash.HashTable
	var hashClient map[string]hash.DirEntry
	err = json.Unmarshal(buffer, &hashClient)
	if err != nil {
		color.Red("Erro ao desserializar os dados para hash.HashTable: %v\n", err)
		return nil, err
	}

	return hashClient, nil
}

func receiveFiles(conn net.Conn, backupPath string, isFirstBackup bool, dirPath string) {
	// =============================================================================
	//                                 HASH CLIENT
	// =============================================================================

	// Pega a hash do cliente
	hashClient, err := getHashClient(conn)
	if err != nil {
		return
	}

	if isFirstBackup {
		hashServer = hashClient
	}

	// envia uma confirmação de 1 byte para o cliente
	_, err = conn.Write([]byte{1})
	if err != nil {
		color.Red("Erro ao enviar confirmação: %v\n", err)
		return
	}

	for {
		// =============================================================================
		//  							 IF NOT FIRST BACKUP
		// =============================================================================

		if !isFirstBackup {
			// faz a comparação entre as duas hash
			toUpdateList := hash.Diff(hashServer, hashClient)

			// envia para o cliente a lista em bytes
			toUpdateListBytes, err := json.Marshal(toUpdateList)
			if err != nil {
				color.Red("Erro ao serializar a lista: %v\n", err)
				return
			}

			// Calcula o tamanho total toUpdateListBytes
			totalSize := len(toUpdateListBytes)

			// Converte o tamanho total em 8 bytes
			totalSizeBytes := make([]byte, 8)
			for i := 0; i < 8; i++ {
				totalSizeBytes[i] = byte(totalSize >> uint(8*(7-i)))
			}

			_, err = conn.Write(totalSizeBytes)
			if err != nil {
				color.Red("Erro ao enviar o tamanho da lista: %v\n", err)
				return
			}

			// Recebe uma confirmação de 1 byte do cliente
			buffer := make([]byte, 1)
			_, err = conn.Read(buffer)
			if err != nil {
				color.Red("Erro ao receber confirmação: %v\n", err)
				return
			}

			_, err = conn.Write(toUpdateListBytes)
			if err != nil {
				color.Red("Erro ao enviar a lista: %v\n", err)
				return
			}

			//  Recebe uma confirmação de 1 byte do cliente
			buffer = make([]byte, 1)
			_, err = conn.Read(buffer)
			if err != nil {
				color.Red("Erro ao receber confirmação: %v\n", err)
				return
			}
		}

		// =============================================================================
		//                                 FILESIZE
		// =============================================================================

		fileSize, err := getFileSize(conn)
		if err != nil {
			return
		}

		// envia uma confirmação de 1 byte para o cliente
		_, err = conn.Write([]byte{1})
		if err != nil {
			color.Red("Erro ao enviar confirmação: %v\n", err)
			return
		}

		// =============================================================================
		//                                 RELATIVE PATH
		// =============================================================================

		fileName, file, err := getRelativePath(conn, backupPath)
		if err != nil {
			return
		}

		fmt.Printf("Recebendo arquivo %s...\n", fileName)

		// envia uma confirmação de 1 byte para o cliente
		_, err = conn.Write([]byte{1})
		if err != nil {
			color.Red("Erro ao enviar confirmação: %v\n", err)
			return
		}

		// =============================================================================
		//                                 FILE CONTENT
		// =============================================================================

		fileContent, err := getFileContent(conn, fileSize, file)
		if err != nil {
			return
		}

		// Escreve o conteúdo do arquivo
		_, err = file.Write(fileContent)
		if err != nil {
			color.Red("Erro ao escrever o arquivo: %v\n", err)
			return
		}

		// Fecha o arquivo
		err = file.Close()
		if err != nil {
			color.Red("Erro ao fechar o arquivo: %v\n", err)
			return
		}

		// Envia uma confirmação ao cliente
		_, err = conn.Write([]byte{1})
		if err != nil {
			color.Red("Erro ao enviar confirmação: %v\n", err)
			return
		}

		fmt.Printf("Arquivo %s recebido com sucesso!\n", fileName)
	}
}

func removeMiddleDots(input string) string {
	parts := strings.Split(input, "/")
	output := []string{}

	for _, part := range parts {
		if part == ".." {
			if len(output) > 0 {
				output = output[:len(output)-1]
			}
		} else {
			output = append(output, part)
		}
	}

	return strings.Join(output, "/")
}

func handleClient(conn net.Conn) {
	color.Green("\nConexão estabelecida!\n\n")
	defer conn.Close()

	for {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Deserializa o JSON de volta para a estrutura Request
		var request Request
		err = json.Unmarshal(buffer[:n], &request)
		if err != nil {
			fmt.Println("Erro ao deserializar a requisição:", err)
			continue
		}

		fmt.Printf("Diretório: %s\nPrimeiro Backup: %v\nSave History: %s\n", request.DirPath, request.IsFirstBackup, request.SaveHistory)

		// Cria o diretório de backup
		backupPath := request.DirPath + "/"
		backupPath = removeMiddleDots(backupPath)
		backupPath = "backups/" + backupPath

		err = os.MkdirAll(backupPath, os.ModePerm)
		if err != nil {
			fmt.Println("Erro ao criar o diretório de backup:", err)
			return
		}

		// Verifica se é o primeiro backup
		if request.IsFirstBackup == "true" {
			// Envia uma confirmação ao cliente
			_, err = conn.Write([]byte{1})
			if err != nil {
				fmt.Println("Erro ao enviar confirmação:", err)
				return
			}

			receiveFiles(conn, backupPath, true, request.DirPath)
		} else {
			fmt.Println("Não é o primeiro backup")
			// Envia uma confirmação ao cliente
			_, err = conn.Write([]byte{1})
			if err != nil {
				fmt.Println("Erro ao enviar confirmação:", err)
				return
			}

			receiveFiles(conn, backupPath, false, request.DirPath)
		}
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

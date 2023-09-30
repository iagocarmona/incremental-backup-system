package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/viper"
)

type Request struct {
	DirPath       string
	IsFirstBackup string
	SaveHistory   string
}

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

func receiveFiles(conn net.Conn, backupPath string) {
	for {
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

			receiveFiles(conn, backupPath)
		} else {
			fmt.Println("Não é o primeiro backup")
			// Envia uma confirmação ao cliente
			_, err = conn.Write([]byte{1})
			if err != nil {
				fmt.Println("Erro ao enviar confirmação:", err)
				return
			}

			receiveFiles(conn, backupPath)
		}
	}
}

func resetConfig() {
	// Configure o Viper para usar o local storage
	viper.SetConfigName("config") // Nome do arquivo de configuração (ex: config.yaml)
	viper.AddConfigPath("../")    // Diretório onde o arquivo de configuração está localizado
	viper.SetConfigType("yaml")   // Tipo do arquivo de configuração (ex: YAML)

	// Ler o arquivo de configuração
	if err := viper.ReadInConfig(); err != nil {
		color.Red("Erro ao ler o arquivo de configuração: %v\n", err)
		return
	}

	// Criar a variável isFirstBackup no local storage
	viper.Set("isFirstBackup", true)

	// Salvar a variável isFirstBackup no local storage
	if err := viper.WriteConfig(); err != nil {
		color.Red("Erro ao salvar o arquivo de configuração: %v\n", err)
		return
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

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	hash "incremental-backup-system/cmd"
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

const (
	CLUSTERSIZE = 1024 // Tamanho do cluster em bytes
)

func verifyIsFirstBackup() bool {
	// Configure o Viper para usar o local storage
	viper.SetConfigName("config") // Nome do arquivo de configuração (ex: config.yaml)
	viper.AddConfigPath("../")    // Diretório onde o arquivo de configuração está localizado
	viper.SetConfigType("yaml")   // Tipo do arquivo de configuração (ex: YAML)

	// Ler o arquivo de configuração
	if err := viper.ReadInConfig(); err != nil {
		color.Red("Erro ao ler o arquivo de configuração: %v\n", err)
		return false
	}

	return viper.GetBool("isFirstBackup")
}

func getHostAndPort() string {
	// Configure o Viper para usar o local storage
	viper.SetConfigName("config") // Nome do arquivo de configuração (ex: config.yaml)
	viper.AddConfigPath("../")    // Diretório onde o arquivo de configuração está localizado
	viper.SetConfigType("yaml")   // Tipo do arquivo de configuração (ex: YAML)

	// Ler o arquivo de configuração
	if err := viper.ReadInConfig(); err != nil {
		color.Red("Erro ao ler o arquivo de configuração: %v\n", err)
		return ""
	}

	return viper.GetString("host") + ":" + viper.GetString("port")
}

func createConfig() {
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
	viper.Set("isFirstBackup", false)

	// Salvar a variável isFirstBackup no local storage
	if err := viper.WriteConfig(); err != nil {
		color.Red("Erro ao salvar o arquivo de configuração: %v\n", err)
		return
	}
}

func printHeader() {
	color.White("\n==================================================================")
	color.White("                  INCREMENTAL BACKUP SYSTEM\n\n")
	color.Blue("> Informe o diretório e se deseja salvar histórico dos arquivos: ")
	color.Green("Exemplo: /home/user/backup true")
	color.White("==================================================================\n\n")
	fmt.Print("> ")
}

func sendHashToServer(request Request, conn net.Conn) error {
	//monta a hash local
	localHash := hash.CreateLocalHash(request.DirPath)

	// Serializa a estrutura Request em JSON
	requestJSON, err := json.Marshal(localHash)
	if err != nil {
		color.Red("Erro ao serializar a requisição em JSON:", err)
		return err
	}

	// Obtém o tamanho da hash serializada em bytes
	hashSize := len(requestJSON)

	// converte o tamanho da hash em 8 bytes
	hashSizeBytes := make([]byte, 8)
	for i := 0; i < 8; i++ {
		hashSizeBytes[i] = byte((hashSize >> uint(8*(7-i))) & 0xff)
	}

	// envia para o servidor o tamanho da hash em 8 bytes
	_, err = conn.Write(hashSizeBytes)
	if err != nil {
		color.Red("Erro ao enviar o tamanho da hash:", err)
		return err
	}

	// servidor responde com 1 byte
	confirmationSize := make([]byte, 1)
	_, err = conn.Read(confirmationSize)
	if err != nil {
		color.Red("Erro ao receber confirmação do servidor:", err)
		return err
	}

	// envia para o servidor a hash
	_, err = conn.Write(requestJSON)
	if err != nil {
		color.Red("Erro ao enviar a requisição para o servidor:", err)
		return err
	}

	// recebo uma confirmação do servidor em 1 byte
	confirmation := make([]byte, 1)
	_, err = conn.Read(confirmation)
	if err != nil {
		color.Red("Erro ao receber confirmação do servidor:", err)
		return err
	}

	return nil
}

func sendIncremenalBackup(request Request, conn net.Conn) error {
	// Envia a hash para o servidor
	err := sendHashToServer(request, conn)
	if err != nil {
		color.Red("Erro ao enviar a hash para o servidor:", err)
		return err
	}

	// Recebe o tamanho da lista de arquivos em 8 bytes
	listSizeBytes := make([]byte, 8)
	_, err = conn.Read(listSizeBytes)
	if err != nil {
		color.Red("Erro ao receber o tamanho da lista de arquivos:", err)
		return err
	}

	// Envia uma confirmação de 1 byte para o servidor
	_, err = conn.Write([]byte{1})
	if err != nil {
		color.Red("Erro ao enviar confirmação:", err)
		return err
	}

	// Converte os 8 bytes para um valor int64 que representa o tamanho da lista
	listSize := int64(0)
	for i := 0; i < 8; i++ {
		listSize += int64(listSizeBytes[i]) << uint(8*(7-i))
	}

	// Recebe a lista de arquivos em JSON
	listBytes := make([]byte, listSize)
	_, err = conn.Read(listBytes)
	if err != nil {
		color.Red("Erro ao receber a lista de arquivos:", err)
		return err
	}

	// Envia uma confirmação de 1 byte para o servidor
	_, err = conn.Write([]byte{1})
	if err != nil {
		color.Red("Erro ao enviar confirmação:", err)
		return err
	}

	// Converte a lista de arquivos em JSON para uma lista de strings
	list := []string{}
	err = json.Unmarshal(listBytes[:listSize], &list)
	if err != nil {
		color.Red("Erro ao converter a lista de arquivos em JSON:", err)
		return err
	}

	// Percorre a lista de arquivos
	for _, fileName := range list {
		if fileName != "" {
			// =============================================================================
			//                                 FILESIZE
			// =============================================================================

			// Abre o arquivo
			file, err := os.Open(fileName)
			if err != nil {
				color.Red("Erro ao abrir o arquivo:", err)
				return err
			}
			defer file.Close()

			// Obtém informações do arquivo
			fileInfo, err := file.Stat()
			if err != nil {
				color.Red("Erro ao obter informações do arquivo:", err)
				return err
			}

			// converte o tamanho do arquivo em 8 bytes
			fileSize := make([]byte, 8)
			for i := 0; i < 8; i++ {
				fileSize[i] = byte((fileInfo.Size() >> uint(8*(7-i))) & 0xff)
			}

			// envia para o servidor o tamanho do arquivo em 8 bytes
			_, err = conn.Write(fileSize)
			if err != nil {
				color.Red("Erro ao enviar o tamanho do arquivo:", err)
				return err
			}

			// servidor responde com 1 byte
			confirmationSize := make([]byte, 1)
			_, err = conn.Read(confirmationSize)
			if err != nil {
				color.Red("Erro ao receber confirmação do servidor:", err)
				return err
			}

			// =============================================================================
			//                          	 RELATIVE PATH
			// =============================================================================

			// Calcula o caminho relativo do arquivo
			relativePath, err := filepath.Rel(request.DirPath, fileName)
			if err != nil {
				color.Red("Erro ao calcular o caminho relativo do arquivo:", err)
				return err
			}

			// envia para o servidor o nome do arquivo
			_, err = conn.Write([]byte(relativePath))
			if err != nil {
				color.Red("Erro ao enviar o nome do arquivo:", err)
				return err
			}

			// aguarda uma confirmação do servidor antes de enviar o conteúdo do arquivo
			confirmationName := make([]byte, 1)
			_, err = conn.Read(confirmationName)
			if err != nil {
				color.Red("Erro ao receber confirmação do servidor:", err)
				return err
			}

			// =============================================================================
			//                          	 FILE CONTENT
			// =============================================================================

			// Lê o conteúdo do arquivo em um slice de bytes
			fileContent, err := io.ReadAll(file)
			if err != nil {
				color.Red("Erro ao ler o conteúdo do arquivo:", err)
				return err
			}

			// envia o conteudo do arquivo a em cluster de 1024 bytes
			for i := 0; i < len(fileContent); i += CLUSTERSIZE {
				end := i + CLUSTERSIZE
				if end > len(fileContent) {
					end = len(fileContent)
				}

				_, err = conn.Write(fileContent[i:end])
				if err != nil {
					color.Red("Erro ao enviar o conteúdo do arquivo:", err)
					return err
				}
			}

			// aguarda uma confirmação do servidor antes de enviar o próximo arquivo
			confirmationContent := make([]byte, 1)
			_, err = conn.Read(confirmationContent)
			if err != nil {
				color.Red("Erro ao receber confirmação do servidor:", err)
				return err
			}

			if confirmationContent[0] != 1 {
				color.Red("Servidor não confirmou o recebimento do arquivo")
				return err
			}

			fmt.Print("Arquivo enviado -> ")
			color.Yellow(relativePath)
		}
	}

	return nil
}

func sendFullBackup(request Request, conn net.Conn) {
	// Envia a hash para o servidor
	err := sendHashToServer(request, conn)
	if err != nil {
		color.Red("Erro ao enviar a hash para o servidor:", err)
		return
	}

	// Percorre recursivamente o diretório e envia os arquivos para o servidor
	err = filepath.Walk(request.DirPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			// =============================================================================
			//                                 FILESIZE
			// =============================================================================

			// Abre o arquivo
			file, err := os.Open(filePath)
			if err != nil {
				color.Red("Erro ao abrir o arquivo:", err)
				return err
			}
			defer file.Close()

			// converte o tamanho do arquivo em 8 bytes
			fileSize := make([]byte, 8)
			for i := 0; i < 8; i++ {
				fileSize[i] = byte((info.Size() >> uint(8*(7-i))) & 0xff)
			}

			// envia para o servidor o tamanho do arquivo em 8 bytes
			_, err = conn.Write(fileSize)
			if err != nil {
				color.Red("Erro ao enviar o tamanho do arquivo:", err)
				return err
			}

			// servidor responde com 1 byte
			confirmationSize := make([]byte, 1)
			_, err = conn.Read(confirmationSize)
			if err != nil {
				color.Red("Erro ao receber confirmação do servidor:", err)
				return err
			}

			// =============================================================================
			//                          	 RELATIVE PATH
			// =============================================================================

			// Calcula o caminho relativo do arquivo
			relativePath, err := filepath.Rel(request.DirPath, filePath)
			if err != nil {
				color.Red("Erro ao calcular o caminho relativo do arquivo:", err)
				return err
			}

			// envia para o servidor o nome do arquivo
			_, err = conn.Write([]byte(relativePath))
			if err != nil {
				color.Red("Erro ao enviar o nome do arquivo:", err)
				return err
			}

			// aguarda uma confirmação do servidor antes de enviar o conteúdo do arquivo
			confirmationName := make([]byte, 1)
			_, err = conn.Read(confirmationName)
			if err != nil {
				color.Red("Erro ao receber confirmação do servidor:", err)
				return err
			}

			// =============================================================================
			//                          	 FILE CONTENT
			// =============================================================================

			// Lê o conteúdo do arquivo em um slice de bytes
			fileContent, err := io.ReadAll(file)
			if err != nil {
				color.Red("Erro ao ler o conteúdo do arquivo:", err)
				return err
			}

			// envia o conteudo do arquivo a em cluster de 1024 bytes
			for i := 0; i < len(fileContent); i += CLUSTERSIZE {
				end := i + CLUSTERSIZE
				if end > len(fileContent) {
					end = len(fileContent)
				}

				_, err = conn.Write(fileContent[i:end])
				if err != nil {
					color.Red("Erro ao enviar o conteúdo do arquivo:", err)
					return err
				}
			}

			// aguarda uma confirmação do servidor antes de enviar o próximo arquivo
			confirmationContent := make([]byte, 1)
			_, err = conn.Read(confirmationContent)
			if err != nil {
				color.Red("Erro ao receber confirmação do servidor:", err)
				return err
			}

			if confirmationContent[0] != 1 {
				color.Red("Servidor não confirmou o recebimento do arquivo")
				return err
			}

			fmt.Print("Arquivo enviado -> ")
			color.Yellow(relativePath)
		}

		return nil
	})

	if err != nil {
		color.Red("Erro ao enviar arquivos:", err)
		return
	}
}

func main() {
	arguments := os.Args
	CONNECT := ""

	if len(arguments) == 1 {
		color.Yellow("host:port não informado. Using from config.yaml")
		CONNECT = getHostAndPort()
	} else {
		CONNECT = arguments[1]
	}

	conn, err := net.Dial("tcp", CONNECT)

	if err != nil {
		fmt.Println(err)
		return
	}

	printHeader()

	for {
		reader := bufio.NewReader(os.Stdin)

		str, _ := reader.ReadString('\n')

		if str == "exit\n" {
			color.White("Saindo...")
			os.Exit(0)
		}

		strArray := strings.Split(str, " ")

		if len(strArray) != 2 {
			color.Red("Erro: Faltam argumentos. Exemplo: <diretório> <bool: salvar_historico>\n\n")
			continue
		} else if len(strArray) > 2 {
			color.Red("Erro: Muitos argumentos. Exemplo: <diretório> <bool: salvar_historico>\n\n")
			continue
		}

		dirPath := strings.TrimSpace(strArray[0])
		saveHistory := strings.TrimSpace(strArray[1])

		if dirPath == ".." {
			dirPath, err = filepath.Abs("..")
			if err != nil {
				color.Red("Erro ao obter o diretório atual:", err)
				continue
			}
		}

		if dirPath == "." {
			dirPath, err = os.Getwd()
			if err != nil {
				color.Red("Erro ao obter o diretório atual:", err)
				continue
			}
		}

		if dirPath == "" {
			color.Red("Erro: Diretório inválido. Exemplo: <diretório> <bool: salvar_historico>\n\n")
			continue
		} else if saveHistory != "true" && saveHistory != "false" {
			color.Red("Erro: Argumento 'salvar_historico' inválido. Exemplo: <diretório> <bool: salvar_historico>\n\n")
			continue
		}

		isFirstBackupString := "false"
		if verifyIsFirstBackup() {
			isFirstBackupString = "true"
		}

		fmt.Print("\nDiretório: " + dirPath)
		fmt.Print("\nSalvar histórico: " + saveHistory)
		fmt.Print("\nPrimeiro backup: " + isFirstBackupString + "\n\n")

		request := Request{
			DirPath:       dirPath,
			IsFirstBackup: isFirstBackupString,
			SaveHistory:   saveHistory,
		}

		if verifyIsFirstBackup() {
			color.Blue("\nPrimeiro backup")
			createConfig()
		} else {
			color.Blue("\nNão é o primeiro backup")
		}

		// Serializa a estrutura Request em JSON
		requestJSON, err := json.Marshal(request)
		if err != nil {
			color.Red("Erro ao serializar a requisição em JSON:", err)
			continue
		}

		_, err = conn.Write(requestJSON)
		if err != nil {
			color.Red("Erro ao enviar a requisição para o servidor:", err)
			continue
		}

		color.White("Requisição enviada para o servidor.")

		// recebo uma confirmação do servidor em 1 byte
		confirmation := make([]byte, 1)
		_, err = conn.Read(confirmation)
		if err != nil {
			color.Red("Erro ao receber confirmação do servidor:", err)
			continue
		}

		if confirmation[0] != 1 {
			color.Red("Servidor não confirmou o recebimento da requisição")
			continue
		}

		color.Green("Servidor confirmou o recebimento da requisição\n\n")

		if request.IsFirstBackup == "true" {
			sendFullBackup(request, conn)
		} else {
			sendIncremenalBackup(request, conn)
		}

		// Se não houve erro, envia uma mensagem de sucesso
		if err != nil {
			color.Red("Erro ao enviar arquivos:", err)
			continue
		} else {
			color.Green("\n -> Backup realizado com sucesso!\n")
			os.Exit(0)
		}
	}
}

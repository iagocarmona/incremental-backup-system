package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/viper"
)

type Request struct {
	DirPath       string
	IsFirstBackup string
	SaveHistory   string
}

func verifyIsFirstBackup() bool {
	// Configure o Viper para usar o local storage
	viper.SetConfigName("config") // Nome do arquivo de configuração (ex: config.yaml)
	viper.AddConfigPath(".")      // Diretório onde o arquivo de configuração está localizado
	viper.SetConfigType("yaml")   // Tipo do arquivo de configuração (ex: YAML)

	// Ler o arquivo de configuração
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Erro ao ler o arquivo de configuração: %v\n", err)
		return false
	}

	return viper.GetBool("isFirstBackup")
}

func getHostAndPort() string {
	// Configure o Viper para usar o local storage
	viper.SetConfigName("config") // Nome do arquivo de configuração (ex: config.yaml)
	viper.AddConfigPath(".")      // Diretório onde o arquivo de configuração está localizado
	viper.SetConfigType("yaml")   // Tipo do arquivo de configuração (ex: YAML)

	// Ler o arquivo de configuração
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Erro ao ler o arquivo de configuração: %v\n", err)
		return ""
	}

	return viper.GetString("host") + ":" + viper.GetString("port")
}

func createConfig() {
	// Configure o Viper para usar o local storage
	viper.SetConfigName("config") // Nome do arquivo de configuração (ex: config.yaml)
	viper.AddConfigPath(".")      // Diretório onde o arquivo de configuração está localizado
	viper.SetConfigType("yaml")   // Tipo do arquivo de configuração (ex: YAML)

	// Ler o arquivo de configuração
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Erro ao ler o arquivo de configuração: %v\n", err)
		return
	}

	// Criar a variável isFirstBackup no local storage
	viper.Set("isFirstBackup", false)

	// Salvar a variável isFirstBackup no local storage
	if err := viper.WriteConfig(); err != nil {
		fmt.Printf("Erro ao salvar o arquivo de configuração: %v\n", err)
		return
	}
}

func printHeader() {
	color.White("\n\n==============================================================\n")
	color.White("Incremental Backup System\n\n")
	color.Blue("Informe o diretório e se deseja salvar histórico dos arquivos: ")
	color.Green("Exemplo: /home/user/backup true\n")
	color.White("==============================================================\n\n")
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
			fmt.Println("Erro ao serializar a requisição em JSON:", err)
			continue
		}

		_, err = conn.Write(requestJSON)
		if err != nil {
			fmt.Println("Erro ao enviar a requisição para o servidor:", err)
			continue
		}

		fmt.Println("Requisição enviada para o servidor.")

		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("->: " + message)

	}
}

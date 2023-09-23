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

func main() {
	// // For status messages
	color.White("Incremental Backup System\n\n") // starting some action/command/function
	// color.Green("Hello World!")             // finished with success
	// color.Red("Hello World!")               // finished with error

	// 1 - Encontrar o diretório (especifiado no comando) na máquina local
	// 2 - Percorrer esse diretório:
	//			Existe algo dentro? (root != nil)
	// 				Percorrer dirEntries[] verificando:

	// 					se (é um file)
	//						se (ainda não existe no servidor) -> copiar para o servidor ( insert_file() )
	//                  	se (existe no servidor) e (foi modificado) -> atualizar no servidor ( update_file() )
	//		                Edge Case: como verificar se um arquivo que foi apagado da máquina local existe no servidor? acho que essa verificação para realizar o delete_file() não ocorre aqui

	//                  se (é um dir)
	// 						se (ainda não existe no servidor) -> criar no servidor ( insert_node() )

	arguments := os.Args

	if len(arguments) == 1 {
		fmt.Println("Please provide host:port.")
		return
	}

	CONNECT := arguments[1]
	conn, err := net.Dial("tcp", CONNECT)

	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		reader := bufio.NewReader(os.Stdin)
		color.Blue("Informe o diretório: ")
		str, _ := reader.ReadString('\n')

		dirPath, saveHistory := strings.Split(str, " ")[0], strings.Split(str, " ")[1]

		fmt.Print("Diretório: " + dirPath)
		fmt.Print("Salvar histórico: " + saveHistory)

		isFirstBackupString := "false"
		if verifyIsFirstBackup() {
			isFirstBackupString = "true"
		}

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

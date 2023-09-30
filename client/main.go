package main

import (
	"fmt"
	"os"

	"bufio"
	"encoding/json"
	"net"
	"strings"

	// hash "incremental-backup-system/cmd"

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

	color.White("Incremental Backup System\n\n")

	color.White("==============================================================\n")
	color.Blue("Informe o diretório e se deseja salvar histórico dos arquivos: ")
	color.Green("Exemplo: /home/user/backup true\n")
	color.White("==============================================================\n\n")

	// 1 - Encontrar o diretório (especificado no comando) na máquina local
	// 2 - Percorrer esse diretório: (func percorrer_diretorio())
	//			Existe algo dentro? (root != nil) (ou dirEntries != nil) - caso base da recursão?
	//				Percorrer dirEntries[] verificando: (loop iterando sobre dirEntries[])

	// 					se (é um file)
	//						se (ainda não existe no servidor) -> copiar para o servidor ( insert_dir_entry() )
	//						se (existe no servidor) e (foi modificado) -> atualizar no servidor ( update_file() )
	//						Edge Case: como verificar se um arquivo que foi apagado da máquina local existe no servidor? acho que essa verificação para realizar o delete_file() não ocorre aqui

	//					se (é um dir) (dirEntry.Dir != nil)
	//						se (ainda não existe no servidor) -> criar no servidor ( insert_dir_entry() e insert_dir_content() )
	//						1, 2 e 3
	// =========================================================================================================================

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

		str, _ := reader.ReadString('\n')
		strArray := strings.Split(str, " ")

		if len(strArray) < 2 {
			color.Red("Erro: Faltam argumentos. Exemplo: <diretório> <bool: salvar_historico>\n\n")
			continue
		}
		if len(strArray) > 2 {
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

			// cria a hash local e envia pro server
			// localHash := hash.CreateLocalHash(dirPath)

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

/*
// Percorre as entradas de diretório do SO
func walkDirLocal() {
	dirPath := "cmd/server" // request.DirPath

	// Verificando se o diretório (especificado no comando) existe na máquina local
	_, err := os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("O diretório '%s' não existe.", dirPath)
		} else {
			log.Fatalf("Erro ao verificar o diretório: %s", err)
		}
	}

	// Percorre o diretório especificado no comando recursivamente de forma préfixa
	err = filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {

		fmt.Println(path, d.Name(), "directory?", d.IsDir())

		// Criando a hash local
		// atribuir 'path' como key da hash
		// atribuir 'd.Name()' e 'd.ModTime()' na dirEntry (valor da hash)

		// Verificando se d é um arquivo
		if !d.IsDir() {

			// Verifiando se o arquivo existe no servidor e/ou foi modificado
			info, err := d.Info()
			if err != nil {
				log.Fatalf("Erro ao verificar o arquivo: %s", err)
			}
			fmt.Print(info)

			// Enviar path e info.ModTime() para o server
			// para ele verificar se o arquivo existe ou foi alterado
			// response := serverDirStatus(path, info.ModTime())

			// Se o arquivo existe e não foi modificado
			// continue

			// Se o arquivo não existe no servidor ou foi modificado {

			// Obtendo o conteúdo do arquivo
			// fileData, err := readFile(path)
			// if err != nil {
			// 	log.Fatalf("Erro ao ler o arquivo: %s", err)
			// }

			// // Codificando o conteúdo do arquivo em base64
			// fileContentBase64 := base64.StdEncoding.EncodeToString(fileData)

			// enviar (path, info.ModTime(), fileContentBase64) para o server

			// }

		}

		// teste
		// if d.Name() == "server" {
		// 	fmt.Println("Encontrado. Iniciando extração dos dados:")

		// 	info, err := d.Info()

		// 	if err != nil {
		// 		fmt.Print(err)
		// 	}

		// 	fmt.Println(info.Name())
		// 	fmt.Println(info.IsDir())
		// 	fmt.Println(info.ModTime())
		// 	fmt.Println(info.Mode())
		// 	fmt.Println(info.Size())
		// 	fmt.Println(info.Sys())
		// }

		return nil
	})
	if err != nil {
		log.Fatalf("impossible to walk directories: %s", err)
	}
}

func readFile(filePath string) ([]byte, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return data, nil
}
*/

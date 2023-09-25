package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

func readFile(filePath string) ([]byte, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func main() {

	// For status messages
	color.Blue("Hello World!")  // starting some action/command/function
	color.Green("Hello World!") // finished with success
	color.Red("Hello World!")   // finished with error

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

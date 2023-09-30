package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// ********* Hash *********
// CreateLocalHash(): criar uma hash, percorrer estrutura de diretórios local com WalkDir() e adicionando chave (path) e valor (dirEntry) à hash
// CreateServerHash(): ?

// Diferenca(localHash, serverHash): lista as dirEntries novas ou alteradas no local e não existem ou estão mais antigas no server (talvez dê pra separar em duas func, mas talvez dê pra fazer tudo junto)
// UpdateServerHash(lista_modificacao, serverHash): insere as novas (e/ou modificadas) dirEntries na hash do server
// DifInversa ou Complemento - n sei se eh complemento (serverHash, localHash): exclui as dirEntries que existem no server e não existem no local (se a flag... ou verificar essa flag antes de permitir a chamada dessa func)

// CreateDirTree(): transforma o hash do server em uma estrutura de diretórios
// ************************

// Entrada de diretório: arquivo ou diretório
type DirEntry struct {
	Name    string    // n sei se é necessário
	ModDate time.Time // modified date
}

type HashTable struct {
	table map[string]DirEntry // key: path, value: DirEntry
	mutex sync.RWMutex
}

func (ht *HashTable) Put(key string, d fs.DirEntry) {
	ht.mutex.Lock()
	defer ht.mutex.Unlock()

	// Convertendo DirEntry para FileInfo para obter data de modificação
	info, err := d.Info()
	if err != nil {
		log.Fatalf("Erro ao converter DirEntry para FileInfo: %s", err)
	}

	// Criando a key na tabela e atribuindo a ela uma DirEntry com nome e data de modificação
	ht.table[key] = DirEntry{
		Name:    d.Name(),
		ModDate: info.ModTime(),
	}
}

// Atualiza o file modificado na máquina local para na árvore do servidor
// func update_file() {}

// Cria a tabela hash da máquina local (client)
func CreateLocalHash(dirPath string) map[string]DirEntry {

	// Verificando se o diretório (especificado no comando) existe na máquina local
	_, err := os.Stat(dirPath)
	if err != nil {

		if os.IsNotExist(err) {
			log.Fatalf("O diretório '%s' não existe.", dirPath)
			// return // alterar retorno
		} else {
			log.Fatalf("Erro ao verificar o diretório: %s", err)
			// return // alterar retorno
		}
	}

	// Criando e inicializando a hash local
	localHash := HashTable{}
	localHash.table = make(map[string]DirEntry)

	fmt.Println("Hash inicial:")
	fmt.Println(localHash.table)

	// Percorre o diretório especificado no comando recursivamente de forma préfixa
	err = filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {

		//fmt.Println(path, d.Name(), "directory?", d.IsDir())

		if d.Name() == ".git" {
			return filepath.SkipDir
		}

		// Inserindo key: path e value: dirEntry na hash
		localHash.Put(path, d)

		// Verificando se d é um arquivo
		// if !d.IsDir() {

		// Verifiando se o arquivo existe no servidor e/ou foi modificado
		// info, err := d.Info()
		// if err != nil {
		// 	log.Fatalf("Erro ao verificar o arquivo: %s", err)
		// }
		// fmt.Print(info)

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

		// }

		return nil
	})
	if err != nil {
		log.Fatalf("impossible to walk directories: %s", err)
	}

	// fmt.Println("Hash final:")
	// fmt.Print(localHash.table)

	return localHash.table
}

func main() {
	local_hash := CreateLocalHash("../../Backup-System")

	fmt.Println("Hash final:")
	fmt.Print(local_hash)

}

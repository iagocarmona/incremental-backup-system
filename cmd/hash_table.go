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

// DiffToUpdate(localHash, serverHash): lista as dirEntries novas ou alteradas no local e não existem ou estão mais antigas no server (talvez dê pra separar em duas func, mas talvez dê pra fazer tudo junto)
// UpdateServerHash(lista_modificacao, serverHash): insere as novas (e/ou modificadas) dirEntries na hash do server
// DiffToDelete (serverHash, localHash): exclui as dirEntries que existem no server e não existem no local (se a flag... ou verificar essa flag antes de permitir a chamada dessa func)

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

	// fmt.Println("Hash inicial:")
	// fmt.Println(localHash.table)

	// Percorre o diretório especificado no comando recursivamente de forma préfixa
	err = filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {

		//fmt.Println(path, d.Name(), "directory?", d.IsDir())

		if d.Name() == ".git" {
			return filepath.SkipDir
		}

		// Verificando se d é um arquivo
		if !d.IsDir() {
			// Inserindo key: path e value: dirEntry na hash
			localHash.Put(path, d)
		}

		return nil
	})
	if err != nil {
		log.Fatalf("impossible to walk directories: %s", err)
	}

	// fmt.Println("Hash final:")
	// fmt.Print(localHash.table)

	return localHash.table
}

// Atualiza a tabela hash do cliente
func (htl *HashTable) UpdateLocalHash(dirPath string) {

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

		if d.Name() == ".git" {
			return filepath.SkipDir
		}

		// Verificando se d é um arquivo
		if !d.IsDir() {

			// se o elemento não existe na Hash Local
			if _, exist := htl.table[path]; !exist {

				// Insere key: path e value: dirEntry na hash
				htl.Put(path, d)
			}
		}

		return nil
	})
	if err != nil {
		log.Fatalf("impossible to walk directories: %s", err)
	}

	// fmt.Println("Hash final:")
	fmt.Print(htl.table)
}

// Compara a hash local com a hash do server buscando o que existe na hash do local mas não existe na hash do server (ou que precisa ser atualizado no server)
func Diferenca(localHash, serverHash map[string]DirEntry) []string {

	// Lista que será enviada para o cliente que fornecerá os arquivos que precisam ser atualizados ou criados no servidor
	toUpdateList := []string{}

	for keyLocal := range localHash {

		// se o elemento que está em localHash também existir em serverHash
		if _, exist := serverHash[keyLocal]; exist {

			// comparar a data de modificação
			// se a data de modificação do elemento em localHash for mais recente que a data de modificação do elemento em serverHash
			if (localHash[keyLocal].ModDate).Before(serverHash[keyLocal].ModDate) {

				// atualizar o elemento que está em localHsash para o serverHash
				// adicionando a keyLocal em toUpdateList
				toUpdateList = append(toUpdateList, keyLocal)
			}
		} else {
			// se o elemento que está em localHash não existir em serverHash
			// Adiciona a keyLocal em toUpdateList
			toUpdateList = append(toUpdateList, keyLocal)
		}
	}

	return toUpdateList
}

// Verifica os arquivos que precisam ser excluídos do server
func DiffToDelete(serverHash, localHash map[string]DirEntry) []string {

	// Lista que o server utilizará para excluir os arquivos
	toDeleteList := []string{}

	for keyServer := range serverHash {

		// se o elemento que está em serverHash não existir em localHash
		if _, exist := localHash[keyServer]; !exist {

			// Adiciona na lista para excluir o path keyServer de serverHash
			toDeleteList = append(toDeleteList, keyServer)
		}
	}
	return toDeleteList
}

func main() {
	hash1 := CreateLocalHash("../../Backup-System")
	fmt.Println("hash backup pronto, ja tenho")
	time.Sleep(20 * time.Second)
	fmt.Println("criando hash2...")
	hash2 := CreateLocalHash("../../Backup-System")

	// simulando hash_loal e hash_server
	fmt.Println(Diferenca(hash1, hash2))
	fmt.Println()
	fmt.Println(DiffToDelete(hash2, hash1))
}

package hash

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

/*******************************************
*              Estruturas                  *
********************************************/
// Entrada de diretório: arquivo ou diretório
type DirEntry struct {
	Name    string
	ModDate time.Time // Data de Modificação
}

// Tabela Hash com Key: Path e Value: DirEntry
type HashTable struct {
	table map[string]DirEntry
	mutex sync.RWMutex
}

/*******************************************
*            Funções Auxiliares            *
********************************************/
// Insere um novo elemento na hash
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

// Remove um elemento da hash
func Remove(serverHash *map[string]DirEntry, key string) {
	delete(*serverHash, key)
}

/*******************************************
*            Funções Principais            *
********************************************/
// Cria a tabela hash da máquina local (client)
func CreateLocalHash(dirPath string) map[string]DirEntry {

	// Verificando se o diretório (especificado no comando) existe na máquina local
	_, err := os.Stat(dirPath)
	if err != nil {

		if os.IsNotExist(err) {
			log.Fatalf("O diretório '%s' não existe.", dirPath)
		} else {
			log.Fatalf("Erro ao verificar o diretório: %s", err)
		}
	}

	// Criando e inicializando a hash local
	localHash := HashTable{}
	localHash.table = make(map[string]DirEntry)

	// Percorre o diretório especificado no comando recursivamente de forma préfixa
	err = filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {

		// ignora os diretórios .git
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
	// fmt.Print(htl.table)
}

// Cria a tabela hash do servidor
func (ht *HashTable) ToMap() map[string]DirEntry {
	ht.mutex.RLock()
	defer ht.mutex.RUnlock()

	result := make(map[string]DirEntry, len(ht.table))

	for key, value := range ht.table {
		result[key] = value
	}

	return result
}

// Compara a hash local com a hash do server buscando o que existe na hash do local mas não existe na hash do server (ou que precisa ser atualizado no server)
func DiffToUpdate(localHash, serverHash map[string]DirEntry) []string {
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

			// Adiciona na lista para exclui o path keyServer de serverHash
			toDeleteList = append(toDeleteList, keyServer)
		}
	}
	return toDeleteList
}

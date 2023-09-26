package main

import (
	"fmt"
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
	Name    string
	ModDate time.Time // modified date
	Content []byte    // if it's a file, else nil
	Dir     *Node     // if it's a directory, else nil
}

// Página/Nó da árvore
type Node struct {
	dirEntries []DirEntry
}

type Tree struct {
	root *Node
	// minItems int
	// maxItems int
}

// Copia o file ou dir (DirEntry) da máquina local para o seu respectivo diretório (Node) na árvore do servidor
func insert_dir_entry(local_dir_entry DirEntry, server_dir *Node) {

	// insere o file ou dir (DirEntry) numa nova posição do array dirEntries do Node
	server_dir.dirEntries = append(server_dir.dirEntries, local_dir_entry)

	// verificando se o dirEntry é um diretório (pasta) não vazio(a)
	if local_dir_entry.Dir != nil {
		// TODO:
		// criar novo node (func create_node())
		// devolver o endereço desse novo node pra server_dir.dirEntries[ultima_posicao].Dir

		// pois o endereço que virá de local_dir_entry.Dir (o node no qual a pasta local_dir_entry está apontando)
		// não será o mesmo do server_dir (até porque esse node ainda não existe no server)

		// Inserir o conteúdo da pasta local_dir_entry em server_dir
		childDirSize := len(local_dir_entry.Dir.dirEntries)

		for childDirSize > 0 {

			server_dir.dirEntries[ultima_posicao].Dir.dirEntries[childDirSize] = local_dir_entry.Dir.dirEntries[childDirSize]
			// mas se encontrar um dirEntry que é uma pasta, deveria entrar numa recursão (recursão de quem? onde? onde começa?)
			childDirSize -= 1
		}

		server_dir.dirEntries = append(server_dir.dirEntries, local_dir_entry.Dir.dirEntries...)
	}
}

// Atualiza o file modificado na máquina local para na árvore do servidor
// func update_file(local_file DirEntry, server_dir *Node) {}

// Inserir Node na árvore (?)
// func insert_node() {} // func interna que vai ser chamada quando um novo dir for criado?

// função para buscar uma entrada de diretório na árvore (procurar num node ou na árvore toda?)
func search_dir_entry(local_dir_entry DirEntry, server_dir *Node) bool {

	// retornar a direntry por parâmetro se encontrar e bool (para o caso de não encontrar)
}

//função para definir o que será feito com o diretório: criar ou atualizar
// func serverDirStatus() {
//  se server_dir for igual a local_dir_entry (se existe na árvore do server)
// 		se o modDate do local_dir_entry é maior que o modDate do server_dir (se houve alteração no local_dir_entry)
// 			atualizar o server_dir
//	se direntry não existe na árvore
// 		inserir direntry
//}

// func create_node() *Node {}

func main() {

	// criando uma árvore
	tree := Tree{}

	// criando um array de DirEntry
	dir_entries_test := []DirEntry{
		{
			Name:    "test original",
			ModDate: time.Now(),
			Content: nil,
			Dir:     nil,
		},
	}

	// criando um nó (com o array de DirEntry criado acima)
	node1 := Node{
		dirEntries: dir_entries_test,
	}

	// criando o root (com o Node criado acima)
	tree.root = &node1 // ponteiro que aponta para o primeiro nó da árvore, isto é, aponta para o diretório que o usuário passar como parâmetro
	// no primeiro backup

	// criando um file para ser inserido pela func insert_file
	file_to_insert := DirEntry{
		Name:    "test novo",
		ModDate: time.Now(),
		Content: nil,
		Dir:     nil,
	}

	// antes de inserir o file
	fmt.Println(node1.dirEntries)

	// inserindo o file na árvore do servidor
	insert_dir_entry(file_to_insert, &node1)

	// depois de inserir o file
	fmt.Println(node1.dirEntries)
}

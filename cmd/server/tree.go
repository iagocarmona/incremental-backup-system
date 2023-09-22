package main

import (
	"fmt"
	"time"
)

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
	tree.root = &node1

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

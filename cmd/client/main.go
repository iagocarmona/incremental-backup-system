package main

import (
	"github.com/fatih/color"
)

func main() {

	// For status messages
	color.Blue("Hello World!")  // starting some action/command/function
	color.Green("Hello World!") // finished with success
	color.Red("Hello World!")   // finished with error

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
}

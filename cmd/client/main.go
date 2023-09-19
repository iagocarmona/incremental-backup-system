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
	// 					se (é um file) e (ainda não existe no servidor) -> copiar para o servidor ( insert_file() )
}

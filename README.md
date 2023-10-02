# Incremental Backup System

Sistema de backup inteligente fornecendo duas formas de realizar o backup, sendo total e incremental. Utilizando a linguagem de Prograação Go e a tecnologia TCP para comunicação entre cliente e servidor, um cliente solicita um diretório para realizar o backup e também pode informar se gostaria de salvar o histórico dos arquivos no backup.

<br>

# Como executar

Neste projeto existem dois tipos de executáveis, sendo um processo executando se comportando como um cliente e outro como servidor. Siga a documentação abaixo para execução dos processos. Por padrão o servidor será executado na porta **6677**.

## Cliente

O cliente é o processo na qual vai receber do usuário o diretório e um boleano indicando se deseja salvar histórico dos arquivos, para em seguida enviar essa informação para o servidor e aguardar a resposta em caso de sucesso ou falha.

Execute as seguintes linhas de comando no terminal:

- `cd client`
- `go run . localhost:<porta>`

por padrão `porta` será **6677**

## Servidor

O servidor é o processo responsável por receber os arquivos do cliente e criar uma estrutura de diretórios que tem como função ser o backup dos arquivos.

Execute as seguintes linhas de comando no terminal:

- `cd server`
- `go run . <porta>`

por padrão `porta` será **6677**

<br>

# Autores

Iago Ortega Carmona

Gustavo Zanzin Guerreiro Martins

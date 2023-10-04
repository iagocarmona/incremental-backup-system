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

<center>
<table>
  <tr>
<td align="center"><a href="https://github.com/iagocarmona">
 <img style="border-radius: 50%;" src="https://avatars.githubusercontent.com/u/69121686?s=400&u=c6fc38d355b96f4abf690ae95912c07e5f057b94&v=4" width="200px;" alt="Avatar Iago"/>
<br />
 <b>Iago Carmona</b>
 </a> <a href="https://github.com/iagocarmona" title="Repositorio Iago"></a>

[![Github Badge](https://img.shields.io/badge/-iagocarmona-000?style=flat-square&logo=Github&logoColor=white&link=https://github.com/iagocarmona)](https://github.com/iagocarmona)</td>

<td align="center"><a href="https://github.com/GustavoMartinx">
 <img style="border-radius: 50%;" src="https://avatars.githubusercontent.com/u/90780907?v=4" width="200px;" alt="Avatar Gustavo"/>
<br />
 <b>Gustavo Zanzin</b>
 </a> <a href="https://github.com/GustavoMartinx" title="Repositorio Gustavo"></a>

[![Github Badge](https://img.shields.io/badge/-GustavoMartinx-000?style=flat-square&logo=Github&logoColor=white&link=https://github.com/GustavoMartinx)](https://github.com/GustavoMartinx)</td>
</tr></table>
</center>

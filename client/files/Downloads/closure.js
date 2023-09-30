function contador() {
  let count = 0 // Esta variável é "capturada" pela função interna
  return function () {
    return count++
  }
}

let minhaContagem = contador()
console.log(minhaContagem()) // 0
console.log(minhaContagem()) // 1

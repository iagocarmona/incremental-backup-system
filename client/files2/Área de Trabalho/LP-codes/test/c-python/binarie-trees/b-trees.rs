use std::env;

// Definição da estrutura do nó da árvore binária
struct TreeNode {
    value: i32,
    left: Option<Box<TreeNode>>,
    right: Option<Box<TreeNode>>,
}

// Função para criar a árvore binária com a profundidade dada
fn create_binary_tree(depth: i32, value: i32) -> Option<Box<TreeNode>> {
    if depth == 0 {
        None
    } else {
        Some(Box::new(TreeNode {
            value,
            left: create_binary_tree(depth - 1, value * 2),
            right: create_binary_tree(depth - 1, value * 2 + 1),
        }))
    }
}

// Função para realizar a travessia em ordem (inorder traversal) na árvore binária
fn inorder_traversal(node: &Option<Box<TreeNode>>) {
    if let Some(inner) = node {
        inorder_traversal(&inner.left);
        println!("{}", inner.value);
        inorder_traversal(&inner.right);
    }
}

fn main() {
    let args: Vec<String> = env::args().collect();

    if args.len() != 2 {
        println!("Uso: {} <profundidade>", args[0]);
        return;
    }

    let depth: i32 = args[1].parse().expect("Profundidade inválida");
    let root = create_binary_tree(depth, 1);

    println!("Árvore binária criada com profundidade {}", depth);
    println!("Resultado da travessia em ordem (inorder traversal):");
    inorder_traversal(&root);
}

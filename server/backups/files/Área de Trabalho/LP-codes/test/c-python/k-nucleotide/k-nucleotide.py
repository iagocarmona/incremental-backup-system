from collections import Counter

def read_sequence(filename):
    with open(filename, 'r') as file:
        return ''.join(line.strip() for line in file if not line.startswith(">"))

def count_kmers(sequence, k):
    kmers = [sequence[i:i+k] for i in range(len(sequence) - k + 1)]
    return Counter(kmers)

def main():
    filename = 'sequence.txt'  # Coloque o nome do arquivo contendo a sequência de DNA
    k = 3  # Tamanho do k-mer

    sequence = read_sequence(filename)
    kmer_counts = count_kmers(sequence, k)

    sorted_kmer_counts = sorted(kmer_counts.items(), key=lambda item: (-item[1], item[0]))

    for kmer, count in sorted_kmer_counts:
        print(f"{kmer}: {count}")

if __name__ == "__main__":
    main()

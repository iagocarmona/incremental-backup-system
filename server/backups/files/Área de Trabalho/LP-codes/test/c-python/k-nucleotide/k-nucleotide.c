#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <ctype.h>

#define MAX_KMER_LENGTH 10

typedef struct {
    char kmer[MAX_KMER_LENGTH + 1];
    int count;
} KmerCount;

int compare_kmer_count(const void *a, const void *b);

void normalize_sequence(char *sequence) {
    for (int i = 0; sequence[i]; i++) {
        sequence[i] = toupper(sequence[i]);
    }
}

void count_kmers(const char *sequence, int k, KmerCount *kmers, int *num_kmers) {
    int seq_length = strlen(sequence);
    *num_kmers = seq_length - k + 1;

    for (int i = 0; i < *num_kmers; i++) {
        strncpy(kmers[i].kmer, sequence + i, k);
        kmers[i].kmer[k] = '\0';
        kmers[i].count = 1;
    }

    qsort(kmers, *num_kmers, sizeof(KmerCount), compare_kmer_count);
    
    for (int i = 1; i < *num_kmers; i++) {
        if (strcmp(kmers[i].kmer, kmers[i - 1].kmer) == 0) {
            kmers[i].count++;
        }
    }

    *num_kmers = (*num_kmers > 10) ? 10 : *num_kmers;
}

int compare_kmer_count(const void *a, const void *b) {
    const KmerCount *k1 = (const KmerCount *)a;
    const KmerCount *k2 = (const KmerCount *)b;

    if (k1->count > k2->count) {
        return -1;
    } else if (k1->count < k2->count) {
        return 1;
    } else {
        return strcmp(k1->kmer, k2->kmer);
    }
}

int main() {
    const char *filename = "sequence.txt";  // Substitua pelo nome do arquivo com a sequência de DNA
    int k = 3;  // Tamanho do k-mer

    FILE *file = fopen(filename, "r");
    if (file == NULL) {
        perror("Error opening file");
        return 1;
    }

    fseek(file, 0, SEEK_END);
    int seq_length = ftell(file);
    fseek(file, 0, SEEK_SET);

    char *sequence = (char *)malloc(seq_length + 1);
    if (sequence == NULL) {
        perror("Error allocating memory");
        fclose(file);
        return 1;
    }

    fread(sequence, 1, seq_length, file);
    fclose(file);

    sequence[seq_length] = '\0';
    normalize_sequence(sequence);

    KmerCount kmers[seq_length - k + 1];
    int num_kmers;

    count_kmers(sequence, k, kmers, &num_kmers);

    for (int i = 0; i < num_kmers; i++) {
        printf("%s: %d\n", kmers[i].kmer, kmers[i].count);
    }

    free(sequence);

    return 0;
}

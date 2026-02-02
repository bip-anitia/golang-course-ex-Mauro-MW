# Esercizio 1: Word Frequency Counter

## Obiettivo
Creare un programma che analizza uno o più file di testo e conta la frequenza di ogni parola.

## Descrizione
Il programma deve leggere file di testo e produrre statistiche sulle parole contenute, contando quante volte ogni parola appare nel testo.

## Requisiti

1. **Lettura File**: Leggere uno o più file di testo specificati come argomenti da linea di comando
2. **Parsing**: Suddividere il testo in parole (gestire punteggiatura e case-insensitive)
3. **Conteggio**: Usare una `map[string]int` per contare le occorrenze
4. **Output**: Stampare le parole ordinate per frequenza (dalla più frequente alla meno frequente)
5. **Gestione Errori**: Gestire correttamente file non esistenti o errori di lettura

## Funzionalità Extra (Opzionali)

- Ignorare parole comuni (stop words) come "il", "la", "di", etc.
- Opzione per salvare i risultati in un file CSV
- Mostrare solo le top N parole più frequenti
- Supporto per encoding diversi (UTF-8, ISO-8859-1, etc.)

## Esempi di Utilizzo

```bash
# Analizza un singolo file
go run main.go testo.txt

# Analizza multipli file
go run main.go file1.txt file2.txt file3.txt

# Con opzioni (da implementare)
go run main.go -top=10 -ignore-case testo.txt
```

## Output Atteso

```
Parole totali: 1523
Parole uniche: 342

Top 10 parole più frequenti:
1. "esempio"    - 45 occorrenze
2. "programma"  - 32 occorrenze
3. "golang"     - 28 occorrenze
...
```

## Concetti Go da Usare

- `os.ReadFile()` o `bufio.Scanner` per leggere file
- `map[string]int` per il conteggio
- `strings` package per manipolazione stringhe
- `sort.Slice()` per ordinare i risultati
- Gestione errori con pattern `if err != nil`

## Suggerimenti

- Usa `strings.Fields()` per dividere in parole
- Usa `strings.ToLower()` per normalizzare le parole
- Converti la map in uno slice di struct per poterla ordinare
- Considera l'uso di `regexp` per rimuovere punteggiatura

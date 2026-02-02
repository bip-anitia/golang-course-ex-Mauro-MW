# Esercizio 4: CLI Tool con Flags

## Obiettivo
Creare un tool da linea di comando completo che accetta vari tipi di flag e argomenti, simile a tool come `grep`, `find`, o `curl`.

## Descrizione
Implementare un programma CLI che processa file con varie opzioni configurabili tramite flags. Il tool deve avere un help chiaro e gestire correttamente gli argomenti.

## Proposta: File Processing Tool

Creare un tool che processa file di testo con diverse opzioni:
- Contare linee, parole, caratteri
- Cercare pattern
- Convertire case
- Filtrare linee

## Requisiti

### 1. Flags da Implementare

```go
// Flags booleani
-verbose, -v       // Output verboso
-quiet, -q         // Output minimo
-help, -h          // Mostra help

// Flags con valori
-output, -o        // File di output (string)
-lines, -n         // Numero di linee da processare (int)
-pattern, -p       // Pattern da cercare (string)
-format, -f        // Formato output: json, text, csv (string)

// Flags multipli
-exclude           // Pattern da escludere (può essere ripetuto)
```

### 2. Subcommands (Opzionale)
```bash
# Diversi comandi
go run main.go count file.txt
go run main.go search -pattern="error" file.txt
go run main.go convert -upper file.txt
go run main.go stats file.txt
```

### 3. Validazione
- Validare che i valori dei flags siano corretti
- Gestire flags incompatibili (es. -verbose e -quiet insieme)
- Verificare che i file specificati esistano

### 4. Help Output
Implementare un help chiaro e ben formattato:

```
Usage: filetools [OPTIONS] COMMAND [FILES...]

A versatile file processing tool.

Commands:
  count       Count lines, words, and characters
  search      Search for patterns in files
  convert     Convert text (uppercase, lowercase)
  stats       Show file statistics

Options:
  -v, -verbose        Enable verbose output
  -q, -quiet          Suppress non-error output
  -o, -output FILE    Write output to FILE instead of stdout
  -n, -lines NUM      Process only first NUM lines
  -p, -pattern TEXT   Pattern to search for
  -f, -format FORMAT  Output format (text, json, csv) [default: text]
  -h, -help           Show this help message

Examples:
  filetools count file.txt
  filetools search -pattern="error" -n=100 app.log
  filetools convert -upper file.txt -output=OUTPUT.txt
  filetools stats -format=json *.txt

For more information, visit: https://github.com/yourusername/filetools
```

## Esempi di Utilizzo

```bash
# Count command
go run main.go count file.txt
# Output: Lines: 150, Words: 1234, Chars: 8432

# Search con pattern
go run main.go search -pattern="error" -verbose app.log

# Con output file
go run main.go convert -upper input.txt -output=upper.txt

# Formato JSON
go run main.go stats -format=json file1.txt file2.txt

# Verbose mode
go run main.go -v count *.txt
```

## Output Attesi

### Count (text format)
```
Processing: file.txt
Lines:      150
Words:      1,234
Characters: 8,432
Size:       8.2 KB
```

### Stats (JSON format)
```json
{
  "files": [
    {
      "name": "file1.txt",
      "lines": 150,
      "words": 1234,
      "size_bytes": 8432
    },
    {
      "name": "file2.txt",
      "lines": 200,
      "words": 2100,
      "size_bytes": 12000
    }
  ],
  "totals": {
    "files": 2,
    "lines": 350,
    "words": 3334,
    "size_bytes": 20432
  }
}
```

## Concetti Go da Usare

- `flag` package per parsing flags
- `flag.String()`, `flag.Int()`, `flag.Bool()` per definire flags
- `flag.Parse()` per parsare gli argomenti
- `flag.Args()` per ottenere argomenti non-flag
- `flag.Usage` per custom help
- `os.Exit()` per codici di uscita appropriati
- Subcommands con pattern dispatcher o librerie come `cobra`

## Struttura Suggerita

```go
package main

import (
    "flag"
    "fmt"
    "os"
)

type Config struct {
    Verbose bool
    Quiet   bool
    Output  string
    Lines   int
    Pattern string
    Format  string
}

func main() {
    config := parseFlags()

    if len(flag.Args()) == 0 {
        flag.Usage()
        os.Exit(1)
    }

    command := flag.Args()[0]
    files := flag.Args()[1:]

    switch command {
    case "count":
        runCount(config, files)
    case "search":
        runSearch(config, files)
    case "convert":
        runConvert(config, files)
    case "stats":
        runStats(config, files)
    default:
        fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
        os.Exit(1)
    }
}

func parseFlags() *Config {
    // TODO: Implementare parsing flags
    return &Config{}
}
```

## Suggerimenti

- Usa `flag.Usage` per customizzare il messaggio di help
- Valida flags mutuamente esclusivi
- Usa `os.Stderr` per messaggi di errore
- Implementa exit codes appropriati (0 = success, 1 = error)
- Testa il tool con vari input edge case
- Considera usare `flag.FlagSet` per subcommands
- Aggiungi colored output se -verbose (usando librerie come `fatih/color`)

## Challenge Extra

- Implementare config file support (es. `.filestoolsrc`)
- Aggiungere completion per bash/zsh
- Progress bar per file grandi
- Supporto per stdin (piping)
- Wildcard expansion cross-platform
- Versioning (`-version` flag)
- Documentazione man page style

# Esercizio 2: Concurrent Web Scraper

## Obiettivo
Creare un web scraper concorrente che scarica contenuti da multipli URL in parallelo utilizzando goroutines e channels.

## Descrizione
Il programma deve accettare una lista di URL, scaricare il contenuto HTML di ogni pagina in modo concorrente, ed estrarre informazioni specifiche (es. titolo, links, lunghezza del contenuto).

## Requisiti

1. **Input URLs**: Leggere una lista di URL da file o da argomenti CLI
2. **Concorrenza**: Usare goroutines per scaricare multipli URL contemporaneamente
3. **Channels**: Usare channels per comunicare i risultati tra goroutines
4. **HTTP Requests**: Fare richieste HTTP GET usando `net/http`
5. **Parsing HTML**: Estrarre informazioni base dall'HTML (titolo, numero di link)
6. **Timeout**: Implementare timeout per le richieste HTTP
7. **Error Handling**: Gestire errori di rete e URL non validi

## Funzionalità da Implementare

### Struttura Risultato
```go
type PageInfo struct {
    URL          string
    Title        string
    StatusCode   int
    ContentSize  int
    LinkCount    int
    Error        error
}
```

### Limiti di Concorrenza
- Massimo N richieste simultanee (es. 5)
- Usare un semaforo con buffered channel

## Esempi di Utilizzo

```bash
# Scraping da lista di URL
go run main.go urls.txt

# Con numero massimo di workers
go run main.go -workers=10 urls.txt
```

## File urls.txt di esempio
```
https://golang.org
https://go.dev
https://github.com/golang/go
https://pkg.go.dev
```

## Output Atteso

```
Scraping 4 URLs con 5 workers...

[OK] https://golang.org
     Status: 200 | Size: 45123 bytes | Links: 87 | Title: "The Go Programming Language"

[OK] https://go.dev
     Status: 200 | Size: 32451 bytes | Links: 45 | Title: "Go.dev"

[ERROR] https://invalid-url.xyz
     Error: Get "https://invalid-url.xyz": dial tcp: lookup invalid-url.xyz: no such host

Completato in 2.3s
Successi: 3/4
```

## Concetti Go da Usare

- `goroutine` per concorrenza
- `chan` per comunicazione tra goroutines
- `sync.WaitGroup` per attendere completamento
- `http.Client` con timeout configurato
- `context.WithTimeout()` per timeout delle richieste
- Buffered channels per limitare concorrenza
- `strings` package o librerie come `golang.org/x/net/html` per parsing

## Suggerimenti

- Crea un worker pool pattern
- Usa `defer resp.Body.Close()` dopo ogni richiesta HTTP
- Imposta un `User-Agent` appropriato
- Per estrarre il titolo cerca il tag `<title>`
- Gestisci redirect HTTP
- Considera rate limiting per evitare di sovraccaricare i server

## Challenge Extra

- Implementare crawling ricorsivo (seguire i link trovati)
- Salvare i risultati in JSON
- Aggiungere retry logic per richieste fallite
- Implementare cache per evitare richieste duplicate

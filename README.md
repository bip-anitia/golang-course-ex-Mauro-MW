# Esercizi Go - Corso Avanzato

Collezione completa di **15 esercizi pratici** per studenti che hanno completato il libro "The Go Programming Language" di Donovan & Kernighan. Ogni esercizio copre concetti fondamentali e avanzati con focus su concorrenza, interfacce, performance, e pattern idiomatici di Go.

## 📊 Panoramica

- **15 esercizi** organizzati per difficoltà crescente
- **README dettagliati** con esempi, output attesi, e best practices
- **File main.go** minimali per iniziare subito
- **Challenge extra** per approfondire
- **Testing examples** per ogni esercizio

## 🚀 Quick Start

```bash
# Clona o scarica gli esercizi
cd esercizi

# Inizia dal primo esercizio
cd esercizio-01-word-frequency
cat README.md          # Leggi le istruzioni
go run main.go         # Testa lo scheletro

# Dopo aver implementato, testa
go run main.go esempio.txt
```

## 📚 Struttura del Corso

### 📘 Livello Base/Intermedio (5 esercizi)

Fondamenta e pattern comuni:

| # | Esercizio | Argomenti Chiave | Difficoltà |
|---|-----------|------------------|------------|
| 01 | [Word Frequency Counter](esercizio-01-word-frequency/) | maps, file I/O, strings | ⭐⭐ |
| 02 | [Concurrent Web Scraper](esercizio-02-web-scraper/) | goroutines, channels, HTTP | ⭐⭐ |
| 03 | [JSON API Server](esercizio-03-json-api/) | REST API, encoding/json, net/http | ⭐⭐ |
| 04 | [CLI Tool con Flags](esercizio-04-cli-flags/) | flag package, CLI design | ⭐⭐ |
| 05 | [Rate Limiter](esercizio-05-rate-limiter/) | channels, time.Ticker, patterns | ⭐⭐⭐ |

### 📗 Livello Intermedio/Avanzato (4 esercizi)

Concorrenza e design patterns:

| # | Esercizio | Argomenti Chiave | Difficoltà |
|---|-----------|------------------|------------|
| 06 | [Worker Pool](esercizio-06-worker-pool/) | concurrency patterns, pooling | ⭐⭐⭐ |
| 07 | [Custom Sort](esercizio-07-custom-sort/) | sort.Interface, comparators | ⭐⭐ |
| 08 | [Context Propagation](esercizio-08-context/) | context, cancellation, timeout | ⭐⭐⭐ |
| 09 | [Interface Design](esercizio-09-interface-design/) | interfaces, SOLID, composition | ⭐⭐⭐ |

### 📕 Livello Avanzato (6 esercizi)

Pattern avanzati e production-ready code:

| # | Esercizio | Argomenti Chiave | Difficoltà |
|---|-----------|------------------|------------|
| 10 | [Pipeline Pattern](esercizio-10-pipeline/) | pipelines, fan-out/fan-in | ⭐⭐⭐⭐ |
| 11 | [Error Wrapping](esercizio-11-error-wrapping/) | error chains, custom errors | ⭐⭐⭐ |
| 12 | [Defer, Panic, Recover](esercizio-12-defer-panic-recover/) | resource cleanup, recovery | ⭐⭐ |
| 13 | [Pomodoro Timer](esercizio-13-pomodoro-timer/) | time, Timer, Ticker, timezone | ⭐⭐ |
| 14 | [Graceful Shutdown](esercizio-14-graceful-shutdown/) | os/signal, graceful stop | ⭐⭐⭐ |
| 15 | [Benchmarking](esercizio-15-benchmarking/) | testing.B, profiling, pprof | ⭐⭐⭐ |

## 🎯 Concetti Chiave Coperti

### Concorrenza
- Goroutines e channels
- sync package (Mutex, RWMutex, WaitGroup)
- Context per cancellazione e timeout
- Pipeline patterns (fan-out/fan-in)
- Worker pools
- Rate limiting

### Design e Architettura
- Interfacce e composizione
- SOLID principles
- Error handling e wrapping
- Custom error types
- Sentinel errors
- Package organization

### Standard Library Essenziale
- **I/O**: `io`, `os`, `bufio`
- **Network**: `net/http`, HTTP client/server
- **Data**: `encoding/json`, `encoding/csv`
- **Text**: `strings`, `regexp`, `strconv`
- **CLI**: `flag`, `os`
- **Concurrency**: `context`, `sync`
- **Errors**: `errors`, `fmt.Errorf`
- **Time**: `time` (Timer, Ticker, Duration, Location)
- **System**: `os/signal`, `syscall`
- **Performance**: `testing`, `runtime/pprof`

### Testing & Performance
- Unit testing
- Table-driven tests
- Benchmarking (`testing.B`)
- Memory profiling
- CPU profiling
- Performance optimization

### Resource Management
- Defer statement
- Panic e recover
- Resource cleanup patterns
- Graceful shutdown
- Signal handling

## 📖 Come Usare Questi Esercizi

### Per gli Studenti

1. **Leggi il README** dell'esercizio per capire obiettivi e requisiti
2. **Implementa la soluzione** nel file `main.go` fornito
3. **Testa il codice** con gli esempi forniti
4. **Verifica l'output** confrontandolo con quello atteso
5. **Sperimenta** con le funzionalità extra e challenge
6. **Scrivi test** se richiesto
7. **Confronta** con le best practice suggerite

### Per gli Istruttori

- Ogni esercizio ha obiettivi chiari di apprendimento
- README dettagliati con teoria e pratica
- Output attesi per facilitare correzione
- Suggerimenti graduali (non soluzioni complete)
- Challenge extra per studenti avanzati
- Criteri di valutazione suggeriti
- Riferimenti a capitoli del libro Donovan & Kernighan

### Tempo Stimato per Esercizio

- **Base**: 2-3 ore
- **Intermedio**: 3-4 ore
- **Avanzato**: 4-6 ore

**Totale corso**: ~50-60 ore di lavoro pratico

## 🗂️ Ordine Consigliato

### Approccio 1: Sequenziale
Segui l'ordine numerico (01 → 15) per progressione naturale di difficoltà.

**Vantaggi**: Costruisce conoscenza gradualmente, ideale per principianti.

### Approccio 2: Per Tema

Raggruppa esercizi per argomento:

**🔄 Concorrenza e Parallelismo**
```
02 (Web Scraper) → 05 (Rate Limiter) → 06 (Worker Pool) → 10 (Pipeline)
```

**📝 I/O e Data Processing**
```
01 (Word Frequency) → 04 (CLI Flags) → 07 (Custom Sort)
```

**🌐 Network e Server**
```
02 (Web Scraper) → 03 (JSON API) → 14 (Graceful Shutdown)
```

**🏗️ Design e Architettura**
```
09 (Interface Design) → 11 (Error Wrapping) → 12 (Defer/Panic/Recover)
```

**⚙️ Context e Controllo**
```
08 (Context) → 10 (Pipeline) → 14 (Graceful Shutdown)
```

**⏱️ Time e Resource Management**
```
12 (Defer/Panic/Recover) → 13 (Pomodoro Timer) → 14 (Graceful Shutdown)
```

**⚡ Performance e Optimization**
```
15 (Benchmarking) - Applicabile a tutti gli esercizi precedenti
```

### Approccio 3: Project-Based

Per progetti reali, combina esercizi:

**Backend API Completo**:
```
03 (JSON API) + 08 (Context) + 11 (Error Wrapping) + 14 (Graceful Shutdown)
```

**Data Processing Pipeline**:
```
01 (Word Frequency) + 06 (Worker Pool) + 10 (Pipeline) + 15 (Benchmarking)
```

**Production-Ready Service**:
```
03 (API) + 05 (Rate Limiter) + 12 (Defer/Panic) + 14 (Graceful Shutdown)
```

## 🔧 Prerequisiti

### Software Richiesto
- **Go 1.18+** installato ([download](https://go.dev/dl/))
- Editor con Go support (VS Code + Go extension, GoLand, vim-go, etc.)
- Terminal/Command line familiarity

### Conoscenze Richieste
- ✅ Completamento libro "The Go Programming Language" (Donovan & Kernighan)
- ✅ Sintassi base Go (variabili, funzioni, struct, slice, map)
- ✅ Goroutines e channels (concetti base)
- ✅ Interfacce (basic understanding)
- ✅ Error handling (`if err != nil` pattern)
- ✅ Package e import

### Verifica Prerequisiti

```bash
# Verifica versione Go
go version  # Dovrebbe essere >= 1.18

# Verifica ambiente
go env

# Test base
cd esercizio-01-word-frequency
go run main.go
```

## 💡 Suggerimenti Generali

### Best Practices da Seguire

1. **Gestione Errori**: Mai ignorare errori, sempre controllare `if err != nil`
   ```go
   data, err := os.ReadFile("file.txt")
   if err != nil {
       return fmt.Errorf("failed to read file: %w", err)
   }
   ```

2. **Concorrenza**: Usa channels per comunicare, non condividere memoria
   ```go
   // ✅ GOOD: communicate via channels
   results := make(chan Result)
   go worker(jobs, results)

   // ❌ BAD: shared memory without sync
   var counter int
   go func() { counter++ }()
   ```

3. **Semplicità**: Scrivi codice semplice e leggibile, evita over-engineering
   ```go
   // ✅ GOOD: simple and clear
   if user == nil {
       return ErrNotFound
   }

   // ❌ BAD: unnecessary complexity
   return ternary(user != nil, nil, ErrNotFound)
   ```

4. **Testing**: Scrivi test per il tuo codice
   ```go
   func TestFunction(t *testing.T) {
       result := myFunction(input)
       if result != expected {
           t.Errorf("got %v, want %v", result, expected)
       }
   }
   ```

5. **Documentazione**: Commenta funzioni esportate
   ```go
   // ProcessData processes the input data and returns the result.
   // Returns an error if the data is invalid.
   func ProcessData(data []byte) (Result, error) {
       // ...
   }
   ```

6. **Formatting**: Usa sempre `gofmt` e `goimports`
   ```bash
   gofmt -w .
   goimports -w .
   ```

### Anti-Pattern da Evitare

- ❌ **Goroutine leaks**: Goroutines che non terminano mai
- ❌ **Race conditions**: Accesso concorrente non sincronizzato
- ❌ **Panic in libraries**: Usa errors invece di panic
- ❌ **Nil pointer dereference**: Controlla sempre nil prima di usare
- ❌ **Ignorare context.Done()**: Rispetta sempre cancellazione
- ❌ **Defer in loop**: Defer dentro loop accumula fino a fine funzione
- ❌ **Premature optimization**: Prima funziona, poi ottimizza
- ❌ **Magic numbers**: Usa costanti con nomi significativi

### Comandi Go Utili

```bash
# Formattazione
go fmt ./...
goimports -w .

# Build
go build
go build -o myapp

# Run
go run main.go
go run .

# Test
go test
go test -v
go test -cover
go test ./...

# Benchmark
go test -bench=.
go test -bench=. -benchmem

# Profiling
go test -cpuprofile=cpu.prof
go test -memprofile=mem.prof
go tool pprof cpu.prof

# Dependencies
go mod init
go mod tidy
go get package

# Vet (static analysis)
go vet ./...

# Race detector
go run -race main.go
go test -race
```

## 📊 Valutazione

Ogni esercizio può essere valutato su:

| Criterio | Peso | Cosa Valutare |
|----------|------|---------------|
| **Correttezza** | 40% | Codice funziona come richiesto, casi edge gestiti |
| **Stile** | 20% | Segue convenzioni Go, codice leggibile, nomi chiari |
| **Error Handling** | 20% | Errori gestiti correttamente, messaggi informativi |
| **Concorrenza** | 10% | Uso corretto goroutines/channels (se applicabile) |
| **Testing** | 10% | Include test appropriati, coverage ragionevole |

### Livelli di Competenza

- **🥉 Base (60-70%)**: Soluzione funzionante, errori base gestiti
- **🥈 Intermedio (70-85%)**: Codice pulito, buona gestione errori, alcuni test
- **🥇 Avanzato (85-95%)**: Best practices, test completi, edge cases, performance considerata
- **🏆 Eccellente (95-100%)**: Codice production-ready, testing completo, documentazione, challenge extra

## 🆘 Supporto e Risorse

### Quando Sei Bloccato

1. **Leggi attentamente** il README dell'esercizio
2. **Consulta** la sezione "Suggerimenti" nel README
3. **Controlla** gli esempi di codice forniti
4. **Rivedi** i capitoli rilevanti del libro Donovan & Kernighan
5. **Sperimenta** in REPL Go ([Go Playground](https://go.dev/play/))

### Risorse Online

- **Go by Example**: https://gobyexample.com/ - Esempi concisi per ogni concetto
- **Effective Go**: https://go.dev/doc/effective_go - Stile e idiomi Go
- **Go Blog**: https://go.dev/blog/ - Articoli approfonditi
- **Go Standard Library**: https://pkg.go.dev/std - Documentazione completa
- **Go Wiki**: https://github.com/golang/go/wiki - Guide e tutorial
- **Go Playground**: https://go.dev/play/ - Test codice online
- **Awesome Go**: https://awesome-go.com/ - Librerie e risorse curate

### Libri Consigliati

- **The Go Programming Language** (Donovan & Kernighan) - *prerequisito*
- **Concurrency in Go** (Katherine Cox-Buday) - Approfondimento concorrenza
- **Learning Go** (Jon Bodner) - Best practices moderne
- **100 Go Mistakes** (Teiva Harsanyi) - Errori comuni da evitare

## 🎓 Certificazione e Completamento

### Criterio di Completamento

Un esercizio è considerato completo quando:
- ✅ Codice compila senza errori
- ✅ Output corrisponde a quello atteso
- ✅ Test passano (se richiesti)
- ✅ Nessun warning da `go vet`
- ✅ Formattato con `gofmt`

### Dopo il Completamento

Una volta completati tutti gli esercizi:
1. **Portfolio**: Aggiungi al tuo GitHub/portfolio
2. **Blog Post**: Scrivi dei concetti appresi
3. **Progetti Personali**: Applica i pattern appresi
4. **Contribuisci**: Open source Go projects
5. **Continua**: Esercizi avanzati (reflection, generics, systems programming)

## 📝 Note per gli Istruttori

### Suggerimenti Didattici

- **Pair Programming**: Fai lavorare studenti in coppia
- **Code Review**: Organizza sessioni di review del codice
- **Live Coding**: Dimostra pattern in tempo reale
- **Debugging Sessions**: Mostra come debuggare problemi comuni
- **Performance Analysis**: Analizza benchmark insieme

### Customizzazione

Gli esercizi possono essere adattati:
- Aggiungere requisiti specifici al contesto
- Modificare domini (es. e-commerce invece di libri)
- Integrare con progetti esistenti
- Creare varianti più semplici/complesse

### Tracking Progress

Template per tracking studente:

```
Studente: _______________
Data Inizio: ___________

Esercizi Completati:
[ ] 01 - Word Frequency Counter
[ ] 02 - Concurrent Web Scraper
...
[ ] 15 - Benchmarking

Note:
___________________________
```

## 📜 Licenza

Questi esercizi sono forniti per **scopi educativi**.

Sei libero di:
- ✅ Usare per apprendimento personale
- ✅ Usare in corsi e workshop
- ✅ Modificare e adattare
- ✅ Condividere con studenti

Si richiede solo di:
- Mantenere attribuzione originale
- Non vendere gli esercizi come prodotto standalone

## 🙏 Ringraziamenti

Esercizi ispirati da:
- "The Go Programming Language" di Donovan & Kernighan
- [Go by Example](https://gobyexample.com/)
- La community Go
- Best practices della Go standard library

---

## 🚀 Ready to Start?

```bash
cd esercizio-01-word-frequency
cat README.md
go run main.go
```

**Buon Coding! 🎉**

*"Simplicity is complicated."* - Rob Pike

*"Clear is better than clever."* - Go Proverbs

---

**Ultima revisione**: 2024
**Versione esercizi**: 1.0
**Go version**: 1.18+

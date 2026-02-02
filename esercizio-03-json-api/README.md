# Esercizio 3: JSON API Server

## Obiettivo
Creare un REST API server che gestisce una risorsa (es. libri, utenti, prodotti) con operazioni CRUD complete, utilizzando JSON per lo scambio dati.

## Descrizione
Implementare un server HTTP che espone endpoint RESTful per gestire una collezione di risorse, con serializzazione/deserializzazione JSON e gestione corretta degli HTTP status codes.

## Requisiti

### 1. Risorsa da Gestire
Scegli una risorsa, ad esempio **Libri**:
```go
type Book struct {
    ID          string    `json:"id"`
    Title       string    `json:"title"`
    Author      string    `json:"author"`
    ISBN        string    `json:"isbn"`
    PublishYear int       `json:"publish_year"`
    CreatedAt   time.Time `json:"created_at"`
}
```

### 2. Endpoint da Implementare

| Metodo | Endpoint        | Descrizione                |
|--------|-----------------|----------------------------|
| GET    | `/books`        | Lista tutti i libri        |
| GET    | `/books/{id}`   | Ottiene un libro per ID    |
| POST   | `/books`        | Crea un nuovo libro        |
| PUT    | `/books/{id}`   | Aggiorna un libro esistente|
| DELETE | `/books/{id}`   | Elimina un libro           |

### 3. Storage in Memoria
- Usare una `map[string]Book` come database in-memory
- Proteggere l'accesso concorrente con `sync.RWMutex`

### 4. Validazione
- Validare input JSON (campi obbligatori, formati)
- Ritornare errori appropriati (400 Bad Request, 404 Not Found, etc.)

### 5. HTTP Status Codes
- `200 OK` - richiesta riuscita
- `201 Created` - risorsa creata
- `400 Bad Request` - input non valido
- `404 Not Found` - risorsa non trovata
- `405 Method Not Allowed` - metodo HTTP non supportato
- `500 Internal Server Error` - errore server

## Esempi di Utilizzo

```bash
# Creare un libro
curl -X POST http://localhost:8080/books \
  -H "Content-Type: application/json" \
  -d '{"title":"The Go Programming Language","author":"Donovan & Kernighan","isbn":"978-0134190440","publish_year":2015}'

# Lista tutti i libri
curl http://localhost:8080/books

# Ottenere un libro specifico
curl http://localhost:8080/books/1

# Aggiornare un libro
curl -X PUT http://localhost:8080/books/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"Updated Title","author":"Donovan & Kernighan","isbn":"978-0134190440","publish_year":2016}'

# Eliminare un libro
curl -X DELETE http://localhost:8080/books/1
```

## Risposte Attese

### GET /books
```json
{
  "books": [
    {
      "id": "1",
      "title": "The Go Programming Language",
      "author": "Donovan & Kernighan",
      "isbn": "978-0134190440",
      "publish_year": 2015,
      "created_at": "2024-01-15T10:30:00Z"
    }
  ]
}
```

### POST /books (201 Created)
```json
{
  "id": "2",
  "title": "The Go Programming Language",
  "author": "Donovan & Kernighan",
  "isbn": "978-0134190440",
  "publish_year": 2015,
  "created_at": "2024-01-15T10:30:00Z"
}
```

### Errore (404 Not Found)
```json
{
  "error": "Book not found"
}
```

## Concetti Go da Usare

- `net/http` package per il server HTTP
- `encoding/json` per Marshal/Unmarshal
- `json.NewEncoder(w).Encode()` per scrivere risposte
- `json.NewDecoder(r.Body).Decode()` per leggere richieste
- `http.HandleFunc()` o custom `http.Handler`
- `sync.RWMutex` per proteggere lo storage
- `http.ResponseWriter` e `http.Request`
- Path parameters extraction (es. usando `strings.Split()` o librerie esterne)

## Suggerimenti

- Crea funzioni helper per rispondere con JSON ed errori
- Usa defer `r.Body.Close()`
- Valida il `Content-Type` header per POST/PUT
- Genera ID univoci (es. con `uuid` o counter atomico)
- Organizza il codice in handlers separati
- Considera logging delle richieste

## Challenge Extra

- Middleware per logging
- Middleware per autenticazione (Basic Auth o JWT)
- Paginazione per GET /books
- Filtering e sorting (es. `?author=Donovan&sort=year`)
- CORS headers
- Persistenza su file/database
- Graceful shutdown
- Rate limiting per prevenire abuse

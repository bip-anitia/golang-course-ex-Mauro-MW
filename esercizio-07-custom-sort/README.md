# Esercizio 7: Custom Sort

## Obiettivo
Implementare ordinamenti personalizzati per tipi custom utilizzando l'interfaccia `sort.Interface` e le funzioni del package `sort`.

## Descrizione
Creare diversi tipi di dati e implementare ordinamenti personalizzati, inclusi ordinamenti multi-campo, ordinamenti personalizzati, e ordinamenti complessi.

## Requisiti

### 1. Implementare sort.Interface

L'interfaccia `sort.Interface` richiede tre metodi:
```go
type Interface interface {
    Len() int
    Less(i, j int) bool
    Swap(i, j int)
}
```

### 2. Tipi di Dati da Ordinare

#### Person
```go
type Person struct {
    Name     string
    Age      int
    Salary   float64
    City     string
    JoinDate time.Time
}
```

#### Movie
```go
type Movie struct {
    Title    string
    Year     int
    Rating   float64
    Duration int // minuti
    Genre    string
}
```

### 3. Ordinamenti da Implementare

Per **Person**:
- Per nome (alfabetico)
- Per età (crescente/decrescente)
- Per salario
- Multi-field: prima per città, poi per età
- Custom: per lunghezza del nome

Per **Movie**:
- Per anno
- Per rating
- Per durata
- Per genere, poi per rating
- Custom: per "score" calcolato (rating * year / duration)

## Implementazioni Richieste

### 1. Sort.Interface Tradizionale

```go
type ByAge []Person

func (a ByAge) Len() int           { /* TODO */ }
func (a ByAge) Less(i, j int) bool { /* TODO */ }
func (a ByAge) Swap(i, j int)      { /* TODO */ }

// Uso
people := []Person{...}
sort.Sort(ByAge(people))
```

### 2. Sort.Slice (più moderno)

```go
// Ordina per età
sort.Slice(people, func(i, j int) bool {
    return people[i].Age < people[j].Age
})

// Ordina per salario (decrescente)
sort.Slice(people, func(i, j int) bool {
    return people[i].Salary > people[j].Salary
})
```

### 3. Multi-field Sort

```go
// Prima per città, poi per età
sort.Slice(people, func(i, j int) bool {
    if people[i].City != people[j].City {
        return people[i].City < people[j].City
    }
    return people[i].Age < people[j].Age
})
```

### 4. Stable Sort

```go
// Mantiene l'ordine relativo di elementi uguali
sort.SliceStable(movies, func(i, j int) bool {
    return movies[i].Rating > movies[j].Rating
})
```

### 5. Sort con Reverse

```go
// Ordine inverso
sort.Sort(sort.Reverse(ByAge(people)))

// O con Slice
sort.Slice(people, func(i, j int) bool {
    return people[i].Age > people[j].Age  // inverte la condizione
})
```

## Esempi di Utilizzo

```bash
# Test di tutti gli ordinamenti
go run main.go

# Test specifico ordinamento
go run main.go -sort=age
go run main.go -sort=salary -reverse

# Con dati custom
go run main.go -input=people.json -sort=city,age
```

## Output Atteso

```
Original data:
  1. Alice    (28 years, $75000, New York)
  2. Bob      (35 years, $95000, Boston)
  3. Charlie  (28 years, $85000, New York)
  4. Diana    (42 years, $110000, Seattle)

Sorted by Age:
  1. Alice    (28 years, $75000, New York)
  2. Charlie  (28 years, $85000, New York)
  3. Bob      (35 years, $95000, Boston)
  4. Diana    (42 years, $110000, Seattle)

Sorted by Salary (descending):
  1. Diana    (42 years, $110000, Seattle)
  2. Bob      (35 years, $95000, Boston)
  3. Charlie  (28 years, $85000, New York)
  4. Alice    (28 years, $75000, New York)

Sorted by City, then Age:
  1. Bob      (35 years, $95000, Boston)
  2. Alice    (28 years, $75000, New York)
  3. Charlie  (28 years, $85000, New York)
  4. Diana    (42 years, $110000, Seattle)

---

Movies sorted by Rating:
  1. The Godfather      (1972, 9.2, 175min, Crime)
  2. The Dark Knight    (2008, 9.0, 152min, Action)
  3. Inception          (2010, 8.8, 148min, Sci-Fi)
  4. Interstellar       (2014, 8.6, 169min, Sci-Fi)

Movies sorted by Custom Score:
  1. The Dark Knight    (2008, 9.0, Score: 119.2)
  2. Inception          (2010, 8.8, Score: 119.5)
  3. The Godfather      (1972, 9.2, Score: 103.6)
  4. Interstellar       (2014, 8.6, Score: 102.4)
```

## Concetti Go da Usare

- `sort.Interface` per implementazioni custom
- `sort.Sort()` per ordinare tipi che implementano Interface
- `sort.Slice()` per ordinamenti inline con closure
- `sort.SliceStable()` per stable sort
- `sort.Reverse()` per ordinamento inverso
- `sort.IsSorted()` per verificare se ordinato
- `sort.Search()` per binary search
- Type conversion per riutilizzare sort types

## Struttura Suggerita

```go
package main

import (
    "fmt"
    "sort"
    "time"
)

type Person struct {
    Name     string
    Age      int
    Salary   float64
    City     string
    JoinDate time.Time
}

// Implementa vari sort types
type ByAge []Person
type BySalary []Person
type ByName []Person

// ByAge implementation
func (a ByAge) Len() int           { return len(a) }
func (a ByAge) Less(i, j int) bool { return a[i].Age < a[j].Age }
func (a ByAge) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// TODO: Implementare BySalary e ByName

// Multi-field sort helper
type PersonSorter struct {
    people []Person
    less   func(p1, p2 *Person) bool
}

func (s PersonSorter) Len() int           { return len(s.people) }
func (s PersonSorter) Swap(i, j int)      { s.people[i], s.people[j] = s.people[j], s.people[i] }
func (s PersonSorter) Less(i, j int) bool { return s.less(&s.people[i], &s.people[j]) }

// Funzioni helper
func SortBy(people []Person, less func(p1, p2 *Person) bool) {
    sort.Sort(PersonSorter{people, less})
}

func main() {
    people := []Person{
        {"Alice", 28, 75000, "New York", time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
        {"Bob", 35, 95000, "Boston", time.Date(2018, 6, 15, 0, 0, 0, 0, time.UTC)},
        {"Charlie", 28, 85000, "New York", time.Date(2019, 3, 10, 0, 0, 0, 0, time.UTC)},
        {"Diana", 42, 110000, "Seattle", time.Date(2015, 9, 1, 0, 0, 0, 0, time.UTC)},
    }

    // TODO: Implementare vari ordinamenti
}
```

## Esercizi Specifici

### Esercizio 1: Basic Sorting
Implementa `ByAge`, `BySalary`, `ByName` usando `sort.Interface`

### Esercizio 2: Reverse Sort
Ordina per età in ordine decrescente usando `sort.Reverse()`

### Esercizio 3: Multi-field Sort
Ordina prima per città (alfabetico), poi per età (crescente)

### Esercizio 4: Custom Comparison
Ordina persone per "seniority score" = Age * YearsInCompany * Salary

### Esercizio 5: Stable Sort
Usa `sort.SliceStable()` per ordinare mantenendo ordine relativo

### Esercizio 6: Search
Dopo ordinamento, usa `sort.Search()` per trovare elementi

### Esercizio 7: Generic Sorter
Crea una struttura generica per ordinare con funzione comparator custom

## Suggerimenti

### Performance
- `sort.Slice()` può essere più lento di `sort.Sort()` per grandi dataset
- `sort.SliceStable()` è più lento ma preserva ordine
- Pre-calcola valori se la funzione `Less` è costosa

### Best Practices
- Per ordinamenti semplici, preferisci `sort.Slice()`
- Per ordinamenti riutilizzabili, implementa `sort.Interface`
- Usa `sort.SliceStable()` quando l'ordine relativo è importante
- Verifica con `sort.IsSorted()` dopo ordinamento

### Pattern Comuni

```go
// Ordinamento con tie-breaking
sort.Slice(items, func(i, j int) bool {
    if items[i].Primary != items[j].Primary {
        return items[i].Primary < items[j].Primary
    }
    return items[i].Secondary < items[j].Secondary
})

// Ordinamento nullable/optional fields
sort.Slice(items, func(i, j int) bool {
    if items[i].Value == nil {
        return false
    }
    if items[j].Value == nil {
        return true
    }
    return *items[i].Value < *items[j].Value
})
```

## Challenge Extra

- **Comparator Chain**: Implementa catena di comparator riutilizzabili
- **Generic Sort**: Usa Go generics (1.18+) per sort generico
- **Locale-aware Sort**: Ordina stringhe con regole locale-specific
- **Parallel Sort**: Implementa merge sort parallelo
- **In-place vs Stable**: Confronta performance
- **Custom Data Structures**: Ordina linked list o tree
- **Top-K Selection**: Implementa partial sort per top K elementi
- **Median Finding**: Usa quickselect algorithm

## Testing

```go
func TestSorting(t *testing.T) {
    people := []Person{
        {"Charlie", 28, 85000, "New York", time.Now()},
        {"Alice", 28, 75000, "New York", time.Now()},
        {"Bob", 35, 95000, "Boston", time.Now()},
    }

    sort.Sort(ByAge(people))

    if !sort.IsSorted(ByAge(people)) {
        t.Error("People not sorted by age")
    }

    if people[0].Name != "Alice" && people[0].Name != "Charlie" {
        t.Error("First person should be Alice or Charlie (both 28)")
    }
}

func BenchmarkSortInterface(b *testing.B) {
    people := generatePeople(1000)
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        sort.Sort(ByAge(people))
    }
}

func BenchmarkSortSlice(b *testing.B) {
    people := generatePeople(1000)
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        sort.Slice(people, func(i, j int) bool {
            return people[i].Age < people[j].Age
        })
    }
}
```

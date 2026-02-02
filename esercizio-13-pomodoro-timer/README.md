# Esercizio 13: Pomodoro Timer

## Obiettivo
Imparare ad usare il package `time` per gestire timer, ticker, durations, parsing, e timezone.

## Descrizione
Creare un timer stile [Pomodoro Technique](https://en.wikipedia.org/wiki/Pomodoro_Technique): 25 minuti di lavoro, 5 minuti di pausa, con countdown visuale e notifiche.

## Requisiti

### 1. Pomodoro Base

Implementare ciclo base work/break:

```go
func pomodoro() {
    workDuration := 25 * time.Minute
    breakDuration := 5 * time.Minute

    fmt.Println("🍅 Work session started!")
    timer := time.NewTimer(workDuration)
    <-timer.C

    fmt.Println("☕ Break time!")
    timer.Reset(breakDuration)
    <-timer.C

    fmt.Println("✅ Pomodoro completed!")
}
```

### 2. Countdown con Ticker

Mostrare countdown in tempo reale:

```go
func countdown(duration time.Duration) {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()

    remaining := duration
    for remaining > 0 {
        fmt.Printf("\rTime remaining: %s ", formatDuration(remaining))
        <-ticker.C
        remaining -= time.Second
    }
    fmt.Println("\n⏰ Time's up!")
}
```

### 3. Time Formatting e Parsing

Gestire formati diversi:

```go
// Parse time da string
start, _ := time.Parse("15:04", "14:30")

// Format time
fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
fmt.Println(time.Now().Format("Mon Jan 2 15:04:05 2006"))
fmt.Println(time.Now().Format(time.RFC3339))

// Custom format
fmt.Println(time.Now().Format("3:04 PM"))
```

### 4. Timezone Handling

Mostrare orari in diverse timezone:

```go
func showWorldTime(t time.Time) {
    locations := map[string]string{
        "New York": "America/New_York",
        "London":   "Europe/London",
        "Tokyo":    "Asia/Tokyo",
        "Sydney":   "Australia/Sydney",
    }

    for city, tz := range locations {
        loc, _ := time.LoadLocation(tz)
        fmt.Printf("%s: %s\n", city, t.In(loc).Format("15:04 MST"))
    }
}
```

## Funzionalità da Implementare

### Feature 1: Pomodoro Completo

```go
type PomodoroConfig struct {
    WorkDuration  time.Duration
    BreakDuration time.Duration
    Cycles        int
}

func runPomodoro(config PomodoroConfig) {
    for i := 1; i <= config.Cycles; i++ {
        fmt.Printf("\n🍅 Pomodoro %d/%d\n", i, config.Cycles)

        // Work session
        fmt.Println("Work session started!")
        countdownWithTicker(config.WorkDuration)

        // Break (tranne ultima iterazione)
        if i < config.Cycles {
            fmt.Println("\n☕ Break time!")
            countdownWithTicker(config.BreakDuration)
        }
    }

    fmt.Println("\n✅ All pomodoros completed!")
}
```

### Feature 2: Stopwatch

Cronometro con start/stop/lap:

```go
type Stopwatch struct {
    start time.Time
    laps  []time.Duration
}

func (s *Stopwatch) Start() {
    s.start = time.Now()
    fmt.Println("⏱️  Stopwatch started")
}

func (s *Stopwatch) Lap() time.Duration {
    lap := time.Since(s.start)
    s.laps = append(s.laps, lap)
    fmt.Printf("Lap %d: %s\n", len(s.laps), lap)
    return lap
}

func (s *Stopwatch) Stop() time.Duration {
    total := time.Since(s.start)
    fmt.Printf("Total time: %s\n", total)
    return total
}
```

### Feature 3: Timer con Pausa/Resume

```go
type PausableTimer struct {
    duration  time.Duration
    remaining time.Duration
    paused    bool
    pauseTime time.Time
}

func (pt *PausableTimer) Pause() {
    if !pt.paused {
        pt.paused = true
        pt.pauseTime = time.Now()
        fmt.Println("⏸️  Timer paused")
    }
}

func (pt *PausableTimer) Resume() {
    if pt.paused {
        pt.paused = false
        pauseDuration := time.Since(pt.pauseTime)
        pt.remaining += pauseDuration
        fmt.Println("▶️  Timer resumed")
    }
}
```

### Feature 4: Meeting Time Calculator

Calcolare orario meeting in diverse timezone:

```go
func scheduleMeeting(timeStr string, fromTZ string) {
    // Parse time in timezone specifica
    loc, _ := time.LoadLocation(fromTZ)
    layout := "2006-01-02 15:04"
    meetingTime, _ := time.ParseInLocation(layout, timeStr, loc)

    fmt.Printf("Meeting scheduled for: %s\n", timeStr)
    fmt.Println("\nEquivalent times:")

    showInTimezone(meetingTime, "America/New_York", "New York")
    showInTimezone(meetingTime, "Europe/London", "London")
    showInTimezone(meetingTime, "Asia/Tokyo", "Tokyo")
}

func showInTimezone(t time.Time, tz string, name string) {
    loc, _ := time.LoadLocation(tz)
    localTime := t.In(loc)
    fmt.Printf("  %s: %s\n", name, localTime.Format("Mon 15:04 MST"))
}
```

## Esempi di Utilizzo

```bash
# Pomodoro standard (25min work, 5min break)
go run main.go

# Pomodoro custom
go run main.go -work=25m -break=5m -cycles=4

# Solo countdown
go run main.go countdown -duration=10m

# Stopwatch
go run main.go stopwatch

# Meeting scheduler
go run main.go meeting -time="2024-03-15 14:00" -timezone="America/New_York"

# World clock
go run main.go worldclock
```

## Output Atteso

### Pomodoro Session
```
🍅 Pomodoro Timer
Starting session: 25 min work, 5 min break

🍅 Pomodoro 1/4
Work session started!
Time remaining: 25:00
Time remaining: 24:59
Time remaining: 24:58
...
Time remaining: 00:01
⏰ Time's up!

☕ Break time!
Time remaining: 5:00
Time remaining: 4:59
...
⏰ Break over!

🍅 Pomodoro 2/4
...

✅ All 4 pomodoros completed!
Total time: 2h 0m
```

### Stopwatch
```
⏱️  Stopwatch
Commands: [l]ap, [s]top

Press Enter to start...
⏱️  Stopwatch started

> l
Lap 1: 5.234s

> l
Lap 2: 10.891s

> s
Stopped
Total time: 15.234s

Lap times:
  Lap 1: 5.234s
  Lap 2: 5.657s
```

### Meeting Scheduler
```
Meeting Time Calculator
Meeting scheduled for: 2024-03-15 14:00 EST

Equivalent times:
  New York:  Fri 14:00 EST
  London:    Fri 19:00 GMT
  Tokyo:     Sat 04:00 JST
  Sydney:    Sat 06:00 AEDT

Duration from now: 2 days 5 hours
```

### World Clock
```
🌍 World Clock - 2024-03-13 10:30:00

Americas:
  New York    10:30 EST  ☀️  Morning
  Los Angeles 07:30 PST  🌅 Early Morning
  São Paulo   11:30 BRT  ☀️  Late Morning

Europe:
  London      15:30 GMT  ☀️  Afternoon
  Paris       16:30 CET  ☀️  Afternoon
  Moscow      18:30 MSK  🌆 Evening

Asia/Pacific:
  Tokyo       00:30 JST  🌙 Night
  Sydney      02:30 AEDT 🌙 Night
  Dubai       19:30 GST  🌆 Evening
```

## Concetti Go da Usare

- `time.Duration` type
- `time.Timer` vs `time.Ticker` - differenze e quando usare
- `time.NewTimer()`, `timer.Reset()`, `timer.Stop()`
- `time.NewTicker()`, `ticker.Stop()`
- `time.Sleep()`, `time.After()`
- `time.Since()`, `time.Until()`
- `time.Parse()` e `time.ParseInLocation()`
- `time.Format()` con layout custom
- `time.LoadLocation()` per timezone
- `time.In()` per conversione timezone
- Duration arithmetic (`+`, `-`, `*`, `/`)
- `select` con time channels

## Struttura Suggerita

```go
package main

import (
    "fmt"
    "time"
)

// Pomodoro configuration
type Config struct {
    WorkDuration  time.Duration
    BreakDuration time.Duration
    LongBreak     time.Duration
    Cycles        int
}

// Default config
func DefaultConfig() Config {
    return Config{
        WorkDuration:  25 * time.Minute,
        BreakDuration: 5 * time.Minute,
        LongBreak:     15 * time.Minute,
        Cycles:        4,
    }
}

// Countdown with ticker
func countdown(duration time.Duration, label string) {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()

    end := time.Now().Add(duration)

    for time.Now().Before(end) {
        remaining := time.Until(end)
        fmt.Printf("\r%s %s   ", label, formatDuration(remaining))
        <-ticker.C
    }

    fmt.Println()
}

// Format duration nicely
func formatDuration(d time.Duration) string {
    d = d.Round(time.Second)
    h := d / time.Hour
    d -= h * time.Hour
    m := d / time.Minute
    d -= m * time.Minute
    s := d / time.Second

    if h > 0 {
        return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
    }
    return fmt.Sprintf("%02d:%02d", m, s)
}

// Run pomodoro session
func runPomodoro(config Config) {
    fmt.Println("🍅 Pomodoro Timer")
    fmt.Printf("Config: %s work, %s break, %d cycles\n\n",
        config.WorkDuration, config.BreakDuration, config.Cycles)

    for i := 1; i <= config.Cycles; i++ {
        fmt.Printf("🍅 Pomodoro %d/%d\n", i, config.Cycles)

        // Work session
        countdown(config.WorkDuration, "Work:")
        fmt.Println("⏰ Work session complete!")

        // Break
        if i < config.Cycles {
            fmt.Println()
            countdown(config.BreakDuration, "Break:")
            fmt.Println("☕ Break complete!\n")
        } else if i == config.Cycles {
            // Long break after last pomodoro
            fmt.Println()
            countdown(config.LongBreak, "Long break:")
            fmt.Println("🎉 Long break complete!")
        }
    }

    fmt.Println("\n✅ All pomodoros completed!")
}

// Simple timer
func simpleTimer(duration time.Duration) {
    fmt.Printf("Timer set for %s\n", duration)

    timer := time.NewTimer(duration)

    fmt.Println("Timer started...")
    <-timer.C
    fmt.Println("⏰ Timer finished!")
}

// Demonstrate Timer vs Ticker
func timerVsTicker() {
    fmt.Println("=== Timer vs Ticker Demo ===\n")

    // Timer: fires once
    fmt.Println("Timer (fires once):")
    timer := time.NewTimer(2 * time.Second)
    <-timer.C
    fmt.Println("Timer fired!")

    fmt.Println()

    // Ticker: fires repeatedly
    fmt.Println("Ticker (fires every second, 5 times):")
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()

    count := 0
    for range ticker.C {
        count++
        fmt.Printf("Tick %d\n", count)
        if count >= 5 {
            break
        }
    }
}

// Timezone demo
func worldClock() {
    locations := []struct {
        city string
        tz   string
    }{
        {"New York", "America/New_York"},
        {"London", "Europe/London"},
        {"Paris", "Europe/Paris"},
        {"Tokyo", "Asia/Tokyo"},
        {"Sydney", "Australia/Sydney"},
    }

    now := time.Now()
    fmt.Printf("🌍 World Clock - %s\n\n", now.Format("2006-01-02 15:04:05"))

    for _, loc := range locations {
        tz, err := time.LoadLocation(loc.tz)
        if err != nil {
            fmt.Printf("Error loading %s: %v\n", loc.city, err)
            continue
        }

        localTime := now.In(tz)
        fmt.Printf("%-12s %s\n", loc.city+":", localTime.Format("15:04 MST"))
    }
}

func main() {
    // Demo 1: Simple timer
    fmt.Println("=== Simple Timer ===")
    simpleTimer(3 * time.Second)
    fmt.Println()

    // Demo 2: Timer vs Ticker
    timerVsTicker()
    fmt.Println()

    // Demo 3: World clock
    worldClock()
    fmt.Println()

    // Demo 4: Pomodoro (shortened for demo)
    config := Config{
        WorkDuration:  5 * time.Second,  // Shortened for demo
        BreakDuration: 2 * time.Second,  // Shortened for demo
        LongBreak:     3 * time.Second,
        Cycles:        2,
    }
    runPomodoro(config)
}
```

## Suggerimenti

### Best Practices

1. **Sempre Stop ticker e timer**
   ```go
   ticker := time.NewTicker(time.Second)
   defer ticker.Stop() // Evita memory leak
   ```

2. **Timer vs Ticker**
   - `Timer`: Usa per "fai X dopo N tempo" (one-shot)
   - `Ticker`: Usa per "fai X ogni N tempo" (periodic)

3. **Reset vs New**
   ```go
   // ✅ Riusa timer
   timer := time.NewTimer(5 * time.Second)
   <-timer.C
   timer.Reset(5 * time.Second)

   // ❌ Crea nuovo timer ogni volta (meno efficiente)
   <-time.NewTimer(5 * time.Second).C
   <-time.NewTimer(5 * time.Second).C
   ```

4. **Duration parsing**
   ```go
   // Parsing da string
   d, err := time.ParseDuration("1h30m45s")

   // Common durations
   time.Millisecond
   time.Second
   time.Minute
   time.Hour
   ```

### Pattern Comuni

#### Timeout con select
```go
select {
case result := <-ch:
    fmt.Println("Got result:", result)
case <-time.After(5 * time.Second):
    fmt.Println("Timeout!")
}
```

#### Periodic task
```go
ticker := time.NewTicker(30 * time.Second)
defer ticker.Stop()

for {
    select {
    case <-ticker.C:
        doPeriodicTask()
    case <-done:
        return
    }
}
```

#### Debouncing
```go
func debounce(interval time.Duration, input chan string, output chan string) {
    var timer *time.Timer
    for {
        select {
        case val := <-input:
            if timer != nil {
                timer.Stop()
            }
            timer = time.AfterFunc(interval, func() {
                output <- val
            })
        }
    }
}
```

## Challenge Extra

- **Pomodoro Stats**: Traccia statistiche (pomodori completati per giorno)
- **Desktop Notification**: Notifiche native OS (usando librerie esterne)
- **Sound Alert**: Riproduci suono quando timer finisce
- **Config File**: Salva/carica configurazione da file JSON
- **Web Interface**: Interfaccia web per controllare timer
- **Multi-Timer**: Gestisci multipli timer simultaneamente
- **Recurring Events**: Scheduler per eventi ricorrenti (giornalieri, settimanali)

## Testing

```go
func TestFormatDuration(t *testing.T) {
    tests := []struct {
        duration time.Duration
        expected string
    }{
        {90 * time.Second, "01:30"},
        {3661 * time.Second, "01:01:01"},
        {59 * time.Second, "00:59"},
    }

    for _, tt := range tests {
        result := formatDuration(tt.duration)
        if result != tt.expected {
            t.Errorf("formatDuration(%v) = %s; want %s",
                tt.duration, result, tt.expected)
        }
    }
}

func TestTimerCompletion(t *testing.T) {
    start := time.Now()
    duration := 100 * time.Millisecond

    timer := time.NewTimer(duration)
    <-timer.C

    elapsed := time.Since(start)
    if elapsed < duration {
        t.Errorf("Timer completed too early: %v < %v", elapsed, duration)
    }
}
```

## Risorse

- [time package documentation](https://pkg.go.dev/time)
- [Go by Example: Timers](https://gobyexample.com/timers)
- [Go by Example: Tickers](https://gobyexample.com/tickers)

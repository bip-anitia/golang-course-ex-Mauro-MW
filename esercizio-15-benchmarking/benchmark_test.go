package main

import (
	"strings"
	"testing"
)

// Esempio: String concatenation benchmarks

// Metodo 1: operatore +=
func concatPlus(n int) string {
	s := ""
	for i := 0; i < n; i++ {
		s += "x"
	}
	return s
}

// Metodo 2: strings.Builder
func concatBuilder(n int) string {
	var builder strings.Builder
	for i := 0; i < n; i++ {
		builder.WriteString("x")
	}
	return builder.String()
}

// Benchmarks
func BenchmarkConcatPlus10(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = concatPlus(10)
	}
}

func BenchmarkConcatPlus100(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = concatPlus(100)
	}
}

func BenchmarkConcatBuilder10(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = concatBuilder(10)
	}
}

func BenchmarkConcatBuilder100(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = concatBuilder(100)
	}
}

// TODO: Aggiungere altri benchmark

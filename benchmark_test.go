package qrprng_test

import (
	"math/rand"
	"testing"

	"github.com/mattwiller/qrprng"
)

var BenchmarkResult uint64

func BenchmarkUInt64(b *testing.B) {
	prime := uint64(9_021_057_379)
	offset := uint64(1_000_014_012)
	intermediateOffset := uint64(2_947_624_585)
	rng, _ := qrprng.New(prime, intermediateOffset, offset)
	for i := 0; i < b.N; i++ {
		n := rng.Uint64()
		BenchmarkResult = n
	}
}

func BenchmarkDefault(b *testing.B) {
	rng := qrprng.Default()
	for i := 0; i < b.N; i++ {
		n := rng.Uint64()
		BenchmarkResult = n
	}
}

func BenchmarkStdLib(b *testing.B) {
	for i := 0; i < b.N; i++ {
		BenchmarkResult = rand.Uint64()
	}
}

func BenchmarkStdLibWithQRPRNGSource(b *testing.B) {
	rng := rand.New(qrprng.Default())
	for i := 0; i < b.N; i++ {
		BenchmarkResult = rng.Uint64()
	}
}

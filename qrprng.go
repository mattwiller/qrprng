package qrprng

import (
	"fmt"
	"math"
	"math/big"
	"math/bits"
)

const (
	INT63_MASK = (1 << 63) - 1
	// Largest prime (3 mod 4) less than 2^64, permutes [0, 2^64-189)
	DEFAULT_PRIME               = uint64(math.MaxUint64 - 188)
	DEFAULT_INTERMEDIATE_OFFSET = 5_577_006_791_947_779_410
)

// QuadraticResiduePRNG is a thread-unsafe PRNG based on Preshing's method using quadratic residues.
// The PRNG has the unique advantage of generating a permutation: when `offset` is 0, the output will cycle through
// all numbers less than `prime` without repeats (until all have been output and the cycle restarts).
// It implements both rand.Source and rand.Source64, and can be used via rand.New() to generate various random data.
type QuadraticResiduePRNG struct {
	prime              uint64
	intermediateOffset uint64
	offset             uint64
	maxMask            uint64
	mask               uint64
	idx                uint64
}

// New creates a new PRNG instance with the given parameters, which are validated for correctness before creation.
// The chosen prime must be 3 mod 4, and the intermediate offset (seed) can be any number less than the prime.  The
// offset will be added to all output, effectively placing a floor on the output values.
func New(prime, intermediateOffset, offset uint64) (*QuadraticResiduePRNG, error) {
	if err := validate(prime, offset, intermediateOffset); err != nil {
		return &QuadraticResiduePRNG{}, err
	}

	maxMask := calculateMaxMask(prime)
	return &QuadraticResiduePRNG{
		prime:              prime,
		intermediateOffset: intermediateOffset,
		offset:             offset,
		maxMask:            maxMask,
		mask:               calculateMask(prime, intermediateOffset, maxMask),
	}, nil
}

// Default returns a new PRNG instance suitable for general-purpose use.
// It uses the largest possible prime to permute 99.999999999999999% of possible uint64 values.
func Default() *QuadraticResiduePRNG {
	prng, err := New(DEFAULT_PRIME, DEFAULT_INTERMEDIATE_OFFSET, 0)
	if err != nil {
		panic(err)
	}
	return prng
}

// Index generates the ith element of the permutation described by the generator.
// If i >= prime, then an error is returned.  However, it can be ignored if desired;
// the sequence will simply cycle.
func (prng *QuadraticResiduePRNG) Index(i uint64) (uint64, error) {
	if i >= prng.prime {
		return i, fmt.Errorf("invalid index %d: must be less than chosen prime", i)
	}

	intermediate := prng.permuteQPR(i) + prng.intermediateOffset
	masked := prng.applyMask(intermediate % prng.prime)
	return prng.offset + prng.permuteQPR(masked), nil
}

func (prng *QuadraticResiduePRNG) applyMask(i uint64) uint64 {
	if i <= prng.maxMask {
		return i ^ prng.mask
	} else {
		return i
	}
}

func (prng *QuadraticResiduePRNG) permuteQPR(i uint64) uint64 {
	residue := (i * i) % prng.prime
	if i <= (prng.prime / 2) {
		return residue
	} else {
		return prng.prime - residue
	}
}

// QuadraticResiduePRNG implements math/rand.Source

func (prng *QuadraticResiduePRNG) Int63() int64 {
	return int64(prng.Uint64() & INT63_MASK)
}

// Seed changes the seed of the PRNG instance and resets the internal state of the generator.
func (prng *QuadraticResiduePRNG) Seed(seed int64) {
	if seed >= 0 {
		prng.intermediateOffset = uint64(seed)
	} else {
		prng.intermediateOffset = math.MaxUint64 - uint64(-1*seed)
	}
	prng.idx = 0
}

// QuadraticResiduePRNG implements math/rand.Source64

func (prng *QuadraticResiduePRNG) Uint64() uint64 {
	n, _ := prng.Index(prng.idx)
	prng.idx++
	return n
}

// ===== Private functions =====

func validate(prime, offset, intermediateOffset uint64) error {
	if prime%4 != 3 {
		return fmt.Errorf("invalid prime %d: must be 3 mod 4", prime)
	} else if intermediateOffset >= prime {
		return fmt.Errorf("invalid intermediate offset %d: must be less than chosen prime", intermediateOffset)
	} else if p := bigIntFromUint64(prime); !p.ProbablyPrime(0) {
		return fmt.Errorf("invalid prime %d: number is not prime", prime)
	}
	return nil
}

func calculateMaxMask(prime uint64) uint64 {
	primeBits := bits.Len64(prime - 1)
	return (1 << (primeBits - 1)) - 1
}

func calculateMask(prime, intermediateOffset, maxMask uint64) uint64 {
	min := uint64(1 << (bits.Len64(maxMask) - 1))
	return min + ((prime + intermediateOffset) % (maxMask - min))
}

func bigIntFromUint64(n uint64) *big.Int {
	var result *big.Int
	if n <= math.MaxInt64 {
		result = big.NewInt(int64(n))
	} else {
		result = big.NewInt(math.MaxInt64)
		result.Add(result, big.NewInt(int64(n-math.MaxInt64)))
	}
	return result
}

Quadratic Residue Pseudo-Random Number Generator
================================================

[![Documentation](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/mattwiller/qrprng)
[![GoReportCard example](https://goreportcard.com/badge/github.com/mattwiller/qrprng)](https://goreportcard.com/report/github.com/mattwiller/qrprng)



A pseudo-random sequence generator for Go using [Preshing's method][preshing], based on computing quadratic residues
modulo some prime number _p_:

n &equiv; x&sup2; (mod p)

The PRNG is reasonably fast: only about [50% slower](#benchmarks) than the default source in `math/rand`.  However, it also has the
added advantage that it produces a permuted sequence with no repeated values, which may be desirable in some contexts.

[preshing]: https://preshing.com/20121224/how-to-generate-a-sequence-of-unique-random-integers/

Installation
------------

```bash
go get -u github.com/mattwiller/qrprng
```

Usage
-----

### As `rand.Source`

For general-purpose use, the PRNG can be used as a source for all the functionality offered by [`math.Rand`][rand]:
```go
rng := rand.New(qrprng.Default())
rng.Seed(0xb0a710ad)
normDistDecimal := rng.NormFloat64()
```

Although the generator only produces `uint64` values directly, this allows it to be used in many different ways.

[rand]: https://pkg.go.dev/math/rand#Rand

### Random sequence access

The generator can calculate any term of the sequence in constant time as a `uint64`:

```go
rng := qrprng.Default()
n, _ := rng.Index(7_648_235_472)
```

### Custom generator

The parameters of the PRNG are fully customizable, and can use any prime _p_ where p &equiv; 3 (mod 4):

```go
p := unit64(11)
rng, _ := qrprng.New(p, 0, 0)

permutation := make([]uint64, p)
for i := uint64(0); i < p; i++ {
	permutation[i], _ = rng.Index(i)
}

fmt.Printf(permutation)
// [8 6 4 7 9 3 2 0 5 1 10]
```

Benchmarks
----------

```
goos: darwin
goarch: amd64
pkg: github.com/mattwiller/qrprng
cpu: Intel(R) Core(TM) i9-9980HK CPU @ 2.40GHz
BenchmarkUInt64-16                      60565976                20.61 ns/op            0 B/op          0 allocs/op
BenchmarkDefault-16                     55240932                21.19 ns/op            0 B/op          0 allocs/op
BenchmarkStdLib-16                      88121001                14.28 ns/op            0 B/op          0 allocs/op
BenchmarkStdLibWithQRPRNGSource-16      54257300                21.04 ns/op            0 B/op          0 allocs/op
PASS
ok      github.com/mattwiller/qrprng    12.468s
```

License
-------

Copyright 2022 Matthew Willer

Licensed under the [MIT License](./LICENSE.txt).  Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
express or implied. See the License for the specific language governing permissions and limitations under the License.

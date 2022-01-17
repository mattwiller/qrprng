Quadratic Residue Pseudo-Random Number Generator
================================================

A pseudo-random sequence generator for Go using [Preshing's method][preshing], based on computing quadratic residues
modulo some prime number _p_:

n &equiv; x&sup2; (mod p)

The PRNG is reasonably fast: only about 50% slower than the default source in `math/rand`.  However, it also has the
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

License
-------

Copyright 2022 Matthew Willer

Licensed under the [MIT License](./LICENSE.txt).  Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
express or implied. See the License for the specific language governing permissions and limitations under the License.

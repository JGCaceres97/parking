package ulid

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

var entropy *rand.Rand

func init() {
	entropy = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func GenerateNewULID() string {
	t := time.Now().UTC()

	// Asegura aleatoriedad y ordenamiento.
	monotonicEntropy := ulid.Monotonic(entropy, 0)

	// Crea ULID a partir del timestamp y la entropía.
	return ulid.MustNew(ulid.Timestamp(t), monotonicEntropy).String()
}

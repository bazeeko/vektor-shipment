package reference

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"
)

const maxRandomNumber = 9999999999999

type Generator struct {
	max *big.Int
}

func New() *Generator {
	return &Generator{
		max: big.NewInt(maxRandomNumber),
	}
}

func (g *Generator) GenerateReferenceNumber() (string, error) {
	randomNumber, err := rand.Int(rand.Reader, g.max)
	if err != nil {
		return "", fmt.Errorf("rand.Int: %w", err)
	}

	refNumberInt64 := time.Now().UnixNano() + randomNumber.Int64()
	refNumber := strings.ToUpper(strconv.FormatInt(refNumberInt64, 36))

	return refNumber, nil
}

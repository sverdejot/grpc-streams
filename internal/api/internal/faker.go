package internal

import (
	"math/rand/v2"

	"github.com/go-faker/faker/v4"
)

const (
	maxBidQty = 200_000_00
)

func FakeBid() Bid {
	return Bid{
		UserID:          faker.FirstName() + " " + faker.LastName(),
		QuantityInCents: int(rand.Int32N(maxBidQty)),
	}
}

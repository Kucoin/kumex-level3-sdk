package str

import (
	"errors"

	"github.com/shopspring/decimal"
)

func Diff(a string, b string) error {
	if a == b {
		return nil
	}

	aF, err := decimal.NewFromString(a)
	if err != nil {
		return err
	}
	bF, err := decimal.NewFromString(b)
	if err != nil {
		return err
	}

	if !aF.Equal(bF) {
		return errors.New("not equal: " + a + " != " + b)
	}

	return nil
}

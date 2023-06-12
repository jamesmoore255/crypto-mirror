package utils

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
)

func ConvertHexToDecimal(hexString string, decimalPlaces ...int32) (decimal.Decimal, error) {
	// Set default decimal places to 0
	var dp int32 = 0
	if len(decimalPlaces) > 0 {
		dp = decimalPlaces[0]
	}

	// Convert hex string to decimal
	bigInt, success := new(big.Int).SetString(hexString[2:], 16)
	if !success {
		return decimal.Decimal{}, fmt.Errorf("invalid hex string: %s", hexString)
	}

	// Set decimal places
	decimalValue := decimal.NewFromBigInt(bigInt, 0).Shift(-dp)

	return decimalValue, nil
}

func HexToDecimal(hex string) (decimal.Decimal, error) {
	bigInt, success := new(big.Int).SetString(hex[2:], 16)
	if !success {
		return decimal.Decimal{}, fmt.Errorf("failed to convert hex to decimal: %s", hex)
	}
	return decimal.NewFromBigInt(bigInt, 0), nil
}

func HexToDecimalHookFunc() mapstructure.DecodeHookFuncType {
	return func(from reflect.Type, to reflect.Type, data any) (any, error) {
		if to == reflect.TypeOf(decimal.Decimal{}) && from == reflect.TypeOf("") {
			hexString := data.(string)
			decimalValue, err := HexToDecimal(hexString)
			if err != nil {
				return nil, err
			}
			return decimalValue, nil
		}
		return data, nil
	}
}

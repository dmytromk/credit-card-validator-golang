package validation

import (
	"card_validator/internal"
	"card_validator/pb"
	"fmt"
	"strconv"
	"unicode"
)

type CardIssuer string

const (
	Mastercard CardIssuer = "Mastercard"
	Visa       CardIssuer = "Visa"
	UnionPay   CardIssuer = "UnionPay"
)

func IsDigit(s *string) *pb.ValidationError {
	if len(*s) == 0 {
		return &pb.ValidationError{Code: 4, Message: "Empty credit card number"}
	}

	for _, r := range *s {
		if !unicode.IsDigit(r) {
			return &pb.ValidationError{Code: 5, Message: "Card number must contain only digits"}
		}
	}

	return nil
}

func IssuerCheck(cardNumber *string) *pb.ValidationError {
	var issuer CardIssuer
	var requiredLength []int

	if len(*cardNumber) < 6 {
		return &pb.ValidationError{Code: 4, Message: "Incorrect card number length"}
	}

	switch {
	case (*cardNumber)[0] == '4':
		issuer = Visa
		requiredLength = []int{13, 16, 19}
	case (*cardNumber)[0:2] >= "51" && (*cardNumber)[0:2] <= "55":
		issuer = Mastercard
		requiredLength = []int{16}
	case (*cardNumber)[0:2] == "62" || (*cardNumber)[0:2] == "81":
		issuer = UnionPay
		requiredLength = []int{16, 19}
	default:
		return &pb.ValidationError{Code: 3, Message: "Unknown credit card issuer"}
	}

	if !internal.CheckLength(cardNumber, &requiredLength) {
		return &pb.ValidationError{Code: 4,
			Message: fmt.Sprintf("%s doesn't have card with length %d", issuer, len(*cardNumber)),
		}
	}

	return nil
}

func LuhnCheck(cardNumber *string) *pb.ValidationError {
	sumResult := 0
	isSecondDigit := false

	for i := len(*cardNumber) - 1; i >= 0; i-- {
		digit, err := strconv.Atoi(string((*cardNumber)[i]))
		if err != nil {
			return &pb.ValidationError{Code: 1, Message: "Encountered invalid character in card number"}
		}

		if isSecondDigit {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sumResult += digit

		isSecondDigit = !isSecondDigit
	}

	if sumResult%10 != 0 {
		return &pb.ValidationError{Code: 2, Message: "Card number failed Luhn check"}
	}

	return nil
}

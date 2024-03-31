package validation

import (
	"card_validator/pb"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type CardIssuer string

const (
	Mastercard CardIssuer = "Mastercard"
	Visa       CardIssuer = "Visa"
	UnionPay   CardIssuer = "UnionPay"
)

func RemoveWhitespace(s *string) string {
	var sb strings.Builder

	for _, r := range *s {
		if !unicode.IsSpace(r) {
			sb.WriteRune(r)
		}
	}

	return sb.String()
}

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

func checkLength(cardNumber *string, possibleLengths *[]int) bool {
	for _, length := range *possibleLengths {
		if len(*cardNumber) == length {
			return true
		}
	}
	return false
}

func IssuerCheck(cardNumber *string) *pb.ValidationError {
	var issuer CardIssuer
	var requiredLength []int

	_, err := strconv.Atoi((*cardNumber)[0:2])
	if err != nil {
		return &pb.ValidationError{Code: 1, Message: "Encountered invalid character in card number"}
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

	if !checkLength(cardNumber, &requiredLength) {
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

func Expiration(expiration_month *string, expiration_year *string) *pb.ValidationError {
	year, year_exception := strconv.Atoi(*expiration_year)
	month, month_exception := strconv.Atoi(*expiration_month)

	if year_exception != nil || year < 1 {
		{
			return &pb.ValidationError{Code: 6, Message: "Incorrect year input"}
		}
	}

	if month_exception != nil || month < 1 || month > 12 {
		{
			return &pb.ValidationError{Code: 7, Message: "Incorrect month input"}
		}
	}

	if year < time.Now().UTC().Year() {
		return &pb.ValidationError{Code: 8, Message: "Card is expired"}
	}

	// Check the expired  year and month
	if year == time.Now().UTC().Year() && month < int(time.Now().UTC().Month()) {
		return &pb.ValidationError{Code: 8, Message: "Card is expired"}
	}

	return nil
}

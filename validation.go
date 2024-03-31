package main

import (
	"card_validator/pb"
	"fmt"
	"strconv"
	"time"
	"unicode"
)

type CardIssuer string

const (
	Mastercard CardIssuer = "Mastercard"
	Visa       CardIssuer = "Visa"
	UnionPay   CardIssuer = "UnionPay"
)

const (
	EmptyField           int32 = 10
	IncorrectFieldFormat       = 20
	UnknownCardIssuer          = 30
	ExpiredCard                = 40
)

func IsDigit(s *string) *pb.ValidationError {
	if len(*s) == 0 {
		return &pb.ValidationError{Code: EmptyField, Message: "Empty credit card number"}
	}

	for _, r := range *s {
		if !unicode.IsDigit(r) {
			return &pb.ValidationError{Code: IncorrectFieldFormat, Message: "Card number must contain only digits"}
		}
	}

	return nil
}

func IssuerCheck(cardNumber *string) *pb.ValidationError {
	var issuer CardIssuer
	var requiredLength []int

	if len(*cardNumber) < 6 {
		return &pb.ValidationError{Code: IncorrectFieldFormat, Message: "Incorrect card number length"}
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
		return &pb.ValidationError{Code: UnknownCardIssuer, Message: "Unknown credit card issuer"}
	}

	if !CheckLength(cardNumber, &requiredLength) {
		return &pb.ValidationError{Code: IncorrectFieldFormat,
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
			return &pb.ValidationError{Code: IncorrectFieldFormat, Message: "Encountered invalid character in card number"}
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

func Expiration(expirationMonth *string, expirationYear *string) *pb.ValidationError {
	year, yearException := strconv.Atoi(*expirationYear)
	month, monthException := strconv.Atoi(*expirationMonth)

	if yearException != nil || year < 1 {
		{
			return &pb.ValidationError{Code: IncorrectFieldFormat, Message: "Incorrect year input"}
		}
	}

	if monthException != nil || month < 1 || month > 12 {
		{
			return &pb.ValidationError{Code: IncorrectFieldFormat, Message: "Incorrect month input"}
		}
	}

	if year < time.Now().UTC().Year() {
		return &pb.ValidationError{Code: ExpiredCard, Message: "Card is expired"}
	}

	if year == time.Now().UTC().Year() && month < int(time.Now().UTC().Month()) {
		return &pb.ValidationError{Code: ExpiredCard, Message: "Card is expired"}
	}

	return nil
}

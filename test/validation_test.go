package test

import (
	"card_validator/pb"
	"card_validator/validation"
	"testing"
)

func TestIsDigit(t *testing.T) {
	testCases := []struct {
		input    string
		expected *pb.ValidationError
	}{
		{"1234567890", nil},
		{"12abc34567890", &pb.ValidationError{Code: validation.IncorrectFieldFormat, Message: "Card number must contain only digits"}},
		{"", &pb.ValidationError{Code: validation.EmptyField, Message: "Empty credit card number"}},
	}

	for _, tc := range testCases {
		result := validation.IsDigit(&tc.input)
		if result == nil && tc.expected != nil {
			t.Errorf("Expected error, got nil for input: %s", tc.input)
		}
		if result != nil && tc.expected == nil {
			t.Errorf("Expected nil, got error for input: %s", tc.input)
		}
		if result != nil && tc.expected != nil && (result.Code != tc.expected.Code || result.Message != tc.expected.Message) {
			t.Errorf("Mismatched error for input: %s, Expected: %+v, Got: %+v", tc.input, tc.expected, result)
		}
	}
}

func TestIssuerCheck(t *testing.T) {
	testCases := []struct {
		input    string
		expected *pb.ValidationError
	}{
		{"4111111111111111", nil},
		{"5211111111111111", nil},
		{"5611111111111111", &pb.ValidationError{Code: validation.UnknownCardIssuer, Message: "Unknown credit card issuer"}},
		{"", &pb.ValidationError{Code: validation.IncorrectFieldFormat, Message: "Incorrect card number length"}},
	}

	for _, tc := range testCases {
		result := validation.IssuerCheck(&tc.input)
		if result == nil && tc.expected != nil {
			t.Errorf("Expected error, got nil for input: %s", tc.input)
		}
		if result != nil && tc.expected == nil {
			t.Errorf("Expected nil, got error for input: %s", tc.input)
		}
		if result != nil && tc.expected != nil && (result.Code != tc.expected.Code || result.Message != tc.expected.Message) {
			t.Errorf("Mismatched error for input: %s, Expected: %+v, Got: %+v", tc.input, tc.expected, result)
		}
	}
}

func TestLuhnCheck(t *testing.T) {
	testCases := []struct {
		input    string
		expected *pb.ValidationError
	}{
		{"4111111111111111", nil},
		{"4111111111111110", &pb.ValidationError{Code: validation.FailedLuhnCheck, Message: "Card number failed Luhn check"}},
	}

	for _, tc := range testCases {
		result := validation.LuhnCheck(&tc.input)
		if result == nil && tc.expected != nil {
			t.Errorf("Expected error, got nil for input: %s", tc.input)
		}
		if result != nil && tc.expected == nil {
			t.Errorf("Expected nil, got error for input: %s", tc.input)
		}
		if result != nil && tc.expected != nil && (result.Code != tc.expected.Code || result.Message != tc.expected.Message) {
			t.Errorf("Mismatched error for input: %s, Expected: %+v, Got: %+v", tc.input, tc.expected, result)
		}
	}
}

func TestExpiration(t *testing.T) {
	testCases := []struct {
		month    string
		year     string
		expected *pb.ValidationError
	}{
		{"12", "2100", nil},
		{"01", "2020", &pb.ValidationError{Code: validation.ExpiredCard, Message: "Card is expired"}},
		{"13", "2025", &pb.ValidationError{Code: validation.IncorrectFieldFormat, Message: "Incorrect month input"}},
		{"01", "0", &pb.ValidationError{Code: validation.IncorrectFieldFormat, Message: "Incorrect year input"}},
	}

	for _, tc := range testCases {
		result := validation.Expiration(&tc.month, &tc.year)
		if result == nil && tc.expected != nil {
			t.Errorf("Expected error, got nil for input: %s/%s", tc.month, tc.year)
		}
		if result != nil && tc.expected == nil {
			t.Errorf("Expected nil, got error for input: %s/%s", tc.month, tc.year)
		}
		if result != nil && tc.expected != nil && (result.Code != tc.expected.Code || result.Message != tc.expected.Message) {
			t.Errorf("Mismatched error for input: %s/%s, Expected: %+v, Got: %+v", tc.month, tc.year, tc.expected, result)
		}
	}
}

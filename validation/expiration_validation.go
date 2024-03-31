package validation

import (
	"card_validator/pb"
	"strconv"
	"time"
)

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

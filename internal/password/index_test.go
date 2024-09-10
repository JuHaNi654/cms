package password

import (
	"encoding/base64"
	"fmt"
	"strings"
	"testing"
)

func TestHash(t *testing.T) {
	isEqual := func(t testing.TB, value, expected any) {
		t.Helper()
		if value != expected {
			t.Errorf("got '%v' :: want '%v'", value, expected)
		}
	}

	isLarger := func(t testing.TB, value, expected int) {
		t.Helper()
		if value <= expected {
			t.Errorf("value (%d) is less or equal than expected value (%d)", value, expected)
		}
	}

	checkValidParsedData := func(t testing.TB, input, format string, expected int) {
		t.Helper()
		if input != fmt.Sprintf("%s%d", format, expected) {
			t.Errorf("expected %s to contain %d", input, expected)
		}
	}

	value := "ExamplePassword123"
	result, err := Hash(value)
	if err != nil {
		t.Fatalf("could not generate valid hash: %v", err)
	}

	values := strings.Split(result, "$")
	info := strings.Split(values[3], ",")

	decodedSalt, err := base64.RawStdEncoding.Strict().Strict().DecodeString(values[4])
	if err != nil {
		t.Fatalf("could not decode salt: %v", err)
	}

	checkValidParsedData(t, info[0], "m=", memory)
	checkValidParsedData(t, info[1], "t=", iterations)
	checkValidParsedData(t, info[2], "p=", parallelism)
	isEqual(t, len(decodedSalt), saltLength)
	isLarger(t, len(values[5]), 0)
}

func TestCompare(t *testing.T) {
	isInvalid := func(t testing.TB, hash, password string) {
		t.Helper()
		valid, err := Compare(password, hash)
		if err != nil {
			t.Fatalf("error occured while comparing passwords: %v", err)
		}

		if valid {
			t.Fatalf("compare functio returned true from invalid password")
		}
	}

	isValid := func(t testing.TB, hash, password string) {
		t.Helper()
		valid, err := Compare(password, hash)
		if err != nil {
			t.Fatalf("error occured while comparing passwords: %v", err)
		}

		if !valid {
			t.Fatalf("compare functio returned false from valid password")
		}
	}

	var (
		value  string
		result string
		err    error
	)
	value = "ExamplePassword123"
	result, err = Hash(value)
	if err != nil {
		t.Fatalf("could not generate valid hash: %v", err)
	}

	isInvalid(t, result, "InvalidPassword")
	isValid(t, result, value)

	value = "Ex5$'Å{][]}!?"
	result, err = Hash(value)
	if err != nil {
		t.Fatalf("could not generate valid hash: %v", err)
	}

	isInvalid(t, result, "InvalidPassword")
	isValid(t, result, value)
}

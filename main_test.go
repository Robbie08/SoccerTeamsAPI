package main

// Simple tests using Go testing
import (
	"log"
	"testing"
)

// This function will test if the GenerateUID function is actually
// 1) returing uint64 values and 2) if the hashing is consistent and
// correct with their respective strings.
func TestHashingAlgorithm(t *testing.T) {
	log.Printf("\n----- Testing GenerateUID(string) -----")
	if GenerateUID("Real Madrid") < 0 {
		t.Error("Value returned is not a valid UID")
	}

	if GenerateUID("Real Madrid") != 13243162659090910742 {
		t.Error("The correct UID sould be: 13243162659090910742")
	}

	if GenerateUID("Juventus") != 10230320473718180607 {
		t.Error("The correct UID should be: 10230320473718180607")
	}
}

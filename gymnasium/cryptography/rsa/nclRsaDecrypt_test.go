package main

import "testing"

func TestNclDecryptRSA(t *testing.T) {
	n, e := 1079, 43
	ciphertext := "996 894 379 631 894 82 379 852 631 677 677 194 893"

	want := "SKY-KRYG-5530"
	got, err := nclDecryptRSA(n, e, ciphertext)
	if want != got || err != nil {
		t.Fatalf("wanted: %s, got: %s; err: %v", want, got, err)
	}
}

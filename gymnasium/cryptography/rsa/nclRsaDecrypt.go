package main

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
)

func leastCommonPrimeFactor(num int) int {
	var min = math.MaxInt
	for res := 2; num > 1; res++ {
		if num%res == 0 {
			x := 0
			for num%res == 0 {
				num /= res
				x++
			}

			if res < min {
				min = res
			}
		}
	}
	return min
}

func extendedGCD(a, b int) (int, int, int) {
	if a == 0 {
		return b, 0, 1
	}
	gcd, xPrime, yPrime := extendedGCD(b%a, a)
	return gcd, yPrime - (b/a)*xPrime, xPrime
}

func inverse(a, m int) (int, error) {
	gcd, x, _ := extendedGCD(a, m)
	if gcd != 1 || m == 0 {
		return 0, fmt.Errorf("no Modular Inverse exists")
	}

	return ((m + (x % m)) % m), nil
}

func calculateThirdPrivateKeyValue(e, p, q int) (int, error) {
	a := (p - 1) * (q - 1)

	invMod, err := inverse(e, a)
	if err != nil {
		return 0, err
	}

	return invMod, nil
}

func stringsToInts(strs []string) ([]int, error) {
	var res []int

	for _, i := range strs {
		num, err := strconv.Atoi(i)
		if err != nil {
			return nil, err
		}

		res = append(res, num)
	}

	return res, nil
}

func nclDecryptRSA(n, e int, ciphertext string) (string, error) {
	q := leastCommonPrimeFactor(n)
	p := n / q

	d, err := calculateThirdPrivateKeyValue(e, p, q)
	if err != nil {
		return "", err
	}

	tokens := strings.Fields(ciphertext)
	encryptedVals, err := stringsToInts(tokens)
	if err != nil {
		return "", err
	}

	exp := big.NewInt(int64(d))
	mod := big.NewInt(int64(n))

	var s strings.Builder
	s.Grow(len(encryptedVals))
	for _, num := range encryptedVals {
		temp := new(big.Int).Exp(big.NewInt(int64(num)), exp, nil)

		s.WriteRune(rune(new(big.Int).Mod(temp, mod).Int64()))
	}

	return s.String(), nil
}

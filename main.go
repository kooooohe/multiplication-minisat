package main

import (
	"fmt"
	"os"
	"strings"
)

// Clause represents a CNF clause as a slice of integers
type Clause []int

// Generate CNF for multiplication of two binary numbers
func generateCnfForMultiplication(n, m int) ([]Clause, int) {
	var clauses []Clause
	counter := 0

	// Helper function to return a new variable
	newVar := func() int {
		counter++
		return counter
	}

	// AND gate CNF for bit multiplication
	andGate := func(a, b, output int) {
		clauses = append(clauses, Clause{-a, -b, output})
		clauses = append(clauses, Clause{a, -output})
		clauses = append(clauses, Clause{b, -output})
	}

	// Full adder CNF // TODO FIX
	fullAdder := func(a, b, cin, s, cout int) {
		sum := newVar()
		clauses = append(clauses, Clause{-a, -b, sum})
		clauses = append(clauses, Clause{a, b, sum})
		clauses = append(clauses, Clause{-sum, -cin, s})
		clauses = append(clauses, Clause{sum, cin, s})
		clauses = append(clauses, Clause{-a, -cin, cout})
		clauses = append(clauses, Clause{-b, -cin, cout})
		clauses = append(clauses, Clause{a, b, cout})
	}

	// Create variables for the inputs and outputs
	inputA := make([]int, n)
	inputB := make([]int, m)
	for i := range inputA {
		inputA[i] = newVar()
	}
	for i := range inputB {
		inputB[i] = newVar()
	}

	// Compute bit products (Step 1: Bitwise multiplication)
	// TODO: この時点でn+m bitでやってよい
	// その後それの足し算を純粋に行う
	productBits := make([][]int, n)
	for i := range productBits {
		productBits[i] = make([]int, m)
		for j := range productBits[i] {
			productBits[i][j] = newVar()
			andGate(inputA[i], inputB[j], productBits[i][j])
		}
	}

	// Add the bit products with shift (Step 2: Shifted addition)
	sumVars := make([]int, n+m-1)
	carry := 0
	for bit := 0; bit < n+m-1; bit++ {
		bitValues := []int{}
		for i := 0; i < n; i++ {
			if j := bit - i; j >= 0 && j < m {
				bitValues = append(bitValues, productBits[i][j])
			}
		}
		if carry != 0 {
			bitValues = append(bitValues, carry)
		}
		if len(bitValues) > 0 {
			sumVars[bit] = newVar()
			newCarry := 0
			if bit < n+m-2 {
				newCarry = newVar()
			}
			fullAdder(bitValues[0], bitValues[1], bitValues[2], sumVars[bit], newCarry)
			carry = newCarry
		}
	}

	return clauses, counter
}

// Convert clauses to string in DIMACS format
func clausesToString(clauses []Clause, varCount int) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("p cnf %d %d\n", varCount, len(clauses)))
	for _, clause := range clauses {
		for _, lit := range clause {
			sb.WriteString(fmt.Sprintf("%d ", lit))
		}
		sb.WriteString("0\n")
	}
	return sb.String()
}

func main() {
	// Example usage: Generate CNF for 4-bit by 3-bit multiplication
	n := 4
	m := 3
	clauses, varCount := generateCnfForMultiplication(n, m)
	cnf := clausesToString(clauses, varCount)

	// Save the CNF to a text file
	filename := "multiplication_cnf.txt"
	if err := os.WriteFile(filename, []byte(cnf), 0644); err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
	fmt.Println("CNF file generated successfully:", filename)
}


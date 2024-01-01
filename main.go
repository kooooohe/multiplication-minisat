package main

import (
	"fmt"
	"os"
	"strings"
)

// Clause represents a CNF clause as a slice of integers
type Clause []int

// Helper function to return a new variable
func intSeq() func() int {
	c := 0
	return func() int {
		c++
		return c
	}
}

// Generate CNF for multiplication of two binary numbers
func generateCnfForMultiplication(n, m int, a, b []byte) ([]Clause, int) {
	clauses := []Clause{}

	newVar := intSeq()

	// AND gate CNF for bit multiplication
	andGate := func(a, b, output int) {
		clauses = append(clauses, Clause{-a, -b, output})
		clauses = append(clauses, Clause{a, -output})
		clauses = append(clauses, Clause{b, -output})
	}

	// Full adder CNF // TODO FIX
	fullAdder := func(a, b, cin, s, cout int) {

		t := newVar()
		t2 := newVar()
		t3 := newVar()
		t4 := newVar()

		tmpClauses := []Clause{
			//for s = a xor b xor x
			//a XOR b = t
			//t XOR cin = s
			{a, b, -t}, {-a, -b, -t}, {a, -b, t}, {-a, b, t},
			{t, cin, -s}, {-t, -cin, -s}, {t, -cin, s}, {-t, cin, s},

			// for c_out (at least two)
			// t2 = a AND b
			{a, -t2}, {b, -t2}, {-a, -b, t2},
			// t3 = a AND c_in
			{a, -t3}, {cin, -t3}, {-a, -cin, t3},
			// t4 = b AND c_in
			{b, -t4}, {cin, -t4}, {-b, -cin, t4},
			//c_out = t2 OR t3 OR t4
			{-cout, t2, t3, t4},
			{cout, -t2}, {cout, -t3}, {cout, -t4},
		}

		clauses = append(clauses, tmpClauses...)

	}

	// Create variables for the inputs and outputs
	inputA := make([]int, n)
	inputB := make([]int, m)
	for i := range inputA {
		inputA[i] = newVar()
		if a[n-i-1] == 1 {
			clauses = append(clauses, Clause{inputA[i]})
		} else {
			clauses = append(clauses, Clause{-inputA[i]})
		}
	}
	for i := range inputB {
		inputB[i] = newVar()
		if b[m-i-1] == 1 {
			clauses = append(clauses, Clause{inputB[i]})
		} else {
			clauses = append(clauses, Clause{-inputB[i]})
		}
	}

	// Compute bit products (Step 1: Bitwise multiplication)
	productBits := make([][]int, n)
	for i := range productBits {
		productBits[i] = make([]int, m)
		for j := range productBits[i] {
			productBits[i][j] = newVar()
			andGate(inputA[i], inputB[j], productBits[i][j])
		}
	}

	productBitsNM := make([][]int, n)
	for i := 0; i < n; i++ {
		productBitsNM[i] = make([]int, n+m)
		for j := 0; j < m; j++ {
			// shift
			productBitsNM[i][j+i] = productBits[i][j]
		}
		// pat space as 0
		for j := 0; j < n+m; j++ {
			if productBitsNM[i][j] == 0 {
				t := newVar()
				productBitsNM[i][j] = t
				clauses = append(clauses, Clause{-t})
			}
		}
	}
	// fmt.Println(productBitsNM)
	// fmt.Println("")

	// Add the bit products with shift (Step 2: Shifted addition)

	pSumVars := make([]int, n+m)
	for i := 0; i < n; i++ {
		sumVars := make([]int, n+m)
		for i := 0; i < len(sumVars); i++ {
			sumVars[i] = newVar()
		}
		carryVars := make([]int, n+m+1)

		for i := 0; i < len(carryVars); i++ {
			carryVars[i] = newVar()
		}

		// carryvars[0] is always false
		clauses = append(clauses, Clause{-carryVars[0]})

		// fmt.Println(sumVars)
		// fmt.Println(pSumVars)
		// fmt.Println(carryVars)

		// sumVars are 0  where i == 0
		allzero := make([]int, n+m)
		if i == 0 {
			for i := 0; i < len(allzero); i++ {
				allzero[i] = newVar()
				clauses = append(clauses, Clause{-allzero[i]})
			}
			// fmt.Println("allzero: ", allzero)
		}
		// fmt.Println("a", productBitsNM[i])
		for j := 0; j < n+m; j++ {
			a := productBitsNM[i][j]
			b := 0
			if i == 0 {
				b = allzero[j]
			} else {
				b = pSumVars[j]
			}
			// fmt.Println("b", b)
			cin := carryVars[j]
			s := sumVars[j]
			cout := carryVars[j+1]
			fullAdder(a, b, cin, s, cout)
			// fmt.Println("s", s)

		}
		copy(pSumVars, sumVars)
		fmt.Println("psum: ", pSumVars)
	}

	// this is the number showing a result
	fmt.Println("result:", pSumVars)
	return clauses, newVar() - 1
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
	a := []byte{0, 1, 1}
	b := []byte{0, 1, 1}
	n := len(a)
	m := len(b)
	clauses, varCount := generateCnfForMultiplication(n, m, a, b)
	cnf := clausesToString(clauses, varCount)

	// Save the CNF to a text file
	filename := "multiplication_cnf.txt"
	if err := os.WriteFile(filename, []byte(cnf), 0644); err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
	fmt.Println("CNF file generated successfully:", filename)
}

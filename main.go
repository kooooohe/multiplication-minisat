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
func generateCnfForMultiplication(n, m int) ([]Clause, int) {
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
			//for s a xor b xor s
			//a XOR b = t
			//t XOR cin = s
			{-a, -b, t}, {a, b, t}, {a, -b, -t}, {-a, b, -t},
			{-t, -cin, s}, {t, cin, s}, {t, -cin, -s}, {-t, cin, -s},

			// for c_out (at least two)
			// t2 = a AND b
			// (-a OR t2) AND (-b OR t2) AND (a OR b OR -t2)
			{-a, t2}, {-b, t2}, {a, b, -t2},
			// t3 = a AND c_in
			// (-a OR t3) AND (-c OR t3) AND (a OR c OR -t3)
			{-a, t3}, {-cin, t3}, {a, cin, -t3},
			// t4 = b AND c_in
			//(-b OR t4) AND (-c OR t4) AND (b OR c OR -t4)
			{-b, t4}, {-cin, t4}, {b, cin, -t4},
			//c_out = t2 OR t3 OR t4
			// (-t2 OR c_out) AND (-t3 OR c_out) AND (-t4 OR c_out) AND (t2 OR t3 OR t4 OR -c_out)
			{-t2, cout}, {-t3, cout}, {-t4, cout}, {t2, t3, t4, -cout},
		}

		clauses = append(clauses, tmpClauses...)
		/*
			clauses = append(clauses, Clause{-a, -b, sum})
			clauses = append(clauses, Clause{a, b, sum})
			clauses = append(clauses, Clause{-sum, -cin, s})
			clauses = append(clauses, Clause{sum, cin, s})
			clauses = append(clauses, Clause{-a, -cin, cout})
			clauses = append(clauses, Clause{-b, -cin, cout})
			clauses = append(clauses, Clause{a, b, cout})
		*/

	}

	// Create variables for the inputs and outputs
	//TODOここに入れる
	inputA := make([]int, n)
	inputB := make([]int, m)
	for i := range inputA {
		inputA[i] = newVar()
	}
	for i := range inputB {
		inputB[i] = newVar()
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

	// return clauses, newVar() -1

	productBitsNM := make([][]int, n)
	for i := 0; i < n; i++ {
		productBitsNM[i] = make([]int, n+m)
		for j := 0; j < m; j++ {
			// shift
			productBitsNM[i][j+i] = productBits[i][j]
		}
		// pat as 0
		// TODO make it 0
		for k := 0; k < i; k++ {
			productBitsNM[i][k] = newVar()
		}

		// pat as 0 上位
		// TODO make it 0
		for k:= n+m-1;k >= 0; k-- {
			productBitsNM[i][k] = newVar()
		}
	}
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

		// TODO carryvars[0]の時は必ずfalseになるようにする
		fmt.Println(sumVars)
		fmt.Println(pSumVars)
		fmt.Println(carryVars)

		for j := 0; j < n+m; j++ {
			a := productBitsNM[i][j]
			b := 0
			if i == 0 {
				// TODO sumVars iが0番目は全部false == 0
				allzero := make([]int, n+m)
				for i := 0; i < len(allzero); i++ {
					allzero[i] = newVar()
				}
				b = allzero[j]
			} else {
				b = pSumVars[j]
			}
			cin := carryVars[j]
			s := sumVars[j]
			cout := carryVars[j+1]
			fullAdder(a, b, cin, s, cout)

		}
		copy(pSumVars, sumVars)
	}

	/*
		sumVars := make([]int, n+m)

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
	*/

	//TODO 最後のsumvarを返す
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

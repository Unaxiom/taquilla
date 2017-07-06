package taquilla

import (
	"fmt"

	"github.com/gonum/matrix/mat64"
)

// computeEquationSolution accepts two matrices A, B, where AX=B and computes X, such that X=A(^-1) x B, and a list of variables
func computeEquationSolution(eqnMatrix *mat64.Dense, resultMatrix *mat64.Dense, ticketTypeList []string) {
	// Calculate inverse of A
	eqnMatrixRows, eqnMatrixCols := eqnMatrix.Caps()
	inv := mat64.NewDense(eqnMatrixRows, eqnMatrixCols, nil)

	err := inv.Inverse(eqnMatrix)
	if err != nil {
		// log.Fatalln(err)
		return
	}

	// Calculate mul of A inverse and B
	mul := mat64.NewDense(eqnMatrixRows, 1, nil)
	mul.Mul(inv, resultMatrix)

	// log.ErrorDump("Solution Matrix Is: \n\n", mul)
	for i, ticketType := range ticketTypeList {
		semGraph.updateMem(ticketType, mul.At(i, 0))

		fmt.Println("===================================================================")
		fmt.Println("Type --> ", ticketType, ", Soln --> ", mul.At(i, 0))
		fmt.Println("Type --> ", semGraph.container[ticketType].name, ", Soln --> ", semGraph.container[ticketType].avgMem)
		fmt.Println("===================================================================")
	}
}

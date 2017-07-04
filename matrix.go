package taquilla

import (
	"github.com/gonum/matrix/mat64"
)

// computeEquationSolution accepts two matrices A, B, where AX=B and computes X, such that X=A(^-1) x B
func computeEquationSolution(eqnMatrix *mat64.Dense, resultMatrix *mat64.Dense) {
	// Calculate inverse of A
	eqnMatrixRows, eqnMatrixCols := eqnMatrix.Caps()
	inv := mat64.NewDense(eqnMatrixRows, eqnMatrixCols, nil)

	err := inv.Inverse(eqnMatrix)
	if err != nil {
		log.Fatalln(err)
	}

	// Calculate mul of A inverse and B
	mul := mat64.NewDense(eqnMatrixRows, 1, nil)
	mul.Mul(inv, resultMatrix)

	log.ErrorDump("Solution Matrix Is: \n\n", mul)
}

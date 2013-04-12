package R

/*

#cgo LDFLAGS: -lm -lR
#cgo CFLAGS: -I /usr/share/R/include/

#include <stdlib.h>
#include <R.h>
#include <Rinternals.h>
#include <Rdefines.h>
#include <R_ext/Parse.h>
#include <Rembedded.h>

void SetNumericVectorElt(SEXP vec, int i, double val) {
    REAL(vec)[i] = val;
}

double NumericVectorElt(SEXP vec, int i) {
    return REAL(vec)[i];
}
*/
import "C"

//import (
//	"fmt"
//	"log"

func boundsCheck(i int, length int) {
	if i >= length || i < 0 {
		panic("Index out of bounds")
	}
}

type Expression interface {
	toSexp() C.SEXP
}

type expression struct {
	expr   C.SEXP
	length int
}

func (this *expression) toSexp() C.SEXP {
	return this.expr
}

func (this *expression) Len() int {
	return this.length
}

func (this *expression) boundsCheck(i int) {
	boundsCheck(i, this.length)
}

type NumericVector struct {
	expression
}

func NewNumericVector(vector []float64) *NumericVector {

	length := len(vector)
	v := NumericVector{}
	v.expr = C.allocVector(C.REALSXP, C.R_len_t(length))
	v.length = length

	v.CopyFrom(vector)

	return &v
}

func (this *NumericVector) Get(i int) float64 {
	this.boundsCheck(i)
	C.Rf_protect(this.expr)
	defer C.Rf_unprotect(1)
	return float64(C.NumericVectorElt(this.expr, C.int(i)))
}

func (this *NumericVector) Set(i int, val float64) {
	this.boundsCheck(i)
	C.Rf_protect(this.expr)
	defer C.Rf_unprotect(1)
	C.SetNumericVectorElt(this.expr, C.int(i), C.double(val))
}

func (this *NumericVector) CopyFrom(src []float64) {
	C.Rf_protect(this.expr)
	defer C.Rf_unprotect(1)
	for i := 0; i < this.length; i++ {
		C.SetNumericVectorElt(this.expr, C.int(i), C.double(src[i]))
	}
}

type Result struct {
	expr C.SEXP
}

func NewResult(expr C.SEXP) *Result {
	return &Result{expr: expr}
}

func (this *Result) IsNumeric() bool {
	return C.Rf_isReal(this.expr) != 0
}

func (this *Result) IsComplex() bool {
	return C.Rf_isComplex(this.expr) != 0
}

func (this *Result) AsComplex() *ComplexVector {
	if !this.IsComplex() {
		panic("Not a complex vector")
	}
	v := ComplexVector{}
	v.length = int(C.Rf_length(this.expr))
	v.expr = this.expr
	return &v

}

func (this *Result) AsNumeric() *NumericVector {
	if !this.IsNumeric() {
		panic("Not a numeric vector")
	}
	v := NumericVector{}
	v.length = int(C.Rf_length(this.expr))
	v.expr = this.expr
	return &v
}

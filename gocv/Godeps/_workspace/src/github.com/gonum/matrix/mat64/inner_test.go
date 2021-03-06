// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat64

import (
	"math"
	"math/rand"
	"testing"

	"github.com/gonum/blas/blas64"
	"gopkg.in/check.v1"
)

func (s *S) TestInner(c *check.C) {
	for i, test := range []struct {
		x []float64
		y []float64
		m [][]float64
	}{
		{
			x: []float64{5},
			y: []float64{10},
			m: [][]float64{{2}},
		},
		{
			x: []float64{5, 6, 1},
			y: []float64{10},
			m: [][]float64{{2}, {-3}, {5}},
		},
		{
			x: []float64{5},
			y: []float64{10, 15},
			m: [][]float64{{2, -3}},
		},
		{
			x: []float64{1, 5},
			y: []float64{10, 15},
			m: [][]float64{
				{2, -3},
				{4, -1},
			},
		},
		{
			x: []float64{2, 3, 9},
			y: []float64{8, 9},
			m: [][]float64{
				{2, 3},
				{4, 5},
				{6, 7},
			},
		},
		{
			x: []float64{2, 3},
			y: []float64{8, 9, 9},
			m: [][]float64{
				{2, 3, 6},
				{4, 5, 7},
			},
		},
	} {
		x := NewDense(1, len(test.x), test.x)
		m := NewDense(flatten(test.m))
		mWant := NewDense(flatten(test.m))
		y := NewDense(len(test.y), 1, test.y)

		mWant.Mul(mWant, y)
		mWant.Mul(x, mWant)

		rm, cm := mWant.Dims()
		c.Check(rm, check.Equals, 1, check.Commentf("Test %v result doesn't have 1 row", i))
		c.Check(cm, check.Equals, 1, check.Commentf("Test %v result doesn't have 1 column", i))

		want := mWant.At(0, 0)

		got := Inner(test.x, m, test.y)
		c.Check(want, check.Equals, got, check.Commentf("Test %v: want %v, got %v", i, want, got))
	}
}

func (s *S) TestInnerSym(c *check.C) {
	n := 10
	x := make([]float64, n)
	y := make([]float64, n)
	data := make([]float64, n*n)
	for i := 0; i < n; i++ {
		x[i] = float64(i)
		y[i] = float64(i)
		for j := i; j < n; j++ {
			data[i*n+j] = float64(i*n + j)
			data[j*n+i] = data[i*n+j]
		}
	}
	m := NewDense(n, n, data)
	ans := Inner(x, m, y)
	sym := NewSymDense(n, data)
	// scramble the lower half of data to ensure it is not used
	for i := 1; i < n; i++ {
		for j := 0; j < i; j++ {
			data[i*n+j] = rand.Float64()
		}
	}

	if math.Abs(Inner(x, sym, y)-ans) > 1e-14 {
		c.Error("inner different symmetric and dense")
	}
}

func benchmarkInner(b *testing.B, m, n int) {
	x := make([]float64, m)
	randomSlice(x)
	y := make([]float64, n)
	randomSlice(y)
	data := make([]float64, m*n)
	randomSlice(data)
	mat := &Dense{mat: blas64.General{Rows: m, Cols: n, Stride: n, Data: data}, capRows: m, capCols: n}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Inner(x, mat, y)
	}
}

func BenchmarkInnerSmSm(b *testing.B) {
	benchmarkInner(b, Sm, Sm)
}

func BenchmarkInnerMedMed(b *testing.B) {
	benchmarkInner(b, Med, Med)
}

func BenchmarkInnerLgLg(b *testing.B) {
	benchmarkInner(b, Lg, Lg)
}

func BenchmarkInnerLgSm(b *testing.B) {
	benchmarkInner(b, Lg, Sm)
}

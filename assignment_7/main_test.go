package main

import (
	"math"
	"testing"
)

func TestCosts(t *testing.T) {
	tests := []struct {
		name         string
		shape        Shape
		expectedCost float64
	}{
		{
			name:         "Square_Cost",
			shape:        square{l: 5},
			expectedCost: 25 * 100,
		},
		{
			name:         "Rectangle_Cost",
			shape:        rectangle{l: 4.222222, b: 1},
			expectedCost: 4.222222 * 20,
		},
		{
			name:         "Square_ZeroArea",
			shape:        square{l: 0},
			expectedCost: 0,
		}, {
			name:         "Rectangle_ZeroArea",
			shape:        rectangle{l: 0, b: 0},
			expectedCost: 0,
		},
		{
			name:         "Square_NegativeSide_Cost",
			shape:        square{l: -5},
			expectedCost: 25 * 100,
		},
		{
			name:         "Rectangle_NegativeDimensions_Cost",
			shape:        rectangle{l: -4, b: -5},
			expectedCost: 20 * 20,
		},
	}

	for _, tt := range tests {
		t.Logf("Running test case: %s \n", tt.name)

		gotCost := costs(tt.shape)
		tol := 1e-6
		currentToll := math.Abs(gotCost - tt.expectedCost)

		if currentToll > tol {
			t.Errorf("FAIL: Test %s For shape %T (Area: %.6f), expected cost %.6f, but got %.6f",
				tt.name, tt.shape, tt.shape.Area(), tt.expectedCost, gotCost)
		} else {
			t.Logf("PASS: Test %s", tt.name)
		}
	}
}

func TestSquareArea(t *testing.T) {
	s1 := square{l: 10}
	expectedArea1 := 100.0
	gotArea1 := s1.Area()
	tol := 1e-9
	if math.Abs(gotArea1-expectedArea1) > tol {
		t.Errorf("FAIL: Test PositiveSide - For side 10, expected area %.2f, but got %.2f", expectedArea1, gotArea1)
	} else {
		t.Logf("PASS: Test PositiveSide - Side 10, Area %.2f", gotArea1)
	}

	s2 := square{l: 0}
	expectedArea2 := 0.0
	gotArea2 := s2.Area()
	if math.Abs(gotArea2-expectedArea2) > tol {
		t.Errorf("FAIL: Test ZeroSide - For side 0, expected area %.2f, but got %.2f", expectedArea2, gotArea2)
	} else {
		t.Logf("PASS: Test ZeroSide - Side 0, Area %.2f", gotArea2)
	}

	s3 := square{l: 2.5}
	expectedArea3 := 6.25
	gotArea3 := s3.Area()
	if math.Abs(gotArea3-expectedArea3) > tol {
		t.Errorf("FAIL: Test FractionalSide - For side 2.5, expected area %.2f, but got %.2f", expectedArea3, gotArea3)
	} else {
		t.Logf("PASS: Test FractionalSide - Side 2.5, Area %.2f", gotArea3)
	}

	s4 := square{l: -5}
	expectedArea4 := 25.0
	gotArea4 := s4.Area()
	if math.Abs(gotArea4-expectedArea4) > tol {
		t.Errorf("FAIL: Test NegativeSide - For side -5, expected area %.2f, but got %.2f", expectedArea4, gotArea4)
	} else {
		t.Logf("PASS: Test NegativeSide - Side -5, Area %.2f", gotArea4)
	}
}

func TestRectangleArea(t *testing.T) {
	r1 := rectangle{b: 10, l: 20}
	expectedArea1 := 200.0
	gotArea1 := r1.Area()
	tol := 1e-9
	if math.Abs(gotArea1-expectedArea1) > tol {
		t.Errorf("FAIL: Test PositiveDimensions - For breadth 10, length 20, expected area %.2f, but got %.2f", expectedArea1, gotArea1)
	} else {
		t.Logf("PASS: Test PositiveDimensions - Breadth 10, Length 20, Area %.2f", gotArea1)
	}

	r2 := rectangle{b: 0, l: 20}
	expectedArea2 := 0.0
	gotArea2 := r2.Area()
	if math.Abs(gotArea2-expectedArea2) > tol {
		t.Errorf("FAIL: Test ZeroBreadth - For breadth 0, length 20, expected area %.2f, but got %.2f", expectedArea2, gotArea2)
	} else {
		t.Logf("PASS: Test ZeroBreadth - Breadth 0, Length 20, Area %.2f", gotArea2)
	}

	r3 := rectangle{b: 10, l: 0}
	expectedArea3 := 0.0
	gotArea3 := r3.Area()
	if math.Abs(gotArea3-expectedArea3) > tol {
		t.Errorf("FAIL: Test ZeroLength - For breadth 10, length 0, expected area %.2f, but got %.2f", expectedArea3, gotArea3)
	} else {
		t.Logf("PASS: Test ZeroLength - Breadth 10, Length 0, Area %.2f", gotArea3)
	}

	r4 := rectangle{b: 0, l: 0}
	expectedArea4 := 0.0
	gotArea4 := r4.Area()
	if math.Abs(gotArea4-expectedArea4) > tol {
		t.Errorf("FAIL: Test BothZero - For breadth 0, length 0, expected area %.2f, but got %.2f", expectedArea4, gotArea4)
	} else {
		t.Logf("PASS: Test BothZero - Breadth 0, Length 0, Area %.2f", gotArea4)
	}

	r5 := rectangle{b: 2.5, l: 4.0}
	expectedArea5 := 10.0
	gotArea5 := r5.Area()
	if math.Abs(gotArea5-expectedArea5) > tol {
		t.Errorf("FAIL: Test FractionalDimensions - For breadth 2.5, length 4.0, expected area %.2f, but got %.2f", expectedArea5, gotArea5)
	} else {
		t.Logf("PASS: Test FractionalDimensions - Breadth 2.5, Length 4.0, Area %.2f", gotArea5)
	}

	r6 := rectangle{b: 5, l: -4}
	expectedArea6 := -20.0
	gotArea6 := r6.Area()
	if math.Abs(gotArea6-expectedArea6) > tol {
		t.Errorf("FAIL: Test NegativeLength - For breadth 5, length -4, expected area %.2f, but got %.2f", expectedArea6, gotArea6)
	} else {
		t.Logf("PASS: Test NegativeLength - Breadth 5, Length -4, Area %.2f", gotArea6)
	}

	r7 := rectangle{b: -5, l: 4}
	expectedArea7 := -20.0
	gotArea7 := r7.Area()
	if math.Abs(gotArea7-expectedArea7) > tol {
		t.Errorf("FAIL: Test NegativeBreadth - For breadth -5, length 4, expected area %.2f, but got %.2f", expectedArea7, gotArea7)
	} else {
		t.Logf("PASS: Test NegativeBreadth - Breadth -5, Length 4, Area %.2f", gotArea7)
	}

	r8 := rectangle{b: -5, l: -4}
	expectedArea8 := 20.0
	gotArea8 := r8.Area()
	if math.Abs(gotArea8-expectedArea8) > tol {
		t.Errorf("FAIL: Test BothNegative - For breadth -5, length -4, expected area %.2f, but got %.2f", expectedArea8, gotArea8)
	} else {
		t.Logf("PASS: Test BothNegative - Breadth -5, Length -4, Area %.2f", gotArea8)
	}
}

func TestCalArea(t *testing.T) {
	tol := 1e-9

	s1 := square{l: 7}
	expectedArea1 := 49.0
	gotArea1 := calArea(s1)
	if math.Abs(gotArea1-expectedArea1) > tol {
		t.Errorf("FAIL: Test CalculateSquareArea - For shape %+v, expected area %.2f, but got %.2f", s1, expectedArea1, gotArea1)
	} else {
		t.Logf("PASS: Test CalculateSquareArea - Shape %+v, Area %.2f", s1, gotArea1)
	}

	r2 := rectangle{b: 5, l: 8}
	expectedArea2 := 40.0
	gotArea2 := calArea(r2)
	if math.Abs(gotArea2-expectedArea2) > tol {
		t.Errorf("FAIL: Test CalculateRectangleArea - For shape %+v, expected area %.2f, but got %.2f", r2, expectedArea2, gotArea2)
	} else {
		t.Logf("PASS: Test CalculateRectangleArea - Shape %+v, Area %.2f", r2, gotArea2)
	}

	s3 := square{l: 0}
	expectedArea3 := 0.0
	gotArea3 := calArea(s3)
	if math.Abs(gotArea3-expectedArea3) > tol {
		t.Errorf("FAIL: Test CalculateSquareZeroArea - For shape %+v, expected area %.2f, but got %.2f", s3, expectedArea3, gotArea3)
	} else {
		t.Logf("PASS: Test CalculateSquareZeroArea - Shape %+v, Area %.2f", s3, gotArea3)
	}

	r4 := rectangle{b: 0, l: 10}
	expectedArea4 := 0.0
	gotArea4 := calArea(r4)
	if math.Abs(gotArea4-expectedArea4) > tol {
		t.Errorf("FAIL: Test CalculateRectangleZeroArea - For shape %+v, expected area %.2f, but got %.2f", r4, expectedArea4, gotArea4)
	} else {
		t.Logf("PASS: Test CalculateRectangleZeroArea - Shape %+v, Area %.2f", r4, gotArea4)
	}

	s5 := square{l: -3}
	expectedArea5 := 9.0
	gotArea5 := calArea(s5)
	if math.Abs(gotArea5-expectedArea5) > tol {
		t.Errorf("FAIL: Test CalculateSquareNegativeSide - For shape %+v, expected area %.2f, but got %.2f", s5, expectedArea5, gotArea5)
	} else {
		t.Logf("PASS: Test CalculateSquareNegativeSide - Shape %+v, Area %.2f", s5, gotArea5)
	}

	r6 := rectangle{b: -2, l: -5}
	expectedArea6 := 10.0
	gotArea6 := calArea(r6)
	if math.Abs(gotArea6-expectedArea6) > tol {
		t.Errorf("FAIL: Test CalculateRectangleNegativeDimensions - For shape %+v, expected area %.2f, but got %.2f", r6, expectedArea6, gotArea6)
	} else {
		t.Logf("PASS: Test CalculateRectangleNegativeDimensions - Shape %+v, Area %.2f", r6, gotArea6)
	}
}

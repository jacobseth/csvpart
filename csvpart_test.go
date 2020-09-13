package main

import "testing"

// linesFromPercs - LFP
func TestLFPHeaderCountToLarge(t *testing.T) {
	if _, err := linesFromPercs(1, []float32{100}, 2, false); err == nil {
		t.Errorf("Expected line sum too large error")
	}
}

func TestLFPCorrectLCForPerc(t *testing.T) {
	// 100 seems too neat - Want those funky rounding errors
	lineCount := 123
	whole := false
	testParams := []struct {
		percs    []float32
		expected []int
	}{
		{
			[]float32{50.0},
			[]int{62},
		},
		{
			[]float32{33.33, 33.33},
			[]int{41, 41},
		},
		{
			[]float32{11, 12, 13.323},
			[]int{14, 15, 16},
		},
	}

	for _, tp := range testParams {
		// If removed, update len checks for RV
		if len(tp.expected) != len(tp.percs) {
			t.Errorf("Testing data sizes don't match; Percentages: %d, Expected: %d", len(tp.percs), len(tp.expected))
		}
		lcps, err := linesFromPercs(lineCount, tp.percs, 0, whole)
		if err != nil {
			t.Error(err)
		}
		if len(lcps) != len(tp.percs) {
			t.Errorf("RV too large; Expected %d, got %d", len(tp.percs), len(lcps))
		}

		for i := 0; i < len(lcps); i++ {
			if lcps[i] != tp.expected[i] {
				t.Errorf("Wrong RV: expected %d, got %d", tp.expected[i], lcps[i])
			}
		}
	}
}

// End -- LFP

// linesFromPercs - LFP with --whole
func TestLFPCorrectLCForPercWithWhole(t *testing.T) {
	// 100 seems too neat - Want those funky rounding errors
	lineCount := 123
	whole := true
	testParams := []struct {
		percs    []float32
		expected []int
	}{
		{
			[]float32{50.0},
			[]int{62, 61},
		},
		{
			[]float32{33.33, 33.33},
			[]int{41, 41, 41},
		},
		{
			[]float32{11, 12, 13.323},
			[]int{14, 15, 16, 78},
		},
	}

	for _, tp := range testParams {
		lcps, err := linesFromPercs(lineCount, tp.percs, 0, whole)
		if err != nil {
			t.Error(err)
		}

		for i := 0; i < len(lcps); i++ {
			if lcps[i] != tp.expected[i] {
				t.Errorf("Wrong RV: expected %d, got %d", tp.expected[i], lcps[i])
			}
		}
	}
}

// End LFP with --whole

// Header test
func TestLFPCorrectLCForPercWithHeader(t *testing.T) {
	// 100 seems too neat - Want those funky rounding errors
	lineCount := 123
	whole := false
	header := 12
	testParams := []struct {
		percs    []float32
		expected []int
	}{
		{
			[]float32{50.0},
			[]int{62},
		},
		{
			[]float32{33.33, 33.33},
			[]int{41, 41},
		},
		{
			[]float32{11, 12, 13.323},
			[]int{14, 15, 16},
		},
	}

	for _, tp := range testParams {
		// If removed, update len checks for RV
		if len(tp.expected) != len(tp.percs) {
			t.Errorf("Testing data sizes don't match; Percentages: %d, Expected: %d", len(tp.percs), len(tp.expected))
		}
		lcps, err := linesFromPercs(lineCount+header, tp.percs, header, whole)
		if err != nil {
			t.Error(err)
		}
		if len(lcps) != len(tp.percs) {
			t.Errorf("RV too large; Expected %d, got %d", len(tp.percs), len(lcps))
		}

		for i := 0; i < len(lcps); i++ {
			if lcps[i] != tp.expected[i] {
				t.Errorf("Wrong RV: expected %d, got %d", tp.expected[i], lcps[i])
			}
		}
	}
}

func TestLFPCorrectLCForPercWithWholeAndHeader(t *testing.T) {
	// 100 seems too neat - Want those funky rounding errors
	lineCount := 123
	whole := true
	header := 12
	testParams := []struct {
		percs    []float32
		expected []int
	}{
		{
			[]float32{50.0},
			[]int{62, 61},
		},
		{
			[]float32{33.33, 33.33},
			[]int{41, 41, 41},
		},
		{
			[]float32{11, 12, 13.323},
			[]int{14, 15, 16, 78},
		},
	}

	for _, tp := range testParams {
		lcps, err := linesFromPercs(lineCount+header, tp.percs, header, whole)
		if err != nil {
			t.Error(err)
		}

		for i := 0; i < len(lcps); i++ {
			if lcps[i] != tp.expected[i] {
				t.Errorf("Wrong RV: expected %d, got %d", tp.expected[i], lcps[i])
			}
		}
	}
}

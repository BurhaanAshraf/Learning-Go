// Package main demonstrates Go's built-in testing and benchmarking toolkit:
// unit tests, subtests, table-driven tests, and memory/CPU benchmarking.
//
// Run tests:        go test -run=^Test -v
// Run benchmarks:   go test -bench=. -run=^$ -benchmem
package main

import (
	"fmt"
	"math/rand"
	"testing"
)

// ---------------------------------------------------------------------------
// CORE FUNCTIONS — the code under test. Kept deliberately simple so focus
// stays on testing/benchmarking mechanics rather than business logic.
// ---------------------------------------------------------------------------

// add returns the sum of two ints. Trivial on purpose — used to demonstrate
// benchmark structure without any real allocation or computation noise.
//
// NOTE (inlining): a function this small is almost guaranteed to be inlined
// by the compiler in a normal build. Inlining merges add's work directly
// into whichever function calls it, so on a CPU profile add may never show
// up as its own frame — its cost gets folded into BenchmarkAddSmallInput
// etc. This is a different gotcha from the escape-analysis one on
// GenerateRandomSlice below: inlining hides *frames* in CPU profiles,
// escape analysis hides *allocations* in memory profiles. See the
// PROFILING GOTCHA note in the benchmarks section for how to work around it.
func add(i, j int) int {
	return i + j
}

// GenerateRandomSlice returns a slice of `size` pseudo-random ints in [0,100).
// NOTE (allocation behavior): make([]int, size) is the only heap allocation
// here under normal use — but if a caller discards the return value (like a
// naive benchmark would), Go's escape analysis can prove the slice never
// escapes and keep it on the stack instead. That means a careless benchmark
// can under-report the real allocation cost. See the "sink" pattern in the
// benchmarks below for the fix.
func GenerateRandomSlice(size int) []int {
	slice := make([]int, size)

	for i := range slice {
		// rand.Intn here uses the package-level global Source, which is
		// mutex-guarded. Fine for a normal benchmark; becomes a contention
		// point if this were ever called from b.RunParallel goroutines.
		slice[i] = rand.Intn(100)
	}
	return slice
}

// SumSlice adds up all elements in a slice. No allocations — pure read +
// arithmetic — so its benchmark should show ~0 B/op and 0 allocs/op.
func SumSlice(slice []int) int {
	sum := 0
	for _, v := range slice {
		sum += v
	}
	return sum
}

// ---------------------------------------------------------------------------
// UNIT TESTS — go test auto-discovers any func TestXxx(t *testing.T).
// Naming rule: must start with "Test", take *testing.T, return nothing.
// ---------------------------------------------------------------------------

// TestGenerateRandomSlice checks one invariant: requested size == actual size.
// NOTE: this does NOT verify the values are in range [0,100) — that would be
// a separate assertion, or a good candidate for a fuzz test (go test -fuzz).
// One test = one clear claim; resist bundling unrelated checks into it.
func TestGenerateRandomSlice(t *testing.T) {
	size := 100
	slice := GenerateRandomSlice(size) 

	if len(slice) != size {
		t.Errorf("Size of Slice = %d | Size Given = %d: SIZE DOES NOT MATCH", len(slice), size)
	}
}

// TestAdd is the simplest possible test shape: call, compare, fail loudly.
// Good for a one-off check; gets unwieldy fast once you have more than a
// couple of cases — that's what the subtest and table-driven patterns below
// are for.
func TestAdd(t *testing.T) {
	result := add(2, 3)
	expectedVal := 5
	if result != expectedVal {
		t.Errorf("Add(2,3) = %d. , Wanted = %d", result, expectedVal)
	}
}

// ---------------------------------------------------------------------------
// SUBTESTS — t.Run groups related cases under one parent test. Each subtest
// gets its own pass/fail line under `go test -v`, and you can target just
// one case directly: go test -run='TestAddSubtests/Add\(2,3\)'
// This is the building block table-driven tests are layered on top of.
// ---------------------------------------------------------------------------
func TestAddSubtests(t *testing.T) {
	tests := []struct{ a, b, expected int }{
		{2, 3, 5}, // fixed from original {2,2,5} — 2+2=4, the expected value didn't match the inputs
		{0, 0, 0},
		{-1, 1, 0},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Add(%d,%d)", test.a, test.b), func(t *testing.T) {
			res := add(test.a, test.b)
			if res != test.expected {
				t.Errorf("Result %d! Wanted = %d", res, test.expected)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// TABLE-DRIVEN TESTS — the idiomatic Go pattern for testing many input/output
// pairs against one function. Adding a new case is one line, no new test
// func needed. Differs from TestAddSubtests above only in that cases aren't
// wrapped in t.Run, so all cases share one pass/fail result instead of
// reporting independently — fine for small tables, but t.Run scales better
// once a table grows or you need to isolate which exact case failed.
// ---------------------------------------------------------------------------
func TestAddTableDriven(t *testing.T) {
	tests := []struct{ a, b, expected int }{
		{2, 3, 5},
		{0, 1, 1}, // fixed from original {0,1,0} — 0+1=1, not 0
		{-1, 1, 0},
	}

	for _, test := range tests {
		res := add(test.a, test.b)
		if res != test.expected {
			t.Errorf("Add(%d , %d) = %d but wanted %d", test.a, test.b, res, test.expected)
		}
	}
}

// ---------------------------------------------------------------------------
// BENCHMARKS — go test discovers func BenchmarkXxx(b *testing.B).
// `for range b.N` lets the framework auto-tune iteration count until timing
// stabilizes (Go 1.22+ range-over-int; equivalent to the classic
// `for i := 0; i < b.N; i++ {}` on older Go versions).
//
//   go test -bench=. -run=^$              run benchmarks only (skip Test funcs)
//   go test -bench=. -run=^$ -benchmem    also report B/op and allocs/op
//
// -run=^$ matches no test names, so TestXxx funcs are skipped entirely —
// keeps test setup/output from polluting benchmark timing.
//
// PROFILING GOTCHA — INLINING:
// Compiler inlining can hide small functions in pprof output by merging
// their work into the caller's frame (see the note on add() above). To
// preserve function boundaries while profiling:
//
//	go test -gcflags='-l' -bench=. -run=^$ -cpuprofile=cpu.pprof
//
// -l disables inlining for that build, so add() shows up as its own line in
// `go tool pprof -list=add cpu.pprof` instead of being folded into its
// caller. To see *which* functions the compiler would normally inline (and
// why), without changing what gets built:
//
//	go test -gcflags='-m' -bench=. -run=^$
//
// Use -gcflags only for debugging/learning profiles. Real benchmarks should
// keep optimizations enabled — disabling inlining changes the actual
// performance characteristics you're trying to measure, so numbers from a
// -l build don't represent production behavior.
// ---------------------------------------------------------------------------

// sinkSlice / sinkInt exist purely to receive benchmark results.
// WHY THIS MATTERS: if a benchmark calls a function and discards the return
// value, the compiler can sometimes prove the result is never observed and
// either skip the heap allocation (via escape analysis) or eliminate work
// entirely — silently making the benchmark measure less than you think.
// Assigning the result to a package-level var defeats that: the compiler
// can't know the var won't be read elsewhere, so it must keep the real work.
var sinkSlice []int
var sinkInt int

func BenchmarkAddSmallInput(b *testing.B) {
	var r int
	for range b.N {
		r = add(2, 3)
	}
	sinkInt = r
}

func BenchmarkAddMediumInput(b *testing.B) {
	var r int
	for range b.N {
		r = add(200, 300)
	}
	sinkInt = r
}

func BenchmarkAddLargeInput(b *testing.B) {
	var r int
	for range b.N {
		r = add(2000, 3000)
	}
	sinkInt = r
}

// BenchmarkGenerateRandomSlice measures the cost of building a 1000-element
// slice per call. Run with -benchmem to see B/op and allocs/op — expect
// ~1 alloc/op (the backing array) now that the result is sunk above instead
// of discarded.
func BenchmarkGenerateRandomSlice(b *testing.B) {
	var r []int
	for range b.N {
		r = GenerateRandomSlice(1000)
	}
	sinkSlice = r
}

// BenchmarkSumSlice isolates *summation* cost from *generation* cost.
// The slice is built once outside the loop as setup, then b.ResetTimer()
// zeroes both the clock and the allocation counters so that one-time setup
// cost isn't charged to the thing actually being measured. Skipping
// ResetTimer here would skew ns/op and B/op, especially at low b.N where
// setup is a larger fraction of total time.
func BenchmarkSumSlice(b *testing.B) {
	slice := GenerateRandomSlice(1000)
	b.ResetTimer()

	var s int
	for range b.N {
		s = SumSlice(slice)
	}
	sinkInt = s
}

// ---------------------------------------------------------------------------
// QUICK REFERENCE (revision notes)
// ---------------------------------------------------------------------------
// go test                            run tests in current package
// go test -v                         verbose: print each test/subtest result
// go test -run=TestAdd               run only tests matching this regex
// go test -bench=.                   run all benchmarks (regex match on name)
// go test -bench=. -run=^$           run benchmarks only, skip all tests
// go test -bench=. -benchmem         add B/op and allocs/op to output
// go test -bench=. -benchtime=3s     run each benchmark ~3s instead of default ~1s
// go test -cpuprofile=cpu.pprof      capture a CPU profile during the run
// go test -memprofile=mem.pprof      capture a heap allocation profile
// go tool pprof -top -cum FILE       inspect a profile, sorted by cumulative cost
// go tool pprof -list=FuncName FILE  line-by-line cost inside one function
// go test -gcflags='-l' ...          disable inlining (debugging profiles only — keeps small funcs as their own frames)
// go test -gcflags='-m' ...          print the compiler's inlining decisions without changing the build

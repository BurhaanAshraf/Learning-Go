package main

import "fmt"

func verbs() {

	// ============================================================
	// 🔹 String Formatting Verbs
	// ============================================================

	// %s   → plain string
	// %q   → double-quoted string
	// %8s  → width 8, right aligned
	// %-8s → width 8, left aligned
	// %x   → hex (byte values)
	// % x  → hex with spaces

	txt := "World"

	fmt.Printf("%s\n", txt)
	fmt.Printf("%q\n", txt)
	fmt.Printf("%8s\n", txt)
	fmt.Printf("%-8s\n", txt)
	fmt.Printf("%x\n", txt)
	fmt.Printf("% x\n", txt)


	// ============================================================
	// 🔹 Boolean Formatting Verbs
	// ============================================================

	// %t → true / false

	t := true
	f := false

	fmt.Printf("%t\n", t)
	fmt.Printf("%t\n", f)


	// ============================================================
	// 🔹 Integer Formatting Verbs
	// ============================================================

	// %b   → binary
	// %d   → decimal
	// %+d  → decimal with sign
	// %o   → octal
	// %O   → octal with 0o
	// %x   → hex (lowercase)
	// %X   → hex (uppercase)
	// %#x  → hex with 0x
	// %4d  → width 4 (right)
	// %-4d → width 4 (left)
	// %04d → pad with zero

	num := 18

	fmt.Printf("%b\n", num)
	fmt.Printf("%d\n", num)
	fmt.Printf("%+d\n", num)
	fmt.Printf("%o\n", num)
	fmt.Printf("%O\n", num)
	fmt.Printf("%x\n", num)
	fmt.Printf("%X\n", num)
	fmt.Printf("%#x\n", num)
	fmt.Printf("%4d\n", num)
	fmt.Printf("%-4d\n", num)
	fmt.Printf("%04d\n", num)


	// ============================================================
	// 🔹 Float Formatting Verbs
	// ============================================================

	// %e   → scientific (e)
	// %f   → decimal
	// %.2f → precision 2
	// %6.2f → width 6, precision 2
	// %g   → compact (auto format)

	flt := 9.18

	fmt.Printf("%e\n", flt)
	fmt.Printf("%f\n", flt)
	fmt.Printf("%.2f\n", flt)
	fmt.Printf("%6.2f\n", flt)
	fmt.Printf("%g\n", flt)
}
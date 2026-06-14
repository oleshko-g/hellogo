package main

import "fmt"

const PGSIZE uint64 = 4096 // bytes per page
const PGSHIFT uint64 = 12  // bits of offset within a page
// const PGROUNDUP = (sz) (((sz) + PGSIZE - 1) & ~(PGSIZE - 1))
// const PGROUNDDOWN(a) (((a)) & ~(PGSIZE - 1))
const PTE_V = 1 << 0 // valid  1L is 0b0000 0000 0000 0000 0000 0000 0000 0000 0000 0001
const PTE_R = 1 << 1
const PTE_W = 1 << 2
const PTE_X = 1 << 3
const PTE_U = 1 << 4        // user can access
const PTE_FLAGS_OFFSET = 10 // bits of offset within a PTE

// shift a physical address to the right place for a PTE.
// Example: The address of the physical frame is
// in hex:
// 0x0000 0000 0000 0123
// or binary:
// 0b0000 0000 0000 0000 0000 0000 0000 0000 0000 0001 0010 0011
//
// Shift it left to free up the space for the PTE flag bits.
// 0b00000000 00000000 00000000 0000 0100 1000 1100 0000 0000
// 0x0000 0000 0004 8C00
// #define PA2PTE(pa) ((((uint64)pa) >> PGSHIFT) << PTE_FLAGS_OFFSET)

// #define PTE2PA(pte) (((pte) >> PTE_FLAGS_OFFSET) << PGSHIFT)

// #define PTE_FLAGS(pte) ((pte) & 0x3FF)

// extract the three 9-bit page table indices from a virtual address.
const PXMASK = 0x1FF // 9 bits

// levels are expected to be:
//  * 2 - "3-rd"
//  * 1 - "2-nd"
//  * 0 - "1-st"
// #define PXSHIFT(level) (PGSHIFT + (9 * (level)))

// PTE index
// #define PX(level, va) ((((uint64)(va)) >> PXSHIFT(level)) & PXMASK)

// one beyond the highest possible virtual address.
// MAXVA is actually one bit less than the max allowed by
// Sv39, to avoid having to sign-extend virtual addresses
// that have the high bit set.
// #define MAXVA (1L << (9 + 9 + 9 + 12 - 1))

func main() {
	va := uint64(0b00000000_00000000_00000000_000000001_000000010_000000101_0001_0010_0101)
	fmt.Printf("%b\n", va)
	fmt.Printf("%b\n", va&(PGSIZE-1))
	fmt.Printf("%b\n", (va & ^(PGSIZE - 1)))
	fmt.Printf("%b\n", ((va & (PGSIZE - 1)) | (va & ^(PGSIZE - 1))))
}

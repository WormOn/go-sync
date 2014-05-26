package index

import (
	"github.com/Redundancy/go-sync/chunks"
	"testing"
)

// Weak checksums must be 4 bytes
var WEAK_A = []byte("aaaa")
var WEAK_B = []byte("bbbb")

func TestMakeIndex(t *testing.T) {
	i := MakeChecksumIndex(
		[]chunks.ChunkChecksum{
			{0, WEAK_A, []byte("b")},
			{1, WEAK_B, []byte("c")},
		},
	)

	if len(i.weakChecksumLookup) != 2 {
		t.Errorf("size of lookup was not expected %v", len(i.weakChecksumLookup))
	}
}

func TestFindWeakInIndex(t *testing.T) {
	i := MakeChecksumIndex(
		[]chunks.ChunkChecksum{
			{0, WEAK_A, []byte("b")},
			{1, WEAK_B, []byte("c")},
			{2, WEAK_B, []byte("d")},
		},
	)

	result := i.FindWeakChecksumInIndex(WEAK_B)

	if result == nil {
		t.Error("Did not find lookfor in the index")
	} else if len(result) != 2 {
		t.Errorf("Wrong number of possible matches found: %v", len(result))
	} else if result[0].ChunkOffset != 1 {
		t.Errorf("Found chunk had offset %v expected 1", result[0].ChunkOffset)
	}
}

func TestFindStrongInIndex(t *testing.T) {
	i := MakeChecksumIndex(
		[]chunks.ChunkChecksum{
			{0, WEAK_A, []byte("b")},
			{1, WEAK_B, []byte("c")},
			{2, WEAK_B, []byte("d")},
		},
	)

	// builds upon TestFindWeakInIndex
	result := i.FindWeakChecksumInIndex(WEAK_B)
	strongs := result.FindStrongChecksum([]byte("c"))

	if len(strongs) != 1 {
		t.Errorf("Incorrect number of strong checksums found: %v", len(strongs))
	} else if strongs[0].ChunkOffset != 1 {
		t.Errorf("Wrong chunk found, had offset %v", strongs[0].ChunkOffset)
	}
}

func TestNotFoundStrongInIndexAtEnd(t *testing.T) {
	i := MakeChecksumIndex(
		[]chunks.ChunkChecksum{
			{0, WEAK_A, []byte("b")},
			{1, WEAK_B, []byte("c")},
			{2, WEAK_B, []byte("d")},
		},
	)

	// builds upon TestFindWeakInIndex
	result := i.FindWeakChecksumInIndex(WEAK_B)
	strongs := result.FindStrongChecksum([]byte("e"))

	if len(strongs) != 0 {
		t.Errorf("Incorrect number of strong checksums found: %v", strongs)
	}
}

func TestNotFoundStrongInIndexInCenter(t *testing.T) {
	// The strong checksum we're looking for is not found
	// but is < another checksum in the strong list

	i := MakeChecksumIndex(
		[]chunks.ChunkChecksum{
			{0, WEAK_A, []byte("b")},
			{1, WEAK_B, []byte("c")},
			{2, WEAK_B, []byte("d")},
			{3, WEAK_B, []byte("f")},
		},
	)

	// builds upon TestFindWeakInIndex
	result := i.FindWeakChecksumInIndex(WEAK_B)
	strongs := result.FindStrongChecksum([]byte("e"))

	if len(strongs) != 0 {
		t.Errorf("Incorrect number of strong checksums found: %v", strongs)
	}
}

func TestFindDuplicatedBlocksInIndex(t *testing.T) {
	i := MakeChecksumIndex(
		[]chunks.ChunkChecksum{
			{0, WEAK_A, []byte("b")},
			{1, WEAK_B, []byte("c")},
			{3, WEAK_B, []byte("c")},
			{2, WEAK_B, []byte("d")},
		},
	)

	// builds upon TestFindWeakInIndex
	result := i.FindWeakChecksumInIndex(WEAK_B)
	strongs := result.FindStrongChecksum([]byte("c"))

	if len(strongs) != 2 {
		t.Fatalf("Incorrect number of strong checksums found: %v", strongs)
	}

	first := strongs[0]
	if first.ChunkOffset != 1 {
		t.Errorf("Wrong chunk found, had offset %v", first.ChunkOffset)
	}

	second := strongs[1]
	if second.ChunkOffset != 3 {
		t.Errorf("Wrong chunk found, had offset %v", second.ChunkOffset)
	}
}

func TestFindTwoDuplicatedBlocksInIndex(t *testing.T) {
	i := MakeChecksumIndex(
		[]chunks.ChunkChecksum{
			{1, WEAK_B, []byte("c")},
			{2, WEAK_B, []byte("c")},
		},
	)

	// builds upon TestFindWeakInIndex
	result := i.FindWeakChecksumInIndex(WEAK_B)
	strongs := result.FindStrongChecksum([]byte("c"))

	if len(strongs) != 2 {
		t.Fatalf("Incorrect number of strong checksums found: %v", strongs)
	}

	first := strongs[0]
	if first.ChunkOffset != 1 {
		t.Errorf("Wrong chunk found, had offset %v", first.ChunkOffset)
	}

	second := strongs[1]
	if second.ChunkOffset != 2 {
		t.Errorf("Wrong chunk found, had offset %v", second.ChunkOffset)
	}
}

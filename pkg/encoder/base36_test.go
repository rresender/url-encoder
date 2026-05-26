package encoder

import "testing"

func TestBase36Encode_Zero(t *testing.T) {
	got := Base36Encode(0)
	if got != "a" {
		t.Fatalf("Base36Encode(0) = %q, want %q", got, "a")
	}
}

func TestDynamicLengthEncode_PadsToMinimum(t *testing.T) {
	// Base36Encode(1) should be "b" with alphabet[0]='a', alphabet[1]='b'.
	got := DynamicLengthEncode(1, 4)
	want := "aaab"
	if got != want {
		t.Fatalf("DynamicLengthEncode(1, 4) = %q, want %q", got, want)
	}
}

func TestDynamicLengthEncode_TruncatesToMinimumLength(t *testing.T) {
	// 36^4 is large enough to produce a Base36 string longer than 4 characters.
	const minLength = 4
	const num = uint64(36 * 36 * 36 * 36)

	got := DynamicLengthEncode(num, minLength)
	if len(got) != minLength {
		t.Fatalf("DynamicLengthEncode(%d, %d) length = %d, want %d; got=%q", num, minLength, len(got), minLength, got)
	}
}

func TestDynamicLengthEncode_NeverShorterThanMinimum(t *testing.T) {
	minLength := 4
	for _, num := range []uint64{0, 1, 10, 123, 999, 10000} {
		got := DynamicLengthEncode(num, minLength)
		if len(got) < minLength {
			t.Fatalf("DynamicLengthEncode(%d, %d) produced %q (len=%d)", num, minLength, got, len(got))
		}
	}
}


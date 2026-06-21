package encoder

const (
	alphabet = "abcdefghijklmnopqrstuvwxyz0123456789"
	length   = uint64(len(alphabet))
)

func Base36Encode(value uint64) string {
	if value == 0 {
		return string(alphabet[0])
	}

	var encoded []byte
	toEncode := value

	for toEncode > 0 {
		encoded = append(encoded, alphabet[toEncode%length])
		toEncode = toEncode / length
	}

	// Reverse the encoded string
	for i, j := 0, len(encoded)-1; i < j; i, j = i+1, j-1 {
		encoded[i], encoded[j] = encoded[j], encoded[i]
	}

	return string(encoded)
}

// DynamicLengthEncode encodes a number into a base36 string of exactly minLength characters.
// It masks the input into the range [0, 36^minLength) so the output is always exactly
// minLength characters (left-padded with 'a'), making the mapping injective within that range.
func DynamicLengthEncode(num uint64, minLength int) string {
	maxVal := uint64(1)
	for i := 0; i < minLength; i++ {
		maxVal *= length
	}
	masked := num % maxVal

	encoded := Base36Encode(masked)
	for len(encoded) < minLength {
		encoded = string(alphabet[0]) + encoded
	}
	return encoded
}

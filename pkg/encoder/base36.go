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

// DynamicLengthEncode encodes a number into a base36 string with a minimum length.
func DynamicLengthEncode(num uint64, minLength int) string {
	encoded := Base36Encode(num)

	if len(encoded) <= minLength {
		return encoded
	}

	// Calculate base size for each part (prefix, middle, suffix)
	partSize := minLength / 3
	remaining := minLength % 3

	// Distribute any remaining characters
	prefixSize := partSize
	middleSize := partSize
	suffixSize := partSize

	switch remaining {
	case 1:
		middleSize++ // Give extra character to middle
	case 2:
		prefixSize++
		suffixSize++
	}

	// Get prefix (first N characters)
	prefix := encoded[:prefixSize]

	// Get suffix (last N characters)
	suffix := encoded[len(encoded)-suffixSize:]

	// Get middle (center characters)
	middleStart := (len(encoded) - middleSize) / 2
	middle := encoded[middleStart : middleStart+middleSize]

	return prefix + middle + suffix
}

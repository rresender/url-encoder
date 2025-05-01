package encoder

const (
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length   = int64(len(alphabet))
)

func Base62Encode(value int64) string {
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

package internal

import (
	"crypto/md5"
	"fmt"
)

// Hash is a simple wrapper around the MD5 algorithm implementation in the
// Go standard library. It takes in data and a salt and returns the hashed
// representation.
func Hash(data string, salt string) string {
	output := md5.Sum([]byte(data + salt))
	return fmt.Sprintf("%x", output)
}

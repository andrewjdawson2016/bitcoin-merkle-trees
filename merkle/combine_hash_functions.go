package merkle

import (
	"crypto/sha256"
	"fmt"
)

func Sha256CombineHashFn(left string, right string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(left+right)))
}

func CommaAppendCombineFn(left string, right string) string {
	return fmt.Sprintf("%v, %v", left, right)
}

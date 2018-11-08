package rand

import (
	"testing"
	"fmt"
)

func TestRand(t *testing.T) {
	fmt.Println(Rand("AB", 13))
	fmt.Println(len(Rand("AB", 13)))
}

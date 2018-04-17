package rand

import (
	"math/rand"
	"time"
	"strconv"
	"strings"
)

func Rand(prefix string, totalLen int) string {
	str := strconv.Itoa(int(time.Now().Nanosecond()))
	bytes := []byte(str)
	data := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < totalLen-2; i++ {
		data = append(data, bytes[r.Intn(len(bytes))])
	}
	result, _ := strconv.Atoi(string(data))
	result = result + rand.Intn(99)
	if len(strconv.Itoa(result)) < totalLen-len(prefix) {
		return prefix + strings.Repeat("0", totalLen-len(prefix)-len(strconv.Itoa(result))) + strconv.Itoa(result)
	}
	return prefix + strconv.Itoa(result)
}

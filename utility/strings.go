package utility

import (
	"fmt"
	"strings"
	"time"
	"math/rand"
)

func RandomString(prefix string) string {
	t := time.Now()
	t.Round(time.Duration(rand.Int63n(8)))
	m := t.UnixNano() / int64(time.Millisecond)
	prefix = strings.ToUpper(prefix)
	return fmt.Sprintf(prefix+"%v", m)
}

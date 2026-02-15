package yasnowid

import (
	"testing"
)

func BenchmarkGenerator(b *testing.B) {
	g, _ := NewGenerator(13)
	for i := 0; i < b.N; i++ {
		_ = g.ID()
	}
}

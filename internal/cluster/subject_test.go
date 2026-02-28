package cluster

import "testing"

func BenchmarkSubject(b *testing.B) {

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		Subject("prefix", "appKind", "appName", "appID")
	}
}

##Original Benchmark results
_8/26/2018 @ 12:26 CDT_

```
goos: windows
goarch: amd64
pkg: github.com/kelbyers/librsync-go
BenchmarkRollsum_Update256-8             3000000               551 ns/op
BenchmarkRollsum_Update1024-8            1000000              2207 ns/op
BenchmarkRollsum_UpdateComplete-8        3000000               551 ns/op
BenchmarkRollsum_Rollin1-8              1000000000               2.51 ns/op
BenchmarkRollsum_Rollin2-8              300000000                5.03 ns/op
BenchmarkRollsum_Rollin5-8              100000000               11.2 ns/op
BenchmarkRollsum_Rollin10-8             100000000               23.8 ns/op
BenchmarkRollsum_Rollin256-8             2000000               646 ns/op
BenchmarkRollsum_RollinComplete-8       100000000               23.7 ns/op
BenchmarkRollsum_Rollout1-8             1000000000               2.54 ns/op
BenchmarkRollsum_Rollout2-8             300000000                4.68 ns/op
BenchmarkRollsum_Rollout5-8             100000000               12.2 ns/op
BenchmarkRollsum_Rollout10-8            50000000                24.8 ns/op
BenchmarkRollsum_Rollout256-8            2000000               643 ns/op
BenchmarkRollsum_RolloutComplete-8      100000000               24.8 ns/op
BenchmarkRollsum_DigestComplete-8       2000000000               0.72 ns/op
```

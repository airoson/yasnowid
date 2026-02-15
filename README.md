## Yet Another Snowflake ID Generator written in Go
This project is a high-performance, lock-free Snowflake ID generator for Go. It uses atomic operations to ensure thread safety and provides built-in protection against system clock drift.
### Structure
Each ID has the same structure:
|63|62-22|21-12|11-0|
|---|---|---|---|
|Always 0|*timestamp*|datacenter ID & machine ID|*counter*|

*timestamp* is a regular UNIX timestamp in millisecond minus 1288834974657. 
*counter* is incremented on each request and drops to 0 every milliseconds.

### Benchmark results for 1/4/12 Cores
```
goos: windows
goarch: amd64
pkg: github.com/airoson/yasnowid
cpu: AMD Ryzen 5 5600H with Radeon Graphics
BenchmarkGenerator               4772202               251.0 ns/op
BenchmarkGenerator-4             4619058               248.9 ns/op
BenchmarkGenerator-12            4767974               254.3 ns/op
PASS
ok      github.com/airoson/yasnowid     4.720s
```
### Sample usage
The usage is straightforward:
```go
gen, err := yasnowid.NewGenerator(101)
var id int64 = gen.ID() // Generate New Snowflake ID
```
You can also create node ID from datacenter + machine ID:
```go
gen, err := yasnowid.NewGenerator(yasnowid.JoinIDs(1, 2)) 
```
Here datacenter ID is 1 and machine id is 2.
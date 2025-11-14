# col-stats

It reads CSV files and executes the following operations:

- sum: It calculates the sum of all values in the specified column.
- avg: It determines the average value of the specified column.

```bash
# Bench marking. `-run ^$` is to skip running any of the tests in the test file while executing the benchmark to prevent impacting the results.
go test -bench . -run ^$

# With `-benchtime=10x`, you can execute benchmarking 10 times.
go test -bench . -benchtime=10x -run ^$

# With `| tee benchresults00.txt`, you can save the result on a Linux/Unix system.
go test -bench . -benchtime=10x -run ^$ | tee benchresults00.txt

# `-cpuprofile cpu00.pprof` enables the CPU profiler. It creates cpu00.pprof and col-stats.test files
go test -bench . -benchtime=10x -run ^$ -cpuprofile cpu00.pprof

# This analyzes the profling results of cpu00.pprof.
# When the profiler is enabled, try the following:
#  top: It shows where the program is spending most of its time.
#  top -cum: It shows th result based on the culmulative time
#  list <function-name> (e.g. list csv2float): It shows the time that takes to execute the specified function
#  web: It generates a relationship graph (This command requires graphviz. If you use Ubuntu, run `sudo apt-get install graphviz`.).
go tool pprof cpu00.pprof

# `-memprofile` shows how much memory the program is allocating, and it creates a memory profile in mem00.pprof.
go test -bench . -benchtime=10x -run ^$ -memprofile mem00.pprof

# This analyzes the profling results of mem00.pprof.
# Try top, top -cum, list <function-name>, and web that are described above.
go tool pprof -alloc_space mem00.pprof


# `-benchmem` displays the memory allocation.
go test -bench . -benchtime=10x -run ^$ -benchmem | tee benchresults00m.txt
```

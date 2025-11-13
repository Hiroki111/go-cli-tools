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
```

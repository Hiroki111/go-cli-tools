## Useful commands

```bash
go build

GOOS=windows go build

go test -v

# Counting words
echo "test test test" | ./word-counter
cat main.go | ./word-counter

# Counting lines (-e tells echo to interpret escape sequences)
echo -e "a\nb\nc" | ./word-counter -l
cat main.go | ./word-counter -l

# Counting bytes (-n means no new line at the end)
echo -n "a ü あ 😀" | ./word-counter -b
echo "a ü あ 😀" | ./word-counter -b
```

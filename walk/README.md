# walk

It walks through a directory and prints files inside it.

## Example usage

```bash
# Print all files in the root of this program
go run .

# Print all files in /tmp/gomisc/ ending with .go
go run . -root /tmp/gomisc/ -ext .go

# Create a folder for archiving files, then archive .go files in /tmp/gomisc/ into /tmp/gomisc_bkp/
mkdir /tmp/gomisc_bkp
go run . -root /tmp/gomisc/ -ext .go -archive /tmp/gomisc_bkp/
```
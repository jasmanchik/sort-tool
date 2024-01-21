### The sort utility

#### The similarity of the sort console utility.

There is utility support for the following keys:

`-k` — specifying the column to sort (words in a row can act as columns, the separator is a space)

`-n` — sort by numeric value

`-r` — sort in reverse order

`-u` — do not output duplicate lines

Example execute: `go run ./cmd/sort -r poems.txt`
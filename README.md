# Coding Challenge #1: Build Your Own wc Tool

This is my solution to the first problem in the [John Crickett's Coding Challenges](https://codingchallenges.fyi/challenges/challenge-wc/).

It's a Go implementation of the Unix command line tool `wc`.


## Setup

1. Clone the repo
1. Run the tool using of the following approaches:

    ```shell
    # Run the tool automatically using the go command
    $ go run main.go . [-clmw] [file... ]

    # Build a binary and run it manually
    $ go build
    $ ./go-wc [-clmw] [file... ]

    # Install the binary in your environment, and run it:
    $ go install
    $ go-wc [-clmw] [file... ]
    ```
1. Done!


## Usage
```shell
Usage of go-wc:
  -c	print the bytes count
  -l	print the lines count
  -m	print the characters count
  -w	print the words count
```

## Examples

```shell
# Default output
$ go-wc test.txt
   7145   58164  342190 test.txt

# Lines Count
$ go-wc -l test.txt
   7145 test.txt

# Words Count
$ go-wc -w test.txt
  58164 test.txt

# Characters Count
$ go-wc -m test.txt
 339292 test.txt

# Bytes Count
$ go-wc -c test.txt
 342190 test.txt

# Using stdin
$ cat test.txt | go-wc
   7145   58164  342190
```

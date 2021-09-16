# RoSe ![build](https://github.com/suzaku/rose/workflows/build/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/suzaku/rose)](https://goreportcard.com/report/github.com/suzaku/rose)

RoSe is a command line tool that allows you to treat files as sets of rows and perform set operations on them.

The name RoSe comes from **Ro**w**Se**t.


## Usage

Currently, 3 set operations are implemented:

1. Intersection
2. Union
3. Subtraction

All operations assume that the files are sorted in alphabetical order,
so that we can work with big files efficiently.

### Intersection

```bash
# What's common between two files?
rose and file1 file2

# List common files in two directories.
rose and <(ls -a /tmp | sort) <(ls -a ~ | sort)
```

### Union

```bash
# Everything that exist in either of the two files
rose or file1 file2
```

### Subtraction

```bash
# What exist in the first file but not in the second?
rose sub file1 file2
```

## Installation

### Go

```bash
go get -u github.com/suzaku/rose
```

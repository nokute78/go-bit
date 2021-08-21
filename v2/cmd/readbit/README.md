# readbit

A command line tool to read bit from STDIN/File.

## Quick Start
```shell
$ printf "\xff\xff" |./readbit -s 4
```

## Options
```
Usage of readbit:
  -B uint
    	offset (in byte)
  -V	show version
  -b uint
    	offset (in bit)
  -s uint
    	read size(in bit)
  -v	verbose mode
```

## Example(STDIN)

```
$ printf "\xff\xff" |./readbit -s 4
0x0f
```

Read 3bit and offset is 8bit.
```
$ printf "\xff\xff" |./readbit -s 3 -b 8
0x07
```

# jlf — a basic JSON line formatter

![Screenshot](https://user-images.githubusercontent.com/50314/83684793-a2f74880-a5b5-11ea-8b3b-c01f6a50f5e2.jpg)

## Installation

Git clone and then `go install .` (requires Go >= 1.13).

## Usage

You can invoke `jlf` with one or more file names, or simply pipe to it. Other flags:

* `-i` or `--include-rest`: Devote the last column to rest of fields that don't have columns specified.
* `-c=COLUMN` or `--column=COLUMN`: Specify a column. Can be specified multiple times. See below.

### Column format

One or more columns must be specified. Each column has this format:

```
--column=myField:100:blue
         |       |   |
         |       |   |
  field name     |   color
                 |   (blank means default)
                 |
               width
         (blank means auto)
```

Some examples:

* `--column=myField`
* `--column=myField:100`
* `--column=myField::blue`
* `--column=myField:100:blue`

Specifying one or more columns without a width will try fill up the full terminal width evenly.

The special column name `...` can be used. It will evaluate to all the other fields that have not been specified, formatted like so:

```
field1: foo
field2: bar
```

## Colors

* `black`
* `red`
* `green`
* `yellow`
* `blue`
* `magenta`
* `cyan`
* `white`
* `hiblack`
* `hired`
* `higreen`
* `hiyellow`
* `hiblue`
* `himagenta`
* `hicyan`
* `hiwhite`

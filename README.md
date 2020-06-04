# jlf — a basic JSON line formatter

![Screenshot](https://user-images.githubusercontent.com/50314/83770876-9de6d780-a64f-11ea-9242-0b1e32fa47b8.jpg)

`jlf` is a simple formatter that takes streams of JSON, formatting is a readable table. For example, the above screenshot is a rendering of [this](https://gist.github.com/atombender/0868a9c1638ab3bb979b5b9a06807a11) data.

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

package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/jessevdk/go-flags"
)

type options struct {
	Columns     []column `short:"c" long:"column" description:"Columns to print." value-name:"NAME[:LENGTH[:COLOR]]" multiple:"true"`
	IncludeRest bool     `short:"i" long:"include-rest" description:"Devote the last column to rest of fields that don't have columns."`
}

func main() {
	var opts options
	parser := flags.NewParser(&opts, flags.HelpFlag|flags.PassDoubleDash)
	parser.Usage = "[FILE]"

	args, err := parser.Parse()
	if err != nil {
		if e, ok := err.(*flags.Error); ok && e.Type == flags.ErrHelp {
			parser.WriteHelp(os.Stdout)
			os.Exit(2)
			return
		}
		fatal(err)
	}

	if len(opts.Columns) == 0 {
		fatal(errors.New("no columns specified"))
	}

	if len(args) == 0 {
		formatStream(opts, os.Stdin)
		return
	}

	for _, a := range args {
		formatFile(opts, a)
	}
}

func formatFile(opts options, fileName string) {
	f, err := os.Open(fileName)
	if err != nil {
		fatal(err)
	}
	defer func() {
		_ = f.Close()
	}()

	formatStream(opts, f)
}

func formatStream(opts options, r io.Reader) {
	fields := make(map[string]bool, len(opts.Columns))
	for _, c := range opts.Columns {
		if c.field != restColField {
			fields[c.field] = true
		}
	}

	restRowCols := make([]column, len(opts.Columns))
	copy(restRowCols, opts.Columns)
	restRowCols[len(restRowCols)-1].colorFunc = color.CyanString
	restRowCols[len(restRowCols)-1].indent = 2

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		var entry map[string]interface{}
		if err := json.Unmarshal([]byte(scanner.Text()), &entry); err != nil {
			logError(err)
			continue
		}

		var cv []interface{}
		for _, c := range opts.Columns {
			var cell interface{}
			if c.field == restColField {
				keys := make([]string, 0, len(entry))
				for k := range entry {
					if !fields[k] {
						keys = append(keys, k)
					}
				}
				sort.Strings(keys)

				var sb strings.Builder
				for i, k := range keys {
					if i > 0 {
						_, _ = sb.WriteString("\n")
					}
					v := entry[k]
					_, _ = sb.WriteString(k + ": ")
					_, _ = sb.WriteString(c.parseString(v))
				}
				cell = sb.String()
			} else if v, ok := entry[c.field]; ok {
				cell = v
			} else {
				cell = ""
			}
			cv = append(cv, cell)
		}

		printRow(cv, opts.Columns)

		if opts.IncludeRest {
			keys := make([]string, 0, len(entry))
			for k := range entry {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			for _, k := range keys {
				v := entry[k]
				kcv := make([]interface{}, len(cv))
				for i := range kcv {
					kcv[i] = ""
				}
				kcv[len(kcv)-1] = fmt.Sprintf("%s: %s", k, opts.Columns[len(opts.Columns)-1].parseString(v))
				printRow(kcv, restRowCols)
			}
		}
	}
}

var sep = " ｜ "

func printRow(row []interface{}, cols []column) {
	if len(row) == 0 {
		return
	}

	termSize := getTerminalSize()

	var (
		totalW      int
		numOpenCols int
	)
	for _, c := range cols {
		if c.width > 0 {
			totalW += c.width + len(sep)
		} else {
			numOpenCols++
		}
	}

	openW := termSize.w - totalW
	if numOpenCols > 0 {
		openW = (openW / numOpenCols) - len(sep)
	}

	li := 0
	for {
		var (
			any   bool
			cells []string
		)
		for i, value := range row {
			col := cols[i]

			cw := col.width
			if cw == 0 {
				cw = openW
			}

			// TODO: Wrapping for each line is expensive and unnecessary
			str := wrapString(col.parseString(value), cw)

			var cell string
			ls := strings.Split(str, "\n")
			if len(ls) > li {
				cell = ls[li]
				any = true
			} else {
				cell = ""
			}

			cell = col.format(pad(cell, cw))

			cells = append(cells, cell)
		}
		if !any {
			break
		}

		line := strings.Join(cells, sep)

		_, err := os.Stdout.WriteString(line + "\n")
		if err != nil {
			fatal(err)
		}

		li++
	}
}

func pad(s string, width int) string {
	// TODO: Optimize
	for len(s) < width {
		s += " "
	}
	return s
}

func wrapString(s string, width int) string {
	s = strings.ReplaceAll(s, "\r", "")

	if width == 0 {
		return s
	}

	var sb strings.Builder
	for {
		var needBreak bool
		i := strings.IndexByte(s, '\n')
		if i < 0 {
			i = len(s)
		} else {
			i++
		}
		if i > width {
			if width > len(s) {
				i = len(s)
			} else {
				i = width
			}
			needBreak = true
		}
		_, _ = sb.WriteString(s[0:i])
		if needBreak {
			_, _ = sb.WriteString("\n")
		}
		s = s[i:]
		if len(s) == 0 {
			break
		}
	}
	return sb.String()
}

const restColField = "..."

type colorFunc func(format string, a ...interface{}) string

type column struct {
	field     string
	width     int
	indent    int
	colorFunc colorFunc
}

func (c *column) format(s string) string {
	for i := 1; i <= c.indent; i++ {
		s = " " + s
	}
	return c.colorFunc(s)
}

func (c *column) parseString(v interface{}) string {
	if str, ok := v.(string); ok {
		if t, err := time.Parse(time.RFC3339, str); err == nil {
			return t.Format("2006-01-02 15:04:05")
		}
		return str
	}
	return fmt.Sprintf("%v", v)
}

var specRegexp = regexp.MustCompile(`([^:]+)(?::(\d*))?(?::([^:]+))?`)

func (c *column) UnmarshalFlag(value string) error {
	m := specRegexp.FindStringSubmatch(value)
	if len(m) <= 1 {
		return fmt.Errorf("invalid column specification: %q", value)
	}

	c.field = m[1]
	if len(m) >= 3 {
		if len(m[2]) > 0 {
			length, err := strconv.Atoi(m[2])
			if err != nil {
				return err
			}
			if length <= 0 {
				return fmt.Errorf("size must be positive and non-zero: %s", m[2])
			}
			c.width = length
		}
	}
	if len(m) >= 4 {
		cf, err := colorFuncFromStr(m[3])
		if err != nil {
			return err
		}
		c.colorFunc = cf
	}
	return nil
}

func logError(err error) {
	s := color.RedString("Error: "+err.Error()) + "\n"
	_, _ = os.Stderr.Write([]byte(s))
}

func fatal(err error) {
	_, _ = os.Stderr.Write([]byte(err.Error() + "\n"))
	os.Exit(1)
}

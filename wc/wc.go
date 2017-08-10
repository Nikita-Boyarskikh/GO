package main

import (
    "os"
    "io"
    "fmt"
    "bufio"
    "bytes"
    "strings"

    "github.com/jessevdk/go-flags"
)

const (
    Version = "1.0.0"
)

type Options struct {
    Bytes bool `short:"c" long:"bytes" description:"print the byte counts"`
    Chars bool `short:"m" long:"chars" description:"print the character counts"`
    Lines bool `short:"l" long:"lines" description:"print the newline counts"`
    Words bool `short:"w" long:"words" description:"print the word counts"`
    MaxLineLength uint `short:"L" long:"max-line-length" description:"print the maximum display width"`
    Files0from flags.Filename `long:"files0from" description:"read input from the files specified by NUL-terminated names in file F; If F is - then read names from standard input" value-name:"FILE"`
    Args struct {
       File flags.Filename `positional-arg-name:"file" description:"filename for reading" value-name:"FILE" default:"-"`
    } `positional-args:"yes"`
    Version func() `short:"v" long:"version" description:"print the version of the program and exit"`
}

func main() {
    var opts Options

    // Set defaults
    opts.Version = func() {
        fmt.Println(Version)
        os.Exit(0)
    }

    // Create and initialize args parser
    p := flags.NewParser(&opts, flags.HelpFlag)
    args, err := p.Parse()

    if len(args) > 0 {
        fatal(p, "Get extra params: " + strings.Join(args, " "), true)
    }

    if err != nil {
        e, ok := err.(*flags.Error)
        if ok {
            fatal(p, err.Error(), e.Type != flags.ErrHelp)
        } else {
            fatal(nil, err.Error(), false)
        }
    }

    // Check options
    if opts.Bytes || opts.Chars || opts.Words || opts.Files0from != "" || opts.MaxLineLength != 0 {
        fatal(nil, "Not implemented yet", false)
    }

    // Set default -l -w -c
    if !(opts.Bytes || opts.Chars || opts.Words || opts.Lines) {
        opts.Lines, opts.Words, opts.Bytes = true, true, true
    }

    // Set default input file to STDOUT and run
    if opts.Args.File == "-" {
        opts.Args.File = ""
    }

    // Run
    result := calcStrs(string(opts.Args.File))
    fmt.Println(result, opts.Args.File)
}

func calcStrs(filename string) uint {
    var file *os.File

    if filename == "" {
        file = os.Stdout
        return countLinesIntoFile(file)
    } else {
        var err error

        file, err = os.Open(filename)
        if err != nil {
            fatal(nil, err.Error(), false)
        }

	defer file.Close()
        return countLinesIntoFile(file)
    }
}

func countLinesIntoFile(file io.Reader) uint {
    var count uint
    scanner := bufio.NewScanner(file)

    count = 0
    for scanner.Scan() {
        count = count + 1
    }

    if err := scanner.Err(); err != nil {
        fatal(nil, err.Error(), false)
    }

    return count
}

func fatal(p *flags.Parser, message string, with_help bool) {
    var b bytes.Buffer

    if with_help {
        p.WriteHelp(&b)
    }

    fmt.Fprintln(os.Stderr, message, "\n", b.String())
    os.Exit(1)
}

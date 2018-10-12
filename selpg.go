package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/spf13/pflag"
)

type selpg_args struct {
	start_page  int
	end_page    int
	in_filename []string
	page_len    int
	page_type   byte
	print_dest  string
}

var progname string

func main() {
	var args selpg_args
	progname = os.Args[0]
	parseArgs(&args)
	processArgs(args)
	processInput(args)
}

func parseArgs(args *selpg_args) {
	pflag.IntVarP(&args.start_page, "start", "s", 1, "Start page number")
	pflag.IntVarP(&args.end_page, "end", "e", 1, "End page number")
	pflag.StringVarP(&args.print_dest, "printdest", "d", "", "Output to destination pipe")
	pflag.IntVarP(&args.page_len, "pagelength", "l", 72, "Line number of each page")
	fword := pflag.Bool("f", false, "Page type")
	pflag.Parse()
	args.page_type = 'l'
	if *fword {
		args.page_type = 'f'
	}
	args.in_filename = pflag.Args()
}

func processArgs(args selpg_args) {
	args_len := len(os.Args)
	if args_len < 3 {
		fmt.Fprintf(os.Stderr, "%s: There should be at least 3 args\n", progname)
		flag.Usage()
		os.Exit(1)
	}
	if os.Args[1] != "-s" {
		fmt.Fprintf(os.Stderr, "%s: The 1st arg should be -s start_page\n", progname)
		flag.Usage()
		os.Exit(2)
	} else {
		if args.start_page < 1 {
			fmt.Fprintf(os.Stderr, "%s: Invalid start page\n", progname)
			flag.Usage()
			os.Exit(3)
		}
	}
	if os.Args[3] != "-e" {
		fmt.Fprintf(os.Stderr, "%s: The 2nd arg should be -e end_page\n", progname)
		flag.Usage()
		os.Exit(4)
	} else {
		if args.end_page < 1 {
			fmt.Fprintf(os.Stderr, "%s: Invalid end page\n", progname)
			flag.Usage()
			os.Exit(5)
		}
		if args.start_page > args.end_page {
			fmt.Fprintf(os.Stderr, "%s: Start page should be smaller than end page\n", progname)
			flag.Usage()
			os.Exit(6)
		}
	}
}

func processInput(args selpg_args) {
	var in *os.File
	var out *os.File
	var cmd *exec.Cmd
	var page_num, line_num int
	if len(args.in_filename) == 0 {
		in = os.Stdin
	} else {
		var err error
		in, err = os.Open(args.in_filename[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not open the file: %s\n", string(args.in_filename[0]))
			return
		}
	}
	if args.print_dest != "" {
		cmd = exec.Command("/user/bin/lp", fmt.Sprintf("-d%s", args.print_dest))
		reader, writer, err := os.Pipe()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not open pipe to %s\n", args.print_dest)
		}
		cmd.Stdin = reader
		out = writer
	} else {
		out = os.Stdout
	}
	reader := bufio.NewReader(in)
	writer := bufio.NewWriter(out)
	page_num = 1
	if args.page_type == 'l' {
		var line []byte
		line_num = 0
		for true {
			var err error
			line, _, err = reader.ReadLine()
			if err != nil {
				if err == io.EOF {
					break
				}
				fmt.Println(err)
				break
			}
			line_num++
			if line_num > args.page_len {
				page_num++
				line_num = 1
			}
			if page_num >= args.start_page && page_num <= args.end_page {
				line = append(line, '\n')
				writer.Write(line)
				writer.Flush()
			}
		}
	} else {
		for true {
			buffer, err := reader.ReadByte()
			if err == io.EOF {
				break
			}
			if buffer == '\f' {
				page_num++
			}
			if page_num >= args.start_page && page_num <= args.end_page {
				writer.WriteByte(buffer)
				writer.Flush()
			}
		}
	}

	if page_num < args.start_page {
		fmt.Fprintf(os.Stderr, "Start page (%d) is greater than total pages (%d)\n", args.start_page, page_num)
	} else if page_num < args.end_page {
		fmt.Fprintf(os.Stderr, "End page (%d) is greater than total pages (%d)\n", args.end_page, page_num)
	}

	if cmd != nil {
		cmd.Run()
	}
	fmt.Println()
}

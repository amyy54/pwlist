package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/amyy54/pwlist/internal/cartesian"
	"github.com/amyy54/pwlist/internal/formatter"
	"github.com/amyy54/pwlist/internal/reader"
)

var (
	Version     string
	VersionLong string
)

func main() {
	var format_str string
	var rtable string
	var files []string

	var v bool

	var version bool

	log.SetFlags(log.Lshortfile)

	flag.StringVar(&format_str, "format", "", "Format to use. See man page for details.")
	flag.StringVar(&rtable, "rtable", "", "Run a hash function and display it alongside the password. [NTLM, MD5, SHA1, SHA256, SHA512]")
	flag.BoolVar(&v, "v", false, "Verbose.")
	flag.BoolVar(&version, "version", false, "Print the version and exit.")
	flag.Parse()
	files = flag.Args()

	if version {
		if v {
			fmt.Printf("pwlist: %s\n", VersionLong)
		} else {
			fmt.Printf("pwlist: %s\n", Version)
		}
		os.Exit(0)
	}

	if v {
		slog.SetLogLoggerLevel(slog.LevelInfo)
	} else {
		slog.SetLogLoggerLevel(slog.LevelWarn)
	}

	var rtable_hash formatter.HashType

	if len(rtable) != 0 {
		switch strings.ToLower(rtable) {
		case "ntlm":
			rtable_hash = formatter.NTLM
		case "md5":
			rtable_hash = formatter.MD5
		case "sha1":
			rtable_hash = formatter.SHA1
		case "sha256":
			rtable_hash = formatter.SHA256
		case "sha512":
			rtable_hash = formatter.SHA512
		default:
			slog.Warn("Hash string invalid, defaulting to MD5")
			rtable_hash = formatter.MD5
		}
	}

	file_contents, err := reader.ReadFiles(files)
	if err != nil {
		log.Fatalf("Could not read one or more files. Error: %s", err)
		os.Exit(1)
	}

	match_list, iter_args := formatter.GenMatchList(format_str, file_contents)

	slog.Info("Starting cartesian iteration")

	c := cartesian.Iter(iter_args...)

	for combine_line := range c {
		line := strings.Clone(format_str)
		for pos, match := range match_list {
			line = strings.Replace(line, match.FormatBlock, combine_line[pos], 1)
		}
		if len(rtable) > 0 {
			fmt.Printf("%s:%s\n", formatter.HashGen(rtable_hash, line), line)
		} else {
			fmt.Println(line)
		}
	}
	os.Exit(0)
}

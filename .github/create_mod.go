// Command create_mod builds a canonical Go module zip from a source directory,
// using the same golang.org/x/mod/zip code path the module proxy uses.
//
// It lives under .github/ so the go tool never treats it as part of the
// golang.org/x/text module (go ignores directories beginning with "."), and it
// is compiled from a throwaway module outside the checkout so that its own
// go.mod/go.sum can never leak into the produced zip.
//
// usage: create_mod <module-path> <version> <src-dir> <out-zip>
package main

import (
	"fmt"
	"os"

	"golang.org/x/mod/module"
	"golang.org/x/mod/zip"
)

func main() {
	if len(os.Args) != 5 {
		fmt.Fprintln(os.Stderr, "usage: create_mod <module-path> <version> <src-dir> <out-zip>")
		os.Exit(2)
	}
	modPath, version, srcDir, outPath := os.Args[1], os.Args[2], os.Args[3], os.Args[4]

	f, err := os.Create(outPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "create:", err)
		os.Exit(1)
	}
	if err := zip.CreateFromDir(f, module.Version{Path: modPath, Version: version}, srcDir); err != nil {
		fmt.Fprintln(os.Stderr, "CreateFromDir:", err)
		os.Exit(1)
	}
	if err := f.Close(); err != nil {
		fmt.Fprintln(os.Stderr, "close:", err)
		os.Exit(1)
	}
	fmt.Printf("wrote %s for %s@%s from %s\n", outPath, modPath, version, srcDir)
}

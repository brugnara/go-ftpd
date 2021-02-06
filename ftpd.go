package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"text/tabwriter"
)

type ftp struct {
	root string
	path string
}

const (
	errorInvalidCommand    = "(╯°□°）╯︵ ┻━┻\nInvalid command"
	errorIncompleteCommand = "ಠ_ಠ\nMissing or invalid argument"
	errorPathInvalid       = "¯\\_(ツ)_/¯\\\nInvalid path"
	errorGeneric           = "¯\\_(ツ)_/¯\\\nAn error occured"
	errorCantOpenDir       = "ಠ_ಠ\nCan't open directories"
	errorFileTooBig        = "¯\\_(ツ)_/¯\\\nThe file is too big"
	wall                   = `┻┳|
┳┻| _
┻┳| •.•)
┳┻|⊂ﾉ
┻┳|`
)

func newFtp(p string) *ftp {
	p = path.Clean(p)
	return &ftp{p, p}
}

func (f ftp) hello(c io.Writer) {
	fmt.Fprintf(c, "\n%s\nWelcome.\nBrugnara made this in 2021.\n", wall)
	f.help(c)
	f.cursor(c)
}

func (f ftp) cursor(c io.Writer) {
	fmt.Fprintf(c, "$ %s > ", f.currentPath())
}

func (f ftp) help(c io.Writer) {
	fmt.Fprintln(c, `Available commands:
    - ls
    - cd <folder>
    - cat <file>`)
}

func (f *ftp) command(c io.Writer, cmd string) {
	defer func() {
		f.cursor(c)
	}()

	log.Println("Executing command:", cmd)
	xc := splitter(cmd)

	if len(xc) == 0 {
		fmt.Fprintf(c, "%s\n", errorInvalidCommand)
		return
	}

	switch xc[0] {
	default:
		fmt.Fprintf(c, "%s\n", errorInvalidCommand)
	case "cd":
		if len(xc) != 2 {
			fmt.Fprintf(c, "%s\n", errorIncompleteCommand)
			return
		}
		f.cd(c, xc[1])
	case "cat":
		f.cat(c, xc[1])
	case "ls":
		f.ls(c)
	case "quit":
		fmt.Fprintln(c, "Bye!")
	}
}

func (f *ftp) cd(c io.Writer, dir string) {
	if p, err := f.wannaBe(dir); err != nil {
		fmt.Fprintln(c, errorPathInvalid)
	} else {
		f.path = p
	}
}

func (f ftp) wannaBe(dir string) (string, error) {
	p := path.Clean(path.Join(f.root, f.currentPath(), dir))
	if !strings.HasPrefix(p, f.root) {
		p = path.Join(f.root, p)
	}
	p = strings.ReplaceAll(p, "..", "")
	// check existance
	if _, err := os.Open(p); err != nil {
		return f.path, err
	}
	return p, nil
}

func (f ftp) ls(c io.Writer) {
	fl, err := os.Open(path.Join(f.root, f.currentPath()))
	if err != nil {
		fmt.Fprintf(c, "%s\n", errorPathInvalid)
		return
	}
	defer fl.Close()

	xf, err := fl.Readdir(-1)
	if err != nil {
		fmt.Fprintf(c, "%s\n", errorPathInvalid)
		return
	}
	tab := tabwriter.NewWriter(c, 0, 8, 2, ' ', 0)
	fmt.Fprintln(tab, "Size\tType\tName")
	fmt.Fprintln(tab, "====\t====\t====")
	for _, ff := range xf {
		var tp string

		if ff.IsDir() {
			tp = "dir"
		} else {
			tp = "file"
		}

		fmt.Fprintln(tab, strings.Join([]string{
			tp,
			toSize(ff.Size()),
			cut(ff.Name(), 40, ".."),
		}, "\t"))
	}
	tab.Flush()
}

func (f ftp) currentPath() string {
	ret := strings.Replace(f.path, f.root, "", -1)
	if ret == "" {
		ret = "/"
	}
	return ret
}

func (f ftp) cat(c io.Writer, file string) {
	quitError := func(err error, output string) {
		log.Println(err)
		fmt.Fprintln(c, output)
	}
	//
	file, err := f.wannaBe(file)
	if err != nil {
		quitError(err, errorCantOpenDir)
		return
	}
	fl, err := os.Open(file)
	if err != nil {
		quitError(err, errorCantOpenDir)
		return
	}
	defer fl.Close()
	stat, err := fl.Stat()
	if err != nil || stat.IsDir() {
		quitError(err, errorCantOpenDir)
		return
	}
	if stat.Size() > pow(1024, 2) {
		quitError(nil, errorFileTooBig)
		return
	}
	io.Copy(c, fl)
}

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

type ftpd struct {
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

func newFtp(p string) *ftpd {
	p = strings.TrimSpace(p)
	p = path.Clean(p)
	if p == "" || p == "." {
		p = "./public"
	}
	return &ftpd{p, p}
}

func (f ftpd) hello(c io.Writer) {
	fmt.Fprintf(c, "\n%s\nWelcome.\nBrugnara made this in 2021.\n", wall)
	f.help(c)
	f.cursor(c)
}

func (f ftpd) cursor(c io.Writer) {
	fmt.Fprintf(c, "$ %s > ", f.currentPath())
}

func (f ftpd) help(c io.Writer) {
	fmt.Fprintln(c, `Available commands:
    - ls
    - cd <folder>
    - cat <file>`)
}

func (f *ftpd) command(c io.Writer, cmd string) bool {
	defer func() {
		f.cursor(c)
	}()

	log.Println("Executing command:", cmd)
	xc := splitter(cmd)

	if len(xc) == 0 {
		fmt.Fprintf(c, "%s - '%s'\n", errorInvalidCommand, "<no command>")
		return false
	}

	switch xc[0] {
	default:
		fmt.Fprintf(c, "%s - '%s'\n", errorInvalidCommand, xc[0])
		return false
	case "cd":
		if len(xc) != 2 {
			fmt.Fprintf(c, "%s\n", errorIncompleteCommand)
			return false
		}
		return f.cd(c, xc[1])
	case "cat":
		if len(xc) != 2 {
			fmt.Fprintf(c, "%s\n", errorIncompleteCommand)
			return false
		}
		return f.cat(c, xc[1])
	case "ls":
		return f.ls(c)
	case "quit":
		fmt.Fprintln(c, "Bye!")
	}
	return true
}

func (f *ftpd) cd(c io.Writer, dir string) bool {
	p, err := f.wannaBe(dir)
	if err != nil {
		fmt.Fprintln(c, errorPathInvalid)
		return false
	}
	f.path = p
	return true

}

func (f ftpd) wannaBe(dir string) (string, error) {
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

func (f ftpd) ls(c io.Writer) bool {
	fl, err := os.Open(path.Join(f.root, f.currentPath()))
	if err != nil {
		fmt.Fprintf(c, "%s\n", errorPathInvalid)
		return false
	}
	defer fl.Close()

	xf, err := fl.Readdir(-1)
	if err != nil {
		fmt.Fprintf(c, "%s\n", errorPathInvalid)
		return false
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
	return true
}

func (f ftpd) currentPath() string {
	ret := strings.Replace(f.path, f.root, "", -1)
	if ret == "" {
		ret = "/"
	}
	return ret
}

func (f ftpd) cat(c io.Writer, file string) bool {
	quitError := func(err error, output string) bool {
		log.Println(err)
		fmt.Fprintln(c, output)
		return false
	}
	//
	file, err := f.wannaBe(file)
	if err != nil {
		return quitError(err, errorCantOpenDir)
	}
	fl, err := os.Open(file)
	if err != nil {
		return quitError(err, errorCantOpenDir)
	}
	defer fl.Close()
	stat, err := fl.Stat()
	if err != nil || stat.IsDir() {
		return quitError(err, errorCantOpenDir)
	}
	if stat.Size() > pow(1024, 2) {
		return quitError(nil, errorFileTooBig)
	}
	_, err = io.Copy(c, fl)
	return err == nil
}

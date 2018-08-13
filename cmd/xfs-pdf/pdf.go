package main

import "io"
import "os"
import "github.com/ledongthuc/pdf"

func main() {
	if len(os.Args) == 1 {
		return
	}
	f, r, err := pdf.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer f.Close()
	b, err := r.GetPlainText()
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(os.Stdout, b)
	if err != nil {
		panic(err)
	}
}

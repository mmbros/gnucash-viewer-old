package main

/*
	References:
	* [Parsing huge xml files with go](http://blog.davidsingleton.org/parsing-huge-xml-files-with-go/)
	* [xml:nested namespace](https://groups.google.com/forum/#!topic/golang-nuts/QFDHM7_VFks)
	* [xml unmarshall example](https://golang.org/pkg/encoding/xml/#example_Unmarshal)

	* [Golang XML Unmarshal and time.Time fields](http://stackoverflow.com/questions/17301149/golang-xml-unmarshal-and-time-time-fields)
*/

import (
	"flag"
	"fmt"
	"time"

	"github.com/mmbros/go/gnucash-go-viewer/model"
	gncxml "github.com/mmbros/go/gnucash-go-viewer/xml"
)

var gnucashPath = flag.String("gnucash-file", "data/data.gnucash", "GnuCash file path")

// --------------------------------------------------------------------------

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("\n--\nfunction %s took %v\n", name, elapsed)
}

func main() {
	defer timeTrack(time.Now(), "task duration:")

	gnc, err := gncxml.ReadFile(*gnucashPath)
	if err != nil {
		panic(err)
	}

	book, err := model.NewBook(gnc)
	if err != nil {
		panic(err)
	}

	fmt.Printf("books: %d\n", len(gnc.Books))
	fmt.Printf("accounts: %d\n", len(book.Accounts.Map))
	fmt.Printf("tansactions: %d\n", len(book.Transactions))

	fmt.Printf("root: %v\n", book.Accounts.Root)
	//	book.Accounts.PrintTree("")

	a := book.Accounts.Map["c623a615013986b49b88d391ce9fd0f1"]
	for j, t := range a.Transactions {
		fmt.Printf("%04d) %v %s\n", j, t.DatePosted, t.Description)
	}
}

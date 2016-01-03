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
	"strings"
	"time"

	"github.com/mmbros/gnucash-viewer/model"
	"github.com/mmbros/gnucash-viewer/numeric"
	gncxml "github.com/mmbros/gnucash-viewer/xml"
)

var gnucashPath = flag.String("gnucash-file", "data/data.gnucash", "GnuCash file path")

// --------------------------------------------------------------------------

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("\n--\nfunction %s took %v\n", name, elapsed)
}

func main2() {
	defer timeTrack(time.Now(), "task duration:")

	a := numeric.New(-3, -2)
	b := numeric.New(5, -7)
	z := numeric.Numeric{}
	fmt.Printf("a = %s\n", a)
	fmt.Printf("b = %s\n", b)
	fmt.Printf("z = %s\n", z)
	z = numeric.Sub(&a, &b)
	fmt.Printf("z = a-b\n")
	fmt.Printf("a = %s\n", a)
	fmt.Printf("b = %s\n", b)
	fmt.Printf("z = %s\n", z)

}

// StringLeft function
func StringLeft(s string, n int) string {
	if n <= 0 {
		return ""
	}
	if L := len([]rune(s)); L < n {
		n = L
	}
	return s[:n]
}

// StringPad function
func StringPad(s string, n int, pad string) string {
	if n <= 0 {
		return ""
	}
	L := len([]rune(s))
	if L >= n {
		return s[:n]
	}
	return s + strings.Repeat(pad, n-L)
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

	//fmt.Printf("root: %v\n", book.Accounts.Root)
	//	book.Accounts.PrintTree("")

	/*
		var tot int
		for _, t := range book.Transactions {
			if len(t.Splits) <= 2 {
				continue
			}
			tot++
			fmt.Printf("%03d) transaction.ID = %s\n", tot, t.ID)
			for _, s := range t.Splits {
				fmt.Printf("    %s %v\n", s.Account.Name, s.Value)
			}
		}
	*/

	//a := book.Accounts.Map["c623a615013986b49b88d391ce9fd0f1"]
	acc := book.Accounts.ByName("Stipendio")
	if acc == nil {
		panic("Account non trovato")
	}
	for j, at := range acc.AccountTransactionList {
		fmt.Printf("%02d) %s %s %5.2f %7.2f %7.0f\n",
			j+1,
			at.Transaction.DatePosted.Format("2006-01-02"),
			StringPad(at.Description(), 41, "."),
			at.PlusValue.Float64(),
			at.MinusValue.Float64(),
			at.Balance.Float64(),
		)
	}

	for _, at := range acc.AccountTransactionList {
		fmt.Printf("'%s', %4.0f\n",
			at.Description()[21:],
			at.MinusValue.Float64(),
		)
	}
	/*
		for j, at := range acc.AccountTransactionList {
			fmt.Printf("%01d) %v\n", j+1, at)
		}
	*/
}

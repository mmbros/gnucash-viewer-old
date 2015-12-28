package model

import (
	"errors"

	gncxml "github.com/mmbros/gnucash-viewer/xml"
)

// Book type
type Book struct {
	Accounts     *Accounts
	Transactions Transactions
}

// NewBook function
func NewBook(gnc *gncxml.Gnc) (*Book, error) {
	if gnc == nil {
		return nil, errors.New("GNC must be not nil")
	}
	switch len(gnc.Books) {
	case 0:
		return nil, errors.New("BOOK not found")
	case 1:
		break
	default:
		return nil, errors.New("Not Implemented: multiple BOOK")
	}
	xmlBook := gnc.Books[0]
	book := Book{}
	var err error

	// init Accounts
	book.Accounts, err = newAccountsFromXML(xmlBook.AccountList)
	if err != nil {
		return nil, err
	}

	// init Transactions
	book.Transactions, err = newTransactionsFromXML(xmlBook.TransactionList, book.Accounts)
	if err != nil {
		return nil, err
	}

	return &book, nil
}

package model

import (
	"errors"
	"fmt"
	"sort"
	"time"

	gncxml "github.com/mmbros/gnucash-viewer/xml"
)

// Transactions type
type Transactions []*Transaction

// Transaction type
type Transaction struct {
	ID          string
	Currency    string
	DatePosted  time.Time
	DateEntered time.Time
	Description string
	Splits      []*Split
}

// Split type
type Split struct {
	ID              string
	ReconciledState string
	ReconcileDate   time.Time
	Value           GncNumeric
	Quantity        GncNumeric
	Account         *Account
}

func timeParse(value string, nullable bool) (time.Time, error) {
	if nullable && len(value) == 0 {
		return time.Time{}, nil
	}
	return time.Parse("2006-01-02 15:04:05 -0700", value)
}

func formatError(object, field, id string, err error) error {
	return fmt.Errorf("Invalid %s in %s (ID=%s): %s", field, object, id, err.Error())
}

func newSplitFromXML(xmlSplit *gncxml.Split, accounts *Accounts) (*Split, error) {
	// check ReconcileDate
	reconcileDate, err := timeParse(xmlSplit.ReconcileDate, true)
	if err != nil {
		return nil, formatError("Split", "ReconcileDate", xmlSplit.ID, err)
	}
	// check Value
	value, err := NewGncNumericFromString(xmlSplit.Value)
	if err != nil {
		return nil, formatError("Split", "Value", xmlSplit.ID, err)
	}
	// check Quantity
	quantity, err := NewGncNumericFromString(xmlSplit.Quantity)
	if err != nil {
		return nil, formatError("Split", "Quantity", xmlSplit.ID, err)
	}
	// check Account
	account, ok := accounts.Map[xmlSplit.AccountID]
	if !ok {
		return nil, formatError("Split", "AccountID", xmlSplit.ID, errors.New("Account not found"))
	}

	// initialize Transaction object
	split := Split{
		ID:              xmlSplit.ID,
		ReconciledState: xmlSplit.ReconciledState,
		ReconcileDate:   reconcileDate,
		Value:           *value,
		Quantity:        *quantity,
		Account:         account,
	}

	return &split, nil
}

func newTransactionFromXML(xmlTransaction *gncxml.Transaction, accounts *Accounts) (*Transaction, error) {
	// check DatePosted
	datePosted, err := timeParse(xmlTransaction.DatePosted, false)
	if err != nil {
		return nil, formatError("Transaction", "DatePosted", xmlTransaction.ID, err)
	}
	// check DateEntered
	dateEntered, err := timeParse(xmlTransaction.DateEntered, false)
	if err != nil {
		return nil, formatError("Transaction", "DateEntered", xmlTransaction.ID, err)
	}
	// initialize Splits
	splits := []*Split{}
	for _, xmlSplit := range xmlTransaction.SplitList {
		split, err := newSplitFromXML(&xmlSplit, accounts)
		if err != nil {
			return nil, err
		}
		splits = append(splits, split)
	}

	// initialize Transaction object
	transaction := Transaction{
		ID:          xmlTransaction.ID,
		Currency:    xmlTransaction.Currency,
		DatePosted:  datePosted,
		DateEntered: dateEntered,
		Description: xmlTransaction.Description,
		Splits:      splits,
	}

	return &transaction, nil
}

func newTransactionsFromXML(xmlTransactionList []gncxml.Transaction, accounts *Accounts) (Transactions, error) {
	// step 0: allocate Transactions object
	transactions := Transactions{}

	// step 1: populate Transactions
	for _, xmlTransaction := range xmlTransactionList {
		t, err := newTransactionFromXML(&xmlTransaction, accounts)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	// step 2: sort Transactions by DatePosted
	sort.Sort(byDatePosted(transactions))

	// step 3: update each Account.transactions field
	// the Account.transactions will be already ordered by DatePosted
	for _, t := range transactions {
		for _, s := range t.Splits {
			a := s.Account
			a.Transactions = append(a.Transactions, t)
		}

	}

	return transactions, nil
}

// used to sort Transactions
type byDatePosted []*Transaction

func (t byDatePosted) Len() int           { return len(t) }
func (t byDatePosted) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t byDatePosted) Less(i, j int) bool { return t[i].DatePosted.Before(t[j].DatePosted) }

package model

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/mmbros/gnucash-viewer/numeric"
	gncxml "github.com/mmbros/gnucash-viewer/xml"
)

// Accounts type
type Accounts struct {
	Root *Account
	Map  map[string]*Account
}

// Account type
type Account struct {
	ID                     string
	Type                   *AccountType
	Name                   string
	Description            string
	Currency               string
	Parent                 *Account
	Children               []*Account
	AccountTransactionList []*AccountTransaction
}

// AccountTransaction type
type AccountTransaction struct {
	Transaction *Transaction
	Split       *Split
	PlusValue   numeric.Numeric
	MinusValue  numeric.Numeric
	Balance     numeric.Numeric
}

func newAccountFromXML(xmlAccount *gncxml.Account) (*Account, error) {
	// check account type
	accType, ok := AccountTypes[xmlAccount.Type]
	if !ok {
		return nil, fmt.Errorf("Invalid AccountType: %s", xmlAccount.Type)
	}

	// initialize Account object
	account := Account{
		ID:          xmlAccount.ID,
		Type:        &accType,
		Name:        xmlAccount.Name,
		Description: xmlAccount.Description,
		Currency:    xmlAccount.Currency,
	}

	return &account, nil
}

func newAccountsFromXML(xmlAccountList []gncxml.Account) (*Accounts, error) {
	// step 0: allocate Accounts object
	a := &Accounts{Map: map[string]*Account{}}

	// step 1: populate Accounts.Map
	for _, xmlAccount := range xmlAccountList {

		// check account unique id
		if _, ok := a.Map[xmlAccount.ID]; ok {
			return nil, fmt.Errorf("Multiple accounts with same ID: %s", xmlAccount.ID)
		}

		// initialize account
		account, err := newAccountFromXML(&xmlAccount)
		if err != nil {
			return nil, err
		}

		// add Account object to Accounts.Map
		a.Map[xmlAccount.ID] = account
	}

	// step 2: initilize root account and parent/children fields
	for _, xmlAccount := range xmlAccountList {
		account := a.Map[xmlAccount.ID]

		if len(xmlAccount.ParentID) == 0 {
			// found root account
			if xmlAccount.Type != "ROOT" {
				return nil, fmt.Errorf("Account of type ROOT can't have parent: Account.ID = %s", xmlAccount.ID)
			}
			if a.Root != nil {
				return nil, errors.New("Not Implemented: multiple ROOT account")
			}
			a.Root = account

		} else {
			// not root account: set parent and children

			parent := a.Map[xmlAccount.ParentID]
			if parent == nil {
				return nil, fmt.Errorf("Parent account not found: ParentID = %s", xmlAccount.ParentID)
			}

			account.Parent = parent
			parent.Children = append(parent.Children, account)
		}
	}

	// step 3: sort each account.children by name
	for _, account := range a.Map {
		sort.Sort(byAccountName(account.Children))
	}

	return a, nil
}

// used to sort each Account.children list
type byAccountName []*Account

func (a byAccountName) Len() int           { return len(a) }
func (a byAccountName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byAccountName) Less(i, j int) bool { return strings.Compare(a[i].Name, a[j].Name) < 0 }

// auxPrintTree is a PrintTree auxiliary function
func auxPrintTree(act *Account, level int, indent string) {
	fmt.Printf("%s[%s] %s (%s)\n", strings.Repeat(indent, level), strings.ToUpper(act.Type.label), act.Name, act.Currency)

	for _, child := range act.Children {
		auxPrintTree(child, level+1, indent)
	}
}

// PrintTree prints account tree
func (accounts *Accounts) PrintTree(indent string) {
	if indent == "" {
		indent = "  "
	}

	if (accounts == nil) || (accounts.Root == nil) {
		fmt.Println("<nil>")
		return
	}

	auxPrintTree(accounts.Root, 0, indent)
}

// ByName return the account with the given name
func (accounts *Accounts) ByName(name string) *Account {
	for _, acc := range accounts.Map {
		if acc.Name == name {
			return acc
		}
	}
	return nil
}

// Balance returns the balance of the account
func (a *Account) Balance() numeric.Numeric {
	return a.AccountTransactionList[len(a.AccountTransactionList)-1].Balance
}

func (accounts *Accounts) postInit(transactions Transactions) {
	// update each Account.AccountTransactionList field
	// the Account.AccountTransactionList will be already ordered by DatePosted
	for _, t := range transactions {
		for _, s := range t.Splits {
			at := AccountTransaction{Transaction: t, Split: s}
			a := s.Account
			a.AccountTransactionList = append(a.AccountTransactionList, &at)
		}
	}
	// initialize account balance
	for _, a := range accounts.Map {
		var balance numeric.Numeric
		for _, at := range a.AccountTransactionList {

			v := at.Split.Value
			balance.AddEqual(v)
			at.Balance.Copy(&balance)

			if v.Sign() >= 0 {
				at.PlusValue.Copy(v)
			} else {
				v.NegEqual()
				at.MinusValue.Copy(v)
			}
		}
	}
}

// Description returns Split.Memo if not null, else Transaction.Description.
func (at *AccountTransaction) Description() string {
	if at == nil || at.Transaction == nil {
		return "<nil>"
	}
	if at.Split != nil && len(at.Split.Memo) > 0 {
		return at.Split.Memo
	}
	return at.Transaction.Description
}

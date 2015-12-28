package model

// AccountType type
type AccountType struct {
	label        string
	root         bool
	invertValues bool
	plusLabel    string
	minusLabel   string
}

// AccountTypes is the map of all AccountType
var AccountTypes = map[string]AccountType{
	"ROOT": AccountType{
		label: "Root",
		root:  true},
	"LIABILITY": AccountType{
		label:        "Liability",
		invertValues: true,
		plusLabel:    "Decrease",
		minusLabel:   "Increase",
	},
	"ASSET": AccountType{
		label:      "Asset",
		plusLabel:  "Increase",
		minusLabel: "Decrease",
	},
	"RECEIVABLE": AccountType{
		label:      "Receivible",
		plusLabel:  "Increase",
		minusLabel: "Decrease",
	},
	"EXPENSE": AccountType{
		label:      "Expense",
		plusLabel:  "Expense",
		minusLabel: "Rebate",
	},
	"INCOME": AccountType{
		label:        "Income",
		invertValues: true,
		plusLabel:    "Charge",
		minusLabel:   "Income",
	},
	"EQUITY": AccountType{
		label:        "Equity",
		invertValues: true,
		plusLabel:    "Decrease",
		minusLabel:   "Increase",
	},
	"BANK": AccountType{
		label:      "Bank",
		plusLabel:  "Deposit",
		minusLabel: "Withdrawal",
	},
	"CASH": AccountType{
		label:      "Cash",
		plusLabel:  "Receive",
		minusLabel: "Spend",
	},
	"CREDIT": AccountType{
		label:      "Credit",
		plusLabel:  "Increase",
		minusLabel: "Decrease",
	},
}

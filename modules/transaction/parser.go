package transaction

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/masudur-rahman/expense-tracker-bot/models"
	"github.com/masudur-rahman/expense-tracker-bot/pkg"
)

/*
	transaction examples: [any keyword can be anywhere, like in natural language]
		transfer 2000 from brac to dbbl on 2020-01-01 note "Bill payment"
		spend 1000 for food-rest on "Jan 13, 2013" from dbbl note "Lunch"
		earn 5000 to brac on 2020-01-01 note "Salary"
		borrow 1000 from user to brac on 2020-01-01
		return 1000 to user from brac on 2020-01-01
		lend 1000 to user from brac on 2020-01-01
		recover 1000 from user to brac on 2020-01-01

	verb keywords: (<keyword> <amount>)
		- transfer
		- expense, spend
		- income, earn
		- borrow
		- return
		- lend
		- recover

	other keywords:
		- from 	[source account (for normal transaction)] [person (for borrow and recover)]
		- to 	[destination account (for normal transaction)] [person (for lend and return)]
		- for 	[subcategory]
		- on 	[date]
		- at	[time]
		- note	[note]
*/

type transactionParser struct {
	txn         models.Transaction
	txnType     models.TransactionType
	amount      string
	from        *string
	to          *string
	fromValue   string
	toValue     string
	subcategory string
	date        string
	time        string
	note        string
}

func ParseTransaction(text string) (models.Transaction, error) {
	opts := transactionParser{}
	kv := pkg.SplitString(text, ' ')
	if len(kv)%2 != 0 {
		return opts.txn, fmt.Errorf("invalid transaction format")
	}

	opts.from = &opts.txn.SrcID
	opts.to = &opts.txn.DstID
	for idx := 0; idx < len(kv); idx += 2 {
		word := strings.ToLower(kv[idx])
		if opts.isVerbKeyword(word) {
			opts.amount = kv[idx+1]
		} else if word == "from" {
			opts.fromValue = kv[idx+1]
		} else if word == "to" {
			opts.toValue = kv[idx+1]
		} else if word == "for" {
			opts.subcategory = kv[idx+1]
		} else if word == "on" {
			opts.date = kv[idx+1]
		} else if word == "at" {
			opts.time = kv[idx+1]
		} else if word == "note" {
			opts.note = kv[idx+1]
		}
	}

	err := opts.parseTransaction()
	return opts.txn, err
}

func (p *transactionParser) isVerbKeyword(keyword string) bool {
	switch keyword {
	case "transfer", "transferred":
		p.txnType = models.TransferTransaction
		p.subcategory = "fin-bank"
	case "withdraw", "withdrew":
		p.txnType = models.TransferTransaction
		p.subcategory = "fin-with"
		p.toValue = "cash"
	case "deposit", "deposited":
		p.txnType = models.TransferTransaction
		p.subcategory = "fin-deposit"
		p.fromValue = "cash"
	case "expense", "spend", "spent":
		p.txnType = models.ExpenseTransaction
	case "giveaway", "donate", "donated":
		p.txnType = models.ExpenseTransaction
		p.subcategory = "misc-give"
	case "income", "earn", "earned":
		p.txnType = models.IncomeTransaction
	case "borrow", "borrowed":
		p.txnType = models.IncomeTransaction
		p.subcategory = models.BorrowSubcategoryID
		p.from = &p.txn.DebtorCreditorName
	case "return", "returned":
		p.txnType = models.ExpenseTransaction
		p.subcategory = models.BorrowReturnSubID
		p.to = &p.txn.DebtorCreditorName
	case "lend", "lent":
		p.txnType = models.ExpenseTransaction
		p.subcategory = models.LoanSubcategoryID
		p.to = &p.txn.DebtorCreditorName
	case "recover", "recovered", "collect", "collected":
		p.txnType = models.IncomeTransaction
		p.subcategory = models.LoanRecoverySubID
		p.from = &p.txn.DebtorCreditorName
	case "flexi":
		p.txnType = models.ExpenseTransaction
		p.subcategory = "fin-flexi"
	default:
		return false
	}
	return true
}

func (p *transactionParser) parseTransaction() error {
	p.txn.Type = p.txnType
	p.txn.SubcategoryID = p.subcategory
	p.txn.Remarks = p.note

	if p.txn.SubcategoryID == "" {
		if p.txn.Type == models.TransferTransaction {
			if p.txn.SrcID == "cash" {
				p.txn.SubcategoryID = "fin-deposit"
			} else if p.txn.DstID == "cash" {
				p.txn.SubcategoryID = "fin-with"
			} else if p.txn.DstID == "credit" {
				p.txn.SubcategoryID = "fin-ccpay"
			}
		} else {
			p.txn.SubcategoryID = "misc-misc"
		}
	}

	p.parseFromTo()
	if err := p.parseAmount(); err != nil {
		return err
	}
	return p.parseTransactionTime()
}

func (p *transactionParser) parseFromTo() {
	if p.fromValue != "" {
		*p.from = p.fromValue
	}
	if p.toValue != "" {
		*p.to = p.toValue
	}
}

func (p *transactionParser) parseTransactionTime() error {
	var year, day, hour, minute, second int
	var month time.Month
	date, err := pkg.ParseDate(p.date)
	if err != nil {
		return err
	}

	tim, err := pkg.ParseTime(p.time)
	if err != nil {
		return err
	}

	year, month, day = date.Date()
	if p.date != "" && p.time == "" { // if date is provided but time is not provided, use 12:00 AM
		hour, minute, second = 0, 0, 0
	} else {
		hour, minute, second = tim.Clock()
	}

	p.txn.Timestamp = time.Date(year, month, day, hour, minute, second, 0, date.Location()).Unix()

	return nil
}

func (p *transactionParser) parseAmount() error {
	var err error
	p.txn.Amount, err = strconv.ParseFloat(p.amount, 64)
	return err
}

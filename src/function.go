package function

import (
	"github.com/credit-to-debit/truelayer"
	"github.com/credit-to-debit/starling"
	"net/http"
)

func Function(http.ResponseWriter, *http.Request) {
	transactions := truelayer.GetTransactions()
	accountID := starling.GetAccount()
	potID := starling.GetPot(accountID)
	for _, transaction := range transactions.TransactionList{
		starling.AddMoneyToGoal(int(transaction.Amount * 100), "GBP", accountID, potID)
	}

}

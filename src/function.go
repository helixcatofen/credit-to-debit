package function

import (
	"github.com/credit-to-debit/truelayer"
	"github.com/credit-to-debit/starling"
	"net/http"
	"fmt"
)

func Function(http.ResponseWriter, *http.Request) {
	transactions := truelayer.GetTransactions()
	fmt.Println("%v", transactions)
	accountID := starling.GetAccount()
	potID := starling.GetPot(accountID)
	for _, transaction := range transactions.TransactionList{
		fmt.Println("Added %f to goal", transaction.Amount)
		starling.AddMoneyToGoal(int(transaction.Amount * 100), "GBP", accountID, potID)
	}

}

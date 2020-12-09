package function

import (
	"github.com/credit-to-debit/truelayer"
	"github.com/credit-to-debit/starling"
	"net/http"
	"fmt"
)

func Function(http.ResponseWriter, *http.Request) {
	transactions := truelayer.GetTransactions()
	fmt.Printf("%v\n", transactions)
	accountID := starling.GetAccount()
	potID := starling.GetPot(accountID)
	for _, transaction := range transactions.TransactionList{
		if(transaction.Amount > 0){
			fmt.Printf("Added %f to goal\n", transaction.Amount)
			starling.AddMoneyToGoal(int(transaction.Amount * 100), "GBP", accountID, potID)
		}
	}

}

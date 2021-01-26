package firebase

import (
	"context"
	"log"

	"github.com/credit-to-debit/truelayer"
	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	// "google.golang.org/api/iterator"
	// "google.golang.org/api/option"
)

func addTransaction(ctx context.Context, client *firestore.Client, transaction truelayer.Transaction) {
	_, err := client.Collection("transactions").Doc(transaction.Id).Set(ctx, transaction)
	if err != nil {
		log.Fatalf("Failed adding %v: %v", transaction, err)
	}
}

// func deleteTransaction(ctx context.Context, client *firestore.Client, id string) {
// 	_, err := client.Collection("transactions").Doc(id).Delete(ctx)
// 	if err != nil {
// 		// Handle any errors in an appropriate way, such as returning them.
// 		log.Printf("Failed deleting: %s, %s", id, err)
// 	}
// }

func transactionExists(ctx context.Context, client *firestore.Client, transaction truelayer.Transaction) bool {
	doc, _:= client.Collection("transactions").Doc(transaction.Id).Get(ctx)
	return doc.Exists()
}

func CheckTransactions(transactions []truelayer.Transaction) []bool{
	isTransactionNew := make([]bool, len(transactions))
	projectID := "credit-to-debit"
	ctx := context.Background()
	conf := &firebase.Config{ProjectID: projectID}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	for i, transaction := range transactions{
		if transactionExists(ctx, client, transaction){
			isTransactionNew[i] = false
		}else{
			addTransaction(ctx, client, transaction)
			isTransactionNew[i] = true
		}
	}
	defer client.Close()
	return isTransactionNew
}

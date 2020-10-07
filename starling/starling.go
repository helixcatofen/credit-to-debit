package starling

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
)

const baseUrl string = "https://api.starlingbank.com"

var starlingKey string = os.Getenv("STARLING_API_KEY")
var bearer string = "Bearer " + starlingKey

type Accounts struct {
	Accounts []Account `json:"accounts"`
}

type Account struct {
	ID   string `json:"accountUid"`
	Type string `json:"accountType"`
}

type Transfer struct {
	Amount struct {
		Currency   string `json:"currency"`
		MinorUnits int    `json:"minorUnits"`
	} `json:"amount"`
}

func listPots(accountId string) {

}

func GetAccount() string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", baseUrl+"/api/v2/accounts", nil)
	req.Header.Add("Authorization", bearer)

	response, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}

	responseData, _ := ioutil.ReadAll(response.Body)

	var accounts Accounts
	json.Unmarshal(responseData, &accounts)
	for _, account := range accounts.Accounts {
		if account.Type == "PRIMARY" {
			return account.ID
		}
	}

	return ""
}

func AddMoneyToGoal(amount int, currency string, accountID string, goalID string) bool {
	var body Transfer
	body.Amount.Currency = currency
	body.Amount.MinorUnits = amount
	reqBody, _ := json.Marshal(&body)
	id, _ := uuid.NewUUID()
	endpoint := fmt.Sprintf("/api/v2/account/%s/savings-goals/%s/add-money/%s", accountID, goalID, id)

	client := &http.Client{}
	req, err := http.NewRequest("PUT", baseUrl+endpoint, bytes.NewBuffer(reqBody))
	req.Header.Add("Authorization", bearer)
	req.Header.Set("Content-Type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}
	if response.StatusCode != 200 {
		return false
	}
	return true
}

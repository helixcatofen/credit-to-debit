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

type Pots struct {
	PotList []struct {
		ID   string `json:"savingsGoalUid"`
		Name string `json:name`
	} `json:"savingsGoalList"`
}

func GetPot(accountId string) string {
	client := &http.Client{}
	endpoint := fmt.Sprintf("%s/api/v2/account/%s/savings-goals", baseUrl, accountId)
	req, err := http.NewRequest("GET", endpoint, nil)
	req.Header.Add("Authorization", bearer)

	response, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}

	if response.StatusCode != 200 {
		fmt.Printf("Starling: %s\n", response.Status)
	}

	responseData, _ := ioutil.ReadAll(response.Body)

	var pots Pots
	json.Unmarshal(responseData, &pots)
	for _, pot := range pots.PotList {
		if pot.Name == "Credit Card" {
			return pot.ID
		}
	}

	return ""
}

func GetAccount() string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", baseUrl+"/api/v2/accounts", nil)
	req.Header.Add("Authorization", bearer)

	response, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}

	if response.StatusCode != 200 {
		fmt.Printf("Starling: %s\n", response.Status)
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

	id, _ := uuid.NewUUID()
	var endpoint string
	if amount > 0 {
		endpoint = fmt.Sprintf("/api/v2/account/%s/savings-goals/%s/add-money/%s", accountID, goalID, id)
	} else {
		endpoint = fmt.Sprintf("/api/v2/account/%s/savings-goals/%s/withdraw-money/%s", accountID, goalID, id)
		amount = -amount
	}
	var body Transfer
	body.Amount.Currency = currency
	body.Amount.MinorUnits = amount
	reqBody, _ := json.Marshal(&body)
	client := &http.Client{}
	req, err := http.NewRequest("PUT", baseUrl+endpoint, bytes.NewBuffer(reqBody))
	req.Header.Add("Authorization", bearer)
	req.Header.Set("Content-Type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}
	if response.StatusCode != 200 {
		fmt.Printf("Starling: %s\n", response.Status)
		return false
	}
	return true
}

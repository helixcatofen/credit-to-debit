package truelayer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

type Cards struct {
	CardList []struct {
		AccountID string `json:"account_id"`
		Name      string `json:"display_name"`
	} `json:"results"`
}

type Transaction struct{
	Description string  `json:"description"`
	Amount      float32 `json:"amount"`
	Id			string	`json:"transaction_id"`
}

type Transactions struct {
	TransactionList []Transaction `json:"results"`
}

const baseUrl string = "https://api.truelayer.com"

func refreshToken() string {
	clientSecret := os.Getenv("TRUELAYER_SECRET")
	clientID := os.Getenv("TRUELAYER_ID")
	token := os.Getenv("TRUELAYER_TOKEN")
	url := "https://auth.truelayer.com/connect/token"
	method := "POST"
	payload := strings.NewReader(fmt.Sprintf("grant_type=refresh_token&client_id=%s&client_secret=%s&refresh_token=%s", clientID, clientSecret, token))
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}
	if response.StatusCode != 200 {
		fmt.Printf("TrueLayer: %s\n", response.Status)
	}
	responseData, _ := ioutil.ReadAll(response.Body)
	var tokenResponse TokenResponse
	json.Unmarshal(responseData, &tokenResponse)

	return tokenResponse.AccessToken
}

func listCards(token string) string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", baseUrl+"/data/v1/cards", nil)
	req.Header.Add("Authorization", "Bearer "+token)

	response, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}
	if response.StatusCode != 200 {
		fmt.Printf("TrueLayer: %s\n", response.Status)
	}

	responseData, _ := ioutil.ReadAll(response.Body)
	var cards Cards
	json.Unmarshal(responseData, &cards)
	return cards.CardList[0].AccountID
}

func GetTransactions() Transactions {
	token := refreshToken()
	cardID := listCards(token)

	var hour time.Duration = -3600 * time.Second
  
	from := time.Now().Add(hour * 168).Format(time.RFC3339)[:19]
	to := time.Now().Add(hour * 1).Format(time.RFC3339)[:19]

	fmt.Println(from)
	fmt.Println(to)
	client := &http.Client{}
	endpoint := fmt.Sprintf("%s/data/v1/cards/%s/transactions?from=%s&to=%s", baseUrl, cardID, from, to)
	req, err := http.NewRequest("GET", endpoint, nil)
	req.Header.Add("Authorization", "Bearer "+token)

	response, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}
	if response.StatusCode != 200 {
		fmt.Printf("TrueLayer: %s\n", response.Status)
	}
	responseData, _ := ioutil.ReadAll(response.Body)
	var transactions Transactions
	json.Unmarshal(responseData, &transactions)
	return transactions
}

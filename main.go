package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

const (
	// URLs
	rawJsonURL      = "https://store-site-backend-static-ipv4.ak.epicgames.com/freeGamesPromotions"
	epicStoreURL    = "https://www.epicgames.com/store/en-US/p/"
	slackWebhookURL = "" // Paste your Slack Webhook URL here
	periodSize      = 24
	periodUnit      = time.Hour
)

type epicJson struct {
	Data data
}

type data struct {
	Catalog Catalog
}

type Catalog struct {
	SearchStore searchStore
}

type searchStore struct {
	Elements []element
}

type element struct {
	Title       string
	Description string
	UrlSlug     string
	Promotions  promotions
}

type promotions struct {
	CurrentPromotionalOffers []promotionalOffers `json:"promotionalOffers"`
}

type promotionalOffers struct {
	PromotionalOffers []offer
}

type offer struct {
	StartDate       string
	EndDate         string
	DiscountSetting discountSetting
}

type discountSetting struct {
	DiscountType       string
	DiscountPercentage int
}

type freeGame struct {
	Title       string
	Description string
	Url         string
	//Picture     string
}

type blocks struct {
	Blocks []section `json:"blocks"`
}

type section struct {
	Type string `json:"type"`
	Text text   `json:"text"`
}

type text struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func main() {
	ticker := time.NewTicker(periodSize * periodUnit)
	for {
		select {
		case <-ticker.C:
			sendFreeGamesToSlack()
		}
	}
}

func sendFreeGamesToSlack() {
	// Get the JSON from EPIC's website
	resp, err := http.Get(rawJsonURL)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	// Unmarshal JSON to struct
	decoder := json.NewDecoder(resp.Body)
	var rawJson epicJson
	decoder.Decode(&rawJson)

	// Identify the free games of the current week and save their relevant data
	var currentFreeGames []freeGame
	for _, element := range rawJson.Data.Catalog.SearchStore.Elements {
		for _, currentpromotionaloffer := range element.Promotions.CurrentPromotionalOffers {
			for _, promotionaloffer := range currentpromotionaloffer.PromotionalOffers {
				if promotionaloffer.DiscountSetting.DiscountType == "PERCENTAGE" && promotionaloffer.DiscountSetting.DiscountPercentage == 0 {
					fG := freeGame{Title: element.Title,
						Description: element.Description,
						Url:         epicStoreURL + element.UrlSlug}
					currentFreeGames = append(currentFreeGames, fG)
				}
			}
		}
	}

	// Pack the data formated for Slack
	var slackBlocks blocks
	slackBlocks.Blocks = append(slackBlocks.Blocks, section{
		Type: "section",
		Text: text{
			Type: "plain_text",
			Text: "The frees game in EPIC Store this week are:"}})
	for _, game := range currentFreeGames {
		slackBlocks.Blocks = append(slackBlocks.Blocks, section{
			Type: "header",
			Text: text{
				Type: "plain_text",
				Text: game.Title}})
		slackBlocks.Blocks = append(slackBlocks.Blocks, section{
			Type: "section",
			Text: text{
				Type: "plain_text",
				Text: game.Description}})
		slackBlocks.Blocks = append(slackBlocks.Blocks, section{
			Type: "section",
			Text: text{
				Type: "mrkdwn",
				Text: ":video_game: <" + game.Url + "|Get the game>"}})
	}

	// Marshal the data into JSON and send it to Slack
	postBody, _ := json.Marshal(slackBlocks)
	responseBody := bytes.NewBuffer(postBody)
	_, err = http.Post(slackWebhookURL, "application/json", responseBody)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()
}

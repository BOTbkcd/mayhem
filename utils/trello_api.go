package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/BOTbkcd/mayhem/entities"
)

type Board struct {
	Id   string
	Name string
	Url  string
}

type List struct {
	Id   string
	Name string
}

type Card struct {
	Id     string
	Name   string
	Desc   string
	Due    string
	Closed string
}

var (
	baseURL     = "https://api.trello.com/"
	new_board   = "/1/boards/"
	new_list    = "/1/boards/{id}/lists"
	board_lists = "/1/boards/{id}/lists"
	list_cards  = "/1/lists/{id}/cards"
	clear_list  = "1//lists/{id}/archiveAllCards"
)

var (
	apiKey   = ""
	apiToken = ""
)

func GenerateNewBoard(apiKey string, apiToken string) Board {
	stacks := entities.FetchStackTitles()

	params := map[string]string{
		"key":   apiKey,
		"token": apiToken,
		"name":  "Mayhem",
	}

	var board Board
	json.Unmarshal(getRawResponse(new_board, "POST", queryParams(params)), &board)

	for _, stack := range stacks {
		params["name"] = stack
		getRawResponse(strings.Replace(new_list, "{id}", board.Id, 1), "POST", queryParams(params))
	}

	return board
}

func SyncBoardData(boardUrl string, apiKey string, apiToken string) {
	params := map[string]string{
		"key":   apiKey,
		"token": apiToken,
	}

	idStart := strings.Index(boardUrl, "/b/") + 3
	idEnd := idStart + 8
	boardId := boardUrl[idStart:idEnd]

	var lists []List
	json.Unmarshal(getRawResponse(strings.Replace(new_list, "{id}", boardId, 1), "GET", queryParams(params)), &lists)

	for _, list := range lists {
		var cards []Card
		json.Unmarshal(getRawResponse(strings.Replace(list_cards, "{id}", list.Id, 1), "GET", queryParams(params)), &cards)

		stack, _ := entities.FetchStackByTitle(list.Name)

		for _, card := range cards {
			newTask := entities.Task{
				StackID:     stack.ID,
				Title:       card.Name,
				Description: card.Desc,
			}

			if card.Due != "" {
				due, _ := time.Parse(time.RFC3339, card.Due)
				newTask.Deadline = due.Local()
			}

			newTask.Save()
			stack.PendingTaskCount++
		}
		stack.Save()

		getRawResponse(strings.Replace(clear_list, "{id}", list.Id, 1), "POST", queryParams(params))
	}

}

func getRawResponse(resource string, method string, params url.Values) []byte {
	u, _ := url.ParseRequestURI(baseURL)
	u.Path = resource
	u.RawQuery = params.Encode()
	urlStr := fmt.Sprintf("%v", u)

	client := &http.Client{}
	req, err := http.NewRequest(method, urlStr, nil)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	return bodyBytes
}

func queryParams(pairs map[string]string) url.Values {
	params := url.Values{}

	for key, val := range pairs {
		params.Add(key, val)
	}

	return params
}

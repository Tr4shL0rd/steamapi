package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func removeDupes(lst [][]string) [][]string {
	uniqueList := make([][]string, 0)
	for _, elem := range lst {
		exists := false
		for _, e := range uniqueList {
			if elem[0] == e[0] && elem[1] == e[1] {
				exists = true
				break
			}
		}
		if !exists {
			uniqueList = append(uniqueList, elem)
		}
	}
	return uniqueList
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	apiKey := os.Getenv("API_KEY")

	url := strings.Join([]string{"https://api.steampowered.com/ISteamApps/GetAppList/v2/?key=", apiKey}, "")
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err)
	}

	header := []string{"appid", "name"}
	appList := make([][]string, 0)
	if appListData, ok := data["applist"]; ok {
		if apps, ok := appListData.(map[string]interface{})["apps"]; ok {
			for _, app := range apps.([]interface{}) {
				appMap := app.(map[string]interface{})
				appIDStr := fmt.Sprint(appMap["appid"])
				// fixes scientific notation
				appID, err := strconv.ParseFloat(appIDStr, 64) // str -> float64
				if err != nil {
					panic(err) // handle error
				}
				if appID < 0 {
					appID = 0
				}
				name := fmt.Sprint(appMap["name"])
				if name == "" {
					name = "None"
				}
				appIDNewStr := fmt.Sprintf("%f", appID)
				appList = append(appList, []string{strings.Split(appIDNewStr, ".")[0], name}) // float64 -> str
			}
		}
	}

	appList = removeDupes(appList)
	file, err := os.Create("apps.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	w := csv.NewWriter(file)
	w.Write(header)
	for _, row := range appList {
		w.Write(row)
	}
	w.Flush()
}

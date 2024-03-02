package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// WordCounter is a structure that contains the count of words in the text and a list of recurring words.
type WordCounter struct {
	WordMap           map[string]int
	RepeatedWordsList []string
}

// TextAnalysis analyzes the given text and returns the WordCounter structure.
func TextAnalysis(text string) WordCounter {
	words := strings.Fields(text)
	wordCounter := WordCounter{
		WordMap:           make(map[string]int),
		RepeatedWordsList: []string{},
	}

	for _, word := range words {
		wordCounter.WordMap[word]++
		if wordCounter.WordMap[word] == 2 {
			wordCounter.RepeatedWordsList = append(wordCounter.RepeatedWordsList, word)
		}
	}

	return wordCounter
}

func analyzeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		http.Error(w, "Only POST or GET requests are supported", http.StatusMethodNotAllowed)
		return
	}

	var text string
	if r.Method == http.MethodPost {
		decoder := json.NewDecoder(r.Body)
		var requestBody struct {
			Text string `json:"text"`
		}
		if err := decoder.Decode(&requestBody); err != nil {
			http.Error(w, "Failed to decode request body", http.StatusBadRequest)
			return
		}
		text = requestBody.Text
	} else if r.Method == http.MethodGet {
		text = r.URL.Query().Get("text")
	}

	analysis := TextAnalysis(text)

	response := struct {
		RepeatedWords map[string]int `json:"repeated_words"`
	}{
		RepeatedWords: make(map[string]int),
	}

	for _, word := range analysis.RepeatedWordsList {
		response.RepeatedWords[word] = analysis.WordMap[word]
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/analyze", analyzeHandler)
	fmt.Println("Starting server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

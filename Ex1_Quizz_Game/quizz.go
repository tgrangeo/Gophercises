package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type ProblemItem struct {
	Question string
	Answer   int
}

func main() {

	//Open the CSV file
	file, err := os.Open("problem.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	//Read the file
	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	//get problems
	var ProblemList []ProblemItem
	for _, record := range records {
		ans, _ := strconv.Atoi(record[1])
		ProblemList = append(ProblemList, ProblemItem{
			Question: record[0],
			Answer:   ans,
		})
	}

	//game begin
	var correct_answer int
	reader := bufio.NewReader(os.Stdin)
	for _, item := range ProblemList {
		fmt.Println(item.Question)
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text) //why atoi in go is different from atoi in c ????
		userOutput, _ := strconv.Atoi(text)
		if userOutput == item.Answer {
			correct_answer++
		}
	}
	fmt.Println("correct answer: ", correct_answer, "/", len(ProblemList))
}

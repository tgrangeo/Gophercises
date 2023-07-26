package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type ProblemItem struct {
	Question string
	Answer   int
}

func shuffle(items []ProblemItem) {
	rand.Seed(time.Now().UnixNano())
	n := len(items)
	for i := n - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		items[i], items[j] = items[j], items[i]
	}
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

	//read
	reader := bufio.NewReader(os.Stdin)

	// Define flags
	var flagVar int
	flag.IntVar(&flagVar, "time", 30, "Duration for the timer")

	var shuffleFlag bool
	flag.BoolVar(&shuffleFlag, "shuffle", false, "Set to true to shuffle questions")
	flag.Parse()
	if shuffleFlag {
		shuffle(ProblemList)
	}

	//timer begin
	fmt.Print("Do you want to begin the quiz ? y/n\n")
	text, _ := reader.ReadString('\n')
	if strings.TrimSpace(text) != "y" {
		fmt.Println("See you next time :)")
		return
	}
	var timer time.Duration = time.Duration(flagVar)
	go func() {
		<-time.After(timer * time.Second)
		fmt.Println(flagVar, "seconds have passed. Program will now exit.")
		os.Exit(1)
	}()

	//game begin
	var correct_answer int
	for _, item := range ProblemList {
		fmt.Println(item.Question)
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		userOutput, _ := strconv.Atoi(text)
		if userOutput == item.Answer {
			correct_answer++
		}
	}
	fmt.Println("correct answer: ", correct_answer, "/", len(ProblemList))
}

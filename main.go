package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

func problemPuller(fileName string) ([]problem, error) {
	if fObj, err := os.Open(fileName); err == nil {

		csvR := csv.NewReader(fObj)

		if cLines, err := csvR.ReadAll(); err == nil {
			return parseProblem(cLines), nil
		} else {
			return nil, fmt.Errorf("error in reading data in csv")
		}

	} else {
		return nil, fmt.Errorf("error in opening file")
	}
}

func main() {

	fName := flag.String("f", "quiz.csv", "path of the csv file")

	timer := flag.Int("t", 30, "timer for the quiz")
	flag.Parse()
	problems, err := problemPuller(*fName)

	if err != nil {
		exit(fmt.Sprintln("something went wrong:%s", err.Error()))
	}

	correctAns := 0

	tObj := time.NewTimer(time.Duration(*timer) * time.Second)

	ansC := make(chan string)

problemLoop:

	for i, p := range problems {
		var answer string
		fmt.Printf("Problems %d: %s=", i+1, p.q)

		go func() {
			fmt.Scanf("%s", &answer)
			ansC <- answer
		}()
		select {
		case <-tObj.C:
			fmt.Println()
			break problemLoop

		case iAns := <-ansC:
			if iAns == p.a {
				correctAns++
			}
			if i == len(problems)-1 {
				close(ansC)
			}
		}

	}

	fmt.Printf("your result is %d out of %d\n", correctAns, len(problems))
	fmt.Printf("press enter to exit")
	<-ansC

}
func parseProblem(lines [][]string) []problem {
	r := make([]problem, len(lines))
	for i := 0; i < len(lines); i++ {
		r[i] = problem{q: lines[i][0], a: lines[i][1]}
	}
	return r
}

type problem struct {
	q string
	a string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

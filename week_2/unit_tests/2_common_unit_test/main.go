package main

import (
	"fmt"

	"github.com/mbakhodurov/examples/week_2/unit_tests/2_common_unit_test/credit_score"
)

func main() {
	client := credit_score.Client{
		Gender:        "male",
		Age:           30,
		Profession:    "engineer",
		Experience:    7,
		AverageSalary: 60000,
	}
	credit_score := credit_score.CalculateCreditScore(client)
	fmt.Printf("Credit score: %d\n", credit_score)
}

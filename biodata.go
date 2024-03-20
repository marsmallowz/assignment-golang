package main

import (
	"fmt"
	"os"
	"strconv"
)

type Friend struct {
	Name                 string
	Address              string
	Job                  string
	ReasonToChooseGolang string
}

var FriendList = []Friend{
	{Name: "AlMan", Address: "Jl. Sample 123", Job: "Programmer", ReasonToChooseGolang: "Want to learn another programming language"},
	{Name: "Jane", Address: "Jl. Sample 456", Job: "Designer", ReasonToChooseGolang: "Want to learn programming"},
	{Name: "Doe", Address: "Jl. Sample 789", Job: "Architect", ReasonToChooseGolang: "Forced by parents"},
}

func ShowData(absen int) {
	if absen <= 0 || absen > len(FriendList) {
		fmt.Println("absen invalid.")
		return
	}

	friend := FriendList[absen-1]
	fmt.Printf("Name: %s\n", friend.Name)
	fmt.Printf("Address: %s\n", friend.Address)
	fmt.Printf("Job: %s\n", friend.Job)
	fmt.Printf("Reason To Choose Golang: %s\n", friend.ReasonToChooseGolang)

}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide the absen number as an argument.")
		return
	}

	absen, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("The absen number must be a number.")
		return
	}

	ShowData(absen)
}

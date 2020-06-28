package main

import (
	"fmt"
	"github.com/DABronskikh/go-lesson-7/pkg/transactions"
	"math/rand"
)

func main() {
	const transitionsPerUser = 10_000
	usersId := []string{"001", "002", "003", "004", "005"}
	users := len(usersId)
	const transactionAmount = 1_00
	mccArray := []string{"5411", "5812"}

	transitions := make([]transactions.Transaction, users*transitionsPerUser)
	for idx := range transitions {
		transitions[idx].Amount = transactionAmount

		MCC := mccArray[0]
		if idx%2 != 0 {
			MCC = mccArray[1]
		}
		transitions[idx].MCC = MCC

		userId := usersId[0]
		if rand.Intn(99)%2 != 0 {
			userId = usersId[1]
		}
		transitions[idx].IdUser = userId
	}

	result := transactions.SumByCategory(usersId[0], transitions)
	fmt.Println("1) ", result)

	result2 := transactions.MutexSumByCategory(usersId[0], transitions)
	fmt.Println("2) ", result2)

	result3 := transactions.ChanSumByCategory(usersId[0], transitions)
	fmt.Println("3) ", result3)

	result4 := transactions.MutexSumByCategoryV2(usersId[0], transitions)
	fmt.Println("4) ", result4)
}

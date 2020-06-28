package transactions

import (
	"reflect"
	"testing"
)

func makeTransactions() []Transaction {
	const users = 10_000
	const transitionsPerUser = 10_000
	const transactionAmount = 1_00
	mccArray := []string{"5411", "5812"}
	transitions := make([]Transaction, users*transitionsPerUser)

	for idx := range transitions {
		transitions[idx].Amount = transactionAmount

		switch idx % 100 {
		case 0:
			// Например, каждая 100-ая транзакция в банке от нашего юзера в категории такой-то
			transitions[idx].IdUser = "001"
			transitions[idx].MCC = mccArray[0]
		case 20:
			// Например, каждая 120-ая транзакция в банке от нашего юзера в категории такой-то
			transitions[idx].IdUser = "001"
			transitions[idx].MCC = mccArray[1]
		default:
			// Транзакции других юзеров, нужны для "общей" массы
			transitions[idx].IdUser = ""
		}
	}
	return transitions
}

func referenceResult() map[string]int64 {
	transitions := make(map[string]int64, 2)
	transitions["5411"] = 1_000_000_00
	transitions["5812"] = 1_000_000_00
	return transitions
}

func TestChanSumByCategory(t *testing.T) {
	type args struct {
		idUser       string
		transactions [] Transaction
	}
	tests := []struct {
		name string
		args args
		want map[string]int64
	}{
		{
			name: "Положительный тест для ChanSumByCategory",
			args: args{
				idUser:       "001",
				transactions: makeTransactions(),
			},
			want: referenceResult(),
		},
	}
	for _, tt := range tests {
		if got := ChanSumByCategory(tt.args.idUser, tt.args.transactions); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("ChanSumByCategory() = %v, want %v", got, tt.want)
		}
	}
}

func TestMutexSumByCategory(t *testing.T) {
	type args struct {
		idUser       string
		transactions [] Transaction
	}
	tests := []struct {
		name string
		args args
		want map[string]int64
	}{
		{
			name: "Положительный тест для MutexSumByCategory",
			args: args{
				idUser:       "001",
				transactions: makeTransactions(),
			},
			want: referenceResult(),
		},
	}
	for _, tt := range tests {
		if got := MutexSumByCategory(tt.args.idUser, tt.args.transactions); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("MutexSumByCategory() = %v, want %v", got, tt.want)
		}
	}
}

func TestMutexSumByCategoryV2(t *testing.T) {
	type args struct {
		idUser       string
		transactions [] Transaction
	}
	tests := []struct {
		name string
		args args
		want map[string]int64
	}{
		{
			name: "Положительный тест для MutexSumByCategoryV2",
			args: args{
				idUser:       "001",
				transactions: makeTransactions(),
			},
			want: referenceResult(),
		},
	}
	for _, tt := range tests {
		if got := MutexSumByCategoryV2(tt.args.idUser, tt.args.transactions); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("MutexSumByCategoryV2() = %v, want %v", got, tt.want)
		}
	}
}

func TestSumByCategory(t *testing.T) {
	type args struct {
		idUser       string
		transactions [] Transaction
	}
	tests := []struct {
		name string
		args args
		want map[string]int64
	}{
		{
			name: "Положительный тест для SumByCategory",
			args: args{
				idUser:       "001",
				transactions: makeTransactions(),
			},
			want: referenceResult(),
		},
	}
	for _, tt := range tests {
		if got := SumByCategory(tt.args.idUser, tt.args.transactions); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("SumByCategory() = %v, want %v", got, tt.want)
		}
	}
}

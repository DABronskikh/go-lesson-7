package transactions

import (
	"sync"
)

type Transaction struct {
	Amount int64
	MCC    string
	IdUser string
}

func Sum(transitions []Transaction) int64 {
	sum := int64(0)
	for _, v := range transitions {
		sum += v.Amount
	}
	return sum
}

/**
1. Обычная функция, которая принимает на вход слайс транзакций и id владельца - возвращает map с категориями
и тратами по ним (сортировать они ничего не должна)
*/
func SumByCategory(idUser string, transactions [] Transaction) map[string]int64 {
	groupByCategory := map[string][] Transaction{}
	for _, transaction := range transactions {
		if transaction.IdUser == idUser {
			groupByCategory[transaction.MCC] = append(groupByCategory[transaction.MCC], transaction)
		}
	}

	result := map[string]int64{}
	for key, transactions := range groupByCategory {
		result[key] = Sum(transactions)
	}

	return result
}

/**
2. Функция с mutex'ом, который защищает любые операции с map, соответственно, её задача: разделить слайс транзакций
на несколько кусков и в отдельных горутинах посчитать map'ы по кускам, после чего собрать всё в один большой map.
Важно: эта функция внутри себя должна вызывать функцию из п.1
*/
func MutexSumByCategory(idUser string, transactions [] Transaction) map[string]int64 {
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	result := make(map[string]int64)

	partsCount := 10
	partsSize := len(transactions) / partsCount

	for i := 0; i < partsCount; i++ {
		wg.Add(1)
		part := transactions[i*partsSize : (i+1)*partsSize]
		go func() {
			m := SumByCategory(idUser, part)

			mu.Lock()

			for mcc, sumTransactions := range m {
				_, ok := result[mcc]
				if !ok {
					result[mcc] = sumTransactions
				} else {
					result[mcc] += sumTransactions
				}
			}

			mu.Unlock()
			wg.Done()
		}()
	}
	wg.Wait()

	return result
}

/**
3. Функция с каналами, соответственно, её задача: разделить слайс транзакций на несколько кусков и в отдельных горутинах
посчитать map'ы по кускам, после чего собрать всё в один большой map (передавайте рассчитанные куски по каналу).
Важно: эта функция внутри себя должна вызывать функцию из п.1
*/
func ChanSumByCategory(idUser string, transactions [] Transaction) map[string]int64 {
	result := make(map[string]int64)
	ch := make(chan map[string]int64)

	partsCount := 10
	partsSize := len(transactions) / partsCount

	for i := 0; i < partsCount; i++ {
		part := transactions[i*partsSize : (i+1)*partsSize]

		go func(ch chan<- map[string]int64) {
			ch <- SumByCategory(idUser, part)
		}(ch)
	}

	finished := 0
	for m := range ch {

		for mcc, sumTransactions := range m {
			_, ok := result[mcc]
			if !ok {
				result[mcc] = sumTransactions
			} else {
				result[mcc] += sumTransactions
			}
		}

		finished++
		if finished == partsCount {
			break
		}
	}

	return result
}

/**
4. Функция с mutex'ом, который защищает любые операции с map, соответственно, её задача: разделить слайс транзакций на
несколько кусков и в отдельных горутинах посчитать, но теперь горутины напрямую пишут в общий map с результатами.
Важно: эта функция внутри себя НЕ должна вызывать функцию из п.1
*/
func MutexSumByCategoryV2(idUser string, transactions [] Transaction) map[string]int64 {
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	result := make(map[string]int64)

	partsCount := 10
	partsSize := len(transactions) / partsCount

	for i := 0; i < partsCount; i++ {
		wg.Add(1)
		part := transactions[i*partsSize : (i+1)*partsSize]
		go func() {
			for _, t := range part {
				if t.IdUser == idUser {
					mu.Lock()
					_, ok := result[t.MCC]
					if !ok {
						result[t.MCC] = t.Amount
					} else {
						result[t.MCC] += t.Amount
					}
					mu.Unlock()
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()

	return result
}

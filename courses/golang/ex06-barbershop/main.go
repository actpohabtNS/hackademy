package main

import "time"

func main() {
	bs := Barbershop{5, make(chan int), 10, 0, 500}
	customerServedChan := make(chan int)
	bsClosingChan := make(chan bool)

	go barber(&bs, customerServedChan, bsClosingChan)
	for i := 0; i < bs.customersPerDay; i++ {
		time.Sleep(time.Duration(500) * time.Millisecond)
		go customer(i, &bs, 1000, customerServedChan)
	}

	<-bsClosingChan
}

package main

import (
	"fmt"
	"time"
)

type Barbershop struct {
	seats           int
	seatsChan       chan int
	customersPerDay int
	customersServed int
	barberDelay     int
}

func barber(bs *Barbershop, customerServedChan chan int, bsClosingChan chan bool) {
	customerId := 0

	for {
		select {
		case customerId = <-bs.seatsChan:
		default:
			fmt.Println("[Barber]: zzz z z Z Z")
			customerId = <-bs.seatsChan
			fmt.Printf("[Barber]: I was woken by Customer #%d\n", customerId)
		}
		serveCustomer(customerId, bs, customerServedChan)
		bs.customersServed++

		if bs.customersServed == bs.customersPerDay {
			break
		}
	}

	fmt.Printf("[Barber]: closing shop, served %d customers!\n", bs.customersServed)
	bsClosingChan <- true
}

func serveCustomer(customerId int, bs *Barbershop, customerServedChan chan int) {
	time.Sleep(time.Duration(bs.barberDelay) * time.Millisecond)
	customerServedChan <- customerId
	fmt.Printf("[Barber]: customer #%d has been served!\n", customerId)
}

func customer(id int, bs *Barbershop, returnDelay int, customerServedChan chan int) {
	serviced := false

	for !serviced {
		select {
		case bs.seatsChan <- id:
			fmt.Printf("[Customer #%d]: entering barbershop!\n", id)
			servedCustomerId := -1
			for servedCustomerId != id {
				servedCustomerId = <-customerServedChan
			}

			fmt.Printf("[Customer #%d]: my hair has been cut!\n", id)
			serviced = true
		default:
			fmt.Printf("[Customer #%d]: barbershop is full :( Waiting outside!\n", id)
			time.Sleep(time.Duration(returnDelay) * time.Millisecond)
		}
	}
}

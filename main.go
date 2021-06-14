package main

import (
	"container/heap"
	"fmt"
	"math"
	"math/rand"

	"github.com/lmicke/golang-des-example/structs"
)

// System entities:  1. Customer-Queue 2. Barista

// System events: Customer Arrival, Customer Departure

// System states: Number of Customers in Queue(0-n), Barista-Status (busy, idle)

// Random Variables: Customer-InterArrival-Time, Barista-Service-Time

// First Example one Barista

type Barista struct {
	Idle            bool
	meanServingTime float64
}

func (b *Barista) ServiceTime() float64 {
	return math.Max((rand.NormFloat64() + b.meanServingTime), float64(0))
}

type CustomerQueue uint

func (cq CustomerQueue) IsEmpty() bool {
	if cq == 0 {
		return true
	}
	return false
}

type Entities struct {
	cq *CustomerQueue
	b  *Barista
}

func newDeparture(customer *structs.Item, pq *structs.PriorityQueue, e *Entities) {
	departureTime := &structs.Item{
		Arrival: false,
	}
	departureTime.Time = customer.Time + e.b.ServiceTime()
	fmt.Println(departureTime)
	e.b.Idle = false
	heap.Push(pq, departureTime)
}

func newArrival(pq *structs.PriorityQueue, lastEventTime float64) {
	time := lastEventTime + float64(5+rand.Intn(15))
	arrivalTime := &structs.Item{
		Time:    time,
		Arrival: true,
	}
	heap.Push(pq, arrivalTime)

}

func step(pq *structs.PriorityQueue, e *Entities) float64 {
	customer := new(structs.Item)
	customer = heap.Pop(pq).(*structs.Item)
	//fmt.Printf("New Customer: %v\n", *customer)
	fmt.Printf("Arrive is %v\t", customer.Arrival)
	fmt.Printf("Time is %v\t", customer.Time)
	fmt.Printf("Queue Length is %v\n", *e.cq)
	if customer.Arrival {
		if e.cq.IsEmpty() && e.b.Idle {
			//fmt.Println("customer arrived, queue empty, barista idle")
			newDeparture(customer, pq, e)
		}
		if !e.b.Idle {
			//fmt.Println("customer arrived, barista busy adding to queue")
			*e.cq = *e.cq + 1
		}
	} else {
		if e.cq.IsEmpty() {
			//fmt.Println("customer departed, queue empty, setting barista to idle")
			e.b.Idle = true
		} else {
			//fmt.Println("Customer departed, there is a queue . getting next one in queue")
			*e.cq = *e.cq - 1
			newDeparture(customer, pq, e)

		}
	}
	newArrival(pq, customer.Time)
	return customer.Time
}

func main() {

	pq := make(structs.PriorityQueue, 5)
	for i := 0; i < 5; i++ {
		pq[i] = &structs.Item{
			Arrival: true,
			Time:    float64(i + rand.Intn(2)),
		}
	}

	heap.Init(&pq)

	b := &Barista{
		Idle:            true,
		meanServingTime: 2,
	}

	var cq CustomerQueue = 0

	e := &Entities{
		b:  b,
		cq: &cq,
	}
	time := 0
	for time < 480 {
		time = int(step(&pq, e))
	}
}

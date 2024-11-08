package models

import (
	"fmt"
	"math/rand"
	"time"
)

type Vehicle struct {
	ID int
}

func (p *ParkingLot) Arrive(vehicle *Vehicle) {
	select {
	case <-p.availableSpots:

		p.mutex.Lock()
		spotIndex := p.findNextAvailableSpot()
		if spotIndex != -1 {
			p.occupiedSpaces[spotIndex] = true
			p.vehicleIDs[spotIndex] = vehicle.ID

		}
		p.mutex.Unlock()

		time.Sleep(time.Duration(3+rand.Intn(3)) * time.Second)
		p.Depart(vehicle)

	default:
		p.queue <- vehicle
	}
}

func (p *ParkingLot) Depart(vehicle *Vehicle) {
	p.mutex.Lock()
	spotFound := false
	spotIndex := -1

	for i := 0; i < p.capacity; i++ {
		if p.vehicleIDs[i] == vehicle.ID {
			spotFound = true
			spotIndex = i
			p.occupiedSpaces[i] = false
			p.vehicleIDs[i] = 0
			break
		}
	}
	p.mutex.Unlock()

	if spotFound {
		fmt.Printf("Vehículo %d saliendo del espacio %d.\n", vehicle.ID, spotIndex+1)
		p.availableSpots <- struct{}{}

		select {
		case nextVehicle := <-p.queue:
			go p.Arrive(nextVehicle)
		default:

		}
	}
}


func (p *ParkingLot) GetVehicleID(index int) int {
	p.mutex.RLock()
   defer p.mutex.RUnlock()
   return p.vehicleIDs[index]
}
package models

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type ParkingLot struct {
    capacity         int
    currentVehicles  int
    mutex       sync.RWMutex
    queue            chan *Vehicle
    availableSpots   chan struct{}
    occupiedSpaces   []bool
    vehicleIDs []int
    nextSpotIndex int
}

func NewParkingLot(capacity int) *ParkingLot {
    parkingLot := &ParkingLot{
        capacity:       capacity,
        queue:          make(chan *Vehicle, capacity),
        availableSpots: make(chan struct{}, capacity),
        occupiedSpaces: make([]bool, capacity),
        vehicleIDs: make([]int, capacity),
        nextSpotIndex: 0,
    }
    for i := 0; i < capacity; i++ {
        parkingLot.availableSpots <- struct{}{}
    }
    return parkingLot
}

func (p *ParkingLot) findNextAvailableSpot() int {
    for i := 0; i < p.capacity; i++ {
        index := (p.nextSpotIndex + i) % p.capacity
        if !p.occupiedSpaces[index] {
            p.nextSpotIndex = (index + 1) % p.capacity
            return index
        }
    }
    return -1
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

func (p *ParkingLot) GetOccupiedSpaces() ([]bool, []int) {
    p.mutex.RLock()
       defer p.mutex.RUnlock()
       occupiedSpacesCopy := make([]bool, len(p.occupiedSpaces))
       vehicleIDsCopy := make([]int, len(p.vehicleIDs))
       copy(occupiedSpacesCopy, p.occupiedSpaces)
       copy(vehicleIDsCopy, p.vehicleIDs)
       return occupiedSpacesCopy, vehicleIDsCopy
}

func (p *ParkingLot) Capacity() int {
    return p.capacity
}

func (p *ParkingLot) GetVehicleID(index int) int {
 	p.mutex.RLock()
    defer p.mutex.RUnlock()
    return p.vehicleIDs[index]
}
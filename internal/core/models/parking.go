package models

import (
	"sync"
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

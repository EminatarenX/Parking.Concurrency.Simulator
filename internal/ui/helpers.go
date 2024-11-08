package ui

import (
	"fmt"
	"math/rand"
	"simulator/internal/core/models"
	"time"

	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
)
func UpdateParkingDisplay(texts []*canvas.Text, carImages []*canvas.Image, parkingLot *models.ParkingLot) {
    occupiedSpaces, vehicleIDs := parkingLot.GetOccupiedSpaces()

    for i := range occupiedSpaces {
        if texts[i] == nil || carImages[i] == nil {
            continue
        }

        if occupiedSpaces[i] {
            text := fmt.Sprintf("\t\t\t#%d  ",vehicleIDs[i])
            texts[i].Color = theme.Color(theme.ColorNameForeground) 
            texts[i].Text = text
            carImages[i].Show()
            
        } else {
            texts[i].Text = fmt.Sprintf("")
            texts[i].Color = theme.Color(theme.ColorNameBackground)
            carImages[i].Hide()
        }

        texts[i].Refresh()
        carImages[i].Refresh()
    }
}

func UpdateScreenAndTiker(texts []*canvas.Text, carImages []*canvas.Image, parkingLot *models.ParkingLot){
 ticker := time.NewTicker(100 * time.Millisecond)
        defer ticker.Stop()

        for range ticker.C {
            UpdateParkingDisplay(texts, carImages, parkingLot)
        }
}

func GenerateCars(totalCars int, duration float64, parkingLot *models.ParkingLot){
	for i := 1; i <= totalCars; i++ {
            time.Sleep(time.Duration(rand.ExpFloat64() * duration) * time.Millisecond)
            vehicle := &models.Vehicle{ID: i}
            go parkingLot.Arrive(vehicle)
        }
}
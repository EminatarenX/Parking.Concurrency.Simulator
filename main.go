package main

import (
    "simulator/internal/core/models"
    "simulator/internal/ui"
    "fyne.io/fyne/v2/app"
)

func main() {
    myApp := app.New()
    parkingLot := models.NewParkingLot(20)
    
  	ui.CreateWindow(myApp, parkingLot, 250, 100)

    myApp.Run()
}

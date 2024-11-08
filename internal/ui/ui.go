package ui

import (
	"fmt"
	"math/rand"
	"simulator/internal/core/models"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"

)

func CreateWindow(app fyne.App, parkingLot *models.ParkingLot, duration float64, totalCars int) fyne.Window {
    myWindow := app.NewWindow("Simulación de Estacionamiento")

    background := canvas.NewImageFromFile("internal/ui/assets/background.png")
    background.FillMode = canvas.ImageFillStretch

    containers := make([]*fyne.Container, parkingLot.Capacity())
    texts := make([]*canvas.Text, parkingLot.Capacity())
    carImages := make([]*canvas.Image, parkingLot.Capacity())

    carColors := []string{
    	"internal/ui/assets/car.png",
	    "internal/ui/assets/blue_car.png",
	    "internal/ui/assets/grey_car.png",
	    "internal/ui/assets/yellow_car.png",
    }

    grid := container.New(layout.NewGridLayoutWithColumns(2))


    for i := 0; i < parkingLot.Capacity(); i++ {
    	randomIndex := rand.Intn(len(carColors))
     	randomCar := carColors[randomIndex]

        texts[i] = canvas.NewText(fmt.Sprintf("Espacio %d: Libre", i+1), theme.Color(theme.ColorNameForeground))
        texts[i].TextSize = 15
        texts[i].Alignment = fyne.TextAlignCenter

        carImages[i] = canvas.NewImageFromFile(randomCar)
        carImages[i].Resize(fyne.NewSize(75, 30))
        carImages[i].Hide()

        spotContainer := container.NewWithoutLayout(
        	texts[i],
         	carImages[i],
        )

        texts[i].Move(fyne.NewPos(100,40))
        carImages[i].Move(fyne.NewPos(70, 40))


        containers[i] = container.New(
            layout.NewPaddedLayout(),
            spotContainer,
        )


        grid.Add(containers[i])
    }

    content := container.NewStack(
	    background,
	    container.NewPadded(grid),
    )


    scrollContainer := container.NewScroll(content)

    myWindow.SetContent(scrollContainer)
    myWindow.Resize(fyne.NewSize(500, 900))
    myWindow.CenterOnScreen()
    myWindow.Show()


    go UpdateScreenAndTiker(texts, carImages, parkingLot)

    go GenerateCars(totalCars, duration, parkingLot)

    return myWindow
}

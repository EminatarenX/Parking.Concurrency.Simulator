package ui

import (
	"fmt"
	"math/rand"
	"simulator/internal/core/models"
	"time"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type CustomLabel struct {
    widget.Label
}



func NewCustomLabel(text string) *CustomLabel {
    label := &CustomLabel{}
    label.ExtendBaseWidget(label)
    label.SetText(text)
    label.Alignment = fyne.TextAlignCenter
    label.TextStyle = fyne.TextStyle{Bold: true}


    text1 := canvas.NewText(text, theme.ForegroundColor())
    text1.TextSize = 25

    return label
}


func (c *CustomLabel) MinSize() fyne.Size {

    return fyne.NewSize(250, 20)
}

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

        texts[i] = canvas.NewText(fmt.Sprintf("Espacio %d: Libre", i+1), theme.ForegroundColor())
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


    go updateScreenAndTiker(texts, carImages, parkingLot)

    go generateCars(totalCars, duration, parkingLot)

    return myWindow
}

func updateParkingDisplay(texts []*canvas.Text, carImages []*canvas.Image, parkingLot *models.ParkingLot) {
    occupiedSpaces, vehicleIDs := parkingLot.GetOccupiedSpaces()

    for i := range occupiedSpaces {
        if texts[i] == nil || carImages[i] == nil {
            continue
        }

        if occupiedSpaces[i] {
            text := fmt.Sprintf("\t\t\t#%d  ",vehicleIDs[i])
            texts[i].Color = theme.ForegroundColor()
            texts[i].Text = text
            carImages[i].Show()
        } else {
            texts[i].Text = fmt.Sprintf("")
            texts[i].Color = theme.SuccessColor()
            carImages[i].Hide()
        }

        texts[i].Refresh()
        carImages[i].Refresh()
    }
}

func updateScreenAndTiker(texts []*canvas.Text, carImages []*canvas.Image, parkingLot *models.ParkingLot){
 ticker := time.NewTicker(100 * time.Millisecond)
        defer ticker.Stop()

        for range ticker.C {
            updateParkingDisplay(texts, carImages, parkingLot)
        }
}

func generateCars(totalCars int, duration float64, parkingLot *models.ParkingLot){
	for i := 1; i <= totalCars; i++ {
            time.Sleep(time.Duration(rand.ExpFloat64() * duration) * time.Millisecond)
            vehicle := &models.Vehicle{ID: i}
            go parkingLot.Arrive(vehicle)
        }
}
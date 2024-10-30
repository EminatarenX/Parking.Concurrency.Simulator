package ui

import (
    "fmt"
    "math/rand"
    "simulator/internal/core/models"
    "time"
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
    "fyne.io/fyne/v2/layout"
    "fyne.io/fyne/v2/canvas"
    "fyne.io/fyne/v2/theme"
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
    

    containers := make([]*fyne.Container, parkingLot.Capacity())
    texts := make([]*canvas.Text, parkingLot.Capacity())
    

    grid := container.New(layout.NewGridLayoutWithColumns(2))
    

    for i := 0; i < parkingLot.Capacity(); i++ {
        // Crear el texto con canvas
        texts[i] = canvas.NewText(fmt.Sprintf("Espacio %d: Libre", i+1), theme.ForegroundColor())
        texts[i].TextSize = 35
        texts[i].Alignment = fyne.TextAlignCenter
        

        textContainer := container.NewPadded(
            container.NewCenter(texts[i]),
        )
        

        containers[i] = container.New(
            layout.NewPaddedLayout(),
            textContainer,
        )
        
   
        grid.Add(containers[i])
    }
    
 
    scrollContainer := container.NewScroll(grid)
    
    myWindow.SetContent(scrollContainer)
    myWindow.Resize(fyne.NewSize(1280, 900))
    myWindow.CenterOnScreen()
    myWindow.Show()
    

    go func() {
        ticker := time.NewTicker(100 * time.Millisecond)
        defer ticker.Stop()
        
        for range ticker.C {
            updateParkingDisplay(texts, parkingLot)
        }
    }()
    

    go func() {
        for i := 1; i <= totalCars; i++ {
            time.Sleep(time.Duration(rand.ExpFloat64() * duration) * time.Millisecond)
            vehicle := &models.Vehicle{ID: i}
            go parkingLot.Arrive(vehicle)
        }
    }()
    
    return myWindow
}

func updateParkingDisplay(texts []*canvas.Text, parkingLot *models.ParkingLot) {
    occupiedSpaces, vehicleIDs := parkingLot.GetOccupiedSpaces()
    
    for i := range occupiedSpaces {
        if texts[i] == nil {
            continue
        }
        
        text := fmt.Sprintf("Espacio %d: ", i+1)
        if occupiedSpaces[i] {
        text = fmt.Sprintf("#%d Ocupado ", i+1)
            text += fmt.Sprintf("por 🚘 #%d", vehicleIDs[i])
            texts[i].Color = theme.ForegroundColor()
        } else {
            text += "Libre"
            texts[i].Color = theme.SuccessColor()
        }
        
        texts[i].Text = text
        texts[i].Refresh()
    }
}
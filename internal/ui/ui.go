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

// CustomLabel es un widget personalizado que extiende el Label básico
type CustomLabel struct {
    widget.Label
}

// NewCustomLabel crea un nuevo label personalizado con texto grande
func NewCustomLabel(text string) *CustomLabel {
    label := &CustomLabel{}
    label.ExtendBaseWidget(label)
    label.SetText(text)
    label.Alignment = fyne.TextAlignCenter
    label.TextStyle = fyne.TextStyle{Bold: true}
    
    // Crear un texto con canvas para hacerlo más grande
    text1 := canvas.NewText(text, theme.ForegroundColor())
    text1.TextSize = 25
    
    return label
}

// MinSize sobreescribe el tamaño mínimo del label
func (c *CustomLabel) MinSize() fyne.Size {
    // Aumentar el tamaño mínimo para acomodar el texto más grande
    return fyne.NewSize(250, 20)
}

func CreateWindow(app fyne.App, parkingLot *models.ParkingLot, duration float64, totalCars int) fyne.Window {
    myWindow := app.NewWindow("Simulación de Estacionamiento")
    
    // Crear contenedores para los textos grandes
    containers := make([]*fyne.Container, parkingLot.Capacity())
    texts := make([]*canvas.Text, parkingLot.Capacity())
    
    // Crear el grid principal
    grid := container.New(layout.NewGridLayoutWithColumns(2))
    
    // Crear los textos y contenedores
    for i := 0; i < parkingLot.Capacity(); i++ {
        // Crear el texto con canvas
        texts[i] = canvas.NewText(fmt.Sprintf("Espacio %d: Libre", i+1), theme.ForegroundColor())
        texts[i].TextSize = 35
        texts[i].Alignment = fyne.TextAlignCenter
        
        // Crear un contenedor con padding usando container.NewPadded
        textContainer := container.NewPadded(
            container.NewCenter(texts[i]),
        )
        
        // Crear un contenedor exterior con más espacio
        containers[i] = container.New(
            layout.NewPaddedLayout(),
            textContainer,
        )
        
        // Añadir el contenedor al grid
        grid.Add(containers[i])
    }
    
    // Envolver el grid en un scroll container
    scrollContainer := container.NewScroll(grid)
    
    myWindow.SetContent(scrollContainer)
    myWindow.Resize(fyne.NewSize(1280, 900))
    myWindow.CenterOnScreen()
    myWindow.Show()
    
    // Actualización periódica de la interfaz
    go func() {
        ticker := time.NewTicker(100 * time.Millisecond)
        defer ticker.Stop()
        
        for range ticker.C {
            updateParkingDisplay(texts, parkingLot)
        }
    }()
    
    // Generación de vehículos
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
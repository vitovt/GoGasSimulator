package main

import (
    "fmt"
    "image/color"
    "math"
    "math/rand"
    "time"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/canvas"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
)

const (
    moleculeCount = 100
    windowWidth   = 800
    windowHeight  = 600
    moleculeSize  = 8
    minSpeed      = 1.0
    maxSpeed      = 5.0
)

type Molecule struct {
    circle *canvas.Circle
    posX   float64
    posY   float64
    velX   float64
    velY   float64
}

func main() {
    // Seed the random number generator
    rand.Seed(time.Now().UnixNano())

    // Create a new application
    myApp := app.New()
    myWindow := myApp.NewWindow("Ideal Gas Simulation")

    // Create a container without layout to control absolute positioning
    moleculeContainer := container.NewWithoutLayout()

    // Initialize molecules
    molecules := make([]*Molecule, moleculeCount)
    for i := 0; i < moleculeCount; i++ {
        // Random initial position avoiding overlap
        posX := rand.Float64()*(windowWidth-moleculeSize*2) + moleculeSize
        posY := rand.Float64()*(windowHeight-moleculeSize*2) + moleculeSize

        // Random initial velocity
        angle := rand.Float64() * 2 * math.Pi
        speed := rand.Float64()*(maxSpeed-minSpeed) + minSpeed
        velX := speed * math.Cos(angle)
        velY := speed * math.Sin(angle)

        // Create a circle for the molecule
        circle := canvas.NewCircle(randomColor())
        circle.Resize(fyne.NewSize(moleculeSize, moleculeSize))
        circle.Move(fyne.NewPos(float32(posX), float32(posY)))

        // Add to container
        moleculeContainer.Add(circle)

        // Add to molecule slice
        molecules[i] = &Molecule{
            circle: circle,
            posX:   posX,
            posY:   posY,
            velX:   velX,
            velY:   velY,
        }
    }

    // Temperature slider
    temperatureSlider := widget.NewSlider(1, 10)
    temperatureSlider.Value = 5 // Starting temperature
    temperatureSlider.Step = 0.1
    temperatureLabel := widget.NewLabel("Temperature: 5.0")

    // Update temperature label when slider changes
    temperatureSlider.OnChanged = func(value float64) {
        temperatureLabel.SetText("Temperature: " + fmt.Sprintf("%.1f", value))
    }

    controls := container.NewVBox(temperatureLabel, temperatureSlider)

    // Animation loop
    go func() {
        ticker := time.NewTicker(16 * time.Millisecond) // Approximately 60 FPS
        defer ticker.Stop()
        for range ticker.C {
            // Update molecule positions and velocities
            for i := 0; i < moleculeCount; i++ {
                m1 := molecules[i]

                // Adjust velocity based on temperature
                speedFactor := temperatureSlider.Value / 5.0 // Normalize to initial temperature
                m1.velX *= speedFactor
                m1.velY *= speedFactor

                // Update position
                m1.posX += m1.velX
                m1.posY += m1.velY

                // Check for collisions with walls
                if m1.posX <= 0 || m1.posX >= windowWidth-moleculeSize {
                    m1.velX *= -1
                }
                if m1.posY <= 0 || m1.posY >= windowHeight-moleculeSize {
                    m1.velY *= -1
                }

                // Keep molecule within bounds
                if m1.posX < 0 {
                    m1.posX = 0
                }
                if m1.posX > windowWidth-moleculeSize {
                    m1.posX = windowWidth - moleculeSize
                }
                if m1.posY < 0 {
                    m1.posY = 0
                }
                if m1.posY > windowHeight-moleculeSize {
                    m1.posY = windowHeight - moleculeSize
                }

                // Check for collisions with other molecules
                for j := i + 1; j < moleculeCount; j++ {
                    m2 := molecules[j]
                    if isColliding(m1, m2) {
                        // Elastic collision response
                        handleCollision(m1, m2)
                    }
                }

                // Move the circle to the new position
                m1.circle.Move(fyne.NewPos(float32(m1.posX), float32(m1.posY)))
            }

            // Refresh the container to update positions
            moleculeContainer.Refresh()
        }
    }()

    // Layout the controls and simulation area
    content := container.NewBorder(controls, nil, nil, nil, moleculeContainer)

    // Set the content and show the window
    myWindow.SetContent(content)
    myWindow.Resize(fyne.NewSize(windowWidth, windowHeight))
    myWindow.ShowAndRun()
}

// Function to generate a random color
func randomColor() color.Color {
    return color.NRGBA{
        R: uint8(rand.Intn(256)),
        G: uint8(rand.Intn(256)),
        B: uint8(rand.Intn(256)),
        A: 255,
    }
}

// Function to check if two molecules are colliding
func isColliding(m1, m2 *Molecule) bool {
    dx := m1.posX - m2.posX
    dy := m1.posY - m2.posY
    distance := math.Sqrt(dx*dx + dy*dy)
    return distance < moleculeSize
}

// Function to handle collision between two molecules
func handleCollision(m1, m2 *Molecule) {
    // Calculate normal vector
    dx := m2.posX - m1.posX
    dy := m2.posY - m1.posY
    distance := math.Sqrt(dx*dx + dy*dy)

    if distance == 0 {
        // Avoid division by zero
        return
    }

    nx := dx / distance
    ny := dy / distance

    // Relative velocity
    dvx := m1.velX - m2.velX
    dvy := m1.velY - m2.velY

    // Dot product of relative velocity and normal vector
    dotProduct := dvx*nx + dvy*ny

    if dotProduct > 0 {
        // Molecules are moving away from each other
        return
    }

    // Exchange velocities (simplified elastic collision)
    m1.velX -= dotProduct * nx
    m1.velY -= dotProduct * ny
    m2.velX += dotProduct * nx
    m2.velY += dotProduct * ny
}


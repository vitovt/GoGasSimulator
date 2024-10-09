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
    "fyne.io/fyne/v2/layout"
    "fyne.io/fyne/v2/widget"
)

// Constants and Variables
const (
    defaultTemperature = 300.0
    moleculeSize   = 8.0
    minSpeed       = 1.0
    maxSpeed       = 5.0
)

var (
    moleculesCount = 100
    windowWidth  = 800.0
    windowHeight  = 600.0 
)

var (
    moleculesColor = color.NRGBA{R: 0, G: 0, B: 255, A: 255} // Blue color for uncharged molecules
    chargedColor   = color.NRGBA{R: 255, G: 0, B: 0, A: 255} // Red color for the charged particle
)

type Molecule struct {
    circle    *canvas.Circle
    posX      float64
    posY      float64
    velX      float64
    velY      float64
    isCharged bool
}

func main() {
    // Seed the random number generator
    rand.Seed(time.Now().UnixNano())

    // Create a new application
    myApp := app.New()
    myWindow := myApp.NewWindow("Ideal Gas Simulation with Electric Field")

    // Create a container without layout to control absolute positioning
    moleculeContainer := container.NewWithoutLayout()

    // Create a visual border
    border := canvas.NewRectangle(color.Transparent)
    border.StrokeColor = color.Black
    border.StrokeWidth = 2
    moleculeContainer.Add(border)

    // Initial temperature value
    var previousTemperature = defaultTemperature // Initial temperature value

    // Declare molecules variable at a higher scope
    var molecules []*Molecule

    // Temperature slider
    temperatureSlider := widget.NewSlider(2, 1000)
    temperatureSlider.Value = previousTemperature
    temperatureSlider.Step = 10.0
    temperatureLabel := widget.NewLabel("Temperature: " + fmt.Sprintf("%.1f", previousTemperature) + "K")

    // Gravity slider
    gravitySlider := widget.NewSlider(0, 20)
    gravitySlider.Value = 0
    gravitySlider.Step = 0.1
    gravityLabel := widget.NewLabel("Gravity: " + fmt.Sprintf("%.1f", gravitySlider.Value) + "g")

    // ElectricX field slider
    electricXFieldSlider := widget.NewSlider(-5, 5)
    electricXFieldSlider.Value = 0 // Starting electric field
    electricXFieldSlider.Step = 0.1
    electricXFieldLabel := widget.NewLabel("ElectricX Field: " + fmt.Sprintf("%.1f", electricXFieldSlider.Value))

    // ElectricY field slider
    electricYFieldSlider := widget.NewSlider(-5, 5)
    electricYFieldSlider.Value = 0 // Starting electric field
    electricYFieldSlider.Step = 0.1
    electricYFieldLabel := widget.NewLabel("ElectricY Field: " + fmt.Sprintf("%.1f", electricYFieldSlider.Value))

     // Create arrow components
    arrowShaft := canvas.NewLine(color.NRGBA{R: 255, G: 0, B: 0, A: 255}) // Red color
    arrowShaft.StrokeWidth = 2

    arrowHeadLeft := canvas.NewLine(color.NRGBA{R: 255, G: 0, B: 0, A: 255})
    arrowHeadLeft.StrokeWidth = 2

    arrowHeadRight := canvas.NewLine(color.NRGBA{R: 255, G: 0, B: 0, A: 255})
    arrowHeadRight.StrokeWidth = 2

    // Add arrow components to the container
    moleculeContainer.Add(arrowShaft)
    moleculeContainer.Add(arrowHeadLeft)
    moleculeContainer.Add(arrowHeadRight)

    // Initial update of the arrow
    updateArrow(arrowShaft, arrowHeadLeft, arrowHeadRight, moleculeContainer, electricXFieldSlider.Value, electricYFieldSlider.Value)

   // Adjust the size of the sliders by wrapping them in containers
    sliderWidth := 200.0

    // Create containers for sliders with fixed size
    temperatureSliderContainer := container.New(
        layout.NewGridWrapLayout(fyne.NewSize(float32(sliderWidth), temperatureSlider.MinSize().Height)),
        temperatureSlider,
    )
    gravitySliderContainer := container.New(
        layout.NewGridWrapLayout(fyne.NewSize(float32(sliderWidth), temperatureSlider.MinSize().Height)),
        gravitySlider,
    )
    electricXFieldSliderContainer := container.New(
        layout.NewGridWrapLayout(fyne.NewSize(float32(sliderWidth), electricXFieldSlider.MinSize().Height)),
        electricXFieldSlider,
    )
    electricYFieldSliderContainer := container.New(
        layout.NewGridWrapLayout(fyne.NewSize(float32(sliderWidth), electricYFieldSlider.MinSize().Height)),
        electricYFieldSlider,
    )

    // Create the Reset button
    resetButton := widget.NewButton("Reset", func() {
        // Reset sliders to default values
        temperatureSlider.SetValue(defaultTemperature)
        gravitySlider.SetValue(0)
        electricXFieldSlider.SetValue(0)
        electricYFieldSlider.SetValue(0)

        // Reset labels
        temperatureLabel.SetText("Temperature: "  + fmt.Sprintf("%.1f", defaultTemperature) + "K")
        gravityLabel.SetText("Gravity: 0.0g")
        electricXFieldLabel.SetText("ElectricX Field: 0.0")
        electricYFieldLabel.SetText("ElectricY Field: 0.0")

        // Reset previousTemperature
        previousTemperature = defaultTemperature

        // Reset molecule velocities
        for _, m := range molecules {
            // Random initial velocity based on default temperature
            angle := rand.Float64() * 2 * math.Pi
            speed := rand.Float64()*(maxSpeed - minSpeed) + minSpeed
            m.velX = speed * math.Cos(angle)
            m.velY = speed * math.Sin(angle)
        }
    })

    // Create the Restart button
    restartButton := widget.NewButton("Restart", func() {
        // Remove existing molecule circles from moleculeContainer
        for _, m := range molecules {
            moleculeContainer.Remove(m.circle)
        }

        // Reinitialize molecules
        molecules = initializeMolecules(moleculeContainer, temperatureSlider.Value)
    })

    // Arrange labels and sliders on the same line
    topControl := container.NewHBox(
        temperatureLabel,
        temperatureSliderContainer,
        gravityLabel,
        gravitySliderContainer,
        resetButton,
        restartButton,
    )
    bottomControl := container.NewHBox(
        electricXFieldLabel,
        electricXFieldSliderContainer,
        electricYFieldLabel,
        electricYFieldSliderContainer,
    )

    // Controls container
    controls := container.NewVBox(
        topControl,
        bottomControl,
    )

    // Layout the controls and simulation area
    content := container.NewBorder(
        controls,            // Top
        nil,                 // Bottom
        nil,                 // Left
        nil,                 // Right
        moleculeContainer,   // Center
    )

    // Set the content and show the window
    myWindow.SetContent(content)
    // Adjust window size to accommodate controls and molecule area
    myWindow.Resize(fyne.NewSize(float32(windowWidth), float32(windowHeight+100)))

    // Initialize molecules after the window and content have been set
    molecules = initializeMolecules(moleculeContainer, temperatureSlider.Value)
    // Handle window resize events
    myWindow.Canvas().SetOnTypedKey(func(event *fyne.KeyEvent) {
        // Optional: Handle keyboard events if needed
    })

    // Update moleculeContainer and border size when the window is resized
    myWindow.SetOnClosed(func() {
        // Clean up if necessary
    })

    // Animation loop
    go func() {
        ticker := time.NewTicker(16 * time.Millisecond) // Approximately 60 FPS
        defer ticker.Stop()
        for range ticker.C {
            // Update windowWidth and windowHeight based on moleculeContainer size
            windowWidth = float64(moleculeContainer.Size().Width)
            windowHeight = float64(moleculeContainer.Size().Height)



            // Update the border size
            border.Resize(moleculeContainer.Size())

            // Update molecule positions and velocities
            for i := 0; i < moleculesCount; i++ {
                m1 := molecules[i]

                // Apply electricX field force to charged particle
                if m1.isCharged && electricXFieldSlider.Value != 0 {
                    // ElectricX field applies a force in the X-direction
                    electricXForce := electricXFieldSlider.Value * 0.1 // Adjust the multiplier as needed
                    m1.velX += electricXForce
                }

                // Apply electricY field force to charged particle
                if m1.isCharged && electricYFieldSlider.Value != 0 {
                    // ElectricY field applies a force in the Y-direction
                    electricYForce := electricYFieldSlider.Value * 0.1 // Adjust the multiplier as needed
                    m1.velY -= electricYForce
                }


                // Apply gravity field force to all particles
                if gravitySlider.Value != 0 {
                    // Gravity field applies a force in the Y-direction
                    gravityForce := gravitySlider.Value * 0.01 // Adjust the multiplier as needed
                    m1.velY += gravityForce
                }

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
                for j := i + 1; j < moleculesCount; j++ {
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

    // Update temperature label and adjust velocities when slider changes
    temperatureSlider.OnChanged = func(value float64) {
        temperatureLabel.SetText("Temperature: " + fmt.Sprintf("%.1f", value) + "K")

        // Calculate the scaling factor
        scalingFactor := math.Sqrt(value / previousTemperature)

        // Adjust velocities of all molecules
        for _, m := range molecules {
            m.velX *= scalingFactor
            m.velY *= scalingFactor
        }

        // Update the previous temperature
        previousTemperature = value
    }

    // Update gravity label when slider changes
    gravitySlider.OnChanged = func(value float64) {
        gravityLabel.SetText("Gravity: " + fmt.Sprintf("%.1f", value) + "g")
    }

    // Update electricX field label when slider changes
    electricXFieldSlider.OnChanged = func(value float64) {
        electricXFieldLabel.SetText("ElectricX Field: " + fmt.Sprintf("%.1f", value))
        updateArrow(arrowShaft, arrowHeadLeft, arrowHeadRight, moleculeContainer, value, electricYFieldSlider.Value)
    }
    // Update electricY field label when slider changes
    electricYFieldSlider.OnChanged = func(value float64) {
        electricYFieldLabel.SetText("ElectricY Field: " + fmt.Sprintf("%.1f", value))
        updateArrow(arrowShaft, arrowHeadLeft, arrowHeadRight, moleculeContainer, electricXFieldSlider.Value, value)
    }

    myWindow.ShowAndRun()
}

// Initialize molecules function
func initializeMolecules(moleculeContainer *fyne.Container, temperature float64) []*Molecule {
    // Get the initial size of the molecule container
    molecules := make([]*Molecule, 0, moleculesCount)
    minDistance := 2.0 * float64(moleculeSize)

    for i := 0; i < moleculesCount; i++ {
        var posX, posY float64
        maxAttempts := 1000
        attempts := 0

        // Loop to find a valid position
        for {
            // Random initial position
            posX = rand.Float64()*(windowWidth-moleculeSize*2) + moleculeSize
            posY = rand.Float64()*(windowHeight-moleculeSize*2) + moleculeSize

            validPosition := true

            // Check against all previously placed molecules
            for _, m := range molecules {
                dx := posX - m.posX
                dy := posY - m.posY
                distance := math.Sqrt(dx*dx + dy*dy)
                if distance < minDistance {
                    validPosition = false
                    break
                }
            }

            if validPosition {
                break
            }

            attempts++
            if attempts >= maxAttempts {
                fmt.Println("Warning: Max attempts reached while placing molecule", i)
                break
            }
        }

        // Random initial velocity
        angle := rand.Float64() * 2 * math.Pi
        speed := rand.Float64()*(maxSpeed-minSpeed) + minSpeed
        speed = speed * temperature/defaultTemperature
        velX := speed * math.Cos(angle)
        velY := speed * math.Sin(angle)

        // Create a circle for the molecule
        var circle *canvas.Circle
        var isCharged bool

        if i == 0 {
            // First molecule is the charged particle
            circle = canvas.NewCircle(chargedColor) // Use chargedColor variable
            isCharged = true
        } else {
            circle = canvas.NewCircle(moleculesColor) // Use moleculesColor variable
            isCharged = false
        }

        circle.Resize(fyne.NewSize(float32(moleculeSize), float32(moleculeSize)))
        circle.Move(fyne.NewPos(float32(posX), float32(posY)))

        // Add to container
        moleculeContainer.Add(circle)

        // Add to molecule slice
        molecule := &Molecule{
            circle:    circle,
            posX:      posX,
            posY:      posY,
            velX:      velX,
            velY:      velY,
            isCharged: isCharged,
        }
        molecules = append(molecules, molecule)
    }

    return molecules
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
        // Apply a small random displacement
        displacement := float64(moleculeSize) * 0.01
        dx = displacement * (rand.Float64()*2 - 1)
        dy = displacement * (rand.Float64()*2 - 1)
        distance = math.Sqrt(dx*dx + dy*dy)
    }

    // Minimum distance between molecules to avoid overlap
    minDistance := float64(moleculeSize)

    // Calculate overlap amount
    overlap := minDistance - distance

    if overlap > 0 {
        // Normalize the collision normal vector
        nx := dx / distance
        ny := dy / distance

        // Push molecules apart based on their masses (assuming equal mass)
        m1.posX -= nx * overlap / 2
        m1.posY -= ny * overlap / 2
        m2.posX += nx * overlap / 2
        m2.posY += ny * overlap / 2

        // Update distance after position correction
        distance = minDistance
        dx = m2.posX - m1.posX
        dy = m2.posY - m1.posY
        nx = dx / distance
        ny = dy / distance

        // Relative velocity
        dvx := m1.velX - m2.velX
        dvy := m1.velY - m2.velY

        // Dot product of relative velocity and normal vector
        dotProduct := dvx*nx + dvy*ny

        // Compute restitution (e = 1 for perfectly elastic collision)
        e := 1.0

        // Impulse scalar
        impulse := -(1 + e) * dotProduct / 2 // Divide by 2 for equal mass

        // Update velocities
        m1.velX += impulse * nx
        m1.velY += impulse * ny
        m2.velX -= impulse * nx
        m2.velY -= impulse * ny
    }
}

// Function to update the arrow's position and orientation
func updateArrow(arrowShaft, arrowHeadLeft, arrowHeadRight *canvas.Line, moleculeContainer *fyne.Container, E_x, E_y float64) {
    // Calculate center point
    centerX := moleculeContainer.Size().Width / 2
    centerY := moleculeContainer.Size().Height / 2

    E_y = -E_y

    // Calculate the magnitude and angle of the electric field
    magnitude := math.Sqrt(E_x*E_x + E_y*E_y)
    angle := math.Atan2(E_y, E_x)

    // Scale the arrow length (adjust scalingFactor as needed)
    scalingFactor := float64(moleculeContainer.Size().Width) / 10 // Adjust this to make the arrow visible but not too big
    arrowLength := float32(magnitude * scalingFactor)

    // Calculate end point of the arrow shaft
    endX := centerX + arrowLength*float32(math.Cos(angle))
    endY := centerY + arrowLength*float32(math.Sin(angle))

    // Update the shaft
    arrowShaft.Position1 = fyne.NewPos(centerX, centerY)
    arrowShaft.Position2 = fyne.NewPos(endX, endY)

    // Calculate points for arrowhead
    headLength := arrowLength * 0.2 // Length of the arrowhead lines
    headAngle := math.Pi / 6        // 30 degrees for arrowhead

    // Left side of arrowhead
    leftAngle := angle + math.Pi - headAngle
    leftX := endX + headLength*float32(math.Cos(leftAngle))
    leftY := endY + headLength*float32(math.Sin(leftAngle))

    // Right side of arrowhead
    rightAngle := angle + math.Pi + headAngle
    rightX := endX + headLength*float32(math.Cos(rightAngle))
    rightY := endY + headLength*float32(math.Sin(rightAngle))

    // Update arrowhead lines
    arrowHeadLeft.Position1 = fyne.NewPos(endX, endY)
    arrowHeadLeft.Position2 = fyne.NewPos(leftX, leftY)

    arrowHeadRight.Position1 = fyne.NewPos(endX, endY)
    arrowHeadRight.Position2 = fyne.NewPos(rightX, rightY)

    // Refresh the lines
    arrowShaft.Refresh()
    arrowHeadLeft.Refresh()
    arrowHeadRight.Refresh()
}

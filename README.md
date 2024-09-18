# GoGasSimulator

A visual simulation of ideal gas molecules with temperature and electric field controls, built in Go using the Fyne GUI toolkit.

![Simulation Screenshot](screenshots/GoGasSimulator.png)

## Overview

**GoGasSimulator** is an interactive application that simulates the behavior of ideal gas molecules within a container. It provides controls to adjust temperature and apply an external electric field to observe their effects on the gas molecules.

## Features

- **100 Ideal Gas Molecules**: Simulates the movement and interaction of 100 gas molecules.
- **Temperature Control**: Adjust the temperature to see how it affects the speed of the molecules.
- **Electric Field Control**: Apply an external electric field that influences a specific charged particle.
- **Electric Field Direction**: Set direction for external electric field, separately for X and Y component.
- **Charged Particle**: A distinguishable molecule affected by the electric field, represented in a different color.
- **Elastic Collisions**: Molecules collide elastically with each other and the container walls.
- **User-Friendly Interface**: Built with Fyne, offering a clean and responsive GUI.
- **Gravitation Field Control**: Apply an external gravitation field that influences to all particles.

## Screenshots

Screenshots are located in the `screenshots` directory.

![Electrical field force](screenshots/ElectricField.png)
![Gravitation force](screenshots/Gravity.png)

## Installation

### Prerequisites

- **Go 1.19 or newer**: [Download Go](https://golang.org/dl/)
- **Fyne Toolkit**: Install the Fyne GUI toolkit for Go.

### Install Go

Ensure Go is installed:

```bash
go version
```

If not installed, download it from the [official website](https://golang.org/dl/).
Or user repository, for example for */Ubuntu it is:

```bash
sudo add-apt-repository ppa:longsleep/golang-backports
sudo apt install golang-go
```

### Install Fyne Dependencies

For **Linux Ubuntu/Kubuntu**:

```bash
sudo apt update
sudo apt install libgl1-mesa-dev xorg-dev
```

### Install Fyne

```bash
go get fyne.io/fyne/v2
go get fyne.io/fyne/v2/cmd/fyne
```

or

```bash
go mod GoGasSimulator
go mod tidy
```


## Building and Running the Application

### Clone the Repository

```bash
git clone https://github.com/yourusername/GoGasSimulator.git
cd GoGasSimulator
```

### Build the Application

```bash
go build -o GoGasSimulator
```

### Run the Application

```bash
./GoGasSimulator
```

Alternatively, you can run it directly:

```bash
go run main.go
```

## Usage

- **Temperature Slider**: Move the slider to increase or decrease the temperature, affecting molecule speeds.
- **Electric Field Slider**: Adjust the slider to apply an electric field, influencing the charged particle.
- **Observation**: Watch the molecules interact and observe the effects of temperature and electric field changes.

## License

This project is licensed under the [MIT License](LICENSE).

## Contributing

Contributions are welcome! Please open issues or submit pull requests for improvements or new features.

## Contact

For questions or suggestions, please contact here or PM in Linkedin (see profile).


# SBQOL (Starbase Quality of Life Features)

**SBQOL** is a collection of quality-of-life features for the game Starbase. This project aims to enhance the player experience by providing additional tools and functionalities not available in the base game.

# All improvements in this tool have been given the verbal OK by Frozenbyte, the developer of Starbase.
Please contact [@Dodgin](https://discord.com/users/112299261734965248) if you have any concerns.

# TL;DR How does it work?
The application does a series of calibration memory scans of your running starbase program to determine the addresses of the floats (a computer data type) that the levers on your ship show the value of. The program then binds your HOTAS or twin-stick setup so that when you input an action from the joystick or throttle, SBQOL inputs your keybind into the game until the desired value is matched. This way, we are reading but not modifying program memory, which is toe-ing the line for most games.

# But isn't Memory Scanning Bad:tm:?
In most cases developers would prefer you not do this. In this case, a few floats are necessary to implement a commonly requested user feature. You can view what parts of the application memory are scanned in [scan.go](scan.go).

## Features

- **Joystick Support**: Control various in-game actions using joysticks.
- **Throttle Control**: Enhanced throttle control for better in-game maneuvering.

## Prerequisites
- SDL2 - This is a pretty new library (1998) for peripheral interaction that you may have heard of. I have included SDL2.dll but you should not trust random dlls in random repositories. Go download SDL2 >>> [HERE](https://wiki.libsdl.org/SDL2/Installation) <<< and compile it or use theirs.
- Go - a c developers fever dream that turned into a web language, also good for scripting one-off apps on occasion.
- A computer - arguable, but suggested.

## Building (Windows)

1. Clone the repository:
   ```sh
   git clone https://github.com/Dodgin/sbqol.git
   cd sbqol
   ```

2. Build the project:
   ```sh
   go build
   ```

3. Run the application:
   ```sh
   ./sbqol
   ```

## Debug Building

### Run npm build in the resources directory
`npm --prefix ./resources run build`
This builds the UI :)

### Run the Go application
```
$env:CGO_CFLAGS="-IC:/SDL2/include"
$env:CGO_LDFLAGS="-LC:/SDL2/lib/x64"
go run .\main.go .\scan.go .\throttle.go .\ui.go .\hotkeys.go .\sdl.go .\throttlecontroller.go .\joy1.go
```
This runs the program

## Usage

1. Connect your joysticks.
2. Launch the application.
3. Follow the on-screen instructions to configure your joystick and other settings.

## Contributing

Contributions are welcome! Please submit a pull request or open an issue to discuss your ideas.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contact

For any questions or suggestions, please contact [@Dodgin](https://discord.com/users/112299261734965248).
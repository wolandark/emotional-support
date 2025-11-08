# Emotional Support Activity Tracker

A cute activity tracker for Linux (X11) that provides emotional support notifications while you code. It tracks your active windows and programs, detects what you're working on, and sends encouraging messages like "wow you've been in vim for an hour! im so proud of you" or "i know java is hard, but you got it!".

## Features

- **Window Tracking**: Monitors active windows and programs on X11
- **Program Detection**: Recognizes popular editors and IDEs (vim, VSCode, Emacs, IntelliJ, etc.)
- **Language Detection**: Attempts to detect programming languages from:
  - File extensions in window titles
  - IDE workspace files (VSCode settings.json)
  - Project files (go.mod, package.json, requirements.txt, etc.)
- **Emotional Support Messages**: 
  - Time-based encouragement (e.g., "You've been coding for 1 hour!")
  - Language-specific support (e.g., "I know Java is hard, but you got it!")
  - Health reminders (stay hydrated, blink your eyes, stretch)
- **State Persistence**: Saves activity history to `~/.config/emotional-support/state.json`

## Requirements

- Linux with X11
- `xdotool` (for window tracking)
- Go 1.21 or later

## Installation

1. Install `xdotool`:
   ```bash
   # On Arch Linux
   sudo pacman -S xdotool
   
   # On Ubuntu/Debian
   sudo apt-get install xdotool
   
   # On Fedora
   sudo dnf install xdotool
   ```

2. Clone or download this repository

3. Install Go dependencies:
   ```bash
   go mod download
   ```

4. Build the program:
   ```bash
   go build -o emotional-support
   ```

## Usage

Run the program:
```bash
./emotional-support
```

The program will:
- Start tracking your active windows
- Send a welcome notification
- Monitor your coding activity
- Send encouraging notifications based on:
  - Time spent coding (every 30 minutes and hourly milestones)
  - Programming language detected
  - Health reminders (every 20 minutes)

## How It Works

1. **Window Tracking**: Uses `xdotool` to get the active window title and process name every 5 seconds
2. **Context Detection**: Analyzes window titles and process names to identify editors/IDEs
3. **Language Detection**: 
   - Extracts file paths from window titles
   - Reads VSCode workspace settings
   - Checks for common project files (go.mod, package.json, etc.)
4. **Message Generation**: Creates contextual messages based on:
   - Time spent in a program
   - Detected programming language
   - Random health reminders
5. **Notifications**: Uses DBus to send desktop notifications

## Customization

You can customize messages by editing `messages.go`:
- `GetTimeBasedMessage()`: Messages for time milestones
- `GetLanguageMessage()`: Language-specific encouragement
- `GetHealthReminder()`: Health and wellness reminders

## State File

Activity history is saved to `~/.config/emotional-support/state.json`. This file tracks:
- Time spent in each window/program
- Activity timestamps
- Total coding time

## Limitations

- Currently only works on X11 (not Wayland)
- Requires `xdotool` to be installed
- Language detection is best-effort and may not always be accurate
- Some window managers may not provide detailed window titles

## Future Improvements

- Wayland support
- More sophisticated language detection
- Configurable notification intervals
- Statistics dashboard
- More editor/IDE support

## License

Feel free to use and modify as you wish! 


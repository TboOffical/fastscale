package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"net"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	//Version of the application
	version = "1.0.0"
)

var (
	//Host
	host            string
	termContent     binding.String
	status          binding.String
	connected       bool
	keyboardEnabled bool
	currentCommand  string
	shift           bool
	snapshot        string
	cursorChar      string = "|"
	notes           binding.String
	connection      net.Conn
)

func termPrint(s string) {
	current, _ := termContent.Get()
	termContent.Set(current + s)
}

func termPrintLn(s string) {
	current, _ := termContent.Get()
	termContent.Set(current + s + "\n")
}

func clearTerm() {
	termContent.Set("\n\n\n\nTerminal Cleared\n\n")
}

func execute(command string) {
	_, err := connection.Write([]byte(command))
	if err != nil {
		return
	}
}

func main() {
	a := app.NewWithID("com.fastscale.terminal")
	w := a.NewWindow("Fastscale Terminal")

	termContent = binding.NewString()
	status = binding.NewString()
	notes = binding.NewString()

	w.SetPadded(false)

	iconRawData, err := os.ReadFile("./fs_term.png")
	if err != nil {
		panic(err)
	}
	logoRawData, err := os.ReadFile("./fastscale-logo.png")
	if err != nil {
		panic(err)
	}

	iconResource := fyne.NewStaticResource("main_icon", iconRawData)
	logoResource := fyne.NewStaticResource("logo", logoRawData)

	w.SetIcon(iconResource)

	w.Resize(fyne.NewSize(800, 600))
	w.CenterOnScreen()

	today := time.Now()
	termContent.Set("\n\n\n\nWelcome to Fastscale Terminal\n" + today.Format("2006-01-02 15:04:05") + "\n\n")

	termText := widget.NewLabelWithStyle("\n\n\n One Moment...", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})
	termText.Bind(termContent)

	term := container.NewScroll(termText)

	termBackground := container.NewStack(canvas.NewRectangle(color.Black), term)
	termBackground.Resize(fyne.NewSize(800, 600))

	image := canvas.NewImageFromResource(logoResource)
	image.SetMinSize(fyne.NewSize(32, 32))

	connectButton := widget.NewButtonWithIcon("", theme.MailForwardIcon(), func() {
		input := widget.NewEntry()
		input.SetPlaceHolder("Enter the host address")
		input.OnChanged = func(s string) {
			host = s
		}

		connectCanvas := container.NewVBox(widget.NewLabel("Enter the host address in the following format <hostname or ip>:<port>"), container.NewPadded(input))
		connectCanvas.Resize(fyne.NewSize(700, 200))

		dialog.ShowCustomConfirm("Connection Details", "Connect", "Back", connectCanvas, func(b bool) {
			if b != true || host == "" || connected {
				return
			}

			termPrintLn("Connecting to " + host + "...")

			connection, err = net.Dial("tcp", host)
			if err != nil {
				termPrintLn(fmt.Sprint("Error in connection: ", err))
				return
			}

			connected = true
			clearTerm()

			//todo handshake

			keyboardEnabled = true
			status.Set(host)
			w.SetTitle("Fastscale Terminal - " + host)
			a.SendNotification(fyne.NewNotification("Connected", "You are now connected to the host "+host))

			go func() {
				for {
					received := make([]byte, 1024)
					_, err = connection.Read(received)
					if err != nil {
						println("Read data failed:", err.Error())
						os.Exit(1)
					}

					var clean []byte

					//remove zeros
					for _, by := range received {
						if by != 0 {
							clean = append(clean, by)
						}
					}

					final := fmt.Sprint(string(clean))
					termPrintLn(final)
				}
			}()
		}, w)
	})

	settingsButton := widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {
		println("settings")
	})

	InfoButton := widget.NewButtonWithIcon("", theme.InfoIcon(), func() {
		icon := canvas.NewImageFromResource(iconResource)
		icon.SetMinSize(fyne.NewSize(100, 100))
		iconContainer := container.NewCenter(icon)

		about := widget.NewLabel("Fastscale Terminal\nVersion: v" + version + "\n\nDeveloped by Callum F")

		aboutCanvas := container.NewVBox(iconContainer, about, widget.NewButton("Open Source Licenses", func() {
			CreditsWindow(a, fyne.NewSize(800, 600)).Show()
		}))
		aboutCanvas.Resize(fyne.NewSize(300, 200))
		dialog.ShowCustom("About", "Close", aboutCanvas, w)
	})

	status := widget.NewLabelWithData(status)

	titleBarInner := container.New(layout.NewHBoxLayout(), image, status, layout.NewSpacer(), connectButton, settingsButton, InfoButton)
	titleBarPadded := container.NewStack(canvas.NewRectangle(color.RGBA{R: 32, G: 32, B: 32, A: 255}), container.NewPadded(titleBarInner))

	content := container.NewStack(termBackground, container.New(layout.NewVBoxLayout(), titleBarPadded, layout.NewSpacer(), container.NewHBox(layout.NewSpacer(), widget.NewLabelWithData(notes))))

	keyboardEnabled = false

	w.Canvas().SetOnTypedKey(func(event *fyne.KeyEvent) {
		if !keyboardEnabled {
			return
		}

		valid, _ := regexp.Match("^[a-zA-Z0-9]*$", []byte(event.Name))
		final := ""

		if valid {
			if shift {
				final = strings.ToUpper(fmt.Sprint(event.Name))
				shift = false
				notes.Set("")
			} else {
				final = strings.ToLower(fmt.Sprint(event.Name))
			}
		}

		tc, _ := termContent.Get()

		switch fmt.Sprint(event.Name) {
		case "BackSpace":
			if len(currentCommand) == 0 {
				return
			}
			currentCommand = currentCommand[:len(currentCommand)-1]
			termContent.Set(tc[:len(tc)-2] + cursorChar)
			break
		case "Return":
			execute(currentCommand)
			currentCommand = ""
			//todo: execute command
			err := termContent.Set(tc[:len(tc)-1] + "\n" + cursorChar)
			if err == nil {
				term.ScrollToBottom()
				term.ScrollToBottom()
			}
			break
		case "Space":
			termContent.Set(tc[:len(tc)-1] + " " + cursorChar)
			currentCommand += " "
			break
		case ".":
			termContent.Set(tc[:len(tc)-1] + "." + cursorChar)
			currentCommand += "."
		case ";":
			termContent.Set(tc[:len(tc)-1] + ";" + cursorChar)
			currentCommand += ";"
		case "-":
			if shift {
				termContent.Set(tc[:len(tc)-1] + "_" + cursorChar)
				currentCommand += "_"
			}
		case "LeftShift":
			shift = true
			notes.Set("shift")
		}

		if valid && len(fmt.Sprint(event.Name)) == 1 {
			currentCommand += final
			termContent.Set(tc[:len(tc)-1] + final + fmt.Sprint(cursorChar))
		}

	})

	w.SetContent(content)
	w.ShowAndRun()
}

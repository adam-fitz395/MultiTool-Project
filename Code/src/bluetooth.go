package main

import (
	"bufio"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"os/exec"
	"regexp"
	"strings"
)

func loadBluetoothMenu() {
	buttons = nil

	btScanButton := tview.NewButton("Scan").
		SetSelectedFunc(func() {
			pages.SwitchToPage("btScan")
		})
	btScanButton.SetBorder(true).
		SetBorderColor(tcell.ColorWhite)

	btDeauthButton := tview.NewButton("Deauth").
		SetSelectedFunc(func() {
			pages.SwitchToPage("btDeauth")
		})
	btDeauthButton.SetBorder(true).
		SetBorderColor(tcell.ColorWhite)

	backButton := tview.NewButton("Back").
		SetSelectedFunc(func() {
			pages.SwitchToPage("main") // Switch back to the main page
		})
	backButton.SetBorder(true).
		SetBorderColor(tcell.ColorWhite)

	bluetoothFlex := tview.NewFlex().
		AddItem(btScanButton, 0, 1, true).
		AddItem(btDeauthButton, 0, 1, true).
		AddItem(backButton, 0, 1, false).
		SetDirection(tview.FlexRow)

	buttons = []*tview.Button{btScanButton, btDeauthButton, backButton}
	pages.AddPage("bluetooth", bluetoothFlex, true, false) // Add the Wi-Fi page to pages
	enableTabFocus(bluetoothFlex, buttons)
}

func loadBluetoothScan() {
	buttons = nil
	var duration int
	btScanText := tview.NewTextView().
		SetDynamicColors(true).
		SetText("[green]Ready to Scan!")

	btScanText.SetBorder(true)

	btScanDuration := tview.NewDropDown().SetLabel("Duration (Seconds): ").SetLabelColor(tcell.ColorWhite)
	btScanDuration.
		AddOption("10",
			func() {
				duration = 10
			}).
		AddOption("20", func() {
			duration = 20
		}).
		AddOption("30", func() {
			duration = 30
		}).AddOption("60", func() {
		duration = 60
	})
	btScanDuration.SetCurrentOption(0)

	btScanButton := tview.NewButton("Scan").
		SetSelectedFunc(func() {
			btScanText.SetText("[white]Scanning in progress...please wait!")
			go func() {
				var devices []string

				// Run bash script with bluetoothctl scan
				cmd := exec.Command("bash", "../scripts/bt_scan.sh", fmt.Sprintf("%d", duration))

				stdout, err := cmd.StdoutPipe()
				if err != nil {
					app.QueueUpdateDraw(func() {
						btScanText.SetText(fmt.Sprintf("[red]Failed to create pipe: %v\n[red]", err))
					})
					return
				}

				// Start the script to initiate scanning
				if err := cmd.Start(); err != nil {
					app.QueueUpdateDraw(func() {
						btScanText.SetText(fmt.Sprintf("[red]Error starting script: %v\n", err))
					})
					return
				}

				// Read script output line by line
				scanner := bufio.NewScanner(stdout)
				for scanner.Scan() {
					line := strings.TrimSpace(scanner.Text())

					re := regexp.MustCompile(`([0-9A-Fa-f:]{17})\s+-\s+(.*)`)
					match := re.FindStringSubmatch(line)

					if len(match) > 0 {
						deviceMAC := match[1]  // MAC address
						deviceName := match[2] // Device name

						// Ensure unique devices (MAC is unique key)

						devices = append(devices, fmt.Sprintf("%s - %s", deviceMAC, deviceName))
					}
				}

				if err := cmd.Wait(); err != nil {
					app.QueueUpdateDraw(func() {
						btScanText.SetText(fmt.Sprintf("[red]Script execution error: %v\n", err))
					})
					return
				}

				if err := scanner.Err(); err != nil {
					app.QueueUpdateDraw(func() {
						btScanText.SetText(fmt.Sprintf("[red]Scanner error: %v\n", err))
					})
					return
				}

				// Update UI with results
				app.QueueUpdateDraw(func() {
					if len(devices) == 0 {
						btScanText.SetText("[green]Scan complete!\n[Red]No devices found.")
					} else {
						btScanText.SetText(fmt.Sprintf("[green]Scan complete!\n[white]Found Devices:\n%s\n[blue]A log file has been created!",
							strings.Join(devices, "\n")))
					}
				})
			}()
		})

	btScanButton.SetBorder(true).
		SetBorderColor(tcell.ColorWhite)

	backButton := tview.NewButton("Back").SetSelectedFunc(func() {
		pages.SwitchToPage("bluetooth")
	})

	backButton.SetBorder(true).
		SetBorderColor(tcell.ColorWhite)

	btScanFlex := tview.NewFlex().
		AddItem(btScanText, 0, 3, false).
		AddItem(btScanDuration, 0, 1, false).
		AddItem(btScanButton, 0, 1, true).
		AddItem(backButton, 0, 1, false)
	btScanFlex.SetDirection(tview.FlexRow)

	buttons = []*tview.Button{btScanButton, backButton}
	pages.AddPage("btScan", btScanFlex, true, false) // Add the Wi-Fi page to pages
	enableTabFocus(btScanFlex, buttons)
}

func loadBluetoothDeauth() {
	buttons = nil

	deauthText := tview.NewTextView().
		SetDynamicColors(true).
		SetText("[green]Ready to Deauth! Please select a device!")

	deauthText.SetBorder(true).
		SetBorderColor(tcell.ColorWhite)

	deviceList := tview.NewDropDown()
	deviceList.SetBorder(true).SetBorderColor(tcell.ColorWhite)
	deviceList.SetLabel("Device HERE!")

	deauthButton := tview.NewButton("Deauth").
		SetSelectedFunc(func() {

		})

	deauthButton.SetBorder(true).
		SetBorderColor(tcell.ColorWhite)

	backButton := tview.NewButton("Back").
		SetSelectedFunc(func() {
			pages.SwitchToPage("bluetooth")
		})
	backButton.SetBorder(true).
		SetBorderColor(tcell.ColorWhite)

	deauthFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(deauthText, 0, 2, false).
		AddItem(deviceList, 0, 1, false).
		AddItem(deauthButton, 0, 1, true).
		AddItem(backButton, 0, 1, false)

	buttons = []*tview.Button{deauthButton, backButton}
	pages.AddPage("btDeauth", deauthFlex, true, false)
	enableTabFocus(deauthFlex, buttons)
}

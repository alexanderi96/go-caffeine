package main

import (
    _ "embed"
    "log"
    "time"

    "github.com/getlantern/systray"
    "github.com/micmonay/keybd_event"
)

var (
    //go:embed assets/awake.ico
    awake []byte

    //go:embed assets/sleepy.ico
    sleepy []byte

    running     bool = false
    runningChan      = make(chan bool)
)

const (
    appName string = "go-caffeine"
    appDesc string = "go-caffeine helps you keep the system awake"
)

func main() {
    systray.Run(onReady, onExit)
}

func onReady() {

    systray.SetIcon(sleepy)
    systray.SetTitle(appName)
    systray.SetTooltip(appDesc)

    triggerCaffeine := systray.AddMenuItem("Start Caffeine", "Caffeine")

    systray.AddSeparator()
    mQuit := systray.AddMenuItem("Quit", "Quit the whole app.")

    k, err := keybd_event.NewKeyBonding()
    if err != nil {
        panic(err)
    }

    // Setting the pressed key to F15
    k.SetKeys(keybd_event.VK_F15)

    for {
        select {
        case <-triggerCaffeine.ClickedCh:
            toggleCaffeine()

        case <-mQuit.ClickedCh:
            systray.Quit()
            return
        }
    }
}

func toggleCaffeine() {
    running = !running

    if running {
        log.Print("Should be started")
        go runCaffeine(k)
        systray.SetIcon(awake)
        triggerCaffeine.SetTitle("Stop Caffeine")
    } else {
        log.Print("Should be stopped")
        runningChan <- !running
        systray.SetIcon(sleepy)
    }
}

func runCaffeine(k keybd_event.KeyBonding) {
    for {
        select {
        case <-runningChan:
            log.Print("Stopping Caffeine")
            return
        default:
            if err := k.Launching(); err != nil {
                panic(err)
            }

            log.Print("Key pressed")

            time.Sleep(5 * time.Second)
        }
    }
}

func killCaffeineIfAlive() {
    if running {
        toggleCaffeine()
    }
}

func onExit() {
    // Cleaning stuff here.
    killCaffeineIfAlive()
    log.Print("Exiting now...")
}

package main

import (
    _ "embed"
    "io/ioutil"
    "log"
    "os"
    "os/user"
    "path/filepath"
    "time"

    "github.com/getlantern/systray"
    "github.com/micmonay/keybd_event"
    "gopkg.in/yaml.v2"
)

type Config struct {
    Key       string `yaml:"key"`
    Time      int    `yaml:"time"`
    Autostart bool   `yaml:"autostart"`
}

var (
    //go:embed assets/awake.ico
    awake []byte

    //go:embed assets/sleepy.ico
    sleepy []byte

    running     bool = false
    runningChan      = make(chan bool)

    configPath string = "/.config/go-caffeine/config.yaml"
    logPath    string = "/.log/go-caffeine.log"
    logFile    *os.File

    icons = map[string][]byte{
        "awake":  awake,
        "sleepy": sleepy,
    }

    config Config
)

func init() {
    u, _ := user.Current()
    configPath = u.HomeDir + configPath
    logPath = u.HomeDir + logPath

    logFile, err := createFileAndDirectory(logPath)
    if err != nil {
        log.Fatalln(err)
    }

    log.SetOutput(logFile)

    createFileAndDirectory(configPath)

    loadConfig()
}

func main() {
    defer logFile.Close()
    log.Print("Program started")
    systray.Run(onReady, onExit)
}

func onReady() {
    systray.SetIcon(icons["sleepy"])
    systray.SetTitle("go-caffeine")
    systray.SetTooltip("go-caffeine helps you keep the system awake")

    toggleItem := systray.AddMenuItem("Start Caffeine ☕", "Start or stop caffeine")
    reloadConfig := systray.AddMenuItem("Reload config ♻️", "Reload settings from config file")

    systray.AddSeparator()
    mQuit := systray.AddMenuItem("Quit 🏃", "Quit the whole app")

    k, err := keybd_event.NewKeyBonding()
    if err != nil {
        log.Fatal(err)
    }
    k.SetKeys(parseKey(config.Key))

    if config.Autostart {
        log.Print("Autostarting")
        toggleCaffeine(k, toggleItem, &config)
    }

    for {
        select {
        case <-toggleItem.ClickedCh:
            toggleCaffeine(k, toggleItem, &config)

        case <-reloadConfig.ClickedCh:
            loadConfig()
            if running {
                toggleCaffeine(k, toggleItem, &config)
                toggleCaffeine(k, toggleItem, &config)
            }

        case <-mQuit.ClickedCh:
            systray.Quit()
            return
        }
    }
}

func toggleCaffeine(k keybd_event.KeyBonding, toggleItem *systray.MenuItem, config *Config) {
    if !running {
        log.Print("Starting caffeine")
        systray.SetIcon(icons["awake"])
        toggleItem.SetTitle("Stop Caffeine 🛑")
        go runCaffeine(k, config.Time)
    } else {
        log.Print("Stopping caffeine")
        systray.SetIcon(icons["sleepy"])
        toggleItem.SetTitle("Start Caffeine ☕")
        runningChan <- running
    }
    running = !running
}

func runCaffeine(k keybd_event.KeyBonding, t int) {
    for {
        select {
        case <-runningChan:
            return
        case <-time.After(time.Duration(t) * time.Second): // I press the desired key each desired seconds
            if err := k.Launching(); err != nil {
                log.Fatal(err)
            }
            log.Print("Still there 👀")
        }
    }
}

func onExit() {
    if running {
        runningChan <- running
    }
    log.Print("Program closed")
}

func loadConfig() {
    config = Config{
        Key:       "F15",
        Time:      290,
        Autostart: true,
    }

    configData, err := ioutil.ReadFile(configPath)
    if err != nil || len(configData) == 0 || yaml.Unmarshal(configData, &config) != nil {
        log.Print("Invalid config file, using defaults")
    }

    log.Print("Config loaded: ", config)
}

func parseKey(key string) int {
    var keyCode int

    switch key {
    case "F1":
        keyCode = keybd_event.VK_F1
    case "F2":
        keyCode = keybd_event.VK_F2
    case "F3":
        keyCode = keybd_event.VK_F3
    case "F4":
        keyCode = keybd_event.VK_F4
    case "F5":
        keyCode = keybd_event.VK_F5
    case "F6":
        keyCode = keybd_event.VK_F6
    case "F7":
        keyCode = keybd_event.VK_F7
    case "F8":
        keyCode = keybd_event.VK_F8
    case "F9":
        keyCode = keybd_event.VK_F9
    case "F10":
        keyCode = keybd_event.VK_F10
    case "F11":
        keyCode = keybd_event.VK_F11
    case "F12":
        keyCode = keybd_event.VK_F12
    case "F13":
        keyCode = keybd_event.VK_F13
    case "F14":
        keyCode = keybd_event.VK_F14
    case "F15":
        keyCode = keybd_event.VK_F15
    case "F16":
        keyCode = keybd_event.VK_F16
    default:
        log.Fatal("Key unavailable")
    }

    return keyCode
}

func createFileAndDirectory(filePath string) (*os.File, error) {
    if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
        return nil, err
    }

    file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
    if err != nil {
        return nil, err
    }

    return file, nil
}

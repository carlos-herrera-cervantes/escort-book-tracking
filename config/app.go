package config

import "os"

type app struct {
    Port string
}

var singletonApp *app

func InitApp() *app {
    if singletonApp != nil {
        return singletonApp
    }

    lock.Lock()
    defer lock.Unlock()

    singletonApp = &app{
        Port: os.Getenv("PORT"),
    }

    return singletonApp
}

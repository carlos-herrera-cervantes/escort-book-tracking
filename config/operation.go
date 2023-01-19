package config

type operationConfig struct {
    NewUser string
}

var singletonOperationConfig *operationConfig

func InitOperationConfig() *operationConfig {
    if singletonOperationConfig != nil {
        return singletonOperationConfig
    }

    lock.Lock()
    defer lock.Unlock()

    singletonOperationConfig = &operationConfig{
        NewUser: "new-user",
    }

    return singletonOperationConfig
}

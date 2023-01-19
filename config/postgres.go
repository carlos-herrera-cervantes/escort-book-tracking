package config

import "os"

type postgresConfig struct {
    Host string
    Databases postgresDatabases
}

type postgresDatabases struct {
    EscortTracking string
    EscortProfile string
}

var singletonPostgresConfig *postgresConfig

func InitPostgresConfig() *postgresConfig {
    if singletonPostgresConfig != nil {
        return singletonPostgresConfig
    }

    lock.Lock()
    defer lock.Unlock()

    singletonPostgresConfig = &postgresConfig{
        Host: os.Getenv("POSTGRES_HOST"),
        Databases: postgresDatabases{
            EscortTracking: os.Getenv("ESCORT_TRACKING_DB"),
            EscortProfile: os.Getenv("ESCORT_PROFILE_DB"),
        },
    }

    return singletonPostgresConfig
}

package db

import (
    "database/sql"
    "fmt"
    "log"
    "sync"

    "escort-book-tracking/config"

    _ "github.com/lib/pq"
)

type PostgresClient struct {
    EscortTrackingDB *sql.DB
    EscortProfileDB *sql.DB
}

var singletonPostgresClient *PostgresClient
var lock = &sync.Mutex{}

func NewPostgresClient() *PostgresClient {
    if singletonPostgresClient != nil {
        return singletonPostgresClient
    }

    lock.Lock()
    defer lock.Unlock()

    escortTrackingDBURI := fmt.Sprintf(
        "%s/%s?sslmode=disable",
        config.InitPostgresConfig().Host,
        config.InitPostgresConfig().Databases.EscortTracking,
    )
    escortTrackingDB, err := sql.Open("postgres", escortTrackingDBURI)

    if err != nil {
        log.Panic(err.Error())
    }

    escortProfileDBURI := fmt.Sprintf(
        "%s/%s?sslmode=disable",
        config.InitPostgresConfig().Host,
        config.InitPostgresConfig().Databases.EscortProfile,
    )
    escortProfileDB, err := sql.Open("postgres", escortProfileDBURI)

    if err != nil {
        log.Panic(err.Error())
    }

    singletonPostgresClient = &PostgresClient{
        EscortProfileDB: escortProfileDB,
        EscortTrackingDB: escortTrackingDB,
    }

    return singletonPostgresClient
}

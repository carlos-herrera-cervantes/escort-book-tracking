package config

import "os"

type s3 struct {
    Endpoint string
    Buckets bucket
}

type bucket struct {
    Profile string
}

var singletonS3 *s3

func InitS3() * s3 {
    if singletonS3 != nil {
        return singletonS3
    }

    lock.Lock()
    defer lock.Unlock()

    singletonS3 = &s3{
        Endpoint: os.Getenv("ENDPOINT"),
        Buckets: bucket{
            Profile: os.Getenv("S3"),
        },
    }

    return singletonS3
}

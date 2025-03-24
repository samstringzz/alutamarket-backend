package utils

import (
    "crypto/rand"
    "encoding/base64"
    "fmt"
)

func main() {
    // Generate a 32-byte random key for SECRET_KEY
    secretKey := make([]byte, 32)
    if _, err := rand.Read(secretKey); err != nil {
        panic(err)
    }
    
    // Generate a 32-byte random key for REFRESH_SECRET_KEY
    refreshKey := make([]byte, 32)
    if _, err := rand.Read(refreshKey); err != nil {
        panic(err)
    }
    
    fmt.Printf("SECRET_KEY=%s\n", base64.StdEncoding.EncodeToString(secretKey))
    fmt.Printf("REFRESH_SECRET_KEY=%s\n", base64.StdEncoding.EncodeToString(refreshKey))
}
package security

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "errors"
    "io"
    "os"
)

func getAESKey() ([]byte, error) {
    k := os.Getenv("AES_KEY") 
    if k == "" {
        k = "0123456789abcdef0123456789abcdef"
    }
    if len(k) == 32 {
        return []byte(k), nil
    }
    b, err := base64.StdEncoding.DecodeString(k)
    if err != nil {
        return nil, err
    }
    if len(b) != 32 {
        return nil, errors.New("AES_KEY debe ser 32 bytes")
    }
    return b, nil
}

func Encrypt(plain []byte) (cipherText, iv string, err error) {
    key, err := getAESKey()
    if err != nil { return "", "", err }

    block, err := aes.NewCipher(key)
    if err != nil { return "", "", err }

    aesgcm, err := cipher.NewGCM(block)
    if err != nil { return "", "", err }

    nonce := make([]byte, aesgcm.NonceSize())
    if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
        return "", "", err
    }

    encrypted := aesgcm.Seal(nil, nonce, plain, nil)
    return base64.StdEncoding.EncodeToString(encrypted), base64.StdEncoding.EncodeToString(nonce), nil
}

func Decrypt(cipherText, iv string) ([]byte, error) {
    key, err := getAESKey()
    if err != nil { return nil, err }

    block, err := aes.NewCipher(key)
    if err != nil { return nil, err }

    aesgcm, err := cipher.NewGCM(block)
    if err != nil { return nil, err }

    ct, err := base64.StdEncoding.DecodeString(cipherText)
    if err != nil { return nil, err }
    nonce, err := base64.StdEncoding.DecodeString(iv)
    if err != nil { return nil, err }

    return aesgcm.Open(nil, nonce, ct, nil)
}

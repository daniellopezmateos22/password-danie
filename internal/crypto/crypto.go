package crypto

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "errors"
    "io"
    "os"
)

var key []byte

func init() {
    k := os.Getenv("ENC_KEY_BASE64")
    b, err := base64.StdEncoding.DecodeString(k)
    if err != nil || len(b) != 32 {
        panic("ENC_KEY_BASE64 inválida: base64 de 32 bytes (AES-256)")
    }
    key = b
}

func Encrypt(plain []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, err
    }

    out := gcm.Seal(nonce, nonce, plain, nil)
    return out, nil
}

func Decrypt(ciphertext []byte) (string, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }
    ns := gcm.NonceSize()
    if len(ciphertext) < ns {
        return "", errors.New("ciphertext demasiado corto")
    }
    nonce := ciphertext[:ns]
    data := ciphertext[ns:]
    plain, err := gcm.Open(nil, nonce, data, nil)
    if err != nil {
        return "", err
    }
    return string(plain), nil
}

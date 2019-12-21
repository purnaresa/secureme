package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"io"
	"log"
)

var (
	masterKey string
	DB        *sql.DB
)

// User is data type for user
type User struct {
	ID             string
	Name           string
	NationalID     string
	CreateTimeUnix string
}

func main() {
	//run unit test
}

func createData(id, name, nationalID, createTimeUnix string) (err error) {
	nationalID, _ = encrypt(nationalID, masterKey) // encryption
	_, err = DB.Exec(`
    INSERT INTO 
        user (id, name, national_id, create_time_unix)
    VALUES 
        ("?", "?", "?", "?")
    `, id, name, nationalID, createTimeUnix)

	return
}

func readData(id string) (user User, err error) {
	row := DB.QueryRow(`
    SELECT
        id, name, national_id, create_time_unix
    FROM 
        user
    WHERE 
        id = "?"`, id)

	err = row.Scan(
		&user.ID,
		&user.Name,
		&user.NationalID,
		&user.CreateTimeUnix)

	// decryption
	user.NationalID, _ = decrypt(user.NationalID, masterKey)
	//
	return
}

func encrypt(plaintext, key string) (ciphertext string, err error) {
	keyByte := []byte(key)
	block, _ := aes.NewCipher(keyByte)
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}
	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return
	}

	ciphertextByte := gcm.Seal(
		nonce,
		nonce,
		[]byte(plaintext),
		nil)
	ciphertext = base64.StdEncoding.EncodeToString(ciphertextByte)

	return
}

func decrypt(cipherText, key string) (plainText string, err error) {
	// prepare cipher
	keyByte := []byte(key)
	block, err := aes.NewCipher(keyByte)
	if err != nil {
		return
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}
	nonceSize := gcm.NonceSize()
	//

	// process ciphertext
	ciphertextByte, _ := base64.StdEncoding.DecodeString(cipherText)
	nonce, ciphertextByteClean := ciphertextByte[:nonceSize], ciphertextByte[nonceSize:]
	plaintextByte, err := gcm.Open(
		nil,
		nonce,
		ciphertextByteClean,
		nil)
	if err != nil {
		log.Println(err)
		return
	}
	plainText = string(plaintextByte)
	//
	return
}

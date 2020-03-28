package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

var me string // name of user that run the service
var myPrivateKey *rsa.PrivateKey
var partnerPublicKey map[string]*rsa.PublicKey

func main() {

	name := flag.String("name", "", "name of user running the app")
	contacts := flag.String("contacts", "", "list of user in contact list separated by comma")
	port := flag.String("port", "8080", "port assigned for the app")
	flag.Parse()

	me = *name
	myPrivateKey = privateKeyReader(*name) // read the user private key
	contactList := strings.Split(*contacts, ",")
	partnerPublicKey = publicKeyReader(contactList) // read the partners public key

	r := gin.Default()
	r.POST("/send", sendData)
	r.POST("/receive", receiveData)
	r.Run(fmt.Sprintf(":%s", *port))
}

func privateKeyReader(name string) (privateKey *rsa.PrivateKey) {
	fileName := fmt.Sprintf("%s-private.pem", me)
	privPEM, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	block, _ := pem.Decode(privPEM)

	if block == nil {
		log.Fatalln("failed to parse PEM block containing the key")
	}

	privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Fatalln(err)
	}
	return
}

func publicKeyReader(names []string) (publicKeys map[string]*rsa.PublicKey) {
	publicKeys = make(map[string]*rsa.PublicKey)

	for _, v := range names {
		fileName := fmt.Sprintf("%s-public.pem", v)

		pubPEM, _ := ioutil.ReadFile(fileName)

		block, _ := pem.Decode(pubPEM)

		if block == nil {
			log.Fatalln("failed to parse PEM block containing the key")
		}

		pub, _ := x509.ParsePKIXPublicKey(block.Bytes)

		switch pub := pub.(type) {
		case *rsa.PublicKey:
			publicKeys[v] = pub
		default:
			log.Fatalln("failed to parse the file to Public Key")
		}
	}
	return
}

func encryptor(plaintext string, pubKey *rsa.PublicKey) (ciphertext string) {

	hash := sha256.New()
	ciphertextByte, _ := rsa.EncryptOAEP(
		hash,
		rand.Reader,
		pubKey,
		[]byte(plaintext),
		[]byte(""),
	)

	ciphertext = base64.StdEncoding.EncodeToString(ciphertextByte)

	return
}

func decipher(ciphertext string, pivKey *rsa.PrivateKey) (plaintext string) {
	ciphertextByte, _ := base64.StdEncoding.DecodeString(ciphertext)

	hash := sha256.New()
	plaintextByte, _ := rsa.DecryptOAEP(
		hash,
		rand.Reader,
		pivKey,
		ciphertextByte,
		[]byte(""),
	)

	plaintext = string(plaintextByte)
	return
}

func signatureWriter(ciphertext string, pivKey *rsa.PrivateKey) (signature string) {
	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto

	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write([]byte(ciphertext))
	hashed := pssh.Sum(nil)
	signatureByte, _ := rsa.SignPSS(
		rand.Reader,
		pivKey,
		newhash,
		hashed,
		&opts,
	)
	signature = string(signatureByte)
	return
}

func signatureVerifier(message, signature string, pubKey *rsa.PublicKey) (err error) {
	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto

	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write([]byte(message))
	hashed := pssh.Sum(nil)
	err = rsa.VerifyPSS(
		pubKey,
		newhash,
		hashed,
		[]byte(signature),
		&opts,
	)
	return
}

func sendData(c *gin.Context) {

	receiver := c.PostForm("receiver")
	host := c.PostForm("host")
	message := c.PostForm("message")

	pubKey := partnerPublicKey[receiver] // get the right public key for receiver

	ciphertext := encryptor(message, pubKey) // encrypt the message into ciphertext

	signature := signatureWriter(ciphertext, myPrivateKey) // write a signature

	formData := url.Values{
		"sender":    {me},
		"message":   {ciphertext},
		"signature": {signature},
	}

	response, _ := http.PostForm(host, formData) // send the ciphertext to target service

	c.JSON(
		response.StatusCode,
		gin.H{"ciphertext": ciphertext}) // response back to the sender

	return
}

func receiveData(c *gin.Context) {
	sender := c.PostForm("sender")
	message := c.PostForm("message")
	signature := c.PostForm("signature")

	pubKey := partnerPublicKey[sender] // get the right public key matching the sender

	errVerify := signatureVerifier(message, signature, pubKey) // verify the signature is come from the right sender

	if errVerify != nil {
		log.Printf("\n[new message - not verified]\n%s : %s\n", sender, message)
	} else {
		plaintext := decipher(message, myPrivateKey)                           // decrypt the message
		log.Printf("\n[new message - verified]\n%s : %s\n", sender, plaintext) // print the plaintext if to signature verified
	}

	c.JSON(http.StatusOK,
		gin.H{"success": "OK"})
	return
}

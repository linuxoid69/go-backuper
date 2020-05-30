package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"encoding/base64"
	"errors"
	"encoding/hex"
	"io"
	"github.com/linuxoid69/go-backuper/config"

)

func encryptText(password string, text []byte) ([]byte, error) {
	key := passwordToHash(password)
	block, err := aes.NewCipher(key)

	if err != nil {
		return nil, err
	}

	b := base64.StdEncoding.EncodeToString(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))

	return ciphertext, nil
}

func decryptText(password string, text []byte) ([]byte, error) {
	key := passwordToHash(password)
	block, err := aes.NewCipher(key)

	if err != nil {
		return nil, err
	}

	if len(text) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := base64.StdEncoding.DecodeString(string(text))

	if err != nil {
		return nil, err
	}

	return data, nil
}

// EncryptFile - function encrypting a file
func EncryptFile(password string, filename string) error {
	data, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Println(err)
		return err
	}

	edata, err := encryptText(password, data)

	if err != nil {
		log.Println(err)
		return err
	}

	newFileName := fmt.Sprintf("%v.enc", filename)
	fi, err := os.Stat(filename)
	if err != nil {
		log.Println(err)
		return err
	}

	ioutil.WriteFile(newFileName, edata, fi.Mode().Perm())

	return nil
}

// DecryptFile - function encrypting a file
func DencryptFile(password string, filename string) error {

	data, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Println(err)
		return err
	}

	ddata, err := decryptText(password, data)

	if err != nil {
		return err
	}

	newFileName := strings.Split(filename, ".enc")

	fi, err := os.Stat(filename)

	if err != nil {
		log.Println(err)
		return err
	}

	ioutil.WriteFile(newFileName[0], ddata, fi.Mode().Perm())

	return nil
}

func passwordToHash(password string) []byte {
	hasher := md5.New()
	hasher.Write([]byte(password))
	return []byte(hex.EncodeToString(hasher.Sum(nil)))[:16]
}

// RunEncrypt - running a process of backup
func RunEncrypt(fileList []string, cfg *config.Config ){
	if (cfg.EncryptBackup == true && cfg.EncryptPassword != ""){
		for _, f := range fileList {
			err := EncryptFile(cfg.EncryptPassword, f)
			if err != nil{
				log.Printf("File %s wasn't encrypted. %v", f, err)
			}
		}
	} else if (cfg.EncryptBackup == true && cfg.EncryptPassword == "") {
		log.Println("You should set a password")
	} else if (cfg.EncryptBackup == false && cfg.EncryptPassword != ""){
		log.Println("You have to enable an option `encrypt_backup` in you config file")
	}
}

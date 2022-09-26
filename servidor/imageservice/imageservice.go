package imageservice

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
)

type ImageService struct {
	path string
}

func (i *ImageService) Open() error {
	_, err := os.Stat(os.Getenv("ImagePath"))

	if os.IsNotExist(err) {
		return errors.New("The path does not exist")
	} else if err != nil {
		return errors.New("Internal error")
	}

	i.path = os.Getenv("jwtSecret")
	return nil
}

func (i *ImageService) Check(hash string) (bool, error) {
	_, err := os.ReadFile(i.path + "/" + hash + ".png")
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, nil
	}
	return true, nil
}

func (i *ImageService) Delete(hash string) (bool, error) {
	err := os.Remove(i.path + "/" + hash + ".png")
	if err == os.ErrNotExist {
		return true, errors.New("The image does not exist")
	} else if err != nil {
		return false, errors.New("Internal error")
	}
	return true, nil
}

func (i *ImageService) Create(img string) (string, error) {

	hash := sha256.Sum256([]byte(img))

	err := os.WriteFile(i.path+"/"+hex.EncodeToString(hash[:])+".png", []byte(img), 0666)
	if err != nil {
		return "", errors.New("Internal error")
	}

	return hex.EncodeToString(hash[:]), nil
}

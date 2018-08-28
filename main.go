package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	valid "github.com/asaskevich/govalidator"
	"github.com/coreos/bbolt"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

var boltDBPath = "/boltdb-data/shortenedURL.db"
var shortUrlBkt = []byte("shortUrlBkt")
var seedChars = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var seedCharsLen = len(seedChars)
var aChar byte = 97
var dbConn *bolt.DB

type Response struct {
	Status int    `json:"status,omitempty"`
	Msg    string `json:"msg,omitempty"`
	Url    string `json:"url,omitempty"`
	Short  string `json:"short,omitempty"`
}

func main() {
	var err error
	dbConn, err = bolt.Open(boltDBPath, 0644, nil)
	if err != nil {
		log.Println(err)
	}

	logFile, err := os.OpenFile("server.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	logrus.SetOutput(logFile)

	router := httprouter.New()
	router.GET("/:code", Redirect)
	router.POST("/url", Create)
	router.POST("/url/", Create)
	router.NotFound = http.HandlerFunc(notFound)
	logrus.Info(http.ListenAndServe(":8080", router))
}

func notFound(w http.ResponseWriter, r *http.Request) {
	var resp *Response
	if r.URL.Path != "/url" && r.URL.Path != "/:code" {
		logrus.WithFields(logrus.Fields{
			"Method":         r.Method,
			"ResponseStatus": http.StatusNotFound,
			"Endpoint":       r.URL.Path,
		}).Error("Unknown endpoint")
		resp = &Response{Status: http.StatusNotFound, Msg: "Unknown endpoint", Url: ""}
		respJson, _ := json.Marshal(resp)
		fmt.Fprint(w, string(respJson))
	}
}

func Create(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var body Response
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	host := r.Host
	baseUrl := host + "/"
	if err != nil {
		panic(err)
	}

	if valid.IsURL(body.Url) == false {
		logrus.WithFields(logrus.Fields{
			"Method":         r.Method,
			"ResponseStatus": http.StatusUnprocessableEntity,
			"URLToShorten":   body.Url,
		}).Error("Invalid input")
		resp := &Response{Status: http.StatusUnprocessableEntity, Msg: "Invalid input", Url: body.Url}
		respJson, _ := json.Marshal(resp)
		fmt.Fprint(w, string(respJson))
		return
	}

	newCode, err := GetNextCode()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Method":         r.Method,
			"Error":          err,
			"ResponseStatus": http.StatusInternalServerError,
			"URLToShorten":   body.Url,
		}).Error("Some error occurred while shortening the URL")
		resp := Response{Status: http.StatusInternalServerError, Msg: "Some error occurred while creating short URL", Url: ""}
		respJson, _ := json.Marshal(resp)
		fmt.Fprint(w, string(respJson))
	}

	byteKey, byteUrl := []byte(newCode), []byte(body.Url)
	err = dbConn.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(shortUrlBkt)
		if err != nil {
			return err
		}

		err = bucket.Put(byteKey, byteUrl)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Method":         r.Method,
			"Error":          err,
			"ResponseStatus": http.StatusInternalServerError,
			"URLToShorten":   body.Url,
		}).Error("Some error occurred while shortening the URL")
		resp := &Response{Status: http.StatusInternalServerError, Msg: "Some error occurred while creating short URL:", Url: body.Url}
		respJson, _ := json.Marshal(resp)
		fmt.Fprint(w, string(respJson))
		return
	}

	shortUrl := baseUrl + newCode
	resp := &Response{Url: body.Url, Short: shortUrl}
	respJson, _ := json.Marshal(resp)
	fmt.Fprint(w, string(respJson))

	logrus.WithFields(logrus.Fields{
		"Method":         r.Method,
		"ResponseStatus": http.StatusOK,
		"URLToShorten":   body.Url,
		"ShortenedURL":   shortUrl,
	}).Info("Successfully Shortened URL")
}

func Redirect(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var resp *Response
	code := ps.ByName("code")
	originalUrl, err := getCodeURL(code)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Method":         r.Method,
			"Error":          "URL not found",
			"ResponseStatus": http.StatusNotFound,
		}).Error()
		resp = &Response{Status: http.StatusNotFound, Msg: "URL not found", Url: ""}
		respJson, _ := json.Marshal(resp)
		fmt.Fprint(w, string(respJson))
		return
	}
	if len(originalUrl) != 0 {
		logrus.WithFields(logrus.Fields{
			"Method":         r.Method,
			"ResponseStatus": http.StatusFound,
			"ShortenedURL":   r.Host + r.RequestURI,
			"RemoteAddress":  r.RemoteAddr,
			"URLToRedirect":  originalUrl,
		}).Info("Redirecting shortened URL to Original URL")
		http.Redirect(w, r, originalUrl, http.StatusMovedPermanently)
	} else {
		logrus.WithFields(logrus.Fields{
			"Method":         r.Method,
			"Error":          "URL not found",
			"ResponseStatus": http.StatusNotFound,
		}).Error()
		resp = &Response{Status: http.StatusNotFound, Msg: "URL not found", Url: ""}
		respJson, _ := json.Marshal(resp)
		fmt.Fprint(w, string(respJson))
		return
	}
	respJson, err := json.Marshal(resp)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Method": r.Method,
			"Error":  err,
		}).Error("Error occurred while creating json response")
		fmt.Fprint(w, "Error occurred while creating json response")
		return
	}
	fmt.Fprint(w, string(respJson))
	logrus.WithFields(logrus.Fields{
		"Method":         r.Method,
		"ResponseStatus": http.StatusOK,
		"URLToRedirect":  originalUrl,
	}).Info("Successfully Redirected URL")
}

func getCodeURL(code string) (string, error) {
	key := []byte(code)
	var originalUrl string

	err := dbConn.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(shortUrlBkt)
		if bucket == nil {
			logrus.WithFields(logrus.Fields{
				"BucketName":    shortUrlBkt,
				"URLToRedirect": originalUrl,
			}).Error("Bucket not found")
			return fmt.Errorf("Bucket %q not found!", shortUrlBkt)
		}

		value := bucket.Get(key)
		originalUrl = string(value)
		fmt.Println(originalUrl)
		return nil
	})

	if err != nil {
		return "", err
	}
	return originalUrl, nil
}

func GetNextCode() (string, error) {
	var newCode string
	err := dbConn.Update(func(tx *bolt.Tx) error {
		// By using locking on db file BoltDB makes sure it will be thread safe operation
		// and no two goroutines can can get same a short code at a time
		bucket, err := tx.CreateBucketIfNotExists(shortUrlBkt)
		if err != nil {
			return err
		}

		existingCodeByteKey := []byte("existingCodeKey")
		existingCode := bucket.Get(existingCodeByteKey)
		newCode, err = GenerateNextCode(string(existingCode))
		if err != nil {
			return err
		}

		err = bucket.Put(existingCodeByteKey, []byte(newCode))
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return "", err
	}
	return newCode, nil
}

/*
	Following method is used to generate alphanumeric incremental code, which will be helpful
	for generating short urls
	this function will return new code like, input > output
	a > b, ax > ay, az > aA, aZ > a1, a9 > ba, 99 > aaa
	it will create shortest alphanumeric code possible for using in url
*/
func GenerateNextCode(code string) (string, error) {
	if code == "" {
		return string(aChar), nil
	}
	codeBytes := []byte(code)
	codeByteLen := len(codeBytes)

	codeCharIndex := -1
	for i := (codeByteLen - 1); i >= 0; i-- {
		codeCharIndex = bytes.IndexByte(seedChars, codeBytes[i])
		if codeCharIndex == -1 || codeCharIndex >= seedCharsLen {
			return "", errors.New("Invalid Exisitng Code.")
		} else if codeCharIndex == (seedCharsLen - 1) {
			codeBytes[i] = aChar
		} else {
			codeBytes[i] = seedChars[(codeCharIndex + 1)]
			return string(codeBytes), nil
		}
	}
	for _, byteVal := range codeBytes {
		if byteVal != aChar {
			return string(codeBytes), nil
		}
	}
	// Prepending "a" for generating new incremental code
	return "a" + string(codeBytes), nil
}


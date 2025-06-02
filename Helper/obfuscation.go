package Helper

import (
	"math/rand"
	"strings"
)

// This simple helper is created to obfuscate sensitive information before passing it to any LLM
// Perhaps it makes more sense to call this obfuscation map?
// An encoding map is created to be able to translate the encoded data back to the user

type EncodeMap struct {
	fromWord string
	toWord   string
}

// randStringBytes will generate a random string of n length
func randStringBytes(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// MapWordEncoding function to obfuscate real emails to the LLM
// Return encoded emails
func MapWordEncoding(inputWord []string, encodedMap *[]EncodeMap) []string {

	newEncodedEmails := []string{}

	for _, email := range inputWord {
		if email == "" || len(email) < 3 {
			continue
		}

		mappedEmail := ""
		// Check if email is already encoded
		for _, emailMapped := range *encodedMap {
			if email == emailMapped.fromWord {
				mappedEmail = emailMapped.toWord
				break
			}
		}

		if mappedEmail == "" {
			var newEmail string

			if strings.Contains(email, "@") {
				newEmail = email[0:2] + randStringBytes(3) + "@redacted.com"
			} else {
				newEmail = "[REDACTED " + randStringBytes(4) + "]"
			}

			*encodedMap = append(*encodedMap, EncodeMap{
				fromWord: email,
				toWord:   newEmail,
			})
			mappedEmail = newEmail
		}

		newEncodedEmails = append(newEncodedEmails, mappedEmail)
	}

	return newEncodedEmails
}

// EncodeWordsInString based on the generated encoded map, find for words in a string to encode
func EncodeWordsInString(inputText string, encodedEmails *[]EncodeMap) string {
	for _, emailObj := range *encodedEmails {
		inputText = strings.Replace(inputText, emailObj.fromWord, emailObj.toWord, -1)
	}

	return inputText
}

// DecodeWordsInString based on the generated encoded map, find for words in a string to decode
func DecodeWordsInString(inputText string, encodedEmails *[]EncodeMap) string {
	for _, emailObj := range *encodedEmails {
		inputText = strings.Replace(inputText, emailObj.toWord, emailObj.fromWord, -1)
	}

	return inputText
}

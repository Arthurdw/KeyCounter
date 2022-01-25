package main

import (
	"github.com/MarinX/keylogger"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const KeysFileBufferSize = 100
const CombinationsFileBufferSize = 10
const FileRoot = "/etc/key-counter"

func readData(file string) map[string]uint64 {
	data, err := os.ReadFile(filepath.Join(FileRoot, file))
	response := map[string]uint64{}

	if err != nil {
		return response
	}

	for _, line := range strings.Split(string(data), "\n") {
		content := strings.Split(line, ",")

		if len(content) != 2 {
			continue
		}

		key := content[0]
		value, _ := strconv.ParseInt(content[1], 10, 64)

		response[key] = uint64(value)
	}
	return response
}

func writeData(fileName string, response map[string]uint64) {
	logrus.Info("Writing data to file")
	file, err := os.Create(filepath.Join(FileRoot, fileName))

	if err != nil {
		logrus.Error(err)
		panic(err)
	}

	defer file.Close()

	for key, value := range response {
		_, err := file.WriteString(key + "," + strconv.FormatInt(int64(value), 10) + "\n")

		if err != nil {
			logrus.Error(err)
			panic(err)
		}
	}
}

func main() {
	keyboard := keylogger.FindKeyboardDevice()
	logrus.SetLevel(logrus.DebugLevel)

	if len(keyboard) == 0 {
		logrus.Error("No keyboard found?")
		return
	}

	if _, err := os.Stat(FileRoot); os.IsNotExist(err) {
		os.Mkdir(FileRoot, 777)
	}

	k, err := keylogger.New(keyboard)
	if err != nil {
		logrus.Error(err)
		return
	}

	defer k.Close()

	pressedKeys := map[string]struct{}{}
	keyCache := readData("./data.csv")
	combinationCache := readData("./combinations.csv")
	lastActionWasUp := false

	// File writing buffers (to prevent file write spamming)
	keysWriteMissedBuffer := 0
	combinationWriteMissedBuffer := 0

	for e := range k.Read() {
		switch e.Type {
		case keylogger.EvKey:
			if !e.KeyPress() && !e.KeyRelease() {
				break
			}

			keyString := e.KeyString()

			if keyString == "," {
				// We have to describe the comma like this for simplifying the csv parsing.
				keyString = "Comma"
			}

			if e.KeyPress() {
				_, ok := pressedKeys[keyString]
				lastActionWasUp = false

				if !ok {
					pressedKeys[keyString] = struct{}{}
				}

				_, ok = keyCache[keyString]
				if ok {
					keyCache[keyString]++
				} else {
					keyCache[keyString] = 1
				}

			} else {
				if !lastActionWasUp && len(pressedKeys) > 1 {
					keys := make([]string, 0, len(pressedKeys))
					for key := range pressedKeys {
						keys = append(keys, key)
					}
					sort.Strings(keys)

					combinationRepresentation := strings.Join(keys, "+")

					_, ok := combinationCache[combinationRepresentation]
					if ok {
						combinationCache[combinationRepresentation]++
					} else {
						combinationCache[combinationRepresentation] = 1
					}

					if combinationWriteMissedBuffer > CombinationsFileBufferSize {
						writeData("./combinations.csv", combinationCache)
						combinationWriteMissedBuffer = 0
					} else {
						combinationWriteMissedBuffer++
					}

				}
				delete(pressedKeys, keyString)
				lastActionWasUp = true
			}

			if keysWriteMissedBuffer > KeysFileBufferSize {
				writeData("./data.csv", keyCache)
				keysWriteMissedBuffer = 0
			} else {
				keysWriteMissedBuffer++
			}
			break
		}
	}
}

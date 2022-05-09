package Utils

import (
	"encoding/gob"
	"encoding/json"
	"log"
	"os"
	"os/exec"
)

func StringInSlice(a string, list []string) bool {
	for _, stringInList := range list {
		if stringInList == a {
			return true
		}
	}
	return false
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func SerializeObject(object interface{}, file_name string) {
	// Create a file for IO
	encodeFile, err := os.Create(file_name)
	if err != nil {
		log.Panic("Error creating file for serialization:", err)
	}
	// Since this is a binary format large parts of it will be unreadable
	encoder := gob.NewEncoder(encodeFile)
	// Write to the file
	if err := encoder.Encode(object); err != nil {
		log.Panic("Error writing serilaized object to file:", err)
	}
	encodeFile.Close()
}

func DeserializeObject(object interface{}, file_name string) {
	// Open a RO file
	decodeFile, err := os.Open(file_name)
	if err != nil {
		log.Panic("Error opening file of serialized object:", err)
	}
	defer decodeFile.Close()
	// Create a decoder
	decoder := gob.NewDecoder(decodeFile)
	// Decode -- We need to pass a pointer otherwise telegram_ids isn't modified
	decoder.Decode(object)
}

func Execute_system_command(program_name string, parameters ...string) (map[string]interface{}, error) {
	cmd := exec.Command(program_name, parameters...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	jsonMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(out), &jsonMap)
	if err != nil {
		return nil, err
	}

	return jsonMap, nil
}

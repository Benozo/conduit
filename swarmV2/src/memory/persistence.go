package memory

import (
	"encoding/json"
	"os"
)

// MemoryData represents the structure of the data to be persisted.
type MemoryData struct {
	SharedData map[string]interface{} `json:"shared_data"`
	ContextData map[string]interface{} `json:"context_data"`
}

// SaveMemory persists the memory data to a file.
func SaveMemory(filePath string, data MemoryData) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(data)
}

// LoadMemory loads the memory data from a file.
func LoadMemory(filePath string) (MemoryData, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return MemoryData{}, err
	}
	defer file.Close()

	var data MemoryData
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	return data, err
}
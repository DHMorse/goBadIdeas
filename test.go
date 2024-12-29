package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"os"
)

func init() {
	gob.Register(map[string]interface{}{})
	gob.Register([]interface{}{})
	gob.Register(string(""))
	gob.Register(float64(0))
	gob.Register(int(0))
	gob.Register(bool(false))
}

func jsonToBinary(jsonFile, binaryFile string) error {
	// Read the JSON file
	jsonData, err := os.ReadFile(jsonFile)
	if err != nil {
		return fmt.Errorf("error reading JSON file: %v", err)
	}

	// Create a map to store the JSON data
	var data interface{}

	// Unmarshal JSON into the map
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	// Create the binary file
	bf, err := os.Create(binaryFile)
	if err != nil {
		return fmt.Errorf("error creating binary file: %v", err)
	}
	defer bf.Close()

	// Create a new encoder and encode the data
	encoder := gob.NewEncoder(bf)
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("error encoding to binary: %v", err)
	}

	fmt.Printf("Successfully converted %s to %s\n", jsonFile, binaryFile)
	return nil
}

func binaryToJson(binaryFile, jsonFile string) error {
	// Open the binary file
	bf, err := os.Open(binaryFile)
	if err != nil {
		return fmt.Errorf("error opening binary file: %v", err)
	}
	defer bf.Close()

	// Create a variable to hold the decoded data as a map
	var data map[string]interface{}

	// Create a decoder and decode the binary data
	decoder := gob.NewDecoder(bf)
	if err := decoder.Decode(&data); err != nil {
		return fmt.Errorf("error decoding binary file: %v", err)
	}

	// Marshal the data into JSON format
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling data to JSON: %v", err)
	}

	// Write the JSON data to the output file
	if err := os.WriteFile(jsonFile, jsonData, 0644); err != nil {
		return fmt.Errorf("error writing JSON file: %v", err)
	}

	fmt.Printf("Successfully converted %s to %s\n", binaryFile, jsonFile)
	return nil
}

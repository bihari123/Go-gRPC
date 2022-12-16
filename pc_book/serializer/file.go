package serializer

import (
	"fmt"
	"io/ioutil"

	"github.com/golang/protobuf/proto"
)

func WriteProtoBufToBinary(message proto.Message, fileName string) error {
	data, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("Cannot marshal proto message to binary: %w", err)
	}
	err = ioutil.WriteFile(fileName, data, 0644)
	if err != nil {
		return fmt.Errorf("cannot write binary data to file: %w", err)
	}

	return nil
}

func ReadProtoBufFromBinaryFile(fileName string, message proto.Message) error {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("Error reading the file %s : %w", fileName, err)
	}

	err = proto.Unmarshal(data, message)
	if err != nil {
		return fmt.Errorf("Error unmarshalling %s : %w", fileName, err)
	}

	return nil
}

func WriteProtobufToJSONFile(message proto.Message, fileName string) error {
	data, err := ProtobufToJson(message)
	if err != nil {
		return fmt.Errorf("Error converting into JSON: %w", err)
	}

	err = ioutil.WriteFile(fileName, []byte(data), 0644)

	if err != nil {
		return fmt.Errorf("Error writing to the file: %w ", err)
	}
	return nil
}

package serializer_test

import (
	pcbook "pcbook/proto"
	"pcbook/sample"
	"pcbook/serializer"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/require"
)

func TestFileSerializer(t *testing.T) {
	t.Parallel()

	binaryFile := "../tmp/laptop.bin"

	laptop1 := sample.NewLaptop()

	err := serializer.WriteProtoBufToBinary(laptop1, binaryFile)

	require.NoError(t, err)

	laptop2 := &pcbook.Laptop{}

	err = serializer.ReadProtoBufFromBinaryFile(binaryFile, laptop2)
	require.NoError(t, err)

	require.True(t, proto.Equal(laptop1, laptop2))

	jsonFile := "../tmp/laptop.json"
	err = serializer.WriteProtobufToJSONFile(laptop2, jsonFile)
	require.NoError(t, err)
}

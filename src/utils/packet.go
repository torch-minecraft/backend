package utils

import (
	"bytes"
	"io"
)

func WritePacket(data *bytes.Buffer, writer io.Writer) error {
	if _, err := WriteVarInt(int32(data.Len()), writer); err != nil {
		return err
	}
	_, err := io.Copy(writer, data)
	return err
}

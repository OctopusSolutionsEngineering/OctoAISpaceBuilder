package compress

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"go.uber.org/zap"
	"io"
)

func CompressByteArray(data []byte) (string, error) {
	// Create a buffer to hold the compressed data
	var compressedBuffer bytes.Buffer

	// Create a gzip writer
	gzipWriter, err := gzip.NewWriterLevel(&compressedBuffer, gzip.BestCompression)

	if err != nil {
		return "", err
	}

	// Write the data to the gzip writer
	_, err = gzipWriter.Write(data)

	if err != nil {
		return "", err
	}

	// Close the writer to flush any pending data
	if err := gzipWriter.Close(); err != nil {
		return "", err
	}

	// Return base64 encoded compressed data
	return base64.StdEncoding.EncodeToString(compressedBuffer.Bytes()), nil
}

func CompressString(data string) (string, error) {
	return CompressByteArray([]byte(data))
}

func DecompressString(compressedData string) (string, error) {
	byteArray, err := DecompressStringToByteArray(compressedData)

	if err != nil {
		return "", err
	}

	return string(byteArray), nil
}

func DecompressStringToByteArray(compressedData string) ([]byte, error) {
	// Decode base64 string
	decoded, err := base64.StdEncoding.DecodeString(compressedData)
	if err != nil {
		return nil, err
	}

	// Create a reader for the compressed data
	compressedReader := bytes.NewReader(decoded)

	// Create a gzip reader
	gzipReader, err := gzip.NewReader(compressedReader)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := gzipReader.Close(); err != nil {
			zap.L().Error("failed to close gzip reader", zap.Error(err))
		}
	}()

	// Read the decompressed data
	decompressed, err := io.ReadAll(gzipReader)
	if err != nil {
		return nil, err
	}

	return decompressed, nil
}

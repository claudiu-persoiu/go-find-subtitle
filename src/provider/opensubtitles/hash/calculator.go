package hash

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

const (
	ChunkSize = 65536 // 64k
)

func BuildHash(path string) (string, error) {

	f, err := os.Open(path)

	if err != nil {
		return "", err
	}

	defer f.Close()

	hash, err := hashFile(f)

	if err != nil {
		return "", err
	}

	return hash, nil
}

// HashFile Generate an OSDB hash for an *os.File.
// The hash was taken from here: https://trac.opensubtitles.org/projects/opensubtitles/wiki/HashSourceCodes#GO
// which in term was referencing: https://github.com/oz/osdb/blob/master/osdb.go
func hashFile(file *os.File) (hash string, err error) {
	fi, err := file.Stat()
	if err != nil {
		return
	}
	if fi.Size() < ChunkSize {
		return "", fmt.Errorf("File is too small")
	}

	// Read head and tail blocks.
	buf := make([]byte, ChunkSize*2)
	err = readChunk(file, 0, buf[:ChunkSize])
	if err != nil {
		return
	}
	err = readChunk(file, fi.Size()-ChunkSize, buf[ChunkSize:])
	if err != nil {
		return
	}

	// Convert to uint64, and sum.
	var nums [(ChunkSize * 2) / 8]uint64
	reader := bytes.NewReader(buf)
	err = binary.Read(reader, binary.LittleEndian, &nums)
	if err != nil {
		return
	}
	var h uint64
	for _, num := range nums {
		h += num
	}

	hash = fmt.Sprintf("%x", h+uint64(fi.Size()))

	return
}

// Read a chunk of a file at `offset` to fill `buf`.
func readChunk(file *os.File, offset int64, buf []byte) (err error) {
	n, err := file.ReadAt(buf, offset)
	if err != nil {
		return
	}
	if n != ChunkSize {
		return fmt.Errorf("Invalid read %v", n)
	}
	return
}

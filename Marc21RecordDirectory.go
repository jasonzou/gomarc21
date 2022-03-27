package gomarc21

import (
	"errors"
	"fmt"
	"strconv"
)

func (dirEntry DirectoryEntry) GetRaw() string {
	// need to make sure 12 bytes
	return string(dirEntry.raw)
}
func (dirEntry DirectoryEntry) GetFieldLength() int {
	return dirEntry.FieldLength
}
func (dirEntry DirectoryEntry) GetStartingPosition() int {
	return dirEntry.StartingPosition
}
func (dirEntry DirectoryEntry) GetTag() string {
	return dirEntry.Tag.GetTag()
}
func NewDirectoryEntry(entryBytes []byte) (DirectoryEntry, error) {
	if len(entryBytes) != 12 {
		msg := fmt.Sprintf("The length of directry entry is not 12, but %d instead", len(entryBytes))
		err := errors.New(msg)
		return DirectoryEntry{}, err
	}

	tag, err := NewTag(entryBytes[0:3])
	if err != nil {
		return DirectoryEntry{}, errors.New("tag is invalid")
	}

	fieldLength, err := strconv.Atoi(string(entryBytes[3:7]))
	if err != nil {
		return DirectoryEntry{}, errors.New("fieldLength is invalid")
	}

	startingPosition, err := strconv.Atoi(string(entryBytes[7:12]))
	if err != nil {
		return DirectoryEntry{}, errors.New("starting position is invalid")
	}

	directoryEntry := DirectoryEntry{
		raw:              entryBytes[:12],
		Tag:              tag,
		FieldLength:      fieldLength,
		StartingPosition: startingPosition,
	}

	return directoryEntry, nil
}

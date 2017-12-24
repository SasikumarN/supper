package media

import (
	"errors"
	"fmt"
	"os"
	"io"
	"path/filepath"

	"github.com/Tympanix/supper/parse"
	"github.com/Tympanix/supper/types"
)

type File struct {
	os.FileInfo
	types.Media
	path string
}

func (f *File) Path() string {
	return f.path
}

// SaveSubtitle saves the subtitle for the given media to disk
func (f *File) SaveSubtitle(s types.Subtitle) error {
	srt, err := s.Download()
	defer srt.Close()

	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("Filename: %s\n", f.FileInfo.Name())
	filename := f.Path()
	extension := filepath.Ext(filename)
	name := filename[0:len(filename)-len(extension)] + ".srt"

	file, err := os.Create(name)

	if err != nil {
		return err
	}

	_, err = io.Copy(file, srt)

	if err != nil {
		return err
	}

	return nil
}

func NewFile(file os.FileInfo, meta types.Metadata, path string) *File {
	return &File{file, NewType(meta), path}
}

type Type struct {
	types.Metadata
}

func (m *Type) Meta() types.Metadata {
	return m.Metadata
}

func (m *Type) TypeMovie() (r types.Movie, ok bool) {
	r, ok = m.Metadata.(types.Movie)
	return
}

func (m *Type) TypeEpisode() (r types.Episode, ok bool) {
	r, ok = m.Metadata.(types.Episode)
	return
}

func NewType(m types.Metadata) *Type {
	return &Type{m}
}

// New parses a file into media attributes
func New(path string) (types.LocalMedia, error) {
	fmt.Println(path)
	filename := parse.Filename(path)
	fmt.Println(filename)
	media, err := NewMetadata(filename)
	fmt.Println(media)

	if err != nil {
		return nil, err
	}

	file, err := os.Stat(path)

	if err != nil {
		return nil, err
	}

	if movie, ok := media.(types.Movie); ok {
		return NewFile(file, movie, path), nil
	} else if episode, ok := media.(types.Episode); ok {
		return NewFile(file, episode, path), nil
	} else {
		return nil, errors.New("Unknown media type")
	}
}

// NewMetadata returns a metadata object parsed from the string
func NewMetadata(str string) (types.Metadata, error) {
	if episodeRegexp.MatchString(str) {
		return NewEpisode(str)
	}
	return NewMovie(str)
}

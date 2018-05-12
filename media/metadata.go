package media

import (
	"encoding/json"
	"strings"

	"github.com/tympanix/supper/meta/codec"
	"github.com/tympanix/supper/meta/misc"
	"github.com/tympanix/supper/meta/quality"
	"github.com/tympanix/supper/meta/source"
	"github.com/tympanix/supper/parse"
)

// Metadata provides release information for media
type Metadata struct {
	group   string
	codec   codec.Tag
	quality quality.Tag
	source  source.Tag
	misc    misc.List
	tags    []string
}

// ParseMetadata generates meta data from a string
func ParseMetadata(tags string) Metadata {
	return Metadata{
		group:   parse.Group(tags),
		codec:   parse.Codec(tags),
		quality: parse.Quality(tags),
		source:  parse.Source(tags),
		tags:    parse.Tags(tags),
		misc:    parse.Miscellaneous(tags),
	}
}

// MarshalJSON returns a JSON representation of metadata
func (m Metadata) MarshalJSON() (b []byte, err error) {
	return json.Marshal(struct {
		Group   string `json:"group"`
		Codec   string `json:"codec"`
		Quality string `json:"quality"`
		Source  string `json:"source"`
	}{
		m.group,
		m.codec.String(),
		m.quality.String(),
		m.source.String(),
	})
}

// String return a description of the metadata
func (m Metadata) String() string {
	return strings.Join([]string{
		m.Group(),
		m.Codec().String(),
		m.Quality().String(),
		m.Source().String(),
	}, ",")
}

// Group returns the release group
func (m Metadata) Group() string {
	return m.group
}

// Codec returns the codec
func (m Metadata) Codec() codec.Tag {
	return m.codec
}

// Quality returns the quality of the media
func (m Metadata) Quality() quality.Tag {
	return m.quality
}

// Source returns the source of the media
func (m Metadata) Source() source.Tag {
	return m.source
}

// AllTags returns all metadata as a list of tags
func (m Metadata) AllTags() []string {
	return m.tags
}

// Misc returns miscellaneous media tags for the media
func (m Metadata) Misc() misc.List {
	return m.misc
}

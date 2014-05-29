package container

import (
	"io"

	"github.com/clipperhouse/typewriter"
)

func init() {
	err := typewriter.Register(NewContainerWriter())
	if err != nil {
		panic(err)
	}
}

type ContainerWriter struct {
	tagsByType map[string]typewriter.Tag // typewriter.Type is not comparable, key by .String()
}

func NewContainerWriter() ContainerWriter {
	return ContainerWriter{
		tagsByType: make(map[string]typewriter.Tag),
	}
}

func (c ContainerWriter) Name() string {
	return "sorted_container"
}

func (c ContainerWriter) Validate(t typewriter.Type) (bool, error) {
	tag, found, err := t.Tags.ByName("containers")
	if found && err == nil {
		c.tagsByType[t.String()] = tag
	}
	return found, err
}

func (c ContainerWriter) WriteHeader(w io.Writer, t typewriter.Type) {
	for _, s := range c.tagsByType[t.String()].Items {
		switch s {
		case "SortedSet":
			license := `// SortedSet is a modification of https://github.com/wfreeman/go-skiplist/sortedset.go
// The MIT License (MIT)
// Copyright (c) 2014 Wes Freeman (freeman.wes@gmail.com)
`
			w.Write([]byte(license))
		}
	}
}

func (c ContainerWriter) Imports(t typewriter.Type) []string {
	return []string{"math", "math/rand"}
}

func (c ContainerWriter) Write(w io.Writer, t typewriter.Type) {
	tag := c.tagsByType[t.String()] // validated above

	for _, s := range tag.Items {
		tmpl, err := containerTemplates.Get(s)
		if err != nil {
			continue
		}
		tmpl.Execute(w, t)
	}

	return
}

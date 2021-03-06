package io3mf

import (
	"io"

	"github.com/qmuntal/opc"
)

type opcRelationship struct {
	rel *opc.Relationship
}

func (o *opcRelationship) Type() string {
	return o.rel.Type
}

func (o *opcRelationship) TargetURI() string {
	return o.rel.TargetURI
}

func newRelationships(rels []*opc.Relationship) []relationship {
	pr := make([]relationship, len(rels))
	for i, r := range rels {
		pr[i] = &opcRelationship{r}
	}
	return pr
}

type opcFile struct {
	r *opc.Reader
	f *opc.File
}

func (o *opcFile) Open() (io.ReadCloser, error) {
	return o.f.Open()
}

func (o *opcFile) Name() string {
	return o.f.Name
}

func (o *opcFile) FindFileFromRel(relType string) (packageFile, bool) {
	name := findOPCFileURIFromRel(relType, o.f.Relationships)
	name = opc.ResolveRelationship(o.f.Name, name)
	return findOPCFileFromName(name, o.r)
}

func (o *opcFile) Relationships() []relationship {
	return newRelationships(o.f.Relationships)
}

type opcReader struct {
	ra   io.ReaderAt
	size int64
	r    *opc.Reader // nil until call Open.
}

func (o *opcReader) Open(f func(r io.Reader) io.ReadCloser) (err error) {
	o.r, err = opc.NewReader(o.ra, o.size)
	if f != nil {
		o.r.SetDecompressor(f)
	}
	return
}

func (o *opcReader) FindFileFromRel(relType string) (packageFile, bool) {
	name := findOPCFileURIFromRel(relType, o.r.Relationships)
	return o.FindFileFromName(name)
}

func (o *opcReader) FindFileFromName(name string) (packageFile, bool) {
	return findOPCFileFromName(name, o.r)
}

func findOPCFileFromName(name string, r *opc.Reader) (packageFile, bool) {
	for _, f := range r.Files {
		if f.Name == name {
			return &opcFile{r, f}, true
		}
	}
	return nil, false
}

func findOPCFileURIFromRel(relType string, rels []*opc.Relationship) string {
	for _, r := range rels {
		if r.Type == relType {
			return r.TargetURI
		}
	}
	return ""
}

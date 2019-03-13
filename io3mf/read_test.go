package io3mf

import (
	"bytes"
	"errors"
	"fmt"
	"image/color"
	"io"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-test/deep"
	"github.com/qmuntal/go3mf"
	"github.com/qmuntal/go3mf/mesh"
	"github.com/qmuntal/go3mf/mesh/meshinfo"
	"github.com/stretchr/testify/mock"
)

type mockRelationship struct {
	mock.Mock
}

func newMockRelationship(relType, targetURI string) *mockRelationship {
	m := new(mockRelationship)
	m.On("Type").Return(relType).Maybe()
	m.On("TargetURI").Return(targetURI).Maybe()
	return m
}

func (m *mockRelationship) Type() string {
	args := m.Called()
	return args.String(0)
}

func (m *mockRelationship) TargetURI() string {
	args := m.Called()
	return args.String(0)
}

type modelBuilder struct {
	str      strings.Builder
	hasModel bool
}

func (m *modelBuilder) withElement(s string) *modelBuilder {
	m.str.WriteString(s)
	m.str.WriteString("\n")
	return m
}

func (m *modelBuilder) addAttr(prefix, name, value string) *modelBuilder {
	if prefix != "" {
		m.str.WriteString(fmt.Sprintf(`%s:`, prefix))
	}
	if name != "" {
		m.str.WriteString(fmt.Sprintf(`%s="%s" `, name, value))
	}
	return m
}

func (m *modelBuilder) withDefaultModel() *modelBuilder {
	m.withModel("millimeter", "en-US")
	return m
}

func (m *modelBuilder) withModel(unit string, lang string) *modelBuilder {
	m.str.WriteString(`<model `)
	m.addAttr("", "unit", unit).addAttr("xml", "lang", lang)
	m.addAttr("", "xmlns", nsCoreSpec).addAttr("xmlns", "m", nsMaterialSpec).addAttr("xmlns", "p", nsProductionSpec)
	m.addAttr("xmlns", "b", nsBeamLatticeSpec).addAttr("xmlns", "s", nsSliceSpec).addAttr("", "requiredextensions", "m p b s")
	m.str.WriteString(">\n")
	m.hasModel = true
	return m
}

func (m *modelBuilder) withEncoding(encode string) *modelBuilder {
	m.str.WriteString(fmt.Sprintf(`<?xml version="1.0" encoding="%s"?>`, encode))
	m.str.WriteString("\n")
	return m
}

func (m *modelBuilder) build() *mockFile {
	if m.hasModel {
		m.str.WriteString("</model>\n")
	}
	f := new(mockFile)
	f.On("Name").Return("3d/3dmodel.model").Maybe()
	f.On("Open").Return(ioutil.NopCloser(bytes.NewBufferString(m.str.String())), nil).Maybe()
	return f
}

type mockFile struct {
	mock.Mock
}

func newMockFile(name string, relationships []relationship, thumb *mockFile, other *mockFile, openErr bool) *mockFile {
	m := new(mockFile)
	m.On("Name").Return(name).Maybe()
	m.On("Relationships").Return(relationships).Maybe()
	m.On("FindFileFromRel", relTypeThumbnail).Return(thumb, thumb != nil).Maybe()
	m.On("FindFileFromRel", mock.Anything).Return(other, other != nil).Maybe()
	var err error
	if openErr {
		err = errors.New("")
	}
	m.On("Open").Return(ioutil.NopCloser(new(bytes.Buffer)), err).Maybe()
	return m
}

func (m *mockFile) Open() (io.ReadCloser, error) {
	args := m.Called()
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func (m *mockFile) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *mockFile) FindFileFromRel(args0 string) (packageFile, bool) {
	args := m.Called(args0)
	return args.Get(0).(packageFile), args.Bool(1)
}

func (m *mockFile) Relationships() []relationship {
	args := m.Called()
	return args.Get(0).([]relationship)
}

type mockPackage struct {
	mock.Mock
}

func newMockPackage(other *mockFile) *mockPackage {
	m := new(mockPackage)
	m.On("FindFileFromRel", mock.Anything).Return(other, other != nil).Maybe()
	m.On("FindFileFromName", mock.Anything).Return(other, other != nil).Maybe()
	return m
}

func (m *mockPackage) FindFileFromName(args0 string) (packageFile, bool) {
	args := m.Called(args0)
	return args.Get(0).(packageFile), args.Bool(1)
}

func (m *mockPackage) FindFileFromRel(args0 string) (packageFile, bool) {
	args := m.Called(args0)
	return args.Get(0).(packageFile), args.Bool(1)
}

func TestReadError_Error(t *testing.T) {
	tests := []struct {
		name string
		e    *ReadError
		want string
	}{
		{"new", new(ReadError), ""},
		{"generic", &ReadError{Message: "generic error"}, "generic error"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Error(); got != tt.want {
				t.Errorf("ReadError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReader_processOPC(t *testing.T) {
	abortReader := &Reader{Model: new(go3mf.Model), r: newMockPackage(newMockFile("/a.model", nil, nil, nil, false))}
	abortReader.SetProgressCallback(callbackFalse, nil)
	thumbFile := newMockFile("/a.png", nil, nil, nil, false)
	thumbErr := newMockFile("/a.png", nil, nil, nil, true)
	tests := []struct {
		name    string
		d       *Reader
		want    *go3mf.Model
		wantErr bool
	}{
		{"noRoot", &Reader{Model: new(go3mf.Model), r: newMockPackage(nil)}, &go3mf.Model{}, true},
		{"abort", abortReader, &go3mf.Model{}, true},
		{"noRels", &Reader{Model: new(go3mf.Model), r: newMockPackage(newMockFile("/a.model", nil, nil, nil, false))}, &go3mf.Model{Path: "/a.model"}, false},
		{"withThumb", &Reader{Model: new(go3mf.Model),
			r: newMockPackage(newMockFile("/a.model", []relationship{newMockRelationship(relTypeThumbnail, "/a.png")}, thumbFile, thumbFile, false)),
		}, &go3mf.Model{
			Path:        "/a.model",
			Thumbnail:   &go3mf.Attachment{RelationshipType: relTypeThumbnail, Path: "/Metadata/thumbnail.png", Stream: new(bytes.Buffer)},
			Attachments: []*go3mf.Attachment{{RelationshipType: relTypeThumbnail, Path: "/a.png", Stream: new(bytes.Buffer)}},
		}, false},
		{"withThumbErr", &Reader{Model: new(go3mf.Model),
			r: newMockPackage(newMockFile("/a.model", []relationship{newMockRelationship(relTypeThumbnail, "/a.png")}, thumbErr, thumbErr, false)),
		}, &go3mf.Model{Path: "/a.model"}, false},
		{"withOtherRel", &Reader{Model: new(go3mf.Model),
			r: newMockPackage(newMockFile("/a.model", []relationship{newMockRelationship("other", "/a.png")}, nil, nil, false)),
		}, &go3mf.Model{Path: "/a.model"}, false},
		{"withModelAttachment", &Reader{Model: new(go3mf.Model),
			r: newMockPackage(newMockFile("/a.model", []relationship{newMockRelationship(relTypeModel3D, "/a.model")}, nil, newMockFile("/a.model", nil, nil, nil, false), false)),
		}, &go3mf.Model{
			Path:                  "/a.model",
			ProductionAttachments: []*go3mf.Attachment{{RelationshipType: relTypeModel3D, Path: "/a.model", Stream: new(bytes.Buffer)}},
		}, false},
		{"withAttRel", &Reader{Model: new(go3mf.Model), AttachmentRelations: []string{"b"},
			r: newMockPackage(newMockFile("/a.model", []relationship{newMockRelationship("b", "/a.xml")}, nil, newMockFile("/a.xml", nil, nil, nil, false), false)),
		}, &go3mf.Model{
			Path:        "/a.model",
			Attachments: []*go3mf.Attachment{{RelationshipType: "b", Path: "/a.xml", Stream: new(bytes.Buffer)}},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.d.processOPC()
			if (err != nil) != tt.wantErr {
				t.Errorf("Reader.processOPC() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := deep.Equal(tt.d.Model, tt.want); diff != nil {
				t.Errorf("Reader.processOPC() = %v", diff)
				return
			}
		})
	}
}

func TestReader_processRootModel_Fail(t *testing.T) {
	abortReader := &Reader{Model: new(go3mf.Model), r: newMockPackage(newMockFile("/a.model", nil, nil, nil, false))}
	abortReader.SetProgressCallback(callbackFalse, nil)
	tests := []struct {
		name    string
		r       *Reader
		wantErr bool
	}{
		{"noRoot", &Reader{Model: new(go3mf.Model), r: newMockPackage(nil)}, true},
		{"abort", abortReader, true},
		{"errOpen", &Reader{Model: new(go3mf.Model), r: newMockPackage(newMockFile("/a.model", nil, nil, nil, true))}, true},
		{"errEncode", &Reader{Model: new(go3mf.Model), r: newMockPackage(new(modelBuilder).withEncoding("utf16").build())}, true},
		{"invalidUnits", &Reader{Model: new(go3mf.Model), r: newMockPackage(new(modelBuilder).withModel("other", "en-US").build())}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.processRootModel(); (err != nil) != tt.wantErr {
				t.Errorf("Reader.processRootModel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestReader_processRootModel(t *testing.T) {
	want := new(go3mf.Model)
	want.Units = go3mf.UnitMillimeter
	want.Language = "en-US"
	baseMaterials := &go3mf.BaseMaterialsResource{ID: 5, ModelPath: "3d/3dmodel.model", Materials: []go3mf.BaseMaterial{
		{Name: "Blue PLA", Color: color.RGBA{0, 0, 85, 255}},
		{Name: "Red ABS", Color: color.RGBA{85, 0, 0, 255}},
	}}
	baseTexture := &go3mf.Texture2DResource{ID: 6, ModelPath: "3d/3dmodel.model", Path: "/3D/Texture/msLogo.png", ContentType: go3mf.PNGTexture, TileStyleU: go3mf.TileWrap, TileStyleV: go3mf.TileMirror, Filter: go3mf.TextureFilterAuto}
	otherSlices := &go3mf.SliceStack{
		BottomZ: 2,
		Slices: []*go3mf.Slice{
			{
				TopZ:     1.2,
				Vertices: []mgl32.Vec2{{1.01, 1.02}, {9.03, 1.04}, {9.05, 9.06}, {1.07, 9.08}},
				Polygons: [][]int{{0, 1, 2, 3, 0}},
			},
		},
	}
	sliceStack := &go3mf.SliceStackResource{ID: 3, ModelPath: "3d/3dmodel.model", SliceStack: &go3mf.SliceStack{
		BottomZ: 1,
		Slices: []*go3mf.Slice{
			{
				TopZ:     0,
				Vertices: []mgl32.Vec2{{1.01, 1.02}, {9.03, 1.04}, {9.05, 9.06}, {1.07, 9.08}},
				Polygons: [][]int{{0, 1, 2, 3, 0}},
			},
			{
				TopZ:     0.1,
				Vertices: []mgl32.Vec2{{1.01, 1.02}, {9.03, 1.04}, {9.05, 9.06}, {1.07, 9.08}},
				Polygons: [][]int{{0, 2, 1, 3, 0}},
			},
		},
	}}
	sliceStackRef := &go3mf.SliceStackResource{ID: 7, ModelPath: "3d/3dmodel.model", SliceStack: otherSlices}
	sliceStackRef.BottomZ = 1.1
	sliceStackRef.UsesSliceRef = true
	sliceStackRef.Slices = append(sliceStackRef.Slices, otherSlices.Slices...)
	meshRes := &go3mf.MeshResource{
		ObjectResource: go3mf.ObjectResource{ID: 1, Name: "Box 1", ModelPath: "3d/3dmodel.model", SliceStackID: 3, SliceResoultion: go3mf.ResolutionLow, PartNumber: "11111111-1111-1111-1111-111111111111"},
		Mesh:           mesh.NewMesh(),
	}
	meshRes.Mesh.Nodes = append(meshRes.Mesh.Nodes, []mesh.Node{
		{Index: 0, Position: mgl32.Vec3{0, 0, 0}},
		{Index: 1, Position: mgl32.Vec3{100, 0, 0}},
		{Index: 2, Position: mgl32.Vec3{100, 100, 0}},
		{Index: 3, Position: mgl32.Vec3{0, 100, 0}},
		{Index: 4, Position: mgl32.Vec3{0, 0, 100}},
		{Index: 5, Position: mgl32.Vec3{100, 0, 100}},
		{Index: 6, Position: mgl32.Vec3{100, 100, 100}},
		{Index: 7, Position: mgl32.Vec3{0, 100, 100}},
	}...)
	meshRes.Mesh.Faces = append(meshRes.Mesh.Faces, []mesh.Face{
		{Index: 0, NodeIndices: [3]uint32{3, 2, 1}},
		{Index: 1, NodeIndices: [3]uint32{1, 0, 3}},
		{Index: 2, NodeIndices: [3]uint32{4, 5, 6}},
		{Index: 3, NodeIndices: [3]uint32{6, 7, 4}},
		{Index: 4, NodeIndices: [3]uint32{0, 1, 5}},
		{Index: 5, NodeIndices: [3]uint32{5, 4, 0}},
		{Index: 6, NodeIndices: [3]uint32{1, 2, 6}},
		{Index: 7, NodeIndices: [3]uint32{6, 5, 1}},
		{Index: 8, NodeIndices: [3]uint32{2, 3, 7}},
		{Index: 9, NodeIndices: [3]uint32{7, 6, 2}},
		{Index: 10, NodeIndices: [3]uint32{3, 0, 4}},
		{Index: 11, NodeIndices: [3]uint32{4, 7, 3}},
	}...)
	handler := meshRes.Mesh.InformationHandler()
	info := handler.AddBaseMaterialInfo(uint32(len(meshRes.Mesh.Faces)))
	*info.FaceData(0).(*meshinfo.BaseMaterial) = meshinfo.BaseMaterial{GroupID: 5}
	*info.FaceData(1).(*meshinfo.BaseMaterial) = meshinfo.BaseMaterial{GroupID: 5}
	*info.FaceData(2).(*meshinfo.BaseMaterial) = meshinfo.BaseMaterial{GroupID: 5, Index: 1}
	*info.FaceData(3).(*meshinfo.BaseMaterial) = meshinfo.BaseMaterial{GroupID: 5, Index: 1}
	*info.FaceData(8).(*meshinfo.BaseMaterial) = meshinfo.BaseMaterial{GroupID: 5}
	*info.FaceData(9).(*meshinfo.BaseMaterial) = meshinfo.BaseMaterial{GroupID: 5}
	*info.FaceData(10).(*meshinfo.BaseMaterial) = meshinfo.BaseMaterial{GroupID: 5}
	*info.FaceData(11).(*meshinfo.BaseMaterial) = meshinfo.BaseMaterial{GroupID: 5}	
	info = handler.AddTextureCoordsInfo(uint32(len(meshRes.Mesh.Faces)))
	*info.FaceData(4).(*meshinfo.TextureCoords) = meshinfo.TextureCoords{TextureID: 6, Coords: [3]mgl32.Vec2{{0.3,0.5}, {0.3,0.8},{0.5, 0.8}}}
	*info.FaceData(5).(*meshinfo.TextureCoords) = meshinfo.TextureCoords{TextureID: 6, Coords: [3]mgl32.Vec2{{0.5,0.5}, {0.3,0.5},{0.5, 0.8}}}
	info = handler.AddNodeColorInfo(uint32(len(meshRes.Mesh.Faces)))
	*info.FaceData(6).(*meshinfo.NodeColor) = meshinfo.NodeColor{Colors: [3]color.RGBA{{R: 85, G: 85, B: 85, A: 255}, {A: 255}, {R: 16, G: 21, B: 103, A: 255}}}
	*info.FaceData(7).(*meshinfo.NodeColor) = meshinfo.NodeColor{Colors: [3]color.RGBA{{R: 16, G: 21, B: 103, A: 255}, {A: 255}, {R: 53, G: 4, B: 80, A: 255}}}

	want.Resources = append(want.Resources, &go3mf.SliceStackResource{ID: 10, ModelPath: "/2D/2Dmodel.model", SliceStack: otherSlices, TimesRefered: 1})
	want.Resources = append(want.Resources, []go3mf.Identifier{baseMaterials, baseTexture, sliceStack, sliceStackRef, meshRes}...)

	got := new(go3mf.Model)
	got.Resources = append(got.Resources, &go3mf.SliceStackResource{ID: 10, ModelPath: "/2D/2Dmodel.model", SliceStack: otherSlices})
	r := &Reader{
		Model: got,
		r: newMockPackage(new(modelBuilder).withDefaultModel().withElement(`
			<resources>
				<basematerials id="5">
					<base name="Blue PLA" displaycolor="#0000FF" />
					<base name="Red ABS" displaycolor="#FF0000" />
				</basematerials>
				<m:texture2d id="6" path="/3D/Texture/msLogo.png" contenttype="image/png" tilestyleu="wrap" tilestylev="mirror" filter="auto" />
				<m:colorgroup id="1">
					<m:color color="#FFFFFF" /> <m:color color="#000000" /> <m:color color="#1AB567" /> <m:color color="#DF045A" />
				</m:colorgroup>
				<m:texture2dgroup id="2" texid="6">
					<m:tex2coord u="0.3" v="0.5" /> <m:tex2coord u="0.3" v="0.8" />	<m:tex2coord u="0.5" v="0.8" />	<m:tex2coord u="0.5" v="0.5" />
				</m:texture2dgroup>
				<s:slicestack id="3" zbottom="1">
					<s:slice ztop="0">
						<s:vertices>
							<s:vertex x="1.01" y="1.02" /> <s:vertex x="9.03" y="1.04" /> <s:vertex x="9.05" y="9.06" /> <s:vertex x="1.07" y="9.08" />
						</s:vertices>
						<s:polygon startv="0">
							<s:segment v2="1"></s:segment> <s:segment v2="2"></s:segment> <s:segment v2="3"></s:segment> <s:segment v2="0"></s:segment>
						</s:polygon>
					</s:slice>
					<s:slice ztop="0.1">
						<s:vertices>
							<s:vertex x="1.01" y="1.02" /> <s:vertex x="9.03" y="1.04" /> <s:vertex x="9.05" y="9.06" /> <s:vertex x="1.07" y="9.08" />
						</s:vertices>
						<s:polygon startv="0"> 
							<s:segment v2="2"></s:segment> <s:segment v2="1"></s:segment> <s:segment v2="3"></s:segment> <s:segment v2="0"></s:segment>
						</s:polygon>
					</s:slice>
				</s:slicestack>
				<s:slicestack id="7" zbottom="1.1">
					<s:sliceref slicestackid="10" slicepath="/2D/2Dmodel.model" />
				</s:slicestack>
				<object id="1" name="Box 1" pid="5" pindex="0" s:meshresolution="lowres" s:slicestackid="3" partnumber="11111111-1111-1111-1111-111111111111" type="model">
					<mesh>
						<vertices>
							<vertex x="0" y="0" z="0" />
							<vertex x="100.00000" y="0" z="0" />
							<vertex x="100.00000" y="100.00000" z="0" />
							<vertex x="0" y="100.00000" z="0" />
							<vertex x="0" y="0" z="100.00000" />
							<vertex x="100.00000" y="0" z="100.00000" />
							<vertex x="100.00000" y="100.00000" z="100.00000" />
							<vertex x="0" y="100.00000" z="100.00000" />
						</vertices>
						<triangles>
							<triangle v1="3" v2="2" v3="1" />
							<triangle v1="1" v2="0" v3="3" />
							<triangle v1="4" v2="5" v3="6" p1="1" />
							<triangle v1="6" v2="7" v3="4" pid="5" p1="1" />
							<triangle v1="0" v2="1" v3="5" pid="2" p1="0" p2="1" p3="2"/>
							<triangle v1="5" v2="4" v3="0" pid="2" p1="3" p2="0" p3="2"/>
							<triangle v1="1" v2="2" v3="6" pid="1" p1="0" p2="1" p3="2"/>
							<triangle v1="6" v2="5" v3="1" pid="1" p1="2" p2="1" p3="3"/>
							<triangle v1="2" v2="3" v3="7" />
							<triangle v1="7" v2="6" v3="2" />
							<triangle v1="3" v2="0" v3="4" />
							<triangle v1="4" v2="7" v3="3" />
						</triangles>
					</mesh>
				</object>
			</resources>`).build()),
	}

	t.Run("base", func(t *testing.T) {
		if err := r.processRootModel(); err != nil {
			t.Errorf("Reader.processRootModel() unexpected error = %v", err)
			return
		}
		deep.CompareUnexportedFields = true
		deep.MaxDepth = 20
		if diff := deep.Equal(r.Model, want); diff != nil {
			t.Errorf("Reader.processRootModel() = %v", diff)
			return
		}
	})
}

func TestReader_namespaceRegistered(t *testing.T) {
	type args struct {
		ns string
	}
	tests := []struct {
		name string
		r    *Reader
		args args
		want bool
	}{
		{"empty", &Reader{namespaces: []string{"http://xml.com"}}, args{""}, false},
		{"exist", &Reader{namespaces: []string{"http://xml.com"}}, args{"http://xml.com"}, true},
		{"noexist", &Reader{namespaces: []string{"http://xml.com"}}, args{"xmls"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.namespaceRegistered(tt.args.ns); got != tt.want {
				t.Errorf("Reader.namespaceRegistered() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_strToMatrix(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    mgl32.Mat4
		wantErr bool
	}{
		{"empty", args{""}, mgl32.Mat4{}, true},
		{"11values", args{"1 1 1 1 1 1 1 1 1 1 1"}, mgl32.Mat4{}, true},
		{"13values", args{"1 1 1 1 1 1 1 1 1 1 1 1 1"}, mgl32.Mat4{}, true},
		{"char", args{"1 1 a 1 1 1 1 1 1 1 1 1"}, mgl32.Mat4{}, true},
		{"base", args{"1 1 1 1 1 1 1 1 1 1 1 1"}, mgl32.Mat4{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 1}, false},
		{"other", args{"0 1 2 10 11 12 20 21 22 30 31 32"}, mgl32.Mat4{0, 10, 20, 30, 1, 11, 21, 31, 2, 12, 22, 32, 0, 0, 0, 1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := strToMatrix(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("strToMatrix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("strToMatrix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_strToSRGB(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		wantC   color.RGBA
		wantErr bool
	}{
		{"empty", args{""}, color.RGBA{}, true},
		{"nohashrgb", args{"101010"}, color.RGBA{}, true},
		{"nohashrgba", args{"10101010"}, color.RGBA{}, true},
		{"invalidChar", args{"#€0101010"}, color.RGBA{}, true},
		{"rgb", args{"#112233"}, color.RGBA{17, 34, 51, 255}, false},
		{"rgb", args{"#000233"}, color.RGBA{0, 2, 51, 255}, false},
		{"rgba", args{"#00023311"}, color.RGBA{0, 2, 51, 17}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotC, err := strToSRGB(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("strToSRGB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotC, tt.wantC) {
				t.Errorf("strToSRGB() = %v, want %v", gotC, tt.wantC)
			}
		})
	}
}
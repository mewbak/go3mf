package stl

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/qmuntal/go3mf/geo"
)

// asciiDecoder can create a Model from a Read stream that is feeded with a ASCII STL.
type asciiDecoder struct {
	r     io.Reader
	units float32
}

func (d *asciiDecoder) decode(ctx context.Context, m *geo.Mesh) (err error) {
	m.StartCreation(geo.CreationOptions{CalculateConnectivity: true})
	defer m.EndCreation()
	position := 0
	nextFaceCheck := checkEveryFaces
	var nodes [3]uint32
	scanner := bufio.NewScanner(d.r)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) == 4 && fields[0] == "vertex" {
			var f [3]float64
			f[0], _ = strconv.ParseFloat(fields[1], 32)
			f[1], _ = strconv.ParseFloat(fields[2], 32)
			f[2], _ = strconv.ParseFloat(fields[3], 32)
			nodes[position] = m.AddNode(geo.Point3D{float32(f[0]), float32(f[1]), float32(f[2])})
			position++

			if position == 3 {
				position = 0
				m.AddFace(nodes[0], nodes[1], nodes[2])
				if len(m.Faces) > nextFaceCheck {
					select {
					case <-ctx.Done():
						err = ctx.Err()
						break
					default: // Default is must to avoid blocking
					}
					nextFaceCheck += checkEveryFaces
				}
			}
		}
		if err != nil {
			return err
		}
	}
	return scanner.Err()
}

type asciiEncoder struct {
	w io.Writer
}

const pstr = "solid\nfacet normal %f %f %f\nouter loop\nvertex %f %f %f\nvertex %f %f %f\nvertex %f %f %f\nendloop\nendfacet\nendsolid\n"

func (e *asciiEncoder) encode(m *geo.Mesh) error {
	for i := range m.Faces {
		n1, n2, n3 := m.FaceNodes(uint32(i))
		n := faceNormal(*n1, *n2, *n3)
		_, err := io.WriteString(e.w, fmt.Sprintf(pstr, n[0], n[1], n[2], n1[0], n1[1], n1[2], n2[0], n2[1], n2[2], n3[0], n3[1], n3[2]))

		if err != nil {
			return err
		}
	}

	return nil
}

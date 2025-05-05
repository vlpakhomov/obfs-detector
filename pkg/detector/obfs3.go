package detector

import (
	"crypto/rand"
	"encoding/binary"
	"gonum.org/v1/gonum/stat"
	"math"
	"obfs-detector/pkg/null"
)

type obfs3 struct {
	alpha float64
}

var Obfs3 obfs3 = obfs3{}

func (o *obfs3) Detect(data []byte) (null.Null[bool], null.Null[string], error) {
	sourceDist := make([]float64, 0)
	for i := range data {
		bits := binary.LittleEndian.Uint64(data[i:i])
		float := math.Float64frombits(bits)
		sourceDist = append(sourceDist, float)
	}

	uniformDistByte := make([]byte, 0, len(sourceDist))
	_, err := rand.Read(uniformDistByte)
	if err != nil {
		return null.NewExplicit(false, false), null.NewExplicit("", false), err
	}

	uniformDist := make([]float64, 0, len(sourceDist))
	for i := range uniformDistByte {
		bits := binary.LittleEndian.Uint64(uniformDistByte[i:i])
		float := math.Float64frombits(bits)
		uniformDist = append(uniformDist, float)
	}

	entropy := stat.KolmogorovSmirnov(sourceDist, []float64{}, uniformDist, []float64{})

	if entropy > o.alpha {
		return null.New(true), null.New("obfs3 detected"), nil
	}

	return null.New(false), null.New("obfs3 undetected"), nil
}

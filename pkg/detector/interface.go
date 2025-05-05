package detector

import "obfs-detector/pkg/null"

type Detector interface {
	Detect(data []byte) (null.Null[bool], null.Null[string], error)
}

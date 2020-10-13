package stage

type Stage string

const (
	Version = "Version"
	Manifest = "Manifest"
	ManifestSignature = "Manifest signature"
	ManifestAugment = "Manifest augment"
	Header = "Header"
	HeaderAugment = "Header Augment"
	Data = "Payload"
)

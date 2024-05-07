package v1alpha1

type TypeMeta struct {
	Kind       string `json:"kind,omitempty" yaml:"kind,omitempty"`
	APIVersion string `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
}

func (tm *TypeMeta) GetKind() string {
	return tm.Kind
}

func (tm *TypeMeta) GetAPIVersion() string {
	return tm.Kind
}

// Discriminator is an interface to discriminate the kind api object.
type Discriminator interface {
	GetKind() string
	GetAPIVersion() string
}

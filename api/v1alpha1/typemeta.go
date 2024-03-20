package v1alpha1

type TypeMeta struct {
	Kind       string `json:"kind,omitempty" yaml:"kind,omitempty"`
	APIVersion string `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
}

func (tm *TypeMeta) GetKind() string {
	return tm.Kind
}

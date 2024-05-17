package v1alpha1

import object "github.com/holos-run/holos/service/gen/holos/object/v1alpha1"

// Form represents a collection of Formly json powered form.
type Form struct {
	TypeMeta `json:",inline" yaml:",inline"`
	Spec     FormSpec `json:"spec" yaml:"spec"`
}

type FormSpec struct {
	Form object.Form `json:"form" yaml:"form"`
}

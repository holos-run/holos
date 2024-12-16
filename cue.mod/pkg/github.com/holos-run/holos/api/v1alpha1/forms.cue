package v1alpha1

// Form represents a collection of Formly json powered form.
#Form: {
	apiVersion: #APIVersion
	kind:       "Form"
}

// #FormBuilder provides a concrete #Form via the Output field.
#FormBuilder: {
	Name: string
	Sections: {[NAME=string]: #FormSection & {name: NAME}}

	Output: #Form & {
		spec: form: fieldConfigs: [for s in Sections {s.wrapper}]
	}
}

// #FormSection represents a configuration section of the front end UI.  The
// wrapper field provides a concrete #FieldConfig for the form section.  The
// fields of the section map to form input fields.
// Refer to: to https://formly.dev/docs/examples/other/nested-formly-forms
#FormSection: {
	name:        string // e.g. "org"
	displayName: string // e.g. "Organization"
	description: string

	expressions: {[string]: string}

	fieldConfigs: {[NAME=string]: #FieldConfig & {key: NAME}}

	let Description = description
	let Expressions = expressions

	// Wrap the fields of the section into one FormlyFieldConfig
	wrapper: #FieldConfig & {
		key: name
		// See our custom wrappers registered in app.config.ts
		wrappers: ["holos-panel"]
		props: label:       displayName
		props: description: Description
		for k, v in Expressions {
			expressions: "\(k)": v
		}

		// Might need to initialize the default value for a fieldGroup
		// https://github.com/ngx-formly/ngx-formly/issues/3667
		fieldGroup: [for fc in fieldConfigs {fc}]
	}
}

// #FieldConfig represents a Formly Field Config.
// Refer to https://formly.dev/docs/api/core#formlyfieldconfig
// Refer to https://formly.dev/docs/api/ui/material/select
#FieldConfig: {
	key: string
	// type is optional, may be a nested form which has no type field
	type?: string | "input" | "select" | "checkbox"
	// For nested forms, refer: to https://formly.dev/docs/examples/other/nested-formly-forms
	wrappers?: [...string]
	// Refer to: https://formly.dev/docs/api/ui/material/select#formlyselectprops
	// and other input field select props.
	props: {
		#FormlySelectProps

		label:        string
		type?:        string
		placeholder?: string
		description:  string
		required?:    *true | false
		pattern?:     string
		minLength?:   number
		maxLength?:   number
	}
	// Refer to: https://github.com/ngx-formly/ngx-formly/blob/v6.3.0/src/core/src/lib/models/fieldconfig.ts#L49-L64
	// We support only the string form.
	validation?: {
		// Note, you can set messages for pattern, minLength, maxLength here.
		messages?: [string]: string
	}

	// Refer to: https://github.com/ngx-formly/ngx-formly/blob/v6.3.0/src/core/src/lib/models/fieldconfig.ts#L66-L71
	// We do not support validators because they must be javascript functions, not data.
	validators?: "not supported"

	// Refer to: https://github.com/ngx-formly/ngx-formly/blob/v6.3.0/src/core/src/lib/models/fieldconfig.ts#L115-L120
	expressions?: [string]: string
	hide?: true | false
	// Required to populate protobuf value.
	resetOnHide:   *true | false
	defaultValue?: _
	className?:    string
	fieldGroup?: [...#FieldConfig]
	focus?: true | *false
	modelOptions?: {
		debounce?: {
			default: number
		}
		updateOn?: "change" | "blur" | "submit"
	}
}

// Refer to https://formly.dev/docs/api/ui/material/select#formlyselectprops
#FormlySelectProps: {
	disableOptionCentering?:    true | false
	multiple?:                  true | false
	panelClass?:                string
	selectAllOption?:           string
	typeaheadDebounceInterval?: number

	options?: [...{value: string | number | bool, label: string, disabled?: true | *false}]

	// These could be used to set different keys for value and label in the
	// options list, but we don't support that level of customization.
	// They're here for documentation purposes only.
	labelProp?: "label"
	valueProp?: "value"
}

package v1alpha1

#Platform: {
	name:        string // short dns label name
	displayName: string // Display name
	description: string // Plaform description

	sections: {[NAME=string]: #ConfigSection & {name: NAME}}

	Form: {
		let Name = name
		apiVersion: "forms.holos.run/v1alpha1"
		kind:       "PlatformForm"
		metadata: name: Name
		spec: #PlatformFormSpec
	}

	let Sections = sections
	Form: spec: sections: [for s in Sections {s.output}]
}

#PlatformFormSpec: {
	sections: [...#ConfigSectionOutput]
}

//  #ConfigSection represents a configuration section of the front end UI.  For
// example, Organization config values.  The fields of the section map to form
// input fields.
#ConfigSection: {
	name:        string // e.g. "org"
	displayName: string // e.g. "Organization"
	description: string
	fieldConfigs: {[NAME=string]: #FieldConfig & {key: NAME}}

	let Name = name
	let DisplayName = displayName
	let Description = description
	let FieldConfigs = fieldConfigs

	output: #ConfigSectionOutput & {
		name:        Name
		displayName: DisplayName
		description: Description
		fieldConfigs: [for fc in FieldConfigs {fc}]
	}
}

#ConfigSectionOutput: {
	name:        string
	displayName: string
	description: string
	fieldConfigs: [...#FieldConfig]
}

// Refer to https://formly.dev/docs/api/core#formlyfieldconfig
// Refer to https://formly.dev/docs/api/ui/material/select
#FieldConfig: {
	key:  string
	type: string | *"input" | "select" | "checkbox"
	// Refer to: https://formly.dev/docs/api/ui/material/select#formlyselectprops
	// and other input field select props.
	props: {
		#FormlySelectProps

		label:        string
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

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
#FieldConfig: {
	key:  string
	type: "input"
	props: {
		label:       string
		placeholder: string
		description: string
		required:    *true | false
	}
}

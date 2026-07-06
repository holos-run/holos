// Authored by hand.  cue get go cannot generate the discriminated union or
// the per-kind input/output cardinality rules, so this file unifies explicit
// per-kind guards into the generated #Task definition per the design decision
// doc/design/v1beta1/schema.md#d5-open-and-closed-structs.
//
// For each kind value: the matching config field is required, every other
// config field is forbidden, and the inputs/output requiredness matches the
// task kinds table in doc/design/v1beta1/schema.md#task-kinds.
package core

#Task: {
	kind: string

	if kind == "Resources" {
		resources!: #Resources
		helm?:      _|_
		file?:      _|_
		kustomize?: _|_
		join?:      _|_
		command?:   _|_
		artifact?:  _|_
		inputs?:    _|_
		output!:    #FileOrDirectoryPath
	}

	if kind == "Helm" {
		helm!:      #Helm
		resources?: _|_
		file?:      _|_
		kustomize?: _|_
		join?:      _|_
		command?:   _|_
		artifact?:  _|_
		inputs?:    _|_
		output!:    #FileOrDirectoryPath
	}

	if kind == "File" {
		file!:      #File
		resources?: _|_
		helm?:      _|_
		kustomize?: _|_
		join?:      _|_
		command?:   _|_
		artifact?:  _|_
		inputs?:    _|_
		output!:    #FileOrDirectoryPath
	}

	if kind == "Kustomize" {
		kustomize!: #Kustomize
		resources?: _|_
		helm?:      _|_
		file?:      _|_
		join?:      _|_
		command?:   _|_
		artifact?:  _|_
		inputs!: [#FileOrDirectoryPath, ...#FileOrDirectoryPath]
		output!: #FileOrDirectoryPath
	}

	if kind == "Join" {
		join!:      #Join
		resources?: _|_
		helm?:      _|_
		file?:      _|_
		kustomize?: _|_
		command?:   _|_
		artifact?:  _|_
		inputs!: [#FileOrDirectoryPath, ...#FileOrDirectoryPath]
		output!: #FileOrDirectoryPath
	}

	if kind == "Command" {
		command!: #Command & {
			// A command without an argument vector cannot execute.
			args!: [string, ...string]
			isStdoutOutput: bool | *false
		}
		resources?: _|_
		helm?:      _|_
		file?:      _|_
		kustomize?: _|_
		join?:      _|_
		artifact?:  _|_

		// Commands may declare zero or more inputs.  Regular with an empty
		// default, rather than optional, so the stdin guards below may
		// reference the field.
		inputs: [...#FileOrDirectoryPath] | *[]

		// A command captures stdout as the task output when isStdoutOutput is
		// true, so an output is required.
		if command.isStdoutOutput {
			output!: #FileOrDirectoryPath
		}

		// Stdin must name one of the task's declared inputs.
		if len(inputs) > 0 {
			command: stdin?: or(inputs)
		}
		if len(inputs) == 0 {
			command: stdin?: _|_
		}
	}

	if kind == "Artifact" {
		artifact!:  #Artifact
		resources?: _|_
		helm?:      _|_
		file?:      _|_
		kustomize?: _|_
		join?:      _|_
		command?:   _|_
		inputs!: [#FileOrDirectoryPath]
		output?: _|_
	}
}

// For validation, see the [Standard constraints](https://github.com/bufbuild/protovalidate/blob/main/docs/standard-constraints.md)
package object

#PlatformConfig: {
	platformModel: {...} @protobuf(1,google.protobuf.Struct,name=platform_model)
}

#Form: {
	fieldConfigs: [...{...}]
}

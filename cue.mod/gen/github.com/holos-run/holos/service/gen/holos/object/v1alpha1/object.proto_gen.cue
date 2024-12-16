// For validation, see the [Standard constraints](https://github.com/bufbuild/protovalidate/blob/main/docs/standard-constraints.md)
package object

import "time"

#Detail: {
	// Created by entity
	createdBy?: #ResourceEditor @protobuf(1,ResourceEditor,name=created_by)

	// Created at timestamp
	createdAt?: time.Time @protobuf(2,google.protobuf.Timestamp,name=created_at,"(buf.validate.field).timestamp.lt_now")

	// Updated by entity
	updatedBy?: #ResourceEditor @protobuf(3,ResourceEditor,name=updated_by)

	// Updated at timestamp
	updatedAt?: time.Time @protobuf(4,google.protobuf.Timestamp,name=updated_at,"(buf.validate.field).timestamp.lt_now")
}

// Subject represents the oidc iss and sub claims which uniquely identify a subject.
#Subject: {
	// iss represents the oidc id token iss claim.  Limits defined at
	// https://openid.net/specs/openid-authentication-1_1.html#limits
	iss?: string @protobuf(1,string,"(buf.validate.field).string=")

	// sub represents the oidc id token sub claim.
	sub?: string @protobuf(2,string,"(buf.validate.field).string=")
}

// UserRef refers to a User by uuid, email, or by the oidc iss and sub claims.
#UserRef: {
	{} | {
		@protobuf(option (buf.validate.oneof).required=true)
	} | {
		userId: string @protobuf(1,string,name=user_id,"(buf.validate.field).string.uuid")
	} | {
		email: string @protobuf(2,string,"(buf.validate.field).string.email")
	} | {
		subject: #Subject @protobuf(3,Subject)
	}
}

// Organization represents the ways in which a organization may be uniquely identified in the system.
#OrganizationRef: {
	{} | {
		@protobuf(option (buf.validate.oneof).required=true)
	} | {
		orgId: string @protobuf(1,string,name=org_id,"(buf.validate.field).string.uuid")
	} | {
		orgName: string @protobuf(2,string,name=org_name,"(buf.validate.field).cel=")
	}
}

// ResourceEditor represents the entity that most recently created or edited a resource.
#ResourceEditor: {
	{} | {
		@protobuf(option (buf.validate.oneof).required=true)
	} | {
		userId: string @protobuf(1,string,name=user_id,"(buf.validate.field).string.uuid")
	}
}

#ResourceOwner: {
	{} | {
		@protobuf(option (buf.validate.oneof).required=true)
	} | {
		orgId: string @protobuf(1,string,name=org_id,"(buf.validate.field).string.uuid")
	} | {
		userId: string @protobuf(2,string,name=user_id,"(buf.validate.field).string.uuid")
	}
}

// Form represents a Formly json powered form.
#Form: {
	// field_configs represents FormlyFieldConfig[] encoded as an array of JSON
	// objects organized by section.
	fieldConfigs?: [...{}] @protobuf(1,google.protobuf.Struct,name=field_configs)
}

// PlatformConfig represents the data passed from the holos cli to CUE when
// rendering configuration.  At present it contains only the platform model from
// the PlatformService, but it is expected to carry additional fields from
// additional data sources.  For this reason, there is a distinction in domain
// language between the "Platform Config" and the "Platform Model"  The config
// is a data transfer object that carries at least the model.  The model
// represents the form values from the PlatformService.
#PlatformConfig: {
	// platform_model represents the form values from the PlatformService.
	platformModel?: {} @protobuf(1,google.protobuf.Struct,name=platform_model)
}

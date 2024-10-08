// @generated by protoc-gen-es v1.10.0 with parameter "target=ts"
// @generated from file holos/organization/v1alpha1/organization.proto (package holos.organization.v1alpha1, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3 } from "@bufbuild/protobuf";
import { Detail } from "../../object/v1alpha1/object_pb.js";

/**
 * @generated from message holos.organization.v1alpha1.Organization
 */
export class Organization extends Message<Organization> {
  /**
   * Unique id assigned by the server.
   *
   * @generated from field: optional string org_id = 1;
   */
  orgId?: string;

  /**
   * Name is the organization name as a dns label.
   *
   * @generated from field: string name = 2;
   */
  name = "";

  /**
   * @generated from field: optional string display_name = 3;
   */
  displayName?: string;

  /**
   * @generated from field: optional holos.object.v1alpha1.Detail detail = 4;
   */
  detail?: Detail;

  constructor(data?: PartialMessage<Organization>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "holos.organization.v1alpha1.Organization";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "org_id", kind: "scalar", T: 9 /* ScalarType.STRING */, opt: true },
    { no: 2, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "display_name", kind: "scalar", T: 9 /* ScalarType.STRING */, opt: true },
    { no: 4, name: "detail", kind: "message", T: Detail, opt: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Organization {
    return new Organization().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Organization {
    return new Organization().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Organization {
    return new Organization().fromJsonString(jsonString, options);
  }

  static equals(a: Organization | PlainMessage<Organization> | undefined, b: Organization | PlainMessage<Organization> | undefined): boolean {
    return proto3.util.equals(Organization, a, b);
  }
}


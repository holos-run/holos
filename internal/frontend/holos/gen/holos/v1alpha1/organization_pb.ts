// @generated by protoc-gen-es v1.8.0 with parameter "target=ts"
// @generated from file holos/v1alpha1/organization.proto (package holos.v1alpha1, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3 } from "@bufbuild/protobuf";
import { Timestamps } from "./timestamps_pb.js";

/**
 * @generated from message holos.v1alpha1.Organization
 */
export class Organization extends Message<Organization> {
  /**
   * Unique id assigned by the server.
   *
   * @generated from field: string id = 1;
   */
  id = "";

  /**
   * @generated from field: string name = 2;
   */
  name = "";

  /**
   * @generated from field: string display_name = 3;
   */
  displayName = "";

  /**
   * @generated from field: holos.v1alpha1.Timestamps timestamps = 4;
   */
  timestamps?: Timestamps;

  constructor(data?: PartialMessage<Organization>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "holos.v1alpha1.Organization";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "display_name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 4, name: "timestamps", kind: "message", T: Timestamps },
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

/**
 * @generated from message holos.v1alpha1.RegisterOrganizationRequest
 */
export class RegisterOrganizationRequest extends Message<RegisterOrganizationRequest> {
  /**
   * @generated from field: optional string name = 1;
   */
  name?: string;

  /**
   * @generated from field: optional string display_name = 2;
   */
  displayName?: string;

  constructor(data?: PartialMessage<RegisterOrganizationRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "holos.v1alpha1.RegisterOrganizationRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */, opt: true },
    { no: 2, name: "display_name", kind: "scalar", T: 9 /* ScalarType.STRING */, opt: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): RegisterOrganizationRequest {
    return new RegisterOrganizationRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): RegisterOrganizationRequest {
    return new RegisterOrganizationRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): RegisterOrganizationRequest {
    return new RegisterOrganizationRequest().fromJsonString(jsonString, options);
  }

  static equals(a: RegisterOrganizationRequest | PlainMessage<RegisterOrganizationRequest> | undefined, b: RegisterOrganizationRequest | PlainMessage<RegisterOrganizationRequest> | undefined): boolean {
    return proto3.util.equals(RegisterOrganizationRequest, a, b);
  }
}

/**
 * @generated from message holos.v1alpha1.RegisterOrganizationResponse
 */
export class RegisterOrganizationResponse extends Message<RegisterOrganizationResponse> {
  /**
   * @generated from field: holos.v1alpha1.Organization organization = 1;
   */
  organization?: Organization;

  /**
   * @generated from field: bool already_exists = 2;
   */
  alreadyExists = false;

  constructor(data?: PartialMessage<RegisterOrganizationResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "holos.v1alpha1.RegisterOrganizationResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "organization", kind: "message", T: Organization },
    { no: 2, name: "already_exists", kind: "scalar", T: 8 /* ScalarType.BOOL */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): RegisterOrganizationResponse {
    return new RegisterOrganizationResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): RegisterOrganizationResponse {
    return new RegisterOrganizationResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): RegisterOrganizationResponse {
    return new RegisterOrganizationResponse().fromJsonString(jsonString, options);
  }

  static equals(a: RegisterOrganizationResponse | PlainMessage<RegisterOrganizationResponse> | undefined, b: RegisterOrganizationResponse | PlainMessage<RegisterOrganizationResponse> | undefined): boolean {
    return proto3.util.equals(RegisterOrganizationResponse, a, b);
  }
}

/**
 * Empty request, claims are pulled from the id token
 *
 * @generated from message holos.v1alpha1.GetOrganizationRequest
 */
export class GetOrganizationRequest extends Message<GetOrganizationRequest> {
  /**
   * name to look up
   *
   * @generated from field: optional string name = 1;
   */
  name?: string;

  constructor(data?: PartialMessage<GetOrganizationRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "holos.v1alpha1.GetOrganizationRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */, opt: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GetOrganizationRequest {
    return new GetOrganizationRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GetOrganizationRequest {
    return new GetOrganizationRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GetOrganizationRequest {
    return new GetOrganizationRequest().fromJsonString(jsonString, options);
  }

  static equals(a: GetOrganizationRequest | PlainMessage<GetOrganizationRequest> | undefined, b: GetOrganizationRequest | PlainMessage<GetOrganizationRequest> | undefined): boolean {
    return proto3.util.equals(GetOrganizationRequest, a, b);
  }
}

/**
 * @generated from message holos.v1alpha1.GetOrganizationResponse
 */
export class GetOrganizationResponse extends Message<GetOrganizationResponse> {
  /**
   * @generated from field: holos.v1alpha1.Organization organization = 1;
   */
  organization?: Organization;

  constructor(data?: PartialMessage<GetOrganizationResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "holos.v1alpha1.GetOrganizationResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "organization", kind: "message", T: Organization },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GetOrganizationResponse {
    return new GetOrganizationResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GetOrganizationResponse {
    return new GetOrganizationResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GetOrganizationResponse {
    return new GetOrganizationResponse().fromJsonString(jsonString, options);
  }

  static equals(a: GetOrganizationResponse | PlainMessage<GetOrganizationResponse> | undefined, b: GetOrganizationResponse | PlainMessage<GetOrganizationResponse> | undefined): boolean {
    return proto3.util.equals(GetOrganizationResponse, a, b);
  }
}


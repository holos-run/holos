// @generated by protoc-gen-es v1.9.0 with parameter "target=ts"
// @generated from file holos/platform/v1alpha1/platform_service.proto (package holos.platform.v1alpha1, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { FieldMask, Message, proto3 } from "@bufbuild/protobuf";
import { Platform } from "./platform_pb.js";

/**
 * @generated from message holos.platform.v1alpha1.CreatePlatformRequest
 */
export class CreatePlatformRequest extends Message<CreatePlatformRequest> {
  /**
   * @generated from field: holos.platform.v1alpha1.Platform platform = 1;
   */
  platform?: Platform;

  constructor(data?: PartialMessage<CreatePlatformRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "holos.platform.v1alpha1.CreatePlatformRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "platform", kind: "message", T: Platform },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): CreatePlatformRequest {
    return new CreatePlatformRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): CreatePlatformRequest {
    return new CreatePlatformRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): CreatePlatformRequest {
    return new CreatePlatformRequest().fromJsonString(jsonString, options);
  }

  static equals(a: CreatePlatformRequest | PlainMessage<CreatePlatformRequest> | undefined, b: CreatePlatformRequest | PlainMessage<CreatePlatformRequest> | undefined): boolean {
    return proto3.util.equals(CreatePlatformRequest, a, b);
  }
}

/**
 * @generated from message holos.platform.v1alpha1.CreatePlatformResponse
 */
export class CreatePlatformResponse extends Message<CreatePlatformResponse> {
  /**
   * @generated from field: holos.platform.v1alpha1.Platform platform = 1;
   */
  platform?: Platform;

  constructor(data?: PartialMessage<CreatePlatformResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "holos.platform.v1alpha1.CreatePlatformResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "platform", kind: "message", T: Platform },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): CreatePlatformResponse {
    return new CreatePlatformResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): CreatePlatformResponse {
    return new CreatePlatformResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): CreatePlatformResponse {
    return new CreatePlatformResponse().fromJsonString(jsonString, options);
  }

  static equals(a: CreatePlatformResponse | PlainMessage<CreatePlatformResponse> | undefined, b: CreatePlatformResponse | PlainMessage<CreatePlatformResponse> | undefined): boolean {
    return proto3.util.equals(CreatePlatformResponse, a, b);
  }
}

/**
 * @generated from message holos.platform.v1alpha1.GetPlatformRequest
 */
export class GetPlatformRequest extends Message<GetPlatformRequest> {
  /**
   * @generated from field: string platform_id = 1;
   */
  platformId = "";

  /**
   * FieldMask represents the response Platform fields to include.
   *
   * @generated from field: google.protobuf.FieldMask field_mask = 2;
   */
  fieldMask?: FieldMask;

  constructor(data?: PartialMessage<GetPlatformRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "holos.platform.v1alpha1.GetPlatformRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "platform_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "field_mask", kind: "message", T: FieldMask },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GetPlatformRequest {
    return new GetPlatformRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GetPlatformRequest {
    return new GetPlatformRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GetPlatformRequest {
    return new GetPlatformRequest().fromJsonString(jsonString, options);
  }

  static equals(a: GetPlatformRequest | PlainMessage<GetPlatformRequest> | undefined, b: GetPlatformRequest | PlainMessage<GetPlatformRequest> | undefined): boolean {
    return proto3.util.equals(GetPlatformRequest, a, b);
  }
}

/**
 * @generated from message holos.platform.v1alpha1.GetPlatformResponse
 */
export class GetPlatformResponse extends Message<GetPlatformResponse> {
  /**
   * @generated from field: holos.platform.v1alpha1.Platform platform = 1;
   */
  platform?: Platform;

  constructor(data?: PartialMessage<GetPlatformResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "holos.platform.v1alpha1.GetPlatformResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "platform", kind: "message", T: Platform },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GetPlatformResponse {
    return new GetPlatformResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GetPlatformResponse {
    return new GetPlatformResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GetPlatformResponse {
    return new GetPlatformResponse().fromJsonString(jsonString, options);
  }

  static equals(a: GetPlatformResponse | PlainMessage<GetPlatformResponse> | undefined, b: GetPlatformResponse | PlainMessage<GetPlatformResponse> | undefined): boolean {
    return proto3.util.equals(GetPlatformResponse, a, b);
  }
}

/**
 * @generated from message holos.platform.v1alpha1.UpdatePlatformResponse
 */
export class UpdatePlatformResponse extends Message<UpdatePlatformResponse> {
  /**
   * @generated from field: holos.platform.v1alpha1.Platform platform = 1;
   */
  platform?: Platform;

  constructor(data?: PartialMessage<UpdatePlatformResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "holos.platform.v1alpha1.UpdatePlatformResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "platform", kind: "message", T: Platform },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): UpdatePlatformResponse {
    return new UpdatePlatformResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): UpdatePlatformResponse {
    return new UpdatePlatformResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): UpdatePlatformResponse {
    return new UpdatePlatformResponse().fromJsonString(jsonString, options);
  }

  static equals(a: UpdatePlatformResponse | PlainMessage<UpdatePlatformResponse> | undefined, b: UpdatePlatformResponse | PlainMessage<UpdatePlatformResponse> | undefined): boolean {
    return proto3.util.equals(UpdatePlatformResponse, a, b);
  }
}

/**
 * @generated from message holos.platform.v1alpha1.UpdatePlatformRequest
 */
export class UpdatePlatformRequest extends Message<UpdatePlatformRequest> {
  /**
   * @generated from field: holos.platform.v1alpha1.Platform platform = 1;
   */
  platform?: Platform;

  /**
   * FieldMask represents the request Platform fields to update.
   *
   * @generated from field: google.protobuf.FieldMask field_mask = 2;
   */
  fieldMask?: FieldMask;

  constructor(data?: PartialMessage<UpdatePlatformRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "holos.platform.v1alpha1.UpdatePlatformRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "platform", kind: "message", T: Platform },
    { no: 2, name: "field_mask", kind: "message", T: FieldMask },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): UpdatePlatformRequest {
    return new UpdatePlatformRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): UpdatePlatformRequest {
    return new UpdatePlatformRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): UpdatePlatformRequest {
    return new UpdatePlatformRequest().fromJsonString(jsonString, options);
  }

  static equals(a: UpdatePlatformRequest | PlainMessage<UpdatePlatformRequest> | undefined, b: UpdatePlatformRequest | PlainMessage<UpdatePlatformRequest> | undefined): boolean {
    return proto3.util.equals(UpdatePlatformRequest, a, b);
  }
}

/**
 * @generated from message holos.platform.v1alpha1.ListPlatformsRequest
 */
export class ListPlatformsRequest extends Message<ListPlatformsRequest> {
  /**
   * @generated from field: string org_id = 1;
   */
  orgId = "";

  /**
   * FieldMask represents the response Platform fields to include.
   *
   * @generated from field: google.protobuf.FieldMask field_mask = 2;
   */
  fieldMask?: FieldMask;

  constructor(data?: PartialMessage<ListPlatformsRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "holos.platform.v1alpha1.ListPlatformsRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "org_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "field_mask", kind: "message", T: FieldMask },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): ListPlatformsRequest {
    return new ListPlatformsRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): ListPlatformsRequest {
    return new ListPlatformsRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): ListPlatformsRequest {
    return new ListPlatformsRequest().fromJsonString(jsonString, options);
  }

  static equals(a: ListPlatformsRequest | PlainMessage<ListPlatformsRequest> | undefined, b: ListPlatformsRequest | PlainMessage<ListPlatformsRequest> | undefined): boolean {
    return proto3.util.equals(ListPlatformsRequest, a, b);
  }
}

/**
 * @generated from message holos.platform.v1alpha1.ListPlatformsResponse
 */
export class ListPlatformsResponse extends Message<ListPlatformsResponse> {
  /**
   * @generated from field: repeated holos.platform.v1alpha1.Platform platforms = 1;
   */
  platforms: Platform[] = [];

  constructor(data?: PartialMessage<ListPlatformsResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "holos.platform.v1alpha1.ListPlatformsResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "platforms", kind: "message", T: Platform, repeated: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): ListPlatformsResponse {
    return new ListPlatformsResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): ListPlatformsResponse {
    return new ListPlatformsResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): ListPlatformsResponse {
    return new ListPlatformsResponse().fromJsonString(jsonString, options);
  }

  static equals(a: ListPlatformsResponse | PlainMessage<ListPlatformsResponse> | undefined, b: ListPlatformsResponse | PlainMessage<ListPlatformsResponse> | undefined): boolean {
    return proto3.util.equals(ListPlatformsResponse, a, b);
  }
}


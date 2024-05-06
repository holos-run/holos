// @generated by protoc-gen-es v1.9.0 with parameter "target=ts"
// @generated from file holos/v1alpha1/system.proto (package holos.v1alpha1, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3 } from "@bufbuild/protobuf";

/**
 * @generated from message holos.v1alpha1.SeedDatabaseRequest
 */
export class SeedDatabaseRequest extends Message<SeedDatabaseRequest> {
  constructor(data?: PartialMessage<SeedDatabaseRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "holos.v1alpha1.SeedDatabaseRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): SeedDatabaseRequest {
    return new SeedDatabaseRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): SeedDatabaseRequest {
    return new SeedDatabaseRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): SeedDatabaseRequest {
    return new SeedDatabaseRequest().fromJsonString(jsonString, options);
  }

  static equals(a: SeedDatabaseRequest | PlainMessage<SeedDatabaseRequest> | undefined, b: SeedDatabaseRequest | PlainMessage<SeedDatabaseRequest> | undefined): boolean {
    return proto3.util.equals(SeedDatabaseRequest, a, b);
  }
}

/**
 * @generated from message holos.v1alpha1.SeedDatabaseResponse
 */
export class SeedDatabaseResponse extends Message<SeedDatabaseResponse> {
  constructor(data?: PartialMessage<SeedDatabaseResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "holos.v1alpha1.SeedDatabaseResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): SeedDatabaseResponse {
    return new SeedDatabaseResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): SeedDatabaseResponse {
    return new SeedDatabaseResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): SeedDatabaseResponse {
    return new SeedDatabaseResponse().fromJsonString(jsonString, options);
  }

  static equals(a: SeedDatabaseResponse | PlainMessage<SeedDatabaseResponse> | undefined, b: SeedDatabaseResponse | PlainMessage<SeedDatabaseResponse> | undefined): boolean {
    return proto3.util.equals(SeedDatabaseResponse, a, b);
  }
}

/**
 * @generated from message holos.v1alpha1.DropTablesRequest
 */
export class DropTablesRequest extends Message<DropTablesRequest> {
  constructor(data?: PartialMessage<DropTablesRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "holos.v1alpha1.DropTablesRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): DropTablesRequest {
    return new DropTablesRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): DropTablesRequest {
    return new DropTablesRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): DropTablesRequest {
    return new DropTablesRequest().fromJsonString(jsonString, options);
  }

  static equals(a: DropTablesRequest | PlainMessage<DropTablesRequest> | undefined, b: DropTablesRequest | PlainMessage<DropTablesRequest> | undefined): boolean {
    return proto3.util.equals(DropTablesRequest, a, b);
  }
}

/**
 * @generated from message holos.v1alpha1.DropTablesResponse
 */
export class DropTablesResponse extends Message<DropTablesResponse> {
  constructor(data?: PartialMessage<DropTablesResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "holos.v1alpha1.DropTablesResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): DropTablesResponse {
    return new DropTablesResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): DropTablesResponse {
    return new DropTablesResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): DropTablesResponse {
    return new DropTablesResponse().fromJsonString(jsonString, options);
  }

  static equals(a: DropTablesResponse | PlainMessage<DropTablesResponse> | undefined, b: DropTablesResponse | PlainMessage<DropTablesResponse> | undefined): boolean {
    return proto3.util.equals(DropTablesResponse, a, b);
  }
}


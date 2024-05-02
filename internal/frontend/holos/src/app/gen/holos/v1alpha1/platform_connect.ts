// @generated by protoc-gen-connect-es v1.4.0 with parameter "target=ts"
// @generated from file holos/v1alpha1/platform.proto (package holos.v1alpha1, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import { AddPlatformRequest, GetPlatformConfigRequest, GetPlatformRequest, GetPlatformResponse, GetPlatformsRequest, GetPlatformsResponse, PlatformConfig, PutPlatformConfigRequest } from "./platform_pb.js";
import { MethodKind } from "@bufbuild/protobuf";

/**
 * @generated from service holos.v1alpha1.PlatformService
 */
export const PlatformService = {
  typeName: "holos.v1alpha1.PlatformService",
  methods: {
    /**
     * @generated from rpc holos.v1alpha1.PlatformService.AddPlatform
     */
    addPlatform: {
      name: "AddPlatform",
      I: AddPlatformRequest,
      O: GetPlatformsResponse,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc holos.v1alpha1.PlatformService.GetPlatforms
     */
    getPlatforms: {
      name: "GetPlatforms",
      I: GetPlatformsRequest,
      O: GetPlatformsResponse,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc holos.v1alpha1.PlatformService.GetPlatform
     */
    getPlatform: {
      name: "GetPlatform",
      I: GetPlatformRequest,
      O: GetPlatformResponse,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc holos.v1alpha1.PlatformService.PutPlatformConfig
     */
    putPlatformConfig: {
      name: "PutPlatformConfig",
      I: PutPlatformConfigRequest,
      O: GetPlatformResponse,
      kind: MethodKind.Unary,
    },
    /**
     * GetConfig provides the unmarshalled config values for use with CUE
     *
     * @generated from rpc holos.v1alpha1.PlatformService.GetConfig
     */
    getConfig: {
      name: "GetConfig",
      I: GetPlatformConfigRequest,
      O: PlatformConfig,
      kind: MethodKind.Unary,
    },
  }
} as const;


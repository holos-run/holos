// @generated by protoc-gen-connect-es v1.4.0 with parameter "target=ts"
// @generated from file holos/platform/v1alpha1/platform_service.proto (package holos.platform.v1alpha1, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import { CreatePlatformRequest, CreatePlatformResponse, GetPlatformRequest, GetPlatformResponse, ListPlatformsRequest, ListPlatformsResponse, UpdatePlatformRequest, UpdatePlatformResponse } from "./platform_service_pb.js";
import { MethodKind } from "@bufbuild/protobuf";

/**
 * @generated from service holos.platform.v1alpha1.PlatformService
 */
export const PlatformService = {
  typeName: "holos.platform.v1alpha1.PlatformService",
  methods: {
    /**
     * @generated from rpc holos.platform.v1alpha1.PlatformService.CreatePlatform
     */
    createPlatform: {
      name: "CreatePlatform",
      I: CreatePlatformRequest,
      O: CreatePlatformResponse,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc holos.platform.v1alpha1.PlatformService.GetPlatform
     */
    getPlatform: {
      name: "GetPlatform",
      I: GetPlatformRequest,
      O: GetPlatformResponse,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc holos.platform.v1alpha1.PlatformService.UpdatePlatform
     */
    updatePlatform: {
      name: "UpdatePlatform",
      I: UpdatePlatformRequest,
      O: UpdatePlatformResponse,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc holos.platform.v1alpha1.PlatformService.ListPlatforms
     */
    listPlatforms: {
      name: "ListPlatforms",
      I: ListPlatformsRequest,
      O: ListPlatformsResponse,
      kind: MethodKind.Unary,
    },
  }
} as const;


import { makeAnyClient, CallOptions, Transport } from '@connectrpc/connect'
import { createAsyncIterable } from '@connectrpc/connect/protocol'
import {
    ServiceType,
    PartialMessage,
    MethodInfoServerStreaming,
    MethodInfo,
    MethodInfoUnary,
    MethodKind,
    Message,
} from '@bufbuild/protobuf'
import { Observable } from 'rxjs'

export type ObservableClient<T extends ServiceType> = {
    [P in keyof T['methods']]: T['methods'][P] extends MethodInfoUnary<
        infer I,
        infer O
    >
        ? UnaryFn<I, O>
        : T['methods'][P] extends MethodInfoServerStreaming<infer I, infer O>
        ? ServerStreamingFn<I, O>
        : never
}

export function createObservableClient<T extends ServiceType>(
    service: T,
    transport: Transport
) {
    return makeAnyClient(service, (method) => {
        switch (method.kind) {
            case MethodKind.Unary:
                return createUnaryFn(transport, service, method)
            case MethodKind.ServerStreaming:
                return createServerStreamingFn(transport, service, method)
            default:
                return null
        }
    }) as ObservableClient<T>
}

type UnaryFn<I extends Message<I>, O extends Message<O>> = (
    request: PartialMessage<I>,
    options?: CallOptions
) => Observable<O>

function createUnaryFn<I extends Message<I>, O extends Message<O>>(
    transport: Transport,
    service: ServiceType,
    method: MethodInfo<I, O>
): UnaryFn<I, O> {
    return function (requestMessage, options) {
        return new Observable<O>((subscriber) => {
            transport
                .unary(
                    service,
                    method,
                    options?.signal,
                    options?.timeoutMs,
                    options?.headers,
                    requestMessage
                )
                .then(
                    (response) => {
                        options?.onHeader?.(response.header)
                        subscriber.next(response.message)
                        options?.onTrailer?.(response.trailer)
                    },
                    (err) => {
                        subscriber.error(err)
                    }
                )
                .finally(() => {
                    subscriber.complete()
                })
        })
    }
}

type ServerStreamingFn<I extends Message<I>, O extends Message<O>> = (
    request: PartialMessage<I>,
    options?: CallOptions
) => Observable<O>

export function createServerStreamingFn<
    I extends Message<I>,
    O extends Message<O>
>(
    transport: Transport,
    service: ServiceType,
    method: MethodInfo<I, O>
): ServerStreamingFn<I, O> {
    return function (input, options) {
        return new Observable<O>((subscriber) => {
            transport
                .stream<I, O>(
                    service,
                    method,
                    options?.signal,
                    options?.timeoutMs,
                    options?.headers,
                    createAsyncIterable([input])
                )
                .then(
                    async (streamResponse) => {
                        options?.onHeader?.(streamResponse.header)
                        for await (const response of streamResponse.message) {
                            subscriber.next(response)
                        }
                        options?.onTrailer?.(streamResponse.trailer)
                    },
                    (err) => {
                        subscriber.error(err)
                    }
                )
                .finally(() => {
                    subscriber.complete()
                })
        })
    }
}

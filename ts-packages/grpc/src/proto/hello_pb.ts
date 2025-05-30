// proto/helloworld.proto

// @generated by protoc-gen-es v2.3.0 with parameter "target=ts"
// @generated from file proto/hello.proto (package helloworld, syntax proto3)
/* eslint-disable */

import type { GenFile, GenMessage, GenService } from "@bufbuild/protobuf/codegenv1";
import { fileDesc, messageDesc, serviceDesc } from "@bufbuild/protobuf/codegenv1";
import type { Message } from "@bufbuild/protobuf";

/**
 * Describes the file proto/hello.proto.
 */
export const file_proto_hello: GenFile = /*@__PURE__*/
  fileDesc("ChFwcm90by9oZWxsby5wcm90bxIKaGVsbG93b3JsZCIcCgxIZWxsb1JlcXVlc3QSDAoEbmFtZRgBIAEoCSIdCgpIZWxsb1JlcGx5Eg8KB21lc3NhZ2UYASABKAkyRwoHR3JlZXRlchI8CghTYXlIZWxsbxIYLmhlbGxvd29ybGQuSGVsbG9SZXF1ZXN0GhYuaGVsbG93b3JsZC5IZWxsb1JlcGx5YgZwcm90bzM");

/**
 * @generated from message helloworld.HelloRequest
 */
export type HelloRequest = Message<"helloworld.HelloRequest"> & {
  /**
   * @generated from field: string name = 1;
   */
  name: string;
};

/**
 * Describes the message helloworld.HelloRequest.
 * Use `create(HelloRequestSchema)` to create a new message.
 */
export const HelloRequestSchema: GenMessage<HelloRequest> = /*@__PURE__*/
  messageDesc(file_proto_hello, 0);

/**
 * @generated from message helloworld.HelloReply
 */
export type HelloReply = Message<"helloworld.HelloReply"> & {
  /**
   * @generated from field: string message = 1;
   */
  message: string;
};

/**
 * Describes the message helloworld.HelloReply.
 * Use `create(HelloReplySchema)` to create a new message.
 */
export const HelloReplySchema: GenMessage<HelloReply> = /*@__PURE__*/
  messageDesc(file_proto_hello, 1);

/**
 * @generated from service helloworld.Greeter
 */
export const Greeter: GenService<{
  /**
   * @generated from rpc helloworld.Greeter.SayHello
   */
  sayHello: {
    methodKind: "unary";
    input: typeof HelloRequestSchema;
    output: typeof HelloReplySchema;
  },
}> = /*@__PURE__*/
  serviceDesc(file_proto_hello, 0);


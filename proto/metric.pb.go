// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: proto/metric.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Metric struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    string  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`         // имя метрики
	MType string  `protobuf:"bytes,2,opt,name=mType,proto3" json:"mType,omitempty"`   // параметр, принимающий значение gauge или counter
	Delta int64   `protobuf:"varint,3,opt,name=delta,proto3" json:"delta,omitempty"`  // значение метрики в случае передачи counter
	Value float64 `protobuf:"fixed64,4,opt,name=value,proto3" json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  string  `protobuf:"bytes,5,opt,name=hash,proto3" json:"hash,omitempty"`     // значение хеш-функции
}

func (x *Metric) Reset() {
	*x = Metric{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_metric_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Metric) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Metric) ProtoMessage() {}

func (x *Metric) ProtoReflect() protoreflect.Message {
	mi := &file_proto_metric_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Metric.ProtoReflect.Descriptor instead.
func (*Metric) Descriptor() ([]byte, []int) {
	return file_proto_metric_proto_rawDescGZIP(), []int{0}
}

func (x *Metric) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Metric) GetMType() string {
	if x != nil {
		return x.MType
	}
	return ""
}

func (x *Metric) GetDelta() int64 {
	if x != nil {
		return x.Delta
	}
	return 0
}

func (x *Metric) GetValue() float64 {
	if x != nil {
		return x.Value
	}
	return 0
}

func (x *Metric) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

type UpdateSingleRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metric *Metric `protobuf:"bytes,1,opt,name=metric,proto3" json:"metric,omitempty"` // метрика для сохранения/обновления
}

func (x *UpdateSingleRequest) Reset() {
	*x = UpdateSingleRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_metric_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateSingleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateSingleRequest) ProtoMessage() {}

func (x *UpdateSingleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_metric_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateSingleRequest.ProtoReflect.Descriptor instead.
func (*UpdateSingleRequest) Descriptor() ([]byte, []int) {
	return file_proto_metric_proto_rawDescGZIP(), []int{1}
}

func (x *UpdateSingleRequest) GetMetric() *Metric {
	if x != nil {
		return x.Metric
	}
	return nil
}

type UpdateSingleResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error string `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"` // ошибка обновления метрики
}

func (x *UpdateSingleResponse) Reset() {
	*x = UpdateSingleResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_metric_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateSingleResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateSingleResponse) ProtoMessage() {}

func (x *UpdateSingleResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_metric_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateSingleResponse.ProtoReflect.Descriptor instead.
func (*UpdateSingleResponse) Descriptor() ([]byte, []int) {
	return file_proto_metric_proto_rawDescGZIP(), []int{2}
}

func (x *UpdateSingleResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type UpdateBatchRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metrics []*Metric `protobuf:"bytes,1,rep,name=metrics,proto3" json:"metrics,omitempty"` // массив метрик для сохранения/обновления
}

func (x *UpdateBatchRequest) Reset() {
	*x = UpdateBatchRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_metric_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateBatchRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateBatchRequest) ProtoMessage() {}

func (x *UpdateBatchRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_metric_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateBatchRequest.ProtoReflect.Descriptor instead.
func (*UpdateBatchRequest) Descriptor() ([]byte, []int) {
	return file_proto_metric_proto_rawDescGZIP(), []int{3}
}

func (x *UpdateBatchRequest) GetMetrics() []*Metric {
	if x != nil {
		return x.Metrics
	}
	return nil
}

type UpdateBatchResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error string `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"` // ошибка обновления массива метрик
}

func (x *UpdateBatchResponse) Reset() {
	*x = UpdateBatchResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_metric_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateBatchResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateBatchResponse) ProtoMessage() {}

func (x *UpdateBatchResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_metric_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateBatchResponse.ProtoReflect.Descriptor instead.
func (*UpdateBatchResponse) Descriptor() ([]byte, []int) {
	return file_proto_metric_proto_rawDescGZIP(), []int{4}
}

func (x *UpdateBatchResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type GetSingleRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metric *Metric `protobuf:"bytes,1,opt,name=metric,proto3" json:"metric,omitempty"` // метрика для запроса на чтение
}

func (x *GetSingleRequest) Reset() {
	*x = GetSingleRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_metric_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetSingleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetSingleRequest) ProtoMessage() {}

func (x *GetSingleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_metric_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetSingleRequest.ProtoReflect.Descriptor instead.
func (*GetSingleRequest) Descriptor() ([]byte, []int) {
	return file_proto_metric_proto_rawDescGZIP(), []int{5}
}

func (x *GetSingleRequest) GetMetric() *Metric {
	if x != nil {
		return x.Metric
	}
	return nil
}

type GetSingleResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metric *Metric `protobuf:"bytes,1,opt,name=metric,proto3" json:"metric,omitempty"` // значение запрошенной метрики
	Error  string  `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`   // ошибка чтения метрики
}

func (x *GetSingleResponse) Reset() {
	*x = GetSingleResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_metric_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetSingleResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetSingleResponse) ProtoMessage() {}

func (x *GetSingleResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_metric_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetSingleResponse.ProtoReflect.Descriptor instead.
func (*GetSingleResponse) Descriptor() ([]byte, []int) {
	return file_proto_metric_proto_rawDescGZIP(), []int{6}
}

func (x *GetSingleResponse) GetMetric() *Metric {
	if x != nil {
		return x.Metric
	}
	return nil
}

func (x *GetSingleResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type ListMetricsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ListMetricsRequest) Reset() {
	*x = ListMetricsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_metric_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListMetricsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListMetricsRequest) ProtoMessage() {}

func (x *ListMetricsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_metric_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListMetricsRequest.ProtoReflect.Descriptor instead.
func (*ListMetricsRequest) Descriptor() ([]byte, []int) {
	return file_proto_metric_proto_rawDescGZIP(), []int{7}
}

type ListMetricsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metrics []*Metric `protobuf:"bytes,1,rep,name=metrics,proto3" json:"metrics,omitempty"` // массив считанных метрик
	Error   string    `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`     // ошибка чтения метрик
}

func (x *ListMetricsResponse) Reset() {
	*x = ListMetricsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_metric_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListMetricsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListMetricsResponse) ProtoMessage() {}

func (x *ListMetricsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_metric_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListMetricsResponse.ProtoReflect.Descriptor instead.
func (*ListMetricsResponse) Descriptor() ([]byte, []int) {
	return file_proto_metric_proto_rawDescGZIP(), []int{8}
}

func (x *ListMetricsResponse) GetMetrics() []*Metric {
	if x != nil {
		return x.Metrics
	}
	return nil
}

func (x *ListMetricsResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type CheckRepoRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CheckRepoRequest) Reset() {
	*x = CheckRepoRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_metric_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckRepoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckRepoRequest) ProtoMessage() {}

func (x *CheckRepoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_metric_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckRepoRequest.ProtoReflect.Descriptor instead.
func (*CheckRepoRequest) Descriptor() ([]byte, []int) {
	return file_proto_metric_proto_rawDescGZIP(), []int{9}
}

type CheckRepoResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error string `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"` // ошибка соединения с базой данных
}

func (x *CheckRepoResponse) Reset() {
	*x = CheckRepoResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_metric_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckRepoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckRepoResponse) ProtoMessage() {}

func (x *CheckRepoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_metric_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckRepoResponse.ProtoReflect.Descriptor instead.
func (*CheckRepoResponse) Descriptor() ([]byte, []int) {
	return file_proto_metric_proto_rawDescGZIP(), []int{10}
}

func (x *CheckRepoResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

var File_proto_metric_proto protoreflect.FileDescriptor

var file_proto_metric_proto_rawDesc = []byte{
	0x0a, 0x12, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x22, 0x6e, 0x0a,
	0x06, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x6d, 0x54, 0x79, 0x70, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6d, 0x54, 0x79, 0x70, 0x65, 0x12, 0x14, 0x0a,
	0x05, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x64, 0x65,
	0x6c, 0x74, 0x61, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x01, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x61, 0x73,
	0x68, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x22, 0x3e, 0x0a,
	0x13, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x53, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x27, 0x0a, 0x06, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2e, 0x4d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x22, 0x2c, 0x0a,
	0x14, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x53, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x3f, 0x0a, 0x12, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x29, 0x0a, 0x07, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2e, 0x4d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x52, 0x07, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x22, 0x2b, 0x0a, 0x13,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x3b, 0x0a, 0x10, 0x47, 0x65, 0x74,
	0x53, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x27, 0x0a,
	0x06, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e,
	0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x06,
	0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x22, 0x52, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x53, 0x69, 0x6e,
	0x67, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x27, 0x0a, 0x06, 0x6d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x6d, 0x65,
	0x74, 0x72, 0x69, 0x63, 0x73, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x06, 0x6d, 0x65,
	0x74, 0x72, 0x69, 0x63, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x14, 0x0a, 0x12, 0x4c, 0x69,
	0x73, 0x74, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x22, 0x56, 0x0a, 0x13, 0x4c, 0x69, 0x73, 0x74, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x29, 0x0a, 0x07, 0x6d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x73, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x07, 0x6d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x12, 0x0a, 0x10, 0x43, 0x68, 0x65, 0x63,
	0x6b, 0x52, 0x65, 0x70, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x29, 0x0a, 0x11,
	0x43, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x65, 0x70, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x32, 0xf2, 0x02, 0x0a, 0x07, 0x4d, 0x65, 0x74, 0x72,
	0x69, 0x63, 0x73, 0x12, 0x4b, 0x0a, 0x0c, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x53, 0x69, 0x6e,
	0x67, 0x6c, 0x65, 0x12, 0x1c, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2e, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x53, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x1d, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2e, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x53, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x48, 0x0a, 0x0b, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x42, 0x61, 0x74, 0x63, 0x68, 0x12,
	0x1b, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x6d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x42, 0x61, 0x74,
	0x63, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x42, 0x0a, 0x09, 0x47, 0x65,
	0x74, 0x53, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x12, 0x19, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63,
	0x73, 0x2e, 0x47, 0x65, 0x74, 0x53, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2e, 0x47, 0x65, 0x74,
	0x53, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x48,
	0x0a, 0x0b, 0x4c, 0x69, 0x73, 0x74, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x12, 0x1b, 0x2e,
	0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x4d, 0x65, 0x74, 0x72,
	0x69, 0x63, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x6d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x73, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x42, 0x0a, 0x09, 0x43, 0x68, 0x65, 0x63,
	0x6b, 0x52, 0x65, 0x70, 0x6f, 0x12, 0x19, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2e,
	0x43, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x65, 0x70, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x1a, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2e, 0x43, 0x68, 0x65, 0x63, 0x6b,
	0x52, 0x65, 0x70, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x15, 0x5a, 0x13,
	0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2d, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_metric_proto_rawDescOnce sync.Once
	file_proto_metric_proto_rawDescData = file_proto_metric_proto_rawDesc
)

func file_proto_metric_proto_rawDescGZIP() []byte {
	file_proto_metric_proto_rawDescOnce.Do(func() {
		file_proto_metric_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_metric_proto_rawDescData)
	})
	return file_proto_metric_proto_rawDescData
}

var file_proto_metric_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_proto_metric_proto_goTypes = []interface{}{
	(*Metric)(nil),               // 0: metrics.Metric
	(*UpdateSingleRequest)(nil),  // 1: metrics.UpdateSingleRequest
	(*UpdateSingleResponse)(nil), // 2: metrics.UpdateSingleResponse
	(*UpdateBatchRequest)(nil),   // 3: metrics.UpdateBatchRequest
	(*UpdateBatchResponse)(nil),  // 4: metrics.UpdateBatchResponse
	(*GetSingleRequest)(nil),     // 5: metrics.GetSingleRequest
	(*GetSingleResponse)(nil),    // 6: metrics.GetSingleResponse
	(*ListMetricsRequest)(nil),   // 7: metrics.ListMetricsRequest
	(*ListMetricsResponse)(nil),  // 8: metrics.ListMetricsResponse
	(*CheckRepoRequest)(nil),     // 9: metrics.CheckRepoRequest
	(*CheckRepoResponse)(nil),    // 10: metrics.CheckRepoResponse
}
var file_proto_metric_proto_depIdxs = []int32{
	0,  // 0: metrics.UpdateSingleRequest.metric:type_name -> metrics.Metric
	0,  // 1: metrics.UpdateBatchRequest.metrics:type_name -> metrics.Metric
	0,  // 2: metrics.GetSingleRequest.metric:type_name -> metrics.Metric
	0,  // 3: metrics.GetSingleResponse.metric:type_name -> metrics.Metric
	0,  // 4: metrics.ListMetricsResponse.metrics:type_name -> metrics.Metric
	1,  // 5: metrics.Metrics.UpdateSingle:input_type -> metrics.UpdateSingleRequest
	3,  // 6: metrics.Metrics.UpdateBatch:input_type -> metrics.UpdateBatchRequest
	5,  // 7: metrics.Metrics.GetSingle:input_type -> metrics.GetSingleRequest
	7,  // 8: metrics.Metrics.ListMetrics:input_type -> metrics.ListMetricsRequest
	9,  // 9: metrics.Metrics.CheckRepo:input_type -> metrics.CheckRepoRequest
	2,  // 10: metrics.Metrics.UpdateSingle:output_type -> metrics.UpdateSingleResponse
	4,  // 11: metrics.Metrics.UpdateBatch:output_type -> metrics.UpdateBatchResponse
	6,  // 12: metrics.Metrics.GetSingle:output_type -> metrics.GetSingleResponse
	8,  // 13: metrics.Metrics.ListMetrics:output_type -> metrics.ListMetricsResponse
	10, // 14: metrics.Metrics.CheckRepo:output_type -> metrics.CheckRepoResponse
	10, // [10:15] is the sub-list for method output_type
	5,  // [5:10] is the sub-list for method input_type
	5,  // [5:5] is the sub-list for extension type_name
	5,  // [5:5] is the sub-list for extension extendee
	0,  // [0:5] is the sub-list for field type_name
}

func init() { file_proto_metric_proto_init() }
func file_proto_metric_proto_init() {
	if File_proto_metric_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_metric_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Metric); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_metric_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateSingleRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_metric_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateSingleResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_metric_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateBatchRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_metric_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateBatchResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_metric_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetSingleRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_metric_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetSingleResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_metric_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListMetricsRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_metric_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListMetricsResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_metric_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckRepoRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_metric_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckRepoResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_metric_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_metric_proto_goTypes,
		DependencyIndexes: file_proto_metric_proto_depIdxs,
		MessageInfos:      file_proto_metric_proto_msgTypes,
	}.Build()
	File_proto_metric_proto = out.File
	file_proto_metric_proto_rawDesc = nil
	file_proto_metric_proto_goTypes = nil
	file_proto_metric_proto_depIdxs = nil
}

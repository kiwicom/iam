// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api/grpc/v1/kiwi_iamapi.proto

package userv1

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type UserRequest struct {
	Email                string   `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UserRequest) Reset()         { *m = UserRequest{} }
func (m *UserRequest) String() string { return proto.CompactTextString(m) }
func (*UserRequest) ProtoMessage()    {}
func (*UserRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1c6f8aa01589961e, []int{0}
}

func (m *UserRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserRequest.Unmarshal(m, b)
}
func (m *UserRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserRequest.Marshal(b, m, deterministic)
}
func (m *UserRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserRequest.Merge(m, src)
}
func (m *UserRequest) XXX_Size() int {
	return xxx_messageInfo_UserRequest.Size(m)
}
func (m *UserRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UserRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UserRequest proto.InternalMessageInfo

func (m *UserRequest) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

type BoocsekAttributes struct {
	Site                 string   `protobuf:"bytes,1,opt,name=site,proto3" json:"site,omitempty"`
	Position             string   `protobuf:"bytes,2,opt,name=position,proto3" json:"position,omitempty"`
	Channel              string   `protobuf:"bytes,3,opt,name=channel,proto3" json:"channel,omitempty"`
	Tier                 string   `protobuf:"bytes,4,opt,name=tier,proto3" json:"tier,omitempty"`
	Team                 string   `protobuf:"bytes,5,opt,name=team,proto3" json:"team,omitempty"`
	TeamManager          string   `protobuf:"bytes,6,opt,name=team_manager,json=teamManager,proto3" json:"team_manager,omitempty"`
	Staff                string   `protobuf:"bytes,7,opt,name=staff,proto3" json:"staff,omitempty"`
	State                string   `protobuf:"bytes,8,opt,name=state,proto3" json:"state,omitempty"`
	KiwibaseId           int32    `protobuf:"varint,9,opt,name=kiwibase_id,json=kiwibaseId,proto3" json:"kiwibase_id,omitempty"`
	Substate             string   `protobuf:"bytes,10,opt,name=substate,proto3" json:"substate,omitempty"`
	Skills               []string `protobuf:"bytes,11,rep,name=skills,proto3" json:"skills,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BoocsekAttributes) Reset()         { *m = BoocsekAttributes{} }
func (m *BoocsekAttributes) String() string { return proto.CompactTextString(m) }
func (*BoocsekAttributes) ProtoMessage()    {}
func (*BoocsekAttributes) Descriptor() ([]byte, []int) {
	return fileDescriptor_1c6f8aa01589961e, []int{1}
}

func (m *BoocsekAttributes) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BoocsekAttributes.Unmarshal(m, b)
}
func (m *BoocsekAttributes) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BoocsekAttributes.Marshal(b, m, deterministic)
}
func (m *BoocsekAttributes) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BoocsekAttributes.Merge(m, src)
}
func (m *BoocsekAttributes) XXX_Size() int {
	return xxx_messageInfo_BoocsekAttributes.Size(m)
}
func (m *BoocsekAttributes) XXX_DiscardUnknown() {
	xxx_messageInfo_BoocsekAttributes.DiscardUnknown(m)
}

var xxx_messageInfo_BoocsekAttributes proto.InternalMessageInfo

func (m *BoocsekAttributes) GetSite() string {
	if m != nil {
		return m.Site
	}
	return ""
}

func (m *BoocsekAttributes) GetPosition() string {
	if m != nil {
		return m.Position
	}
	return ""
}

func (m *BoocsekAttributes) GetChannel() string {
	if m != nil {
		return m.Channel
	}
	return ""
}

func (m *BoocsekAttributes) GetTier() string {
	if m != nil {
		return m.Tier
	}
	return ""
}

func (m *BoocsekAttributes) GetTeam() string {
	if m != nil {
		return m.Team
	}
	return ""
}

func (m *BoocsekAttributes) GetTeamManager() string {
	if m != nil {
		return m.TeamManager
	}
	return ""
}

func (m *BoocsekAttributes) GetStaff() string {
	if m != nil {
		return m.Staff
	}
	return ""
}

func (m *BoocsekAttributes) GetState() string {
	if m != nil {
		return m.State
	}
	return ""
}

func (m *BoocsekAttributes) GetKiwibaseId() int32 {
	if m != nil {
		return m.KiwibaseId
	}
	return 0
}

func (m *BoocsekAttributes) GetSubstate() string {
	if m != nil {
		return m.Substate
	}
	return ""
}

func (m *BoocsekAttributes) GetSkills() []string {
	if m != nil {
		return m.Skills
	}
	return nil
}

type UserResponse struct {
	EmployeeNumber       int64              `protobuf:"varint,1,opt,name=employee_number,json=employeeNumber,proto3" json:"employee_number,omitempty"`
	Email                string             `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty"`
	FirstName            string             `protobuf:"bytes,3,opt,name=first_name,json=firstName,proto3" json:"first_name,omitempty"`
	LastName             string             `protobuf:"bytes,4,opt,name=last_name,json=lastName,proto3" json:"last_name,omitempty"`
	Position             string             `protobuf:"bytes,5,opt,name=position,proto3" json:"position,omitempty"`
	Department           string             `protobuf:"bytes,6,opt,name=department,proto3" json:"department,omitempty"`
	Location             string             `protobuf:"bytes,7,opt,name=location,proto3" json:"location,omitempty"`
	IsVendor             bool               `protobuf:"varint,8,opt,name=is_vendor,json=isVendor,proto3" json:"is_vendor,omitempty"`
	Manager              string             `protobuf:"bytes,9,opt,name=manager,proto3" json:"manager,omitempty"`
	TeamMembership       []string           `protobuf:"bytes,10,rep,name=team_membership,json=teamMembership,proto3" json:"team_membership,omitempty"`
	Boocsek              *BoocsekAttributes `protobuf:"bytes,11,opt,name=boocsek,proto3" json:"boocsek,omitempty"`
	Permissions          []string           `protobuf:"bytes,12,rep,name=permissions,proto3" json:"permissions,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *UserResponse) Reset()         { *m = UserResponse{} }
func (m *UserResponse) String() string { return proto.CompactTextString(m) }
func (*UserResponse) ProtoMessage()    {}
func (*UserResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1c6f8aa01589961e, []int{2}
}

func (m *UserResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserResponse.Unmarshal(m, b)
}
func (m *UserResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserResponse.Marshal(b, m, deterministic)
}
func (m *UserResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserResponse.Merge(m, src)
}
func (m *UserResponse) XXX_Size() int {
	return xxx_messageInfo_UserResponse.Size(m)
}
func (m *UserResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_UserResponse.DiscardUnknown(m)
}

var xxx_messageInfo_UserResponse proto.InternalMessageInfo

func (m *UserResponse) GetEmployeeNumber() int64 {
	if m != nil {
		return m.EmployeeNumber
	}
	return 0
}

func (m *UserResponse) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *UserResponse) GetFirstName() string {
	if m != nil {
		return m.FirstName
	}
	return ""
}

func (m *UserResponse) GetLastName() string {
	if m != nil {
		return m.LastName
	}
	return ""
}

func (m *UserResponse) GetPosition() string {
	if m != nil {
		return m.Position
	}
	return ""
}

func (m *UserResponse) GetDepartment() string {
	if m != nil {
		return m.Department
	}
	return ""
}

func (m *UserResponse) GetLocation() string {
	if m != nil {
		return m.Location
	}
	return ""
}

func (m *UserResponse) GetIsVendor() bool {
	if m != nil {
		return m.IsVendor
	}
	return false
}

func (m *UserResponse) GetManager() string {
	if m != nil {
		return m.Manager
	}
	return ""
}

func (m *UserResponse) GetTeamMembership() []string {
	if m != nil {
		return m.TeamMembership
	}
	return nil
}

func (m *UserResponse) GetBoocsek() *BoocsekAttributes {
	if m != nil {
		return m.Boocsek
	}
	return nil
}

func (m *UserResponse) GetPermissions() []string {
	if m != nil {
		return m.Permissions
	}
	return nil
}

func init() {
	proto.RegisterType((*UserRequest)(nil), "kiwi.iam.user.v1.UserRequest")
	proto.RegisterType((*BoocsekAttributes)(nil), "kiwi.iam.user.v1.BoocsekAttributes")
	proto.RegisterType((*UserResponse)(nil), "kiwi.iam.user.v1.UserResponse")
}

func init() { proto.RegisterFile("api/grpc/v1/kiwi_iamapi.proto", fileDescriptor_1c6f8aa01589961e) }

var fileDescriptor_1c6f8aa01589961e = []byte{
	// 551 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x93, 0x4f, 0x6b, 0xdb, 0x30,
	0x18, 0xc6, 0x49, 0xd2, 0xfc, 0xf1, 0xeb, 0xd2, 0x76, 0xa2, 0x0c, 0xd1, 0xd1, 0x2e, 0x6b, 0x0f,
	0xcb, 0xc9, 0x25, 0xdd, 0x79, 0x87, 0x04, 0x76, 0x30, 0xa5, 0x25, 0x78, 0x34, 0x87, 0x91, 0x61,
	0xe4, 0xe4, 0x4d, 0x2b, 0x62, 0x59, 0x9e, 0xa4, 0xa4, 0xec, 0xeb, 0xf4, 0xb6, 0x7d, 0x94, 0x7d,
	0xaa, 0x21, 0xc9, 0xce, 0x42, 0xc3, 0x4e, 0xf6, 0xf3, 0x7b, 0xf5, 0x87, 0xf7, 0x79, 0xf4, 0xc2,
	0x39, 0x2b, 0xf9, 0xf5, 0xa3, 0x2a, 0xe7, 0xd7, 0x9b, 0xe1, 0xf5, 0x8a, 0x3f, 0xf3, 0x94, 0x33,
	0xc1, 0x4a, 0x1e, 0x95, 0x4a, 0x1a, 0x49, 0x4e, 0x2c, 0x8a, 0x38, 0x13, 0xd1, 0x5a, 0xa3, 0x8a,
	0x36, 0xc3, 0xcb, 0x2b, 0x08, 0x1f, 0x34, 0xaa, 0x04, 0x7f, 0xac, 0x51, 0x1b, 0x72, 0x0a, 0x6d,
	0x14, 0x8c, 0xe7, 0xb4, 0xd1, 0x6f, 0x0c, 0x82, 0xc4, 0x8b, 0xcb, 0x5f, 0x4d, 0x78, 0x33, 0x96,
	0x72, 0xae, 0x71, 0x35, 0x32, 0x46, 0xf1, 0x6c, 0x6d, 0x50, 0x13, 0x02, 0x07, 0x9a, 0x1b, 0xac,
	0x96, 0xba, 0x7f, 0x72, 0x06, 0xbd, 0x52, 0x6a, 0x6e, 0xb8, 0x2c, 0x68, 0xd3, 0xf1, 0xad, 0x26,
	0x14, 0xba, 0xf3, 0x27, 0x56, 0x14, 0x98, 0xd3, 0x96, 0x2b, 0xd5, 0xd2, 0x9e, 0x64, 0x38, 0x2a,
	0x7a, 0xe0, 0x4f, 0xb2, 0xff, 0x8e, 0x21, 0x13, 0xb4, 0x5d, 0x31, 0x64, 0x82, 0x7c, 0x80, 0x43,
	0xfb, 0x4d, 0x05, 0x2b, 0xd8, 0x23, 0x2a, 0xda, 0x71, 0xb5, 0xd0, 0xb2, 0x3b, 0x8f, 0x6c, 0x03,
	0xda, 0xb0, 0xe5, 0x92, 0x76, 0x7d, 0x03, 0x4e, 0x54, 0xd4, 0x20, 0xed, 0x6d, 0xa9, 0x41, 0xf2,
	0x1e, 0x42, 0xeb, 0x47, 0xc6, 0x34, 0xa6, 0x7c, 0x41, 0x83, 0x7e, 0x63, 0xd0, 0x4e, 0xa0, 0x46,
	0xf1, 0xc2, 0x76, 0xa3, 0xd7, 0x99, 0xdf, 0x09, 0xbe, 0x9b, 0x5a, 0x93, 0xb7, 0xd0, 0xd1, 0x2b,
	0x9e, 0xe7, 0x9a, 0x86, 0xfd, 0xd6, 0x20, 0x48, 0x2a, 0x75, 0xf9, 0xd2, 0x82, 0x43, 0xef, 0xa8,
	0x2e, 0x65, 0xa1, 0x91, 0x7c, 0x84, 0x63, 0x14, 0x65, 0x2e, 0x7f, 0x22, 0xa6, 0xc5, 0x5a, 0x64,
	0xa8, 0x9c, 0x63, 0xad, 0xe4, 0xa8, 0xc6, 0xf7, 0x8e, 0xfe, 0xf3, 0xbe, 0xb9, 0xe3, 0x3d, 0x39,
	0x07, 0x58, 0x72, 0xa5, 0x4d, 0x5a, 0x30, 0x81, 0x95, 0x71, 0x81, 0x23, 0xf7, 0x4c, 0x20, 0x79,
	0x07, 0x41, 0xce, 0xea, 0xaa, 0xf7, 0xaf, 0x67, 0x81, 0x2b, 0xee, 0xa6, 0xd1, 0x7e, 0x95, 0xc6,
	0x05, 0xc0, 0x02, 0x4b, 0xa6, 0x8c, 0xc0, 0xc2, 0x54, 0x4e, 0xee, 0x10, 0xbb, 0x37, 0x97, 0x73,
	0xe6, 0xf6, 0x76, 0xab, 0x73, 0x2b, 0x6d, 0x2f, 0xe5, 0x3a, 0xdd, 0x60, 0xb1, 0x90, 0xca, 0x59,
	0xda, 0x4b, 0x7a, 0x5c, 0x4f, 0x9d, 0xb6, 0x31, 0xd7, 0xf9, 0x04, 0x3e, 0xe6, 0x4a, 0x5a, 0x27,
	0x7c, 0x7c, 0x68, 0xfb, 0xd5, 0x4f, 0xbc, 0xa4, 0xe0, 0xbc, 0x3b, 0x72, 0x09, 0x6e, 0x29, 0xf9,
	0x0c, 0xdd, 0xcc, 0x3f, 0x37, 0x1a, 0xf6, 0x1b, 0x83, 0xf0, 0xe6, 0x2a, 0x7a, 0xfd, 0x70, 0xa3,
	0xbd, 0xf7, 0x98, 0xd4, 0x7b, 0x48, 0x1f, 0xc2, 0x12, 0x95, 0xe0, 0x5a, 0x73, 0x59, 0x68, 0x7a,
	0xe8, 0xee, 0xd8, 0x45, 0x37, 0x5f, 0x01, 0x6e, 0xf9, 0x33, 0x8f, 0x47, 0x77, 0xa3, 0x49, 0x4c,
	0xbe, 0xc0, 0x81, 0x4d, 0x8c, 0x9c, 0xef, 0xdf, 0xb2, 0x33, 0x1b, 0x67, 0x17, 0xff, 0x2b, 0xfb,
	0xa0, 0xc7, 0xdf, 0xe1, 0x74, 0x2e, 0xc5, 0xde, 0xa2, 0xf1, 0xb1, 0xbb, 0xca, 0x8d, 0xe1, 0xc4,
	0x4e, 0xe1, 0xa4, 0xf1, 0xad, 0x63, 0x6b, 0x9b, 0xe1, 0x4b, 0xb3, 0x75, 0x1b, 0x3f, 0xfc, 0x6e,
	0x9e, 0xd8, 0x15, 0x51, 0xcc, 0x84, 0x3b, 0x30, 0x9a, 0x0e, 0xff, 0x78, 0x34, 0x8b, 0x99, 0x98,
	0x59, 0x34, 0x9b, 0x0e, 0xb3, 0x8e, 0x1b, 0xe1, 0x4f, 0x7f, 0x03, 0x00, 0x00, 0xff, 0xff, 0xb2,
	0x97, 0x29, 0xc3, 0xe3, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// KiwiIAMAPIClient is the client API for KiwiIAMAPI service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type KiwiIAMAPIClient interface {
	// User retrieves a Kiwi user information from OKTA.
	User(ctx context.Context, in *UserRequest, opts ...grpc.CallOption) (*UserResponse, error)
}

type kiwiIAMAPIClient struct {
	cc *grpc.ClientConn
}

func NewKiwiIAMAPIClient(cc *grpc.ClientConn) KiwiIAMAPIClient {
	return &kiwiIAMAPIClient{cc}
}

func (c *kiwiIAMAPIClient) User(ctx context.Context, in *UserRequest, opts ...grpc.CallOption) (*UserResponse, error) {
	out := new(UserResponse)
	err := c.cc.Invoke(ctx, "/kiwi.iam.user.v1.KiwiIAMAPI/User", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// KiwiIAMAPIServer is the server API for KiwiIAMAPI service.
type KiwiIAMAPIServer interface {
	// User retrieves a Kiwi user information from OKTA.
	User(context.Context, *UserRequest) (*UserResponse, error)
}

// UnimplementedKiwiIAMAPIServer can be embedded to have forward compatible implementations.
type UnimplementedKiwiIAMAPIServer struct {
}

func (*UnimplementedKiwiIAMAPIServer) User(ctx context.Context, req *UserRequest) (*UserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method User not implemented")
}

func RegisterKiwiIAMAPIServer(s *grpc.Server, srv KiwiIAMAPIServer) {
	s.RegisterService(&_KiwiIAMAPI_serviceDesc, srv)
}

func _KiwiIAMAPI_User_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KiwiIAMAPIServer).User(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kiwi.iam.user.v1.KiwiIAMAPI/User",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KiwiIAMAPIServer).User(ctx, req.(*UserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _KiwiIAMAPI_serviceDesc = grpc.ServiceDesc{
	ServiceName: "kiwi.iam.user.v1.KiwiIAMAPI",
	HandlerType: (*KiwiIAMAPIServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "User",
			Handler:    _KiwiIAMAPI_User_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/grpc/v1/kiwi_iamapi.proto",
}

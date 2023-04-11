// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: profile.proto

package __

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ProfileClient is the client API for Profile service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ProfileClient interface {
	GetUserProfile(ctx context.Context, in *UserID, opts ...grpc.CallOption) (*ProfileData, error)
	EditProfile(ctx context.Context, in *EditProfileData, opts ...grpc.CallOption) (*Empty, error)
	EditAvatar(ctx context.Context, in *EditAvatarData, opts ...grpc.CallOption) (*Empty, error)
	UploadAvatar(ctx context.Context, in *UploadInputFile, opts ...grpc.CallOption) (*FileName, error)
	GetAvatar(ctx context.Context, in *UserID, opts ...grpc.CallOption) (*FileName, error)
	AcceptInvitationToFamily(ctx context.Context, in *AddToFamily, opts ...grpc.CallOption) (*Empty, error)
	CreateFamily(ctx context.Context, in *UserID, opts ...grpc.CallOption) (*Empty, error)
	DeleteFamily(ctx context.Context, in *UserID, opts ...grpc.CallOption) (*Empty, error)
	DeleteFromFamily(ctx context.Context, in *Delete, opts ...grpc.CallOption) (*Empty, error)
	DeleteMember(ctx context.Context, in *Delete, opts ...grpc.CallOption) (*Empty, error)
	AddMember(ctx context.Context, in *MemberData, opts ...grpc.CallOption) (*Empty, error)
	GetFamily(ctx context.Context, in *UserID, opts ...grpc.CallOption) (*ResponseMemberDataArr, error)
	HasFamily(ctx context.Context, in *UserID, opts ...grpc.CallOption) (*HasFamilyResp, error)
	UserExists(ctx context.Context, in *EmailData, opts ...grpc.CallOption) (*Exists, error)
	AddMedicine(ctx context.Context, in *AddMed, opts ...grpc.CallOption) (*Empty, error)
	DeleteMedicine(ctx context.Context, in *DeleteMed, opts ...grpc.CallOption) (*Empty, error)
	GetMedicine(ctx context.Context, in *UserID, opts ...grpc.CallOption) (*MedicineArr, error)
}

type profileClient struct {
	cc grpc.ClientConnInterface
}

func NewProfileClient(cc grpc.ClientConnInterface) ProfileClient {
	return &profileClient{cc}
}

func (c *profileClient) GetUserProfile(ctx context.Context, in *UserID, opts ...grpc.CallOption) (*ProfileData, error) {
	out := new(ProfileData)
	err := c.cc.Invoke(ctx, "/profile.Profile/GetUserProfile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileClient) EditProfile(ctx context.Context, in *EditProfileData, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/profile.Profile/EditProfile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileClient) EditAvatar(ctx context.Context, in *EditAvatarData, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/profile.Profile/EditAvatar", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileClient) UploadAvatar(ctx context.Context, in *UploadInputFile, opts ...grpc.CallOption) (*FileName, error) {
	out := new(FileName)
	err := c.cc.Invoke(ctx, "/profile.Profile/UploadAvatar", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileClient) GetAvatar(ctx context.Context, in *UserID, opts ...grpc.CallOption) (*FileName, error) {
	out := new(FileName)
	err := c.cc.Invoke(ctx, "/profile.Profile/GetAvatar", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileClient) AcceptInvitationToFamily(ctx context.Context, in *AddToFamily, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/profile.Profile/AcceptInvitationToFamily", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileClient) CreateFamily(ctx context.Context, in *UserID, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/profile.Profile/CreateFamily", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileClient) DeleteFamily(ctx context.Context, in *UserID, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/profile.Profile/DeleteFamily", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileClient) DeleteFromFamily(ctx context.Context, in *Delete, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/profile.Profile/DeleteFromFamily", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileClient) DeleteMember(ctx context.Context, in *Delete, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/profile.Profile/DeleteMember", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileClient) AddMember(ctx context.Context, in *MemberData, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/profile.Profile/AddMember", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileClient) GetFamily(ctx context.Context, in *UserID, opts ...grpc.CallOption) (*ResponseMemberDataArr, error) {
	out := new(ResponseMemberDataArr)
	err := c.cc.Invoke(ctx, "/profile.Profile/GetFamily", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileClient) HasFamily(ctx context.Context, in *UserID, opts ...grpc.CallOption) (*HasFamilyResp, error) {
	out := new(HasFamilyResp)
	err := c.cc.Invoke(ctx, "/profile.Profile/HasFamily", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileClient) UserExists(ctx context.Context, in *EmailData, opts ...grpc.CallOption) (*Exists, error) {
	out := new(Exists)
	err := c.cc.Invoke(ctx, "/profile.Profile/UserExists", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileClient) AddMedicine(ctx context.Context, in *AddMed, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/profile.Profile/AddMedicine", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileClient) DeleteMedicine(ctx context.Context, in *DeleteMed, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/profile.Profile/DeleteMedicine", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileClient) GetMedicine(ctx context.Context, in *UserID, opts ...grpc.CallOption) (*MedicineArr, error) {
	out := new(MedicineArr)
	err := c.cc.Invoke(ctx, "/profile.Profile/GetMedicine", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProfileServer is the server API for Profile service.
// All implementations should embed UnimplementedProfileServer
// for forward compatibility
type ProfileServer interface {
	GetUserProfile(context.Context, *UserID) (*ProfileData, error)
	EditProfile(context.Context, *EditProfileData) (*Empty, error)
	EditAvatar(context.Context, *EditAvatarData) (*Empty, error)
	UploadAvatar(context.Context, *UploadInputFile) (*FileName, error)
	GetAvatar(context.Context, *UserID) (*FileName, error)
	AcceptInvitationToFamily(context.Context, *AddToFamily) (*Empty, error)
	CreateFamily(context.Context, *UserID) (*Empty, error)
	DeleteFamily(context.Context, *UserID) (*Empty, error)
	DeleteFromFamily(context.Context, *Delete) (*Empty, error)
	DeleteMember(context.Context, *Delete) (*Empty, error)
	AddMember(context.Context, *MemberData) (*Empty, error)
	GetFamily(context.Context, *UserID) (*ResponseMemberDataArr, error)
	HasFamily(context.Context, *UserID) (*HasFamilyResp, error)
	UserExists(context.Context, *EmailData) (*Exists, error)
	AddMedicine(context.Context, *AddMed) (*Empty, error)
	DeleteMedicine(context.Context, *DeleteMed) (*Empty, error)
	GetMedicine(context.Context, *UserID) (*MedicineArr, error)
}

// UnimplementedProfileServer should be embedded to have forward compatible implementations.
type UnimplementedProfileServer struct {
}

func (UnimplementedProfileServer) GetUserProfile(context.Context, *UserID) (*ProfileData, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserProfile not implemented")
}
func (UnimplementedProfileServer) EditProfile(context.Context, *EditProfileData) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EditProfile not implemented")
}
func (UnimplementedProfileServer) EditAvatar(context.Context, *EditAvatarData) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EditAvatar not implemented")
}
func (UnimplementedProfileServer) UploadAvatar(context.Context, *UploadInputFile) (*FileName, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UploadAvatar not implemented")
}
func (UnimplementedProfileServer) GetAvatar(context.Context, *UserID) (*FileName, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAvatar not implemented")
}
func (UnimplementedProfileServer) AcceptInvitationToFamily(context.Context, *AddToFamily) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AcceptInvitationToFamily not implemented")
}
func (UnimplementedProfileServer) CreateFamily(context.Context, *UserID) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateFamily not implemented")
}
func (UnimplementedProfileServer) DeleteFamily(context.Context, *UserID) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteFamily not implemented")
}
func (UnimplementedProfileServer) DeleteFromFamily(context.Context, *Delete) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteFromFamily not implemented")
}
func (UnimplementedProfileServer) DeleteMember(context.Context, *Delete) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteMember not implemented")
}
func (UnimplementedProfileServer) AddMember(context.Context, *MemberData) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddMember not implemented")
}
func (UnimplementedProfileServer) GetFamily(context.Context, *UserID) (*ResponseMemberDataArr, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFamily not implemented")
}
func (UnimplementedProfileServer) HasFamily(context.Context, *UserID) (*HasFamilyResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HasFamily not implemented")
}
func (UnimplementedProfileServer) UserExists(context.Context, *EmailData) (*Exists, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserExists not implemented")
}
func (UnimplementedProfileServer) AddMedicine(context.Context, *AddMed) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddMedicine not implemented")
}
func (UnimplementedProfileServer) DeleteMedicine(context.Context, *DeleteMed) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteMedicine not implemented")
}
func (UnimplementedProfileServer) GetMedicine(context.Context, *UserID) (*MedicineArr, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMedicine not implemented")
}

// UnsafeProfileServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ProfileServer will
// result in compilation errors.
type UnsafeProfileServer interface {
	mustEmbedUnimplementedProfileServer()
}

func RegisterProfileServer(s grpc.ServiceRegistrar, srv ProfileServer) {
	s.RegisterService(&Profile_ServiceDesc, srv)
}

func _Profile_GetUserProfile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServer).GetUserProfile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/profile.Profile/GetUserProfile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServer).GetUserProfile(ctx, req.(*UserID))
	}
	return interceptor(ctx, in, info, handler)
}

func _Profile_EditProfile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EditProfileData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServer).EditProfile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/profile.Profile/EditProfile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServer).EditProfile(ctx, req.(*EditProfileData))
	}
	return interceptor(ctx, in, info, handler)
}

func _Profile_EditAvatar_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EditAvatarData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServer).EditAvatar(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/profile.Profile/EditAvatar",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServer).EditAvatar(ctx, req.(*EditAvatarData))
	}
	return interceptor(ctx, in, info, handler)
}

func _Profile_UploadAvatar_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UploadInputFile)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServer).UploadAvatar(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/profile.Profile/UploadAvatar",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServer).UploadAvatar(ctx, req.(*UploadInputFile))
	}
	return interceptor(ctx, in, info, handler)
}

func _Profile_GetAvatar_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServer).GetAvatar(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/profile.Profile/GetAvatar",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServer).GetAvatar(ctx, req.(*UserID))
	}
	return interceptor(ctx, in, info, handler)
}

func _Profile_AcceptInvitationToFamily_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddToFamily)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServer).AcceptInvitationToFamily(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/profile.Profile/AcceptInvitationToFamily",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServer).AcceptInvitationToFamily(ctx, req.(*AddToFamily))
	}
	return interceptor(ctx, in, info, handler)
}

func _Profile_CreateFamily_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServer).CreateFamily(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/profile.Profile/CreateFamily",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServer).CreateFamily(ctx, req.(*UserID))
	}
	return interceptor(ctx, in, info, handler)
}

func _Profile_DeleteFamily_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServer).DeleteFamily(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/profile.Profile/DeleteFamily",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServer).DeleteFamily(ctx, req.(*UserID))
	}
	return interceptor(ctx, in, info, handler)
}

func _Profile_DeleteFromFamily_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Delete)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServer).DeleteFromFamily(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/profile.Profile/DeleteFromFamily",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServer).DeleteFromFamily(ctx, req.(*Delete))
	}
	return interceptor(ctx, in, info, handler)
}

func _Profile_DeleteMember_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Delete)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServer).DeleteMember(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/profile.Profile/DeleteMember",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServer).DeleteMember(ctx, req.(*Delete))
	}
	return interceptor(ctx, in, info, handler)
}

func _Profile_AddMember_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MemberData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServer).AddMember(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/profile.Profile/AddMember",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServer).AddMember(ctx, req.(*MemberData))
	}
	return interceptor(ctx, in, info, handler)
}

func _Profile_GetFamily_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServer).GetFamily(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/profile.Profile/GetFamily",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServer).GetFamily(ctx, req.(*UserID))
	}
	return interceptor(ctx, in, info, handler)
}

func _Profile_HasFamily_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServer).HasFamily(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/profile.Profile/HasFamily",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServer).HasFamily(ctx, req.(*UserID))
	}
	return interceptor(ctx, in, info, handler)
}

func _Profile_UserExists_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmailData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServer).UserExists(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/profile.Profile/UserExists",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServer).UserExists(ctx, req.(*EmailData))
	}
	return interceptor(ctx, in, info, handler)
}

func _Profile_AddMedicine_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddMed)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServer).AddMedicine(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/profile.Profile/AddMedicine",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServer).AddMedicine(ctx, req.(*AddMed))
	}
	return interceptor(ctx, in, info, handler)
}

func _Profile_DeleteMedicine_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteMed)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServer).DeleteMedicine(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/profile.Profile/DeleteMedicine",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServer).DeleteMedicine(ctx, req.(*DeleteMed))
	}
	return interceptor(ctx, in, info, handler)
}

func _Profile_GetMedicine_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServer).GetMedicine(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/profile.Profile/GetMedicine",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServer).GetMedicine(ctx, req.(*UserID))
	}
	return interceptor(ctx, in, info, handler)
}

// Profile_ServiceDesc is the grpc.ServiceDesc for Profile service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Profile_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "profile.Profile",
	HandlerType: (*ProfileServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetUserProfile",
			Handler:    _Profile_GetUserProfile_Handler,
		},
		{
			MethodName: "EditProfile",
			Handler:    _Profile_EditProfile_Handler,
		},
		{
			MethodName: "EditAvatar",
			Handler:    _Profile_EditAvatar_Handler,
		},
		{
			MethodName: "UploadAvatar",
			Handler:    _Profile_UploadAvatar_Handler,
		},
		{
			MethodName: "GetAvatar",
			Handler:    _Profile_GetAvatar_Handler,
		},
		{
			MethodName: "AcceptInvitationToFamily",
			Handler:    _Profile_AcceptInvitationToFamily_Handler,
		},
		{
			MethodName: "CreateFamily",
			Handler:    _Profile_CreateFamily_Handler,
		},
		{
			MethodName: "DeleteFamily",
			Handler:    _Profile_DeleteFamily_Handler,
		},
		{
			MethodName: "DeleteFromFamily",
			Handler:    _Profile_DeleteFromFamily_Handler,
		},
		{
			MethodName: "DeleteMember",
			Handler:    _Profile_DeleteMember_Handler,
		},
		{
			MethodName: "AddMember",
			Handler:    _Profile_AddMember_Handler,
		},
		{
			MethodName: "GetFamily",
			Handler:    _Profile_GetFamily_Handler,
		},
		{
			MethodName: "HasFamily",
			Handler:    _Profile_HasFamily_Handler,
		},
		{
			MethodName: "UserExists",
			Handler:    _Profile_UserExists_Handler,
		},
		{
			MethodName: "AddMedicine",
			Handler:    _Profile_AddMedicine_Handler,
		},
		{
			MethodName: "DeleteMedicine",
			Handler:    _Profile_DeleteMedicine_Handler,
		},
		{
			MethodName: "GetMedicine",
			Handler:    _Profile_GetMedicine_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "profile.proto",
}

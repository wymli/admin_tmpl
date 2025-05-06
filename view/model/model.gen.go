package model

type ViewModelI interface {
	StructName() string
}

func (v LoginUserRsp) StructName() string {
	return "LoginUserRsp"
}

func (v GetUserRsp) StructName() string {
	return "GetUserRsp"
}

func (v PingReq) StructName() string {
	return "PingReq"
}

func (v EchoReq) StructName() string {
	return "EchoReq"
}

func (v LoginUserReq) StructName() string {
	return "LoginUserReq"
}

func (v GetUserReq) StructName() string {
	return "GetUserReq"
}

func (v AppContext) StructName() string {
	return "AppContext"
}

func (v PingRsp) StructName() string {
	return "PingRsp"
}

func (v EchoRsp) StructName() string {
	return "EchoRsp"
}

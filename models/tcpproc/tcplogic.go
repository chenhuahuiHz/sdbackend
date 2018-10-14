package tcpproc

import (
	"github.com/astaxie/beego"
	"prototcp/protoserver"
	"prototcp/typedefs"
	"prototcp/protomsg"
)

var IMConn *protoserver.ProtoClient
var ip string
var port int

func InitIMConn() {
	ip = beego.AppConfig.String("IMConn::ip")
	port = typedefs.MyAtoi(beego.AppConfig.String("IMConn::port"))

	IMConn = &protoserver.ProtoClient{
		IP:"",
		Port:0,
		Session:nil,
		State:0,
		PBMsgMgr:protomsg.NewPBMsgMgr(),
		}
	IMConn.Connect(ip, port)
}

func SendCommNotify(aimid uint32, aimtype int32, opttype int32, showtype int32, content string) (result string){

	if nil == IMConn {
		beego.Error("SendCommNotify failed: IMConn not ready")
		result = "IMConn not ready"
		return
	}

	if IMConn.State != protoserver.Connected {
		beego.Info("SendCommNotify: try to connect")
		IMConn.Connect(ip, port)
	}

	if IMConn.State != protoserver.Connected {
		beego.Error("SendCommNotify failed: IM not connected")
		result = "IM not connected"
		return
	}

	msg := &protomsg.PBIMComNotify{
		AimId:aimid,
		AimType:aimtype,
		OptType:opttype,
		ShowType:showtype,
		Content:content,
	}

	IMConn.Session.SendPBMsg(msg, uint32(aimid), 0)
	result = "OK"
	return
}
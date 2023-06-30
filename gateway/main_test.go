package main

import (
	"fmt"
	"io"
	"net"
	"service/gateway2/gutil"
	"service/gateway2/tcp_impl"
	"testing"
	"time"
	"utils/config"
)

func TestTCPServer(t *testing.T) {
	config.InitDBConfig()
	config.InitDedicateConf(&config.Cfg.Gateway, "gateway.json")

	gutil.MustInit()

	t.Log("Client Test ... start")

	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		t.Log("client start err, exit!")
		return
	}

	mid := uint32(0)
	for {
		//发封包message消息
		dp := tcp_impl.NewDataPack()
		msg, _ := dp.Pack(tcp_impl.NewMsgPackage(0, []byte(fmt.Sprintf("Lun V0.5 Client Test Message=%d", mid))))
		_, err := conn.Write(msg)
		if err != nil {
			t.Log("write error err ", err)
			return
		}
		mid++
		//先读出流中的head部分
		headData := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, headData) // ReadFull 会把msg填充满为止
		if err != nil {
			t.Log("read head error")
			break
		}
		//将headData字节流 拆包到msg中
		msgHead, err := dp.Unpack(headData)
		if err != nil {
			t.Log("server unpack err:", err)
			return
		}
		println("222", msgHead.GetMsgId(), msgHead.GetDataLen())
		if msgHead.GetDataLen() > 0 {
			//msg 是有data数据的，需要再次读取data数据
			msg := msgHead.(*tcp_impl.Message)
			msg.Data = make([]byte, msg.GetDataLen())

			//根据dataLen从io中读取字节流
			_, err := io.ReadFull(conn, msg.Data)
			if err != nil {
				t.Log("server unpack data err:", err)
				return
			}

			println("==> Recv Msg: ID=", msg.Id, ", len=", msg.DataLen, ", data=", string(msg.Data))
		}

		time.Sleep(time.Second * 30)
	}
}

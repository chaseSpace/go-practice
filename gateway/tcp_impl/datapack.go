package tcp_impl

import (
	"bytes"
	"encoding/binary"
	"service/gateway2/tcp_iface"
)

type DataPack struct{}

func NewDataPack() *DataPack {
	return &DataPack{}
}

const (
	DataSizeSelf = 4
	MsgIdSize    = 4
)

func (dp *DataPack) GetHeadLen() uint32 {
	//Id uint32(4字节) +  DataLen uint32(4字节)
	return DataSizeSelf + MsgIdSize
}

func (dp *DataPack) Pack(msg tcp_iface.IMessage) ([]byte, error) {
	dataBuff := bytes.NewBuffer([]byte{})

	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}

	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

func (dp *DataPack) Unpack(binaryData []byte) (tcp_iface.IMessage, error) {
	dataBuff := bytes.NewReader(binaryData)

	//只解压head的信息，得到dataLen和msgID
	msg := &Message{}

	//读dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	//读msgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	//这里只需要把head的数据拆包出来就可以了，然后再通过head的长度，再从conn读取一次数据
	return msg, nil
}

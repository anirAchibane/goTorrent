package peer

import (
	"encoding/binary"
	"io"

)

type messageID uint8

const (
    MsgChoke         messageID = 0
    MsgUnchoke       messageID = 1
    MsgInterested    messageID = 2
    MsgNotInterested messageID = 3
    MsgHave          messageID = 4
    MsgBitfield      messageID = 5
    MsgRequest       messageID = 6
    MsgPiece         messageID = 7
    MsgCancel        messageID = 8
)

type Message struct {
    Length int
    ID      messageID
    Payload []byte
}

func Parse_message(r io.Reader) (*Message,error){
    length_buff := make([]byte,4)
    _,err := io.ReadFull(r,length_buff)

    if err != nil{
        return nil, err
    }
    
    length := binary.BigEndian.Uint32(length_buff)
    if length == 0{
        return nil,nil // keep-alive message
    }

    message_buff := make([]byte,length)
    
    _ , err = io.ReadFull(r, message_buff)
    if err != nil{
        return nil, err
    }

    return &Message{
        Length: int(length),
        ID: messageID(message_buff[0]),
        Payload: message_buff[1:],
    }, nil 
}

func Serialize_message(m *Message) ([]byte,error){
    if m == nil{
        return make([]byte, 4), nil
    }

    buff := make([]byte,len(m.Payload)+5)
    
    binary.BigEndian.PutUint32(buff[0:4], uint32(m.Length))
    buff[4] = byte(m.ID)
    copy(buff[5:],m.Payload)

    return buff, nil
}
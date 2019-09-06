package datapackage

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

var (
	ERR_EOF  = errors.New("EOF")
	ERR_EXIT = errors.New("EXIT")
)

type Package struct {
	Version        [2]byte // 协议版本，暂定V1
	Length         int16   // 数据部分长度
	Timestamp      int64   // 时间戳
	HostnameLength int16   // 主机名长度
	Hostname       []byte  // 主机名
	TagLength      int16   // 标签长度
	Tag            []byte  // 标签
	Msg            []byte  // 日志数据
}

func init() {
	//fmt.Println("package datapackage init()")
}

func (p *Package) Pack(writer io.Writer) error {
	var err error
	err = binary.Write(writer, binary.BigEndian, &p.Version)
	err = binary.Write(writer, binary.BigEndian, &p.Length)
	err = binary.Write(writer, binary.BigEndian, &p.Timestamp)
	err = binary.Write(writer, binary.BigEndian, &p.HostnameLength)
	err = binary.Write(writer, binary.BigEndian, &p.Hostname)
	err = binary.Write(writer, binary.BigEndian, &p.TagLength)
	err = binary.Write(writer, binary.BigEndian, &p.Tag)
	err = binary.Write(writer, binary.BigEndian, &p.Msg)
	return err
}

func (p *Package) Unpack(reader io.Reader) error {
	var err error
	err = binary.Read(reader, binary.BigEndian, &p.Version)
	err = binary.Read(reader, binary.BigEndian, &p.Length)
	err = binary.Read(reader, binary.BigEndian, &p.Timestamp)
	err = binary.Read(reader, binary.BigEndian, &p.HostnameLength)
	p.Hostname = make([]byte, p.HostnameLength)
	err = binary.Read(reader, binary.BigEndian, &p.Hostname)
	err = binary.Read(reader, binary.BigEndian, &p.TagLength)
	p.Tag = make([]byte, p.TagLength)
	err = binary.Read(reader, binary.BigEndian, &p.Tag)
	p.Msg = make([]byte, p.Length-8-2-p.HostnameLength-2-p.TagLength)
	err = binary.Read(reader, binary.BigEndian, &p.Msg)
	return err
}

func (p *Package) String() string {
	return fmt.Sprintf("version:%s, length:%d, timestamp:%d, hostname:%s, tag:%s, msg:[%s]",
		p.Version,
		p.Length,
		p.Timestamp,
		p.Hostname,
		p.Tag,
		p.Msg,
	)
}

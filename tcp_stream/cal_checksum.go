package tcp_stream

import "fmt"

/*
1. 使用一个4字节变量（uint32）存储TCP报文段中的所有字节（包括首部和数据部分）相加的和。
2. 如果和的高16位不为0，则将其与低16位相加，直到高16位为0为止。
3. 将最终得到的和按位取反，得到校验和。
*/
func calculateChecksum(data []byte) uint16 {
	var sum uint32
	for i := 0; i < len(data); i++ {
		sum += uint32(data[i])
	}
	for (sum >> 16) > 0 {
		sum = (sum & 0xFFFF) + (sum >> 16)
	}
	return ^uint16(sum)
}

func runExampleChecksum() {
	// tcp字节流 示例
	data := []byte{
		0x00, 0x3c, 0x46, 0x00,
		0x06, 0xf4, 0xa8, 0x01,
	}

	checksum := calculateChecksum(data)

	fmt.Printf("Checksum: %04X\n", checksum) // 16b=2B=4个hex字符（每个hex字符表示4b）
}

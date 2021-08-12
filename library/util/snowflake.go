package snowflake

import (
	"fmt"
	"time"
)

// snowflake 雪花算法
// 1 位，不用。
// 二进制中最高位为符号位，我们生成的 ID 一般都是正整数，所以这个最高位固定是 0。
// 41 位，用来记录时间戳（毫秒）。
// 41 位 可以表示 2^41 - 1 个数字。
// 也就是说 41 位 可以表示 2^41 - 1 个毫秒的值，转化成单位年则是 (2^41 - 1) / (1000 * 60 * 60 * 24 * 365) 约为 69 年。
// 10 位，用来记录工作机器 ID。
// 可以部署在 2^10 共 1024 个节点，包括 5 位 DatacenterId 和 5 位 WorkerId。
// 12 位，序列号，用来记录同毫秒内产生的不同 id。
// 12 位 可以表示的最大正整数是 2^12 - 1 共 4095 个数字，来表示同一机器同一时间截（毫秒)内产生的 4095 个 ID 序号。
// Snowflake 可以保证：
// 所有生成的 ID 按时间趋势递增。


// EPOCH_OFF_SET 开始时间
const EPOCH_OFF_SET = 1293811200000
// SIGN_BITS 首位
const SIGN_BITS = 1
// TIMESTAMP_BITS 时间戳相减占用的位数
const TIMESTAMP_BITS = 41
// DATA_CENTER_BITS  数据中心位数
const DATA_CENTER_BITS = 5
// MACHINE_ID_BITS 机器位数
const MACHINE_ID_BITS = 5
// SEQUENCE_BITS 同一毫秒内自增的位数
const SEQUENCE_BITS = 12

var lastTimestamp int64
var lastId int64

// sequence 步长
var sequence = 1
// signLeftShift 标志位需要位移的长度
var signLeftShift = TIMESTAMP_BITS + DATA_CENTER_BITS + MACHINE_ID_BITS + SEQUENCE_BITS
// timestampLeftShift  时间差值需要位移的长度
var timestampLeftShift = DATA_CENTER_BITS + MACHINE_ID_BITS + SEQUENCE_BITS
// dataCenterLeftShift 数据中心需要位移的长度
var dataCenterLeftShift = MACHINE_ID_BITS + SEQUENCE_BITS
// machineLeftShift 机器需要位移的长度
var machineLeftShift = SEQUENCE_BITS
//最大自增数
//下边语法执行效果 等于 (1 << 12) - 1;
//-1 << self::SEQUENCE_BITS = 1111111111111111111111111111111111111111111111111111000000000000
var maxSequenceId int = (1 << SEQUENCE_BITS) -1
//最大机器数
var maxMachineId  = (1 << MACHINE_ID_BITS) - 1
//最大数据中心机器数
var maxDataCenterId = (1 << DATA_CENTER_BITS) - 1

func GetId(dataCenterId int, machineId int,consistency bool) int64 {
	nowTime := time.Now().UnixNano()
	if nowTime < lastTimestamp {
		fmt.Println(" 时间有回滚,同步时间范围")
		if consistency {
			for nowTime < lastTimestamp {
				nowTime = time.Now().UnixNano()
			}
		} else {

		}
	}
    if lastTimestamp == nowTime {
    	if sequence > maxSequenceId {
    		return GetId(dataCenterId,machineId,consistency)
		}
		sequence++
	} else {
		sequence = 0
	}
	lastTimestamp = nowTime
	timeDiff := lastTimestamp - EPOCH_OFF_SET
	lastId = int64(0 << signLeftShift) | int64(timeDiff << timestampLeftShift) | int64(dataCenterId << dataCenterLeftShift) | int64(machineId << machineLeftShift) | int64(sequence)
	return lastId
}

func Parse(id int64) {
    se :=  id | (0 << SEQUENCE_BITS)
    fmt.Println(se)
}
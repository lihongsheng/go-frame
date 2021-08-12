package log


import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

var LogOut LogInterface

func init()  {
	LogOut = NewZapFileLog("./","log")
}

/**
 * 对外提供接口,用于后续更改日志组件的时候，外部代码不用修改
 */
type LogInterface interface {
	Debug(message string, content ...interface{})
	Info(message string, content ...interface{})
	Waring(message string, content ...interface{})
	Error(message string, content ...interface{})
	Panic(message string, content ...interface{})
}

// 文件按日期分割写入，写入源
type rotatelogs struct {
	path string  // 日志路径
	fileTime int64 // 当前日志时间
	filePrefixName string
	writer *os.File
	lock sync.Mutex
}

// 按照日期分割文件记录日志
func (r *rotatelogs) fileName() string {
	if r.fileTime == 0 || (time.Now().Unix() - r.fileTime) >= 864000 {
		currentTime := time.Now()
		r.fileTime = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location()).Unix()
	}
	var sb strings.Builder
	sb.WriteString(r.path)
	sb.WriteString(r.filePrefixName)
	sb.WriteString("-")
	sb.WriteString(fmt.Sprintf("%d-%d-%d.log",time.Now().Year(),time.Now().Month(),time.Now().Day()))
	return sb.String()
}

// 判断是否需要打开新的文件句柄
func (r *rotatelogs) openExistingOrNew() error {
	fn := r.fileName()
	// 获取文件信息
	_,err := os.Stat(fn)
	// 文件不存在
	if os.IsNotExist(err) {
		err := os.MkdirAll(r.path, 0755)
		if err != nil {
			return fmt.Errorf("can't make directories for new logfile: %s", err)
		}
		//os.O_APPEND|os.O_WRONLY, 0644)
		f, err := os.OpenFile(fn, os.O_CREATE|os.O_APPEND|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return fmt.Errorf("can't open new logfile: %s", err)
		}
		r.writer = f
	} else if r.writer == nil {
		f, err := os.OpenFile(fn, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("can't open new logfile: %s", err)
		}
		r.writer = f
	}
	return nil
}

// 加锁写入，防止日志写串
func (r *rotatelogs) Write(p []byte) (n int, err error)  {
	defer r.lock.Unlock()
	r.lock.Lock()
	//创建新的日志文件，后续参考 https://github.com/natefinch/lumberjack/blob/v2.0/lumberjack.go 实现一个
	fmt.Println(string(p))
	if err := r.openExistingOrNew(); err != nil {
		return 0,err
	}

	l, e := r.writer.Write(p)
	return l,e
}

/**
 * zap 文件日志
 */
type ZapLog struct {
	lg *zap.Logger
}

// kafka 写入源
type KafkaLog struct {

}

func (k *KafkaLog) Write(p []byte) (n int, err error)  {
	// kfaka生产者发送日志
	return 0,nil
}

//// zap 支持多输出源，kafak 日志源
//func NewZapKfakaLog() *ZapLog {
//	zl := &ZapLog{}
//	ws := zapcore.AddSync(&KafkaLog{})
//	ec := zap.NewProductionEncoderConfig()
//	ec.EncodeTime = zapcore.ISO8601TimeEncoder
//	en := zapcore.NewJSONEncoder(ec)
//	core := zapcore.NewCore(en, ws, zapcore.DebugLevel)
//	logger := zap.New(core)
//	zl.lg = logger
//	return zl
//}

// zap 支持多输出源，这里先默认一个
func NewZapFileLog(p string,filePrefixName string) *ZapLog {
	zl := &ZapLog{}
	ws := zapcore.AddSync(&rotatelogs{
		path: p,
		filePrefixName: filePrefixName,
	})

	ec := zap.NewProductionEncoderConfig()
	ec.EncodeTime = zapcore.ISO8601TimeEncoder
	en := zapcore.NewJSONEncoder(ec)
	core := zapcore.NewCore(en, ws, zapcore.DebugLevel)
	logger := zap.New(core,zap.AddCaller())
	zl.lg = logger
	return zl
}

// 调用ZAP的核心log方法
//func (z *ZapLog) log(level Level,message string,content[] interface{})  {
//}

func (z *ZapLog) Debug(message string,content ...interface{})  {
	_, file, line, _ := runtime.Caller(2)
	z.lg.Sugar().Debugw(message,"content",content,file,"line",line)
}

// 目前先用 zap的Sugar()写入
func (z *ZapLog) Info(message string,content ...interface{})  {
	_, file, line, _ := runtime.Caller(2)
	z.lg.Sugar().Infow(message,"content",content,file,"line",line)
}

func (z *ZapLog) Waring(message string,content ...interface{})  {
	_, file, line, _ := runtime.Caller(2)
	z.lg.Sugar().Warnw(message,"content",content,file,"line",line)
}

func (z *ZapLog) Error(message string,content ...interface{})  {
	_, file, line, _ := runtime.Caller(2)
	z.lg.Sugar().Errorw(message,"content",content,"file",file,"line",line)
}

func (z *ZapLog) Panic(message string,content ...interface{})  {
	_, file, line, _ := runtime.Caller(2)
	z.lg.Sugar().Panicw(message,"content",content,"file",file,"line",line)
}

func getLogFile() ( string , int ) {
	_, file, line, _ := runtime.Caller(2)
	return file,line
}

/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/6 0006
 * Time: 10:19
 */
package corp

const (
    // 消息
    MessageTypeNews       = "news"               // 类型-图文
    MessageTypeMpNews     = "mpnews"             // 类型-图文
    MessageTypeText       = "text"               // 类型-文本
    MessageTypeVoice      = "voice"              // 类型-语音
    MessageTypeImage      = "image"              // 类型-图片
    MessageTypeVideo      = "video"              // 类型-视频
    MessageTypeFile       = "file"               // 类型-文件
    MessageTypeTextCard   = "textcard"           // 类型-文本卡片
    MessageTypeMarkdown   = "markdown"           // 类型-markdown
    MessageTypeMiniNotice = "miniprogram_notice" // 类型-小程序通知

    // 发票
    InvoiceReimburseStatusInit    = "INVOICE_REIMBURSE_INIT"    // 报销状态-未锁定
    InvoiceReimburseStatusLock    = "INVOICE_REIMBURSE_LOCK"    // 报销状态-已锁定
    InvoiceReimburseStatusClosure = "INVOICE_REIMBURSE_CLOSURE" // 报销状态-已核销
)

var (
    MessageTypes               map[string]string
    InvoiceReimburseStatusList map[string]string
)

func init() {
    MessageTypes = make(map[string]string)
    MessageTypes[MessageTypeNews] = "图文"
    MessageTypes[MessageTypeMpNews] = "图文"
    MessageTypes[MessageTypeText] = "文本"
    MessageTypes[MessageTypeVoice] = "语音"
    MessageTypes[MessageTypeImage] = "图片"
    MessageTypes[MessageTypeVideo] = "视频"
    MessageTypes[MessageTypeFile] = "文件"
    MessageTypes[MessageTypeTextCard] = "文本卡片"
    MessageTypes[MessageTypeMarkdown] = "markdown"
    MessageTypes[MessageTypeMiniNotice] = "小程序通知"

    InvoiceReimburseStatusList = make(map[string]string)
    InvoiceReimburseStatusList[InvoiceReimburseStatusInit] = "未锁定"
    InvoiceReimburseStatusList[InvoiceReimburseStatusLock] = "已锁定"
    InvoiceReimburseStatusList[InvoiceReimburseStatusClosure] = "已核销"
}

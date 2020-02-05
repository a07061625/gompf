/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/2 0002
 * Time: 19:40
 */
package studio

import (
    "strconv"

    "github.com/a07061625/gompf/mpf/api/mpiot"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 新建产品
type productCreate struct {
    mpiot.BaseTencent
    productName    string // 产品名称
    productType    int    // 产品类型
    productDesc    string // 产品描述
    categoryId     int    // 分组模板ID
    encryptionType string // 加密类型
    netType        string // 连接类型
    dataProtocol   int    // 数据协议
    projectId      int    // 项目ID
}

func (pc *productCreate) SetProductName(productName string) {
    if len(productName) > 0 {
        pc.productName = productName
    } else {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "产品名称不合法", nil))
    }
}

func (pc *productCreate) SetProductType(productType int) {
    if productType > 0 {
        pc.productType = productType
    } else {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "产品类型不合法", nil))
    }
}

func (pc *productCreate) SetProductDesc(productDesc string) {
    if len(productDesc) > 0 {
        pc.productDesc = productDesc
    } else {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "产品描述不合法", nil))
    }
}

func (pc *productCreate) SetCategoryId(categoryId int) {
    if categoryId > 0 {
        pc.categoryId = categoryId
    } else {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "分组模板ID不合法", nil))
    }
}

func (pc *productCreate) SetEncryptionType(encryptionType string) {
    if len(encryptionType) > 0 {
        pc.encryptionType = encryptionType
    } else {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "加密类型不合法", nil))
    }
}

func (pc *productCreate) SetNetType(netType string) {
    if len(netType) > 0 {
        pc.netType = netType
    } else {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "连接类型不合法", nil))
    }
}

func (pc *productCreate) SetDataProtocol(dataProtocol int) {
    if dataProtocol > 0 {
        pc.dataProtocol = dataProtocol
    } else {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "数据协议不合法", nil))
    }
}

func (pc *productCreate) SetProjectId(projectId int) {
    if projectId > 0 {
        pc.projectId = projectId
    } else {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "项目ID不合法", nil))
    }
}

func (pc *productCreate) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(pc.productName) == 0 {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "产品名称不能为空", nil))
    }
    if pc.productType <= 0 {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "产品类型不能为空", nil))
    }
    if len(pc.productDesc) == 0 {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "产品描述不能为空", nil))
    }
    if pc.categoryId <= 0 {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "分组模板ID不能为空", nil))
    }
    if len(pc.encryptionType) == 0 {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "加密类型不能为空", nil))
    }
    if len(pc.netType) == 0 {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "连接类型不能为空", nil))
    }
    if pc.dataProtocol <= 0 {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "数据协议不能为空", nil))
    }
    if pc.projectId <= 0 {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "项目ID不能为空", nil))
    }
    pc.ReqData["ProjectName"] = pc.productName
    pc.ReqData["ProductType"] = strconv.Itoa(pc.productType)
    pc.ReqData["ProductDesc"] = pc.productDesc
    pc.ReqData["CategoryId"] = strconv.Itoa(pc.categoryId)
    pc.ReqData["EncryptionType"] = pc.encryptionType
    pc.ReqData["NetType"] = pc.netType
    pc.ReqData["DataProtocol"] = strconv.Itoa(pc.dataProtocol)
    pc.ReqData["ProjectId"] = strconv.Itoa(pc.projectId)

    return pc.GetRequest()
}

func NewProductCreate() *productCreate {
    pc := &productCreate{mpiot.NewBaseTencent(), "", 0, "", 0, "", "", 0, 0}
    pc.ReqHeader["X-TC-Action"] = "CreateStudioProduct"
    return pc
}

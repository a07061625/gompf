package mpf

import (
    "encoding/xml"
    "io"
    "strings"
)

const (
    XmlKeyName  = "_name"
    XmlKeyCData = "_cdata"
)

type XmlMap map[string]string

type xmlEntry struct {
    XMLName xml.Name
    Value   string `xml:",innerxml"`
}

// 重写xml编码方法
// 用法:
// <pre>
// xmlMap := map[string]string{
//     "key1": "One",
//     "key2": "Two",
// }
// mapStr, _ := xml.Marshal(tool.XmlMap(xmlMap))
// </pre>
func (m XmlMap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if len(m) == 0 {
        return nil
    }

    startName, ok := m[XmlKeyName]
    if ok {
        start.Name.Local = startName
    } else {
        start.Name.Local = "xml"
    }
    err := e.EncodeToken(start)
    if err != nil {
        return err
    }

    xmlCData, ok := m[XmlKeyCData]
    delete(m, XmlKeyName)
    delete(m, XmlKeyCData)

    if ok && (xmlCData == "N") {
        for k, v := range m {
            e.Encode(xmlEntry{XMLName: xml.Name{Local: k}, Value: v})
        }
    } else {
        for k, v := range m {
            e.Encode(xmlEntry{XMLName: xml.Name{Local: k}, Value: "<![CDATA[" + v + "]]>"})
        }
    }

    return e.EncodeToken(start.End())
}

// 重写xml解码方法
// 用法:
// <pre>
// newMap := make(map[string]string)
// xml.Unmarshal(mapStr, (*tool.XmlMap)(&newMap))
// </pre>
func (m *XmlMap) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    *m = XmlMap{}
    for {
        var e xmlEntry

        err := d.Decode(&e)
        if err == io.EOF {
            break
        } else if err != nil {
            return err
        }

        eVal := ""
        if (strings.HasPrefix(e.Value, "<![CDATA[")) && (strings.HasSuffix(e.Value, "]]>")) {
            valTmp := strings.TrimPrefix(e.Value, "<![CDATA[")
            eVal = strings.TrimSuffix(valTmp, "]]>")
        } else {
            eVal = e.Value
        }
        (*m)[e.XMLName.Local] = eVal
    }
    return nil
}

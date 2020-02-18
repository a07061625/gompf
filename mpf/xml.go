package mpf

import (
    "encoding/xml"
    "io"
    "strings"
)

const (
    // XMLKeyName KeyName
    XMLKeyName = "_name"
    // XMLKeyCData KeyCData
    XMLKeyCData = "_cdata"
)

// XMLMap XMLMap
type XMLMap map[string]string

type xmlEntry struct {
    XMLName xml.Name
    Value   string `xml:",innerxml"`
}

// MarshalXML 重写xml编码方法
// 用法:
// <pre>
// xmlMap := map[string]string{
//     "key1": "One",
//     "key2": "Two",
// }
// mapStr, _ := xml.Marshal(tool.XmlMap(xmlMap))
// </pre>
func (m XMLMap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if len(m) == 0 {
        return nil
    }

    startName, ok := m[XMLKeyName]
    if ok {
        start.Name.Local = startName
    } else {
        start.Name.Local = "xml"
    }
    err := e.EncodeToken(start)
    if err != nil {
        return err
    }

    xmlCData, ok := m[XMLKeyCData]
    delete(m, XMLKeyName)
    delete(m, XMLKeyCData)

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

// UnmarshalXML 重写xml解码方法
// 用法:
// <pre>
// newMap := make(map[string]string)
// xml.Unmarshal(mapStr, (*tool.XmlMap)(&newMap))
// </pre>
func (m *XMLMap) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    *m = XMLMap{}
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

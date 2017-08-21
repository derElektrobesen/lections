package main

import (
	"encoding/xml"
	"strconv"
)

type jigurdaConfigXML struct {
	XMLName      xml.Name
	JigurdaName  string                  `xml:"name,attr"`
	MobileConfig *mobileJigurdaConfigXML `xml:"mobile,omitempty"`
	WebConfig    *disableableConfigXML   `xml:"web,omitempty"`
	Phones       []string                `xml:"phones>phone"`
	Permille     uint32                  `xml:"permille"`
}

type disableableConfigXML struct {
	Disabled bool `xml:"disabled,attr,omitempty"`
}

type versionTypeXML struct {
	version []int
}

type mobileVersionJigurdaConfigXML struct {
	disableableConfigXML
	MinVersion *versionTypeXML `xml:"min_version,attr,omitempty"`
	MaxVersion *versionTypeXML `xml:"max_version,attr,omitempty"`
}

type mobileJigurdaConfigXML struct {
	disableableConfigXML
	AndroidConfig *mobileVersionJigurdaConfigXML `xml:"android,omitempty"`
	IOSConfig     *mobileVersionJigurdaConfigXML `xml:"ios,omitempty"`
}

func (c disableableConfig) toXML() disableableConfigXML {
	return disableableConfigXML{
		Disabled: c.Disabled,
	}
}

func (c mobileJigurdaConfig) toXML() *mobileJigurdaConfigXML {
	m := mobileJigurdaConfigXML{
		disableableConfigXML: c.disableableConfig.toXML(),
	}

	if android := c.AndroidConfig; android != nil {
		m.AndroidConfig = android.toXML()
	}

	if ios := c.IOSConfig; ios != nil {
		m.IOSConfig = ios.toXML()
	}

	return &m
}

func (c mobileVersionJigurdaConfig) toXML() *mobileVersionJigurdaConfigXML {
	var minV, maxV *versionTypeXML

	if v := []int(c.MinVersion.versionType); len(v) > 0 {
		minV = &versionTypeXML{version: v}
	}

	if v := []int(c.MaxVersion.versionType); len(v) > 0 {
		maxV = &versionTypeXML{version: v}
	}

	return &mobileVersionJigurdaConfigXML{
		disableableConfigXML: c.disableableConfig.toXML(),
		MinVersion:           minV,
		MaxVersion:           maxV,
	}
}

func (v versionTypeXML) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	var ver string
	for _, x := range v.version {
		if len(ver) > 0 {
			ver += "."
		}

		ver += strconv.Itoa(x)
	}

	return xml.Attr{name, ver}, nil
}

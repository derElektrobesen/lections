package main

import (
	"fmt"
	"hash/crc32"
	"log"

	"gitlab.corp.mail.ru/ride-share/backend/types"

	"strconv"
	"strings"
)

// PlatformType declares application platform type
type PlatformType string

const (
	PlatformAndroid PlatformType = "android"
	PlatformIOS     PlatformType = "ios"
	PlatformWeb     PlatformType = "web"
)

type DeviceInfo struct {
	Id       string // Device id (it could be changed for the same device)
	Platform PlatformType
	OsVer    string
	AppVer   string
	AppBuild string
}

func (p PlatformType) IsMobile() bool {
	return p == PlatformAndroid || p == PlatformIOS
}

func (i DeviceInfo) IsMobileApp() bool {
	return i.Platform.IsMobile()
}

type disableableConfig struct {
	Disabled bool `yaml:"disabled"`
}

func (dc disableableConfig) disabled() bool {
	return dc.Disabled
}

type versionType []int

type minVersionType struct {
	versionType
}

// inJigurda returns true if version > minVersion or version is incorrect (for example "1.2-alpha", "")
// or minVersion is not specified
func (v minVersionType) inJigurda(version string) bool {
	if len(v.versionType) == 0 {
		return true
	}

	versionSlice, err := convertVersionToVersionTypeWithoutLastZeros(version)
	if err != nil {
		log.Printf("unexpected device version %q found in jigurda check (minVersion): %s", version, err)
		return true
	}

	result := v.compareWithFixedVersion(versionSlice)
	switch result {
	case versionLess:
		// version < minVersion
		return false
	case versionGreat:
		// version >= minVersion
		fallthrough
	case versionEquals:
		return true
	default:
		log.Printf("unexpected compareWithMaxResult = %d returned compareWithFixedVersion function (minVersion)",
			result)
		return false
	}
}

type compareWithMaxResult int

const (
	versionLess compareWithMaxResult = iota + 1
	versionGreat
	versionEquals
)

// compareWithFixedVersion compares versionSlice and v, where v is fixedVersion
// For example: v = [1 2], versionSlice = [1 3], result = versionGreat
func (v versionType) compareWithFixedVersion(versionSlice versionType) compareWithMaxResult {
	for i := range v {
		// for example v = 1.2; versionSlice = 1
		if i >= len(versionSlice) {
			return versionLess
		}

		if versionSlice[i] < v[i] {
			return versionLess
		} else if versionSlice[i] > v[i] {
			return versionGreat
		}
	}

	// for example v = 1.2.3; version = 1.2.3.4
	if len(versionSlice) > len(v) {
		return versionGreat
	}

	return versionEquals
}

type maxVersionType struct {
	versionType
}

// inJigurda returns true if version < maxVersion or version is incorrect (for example "1.2-alpha", "")
// or max is not specified
func (v maxVersionType) inJigurda(version string) bool {
	if len(v.versionType) == 0 {
		return true
	}

	versionSlice, err := convertVersionToVersionTypeWithoutLastZeros(version)
	if err != nil {
		log.Printf("unexpected device version %q found in jigurda check (maxVersion): %s", version)
		return true
	}

	result := v.compareWithFixedVersion(versionSlice)
	switch result {
	case versionLess:
		// version <= versionMax
		fallthrough
	case versionEquals:
		return true
	case versionGreat:
		return false
	default:
		log.Printf("unexpected compareWithMaxResult = %d returned compareWithFixedVersion function (maxVersion)",
			result)
		return false
	}
}

type mobileVersionJigurdaConfig struct {
	disableableConfig `yaml:",inline"`
	MinVersion        minVersionType `yaml:"min_version"`
	MaxVersion        maxVersionType `yaml:"max_version"`
}

type arrToString map[string]bool

type webJigurdaConfig struct {
	disableableConfig `yaml:",inline"`
}

type mobileJigurdaConfig struct {
	disableableConfig `yaml:",inline"`
	AndroidConfig     *mobileVersionJigurdaConfig `yaml:"android"`
	IOSConfig         *mobileVersionJigurdaConfig `yaml:"ios"`
}

type phonesJigurdaConfig map[string]bool

type jigurdaConfig struct {
	MobileConfig *mobileJigurdaConfig `yaml:"mobile"`
	WebConfig    *webJigurdaConfig    `yaml:"web"`
	Phones       phonesJigurdaConfig  `yaml:"phones"`
	Permille     uint32               `yaml:"permille"`
}

func (wc *webJigurdaConfig) inJigurda(userAgent string) bool {
	if wc == nil {
		return true
	}

	if wc.disabled() {
		return false
	}

	return true
}

func (mc *mobileJigurdaConfig) inJigurda(version string, platform PlatformType) bool {
	if mc == nil {
		return true
	}

	if mc.disabled() {
		return false
	}

	if platform == PlatformAndroid {
		return mc.AndroidConfig.inJigurda(version)
	}

	if platform == PlatformIOS {
		return mc.IOSConfig.inJigurda(version)
	}

	return true
}

func (mc *mobileVersionJigurdaConfig) inJigurda(version string) bool {
	// we return true when jigurda config not found
	if mc == nil {
		return true
	}
	return !mc.Disabled && mc.MinVersion.inJigurda(version) && mc.MaxVersion.inJigurda(version)
}
func InJigurda(jigurdaConfigs map[string]jigurdaConfig,
	name string,
	info *DeviceInfo,
	userID types.UserID,
	phone string,
) bool {
	if info == nil {
		return true
	}

	cfg, ok := jigurdaConfigs[name]
	if !ok {
		return true
	}

	if cfg.Phones.inJigurda(phone) {
		return true
	}

	if info.IsMobileApp() {
		if !cfg.MobileConfig.inJigurda(info.AppVer, info.Platform) {
			return false
		}
	} else {
		// HACK: in controller/device.go we set deviceID = userAgent for web
		if !cfg.WebConfig.inJigurda(info.Id) {
			return false
		}
	}

	return inJigurdaByUserID(userID, name, cfg.Permille)
}

func inJigurdaByUserID(userID types.UserID, name string, permille uint32) bool {
	var maxPermille uint32 = 1000
	if permille == maxPermille {
		return true
	}

	if userID == 0 {
		// in some cases we haven't got user id (and don't need it)
		// Don't include them in jigurda
		return false
	}

	data := fmt.Sprintf("%s_%d", name, userID)
	hashSum := crc32.ChecksumIEEE([]byte(data))

	return hashSum%maxPermille < permille
}

func (p *phonesJigurdaConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if *p == nil {
		*p = phonesJigurdaConfig{}
	}

	var phones []string
	if err := unmarshal(&phones); err != nil {
		return fmt.Errorf("can't unmarshal phones list for jigurda: %s", err)
	}

	for _, phone := range phones {
		(*p)[phone] = true
	}

	return nil
}

func (p phonesJigurdaConfig) inJigurda(phone string) bool {
	return p != nil && phone != "" && p[phone]
}

func (v *versionType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var version string
	if err := unmarshal(&version); err != nil {
		return fmt.Errorf("can't unmarshal version: %s", err)
	}

	if version == "" {
		return nil
	}

	versionSlice, err := convertVersionToVersionTypeWithoutLastZeros(version)
	if err != nil {
		return fmt.Errorf("can't convert %q to slice of ints: %s", version, err)
	}

	*v = versionSlice
	return nil
}

func convertVersionToVersionTypeWithoutLastZeros(version string) (versionType, error) {
	versionParts := strings.Split(version, ".")
	versionSlice := make(versionType, len(versionParts))
	for i, versionPart := range versionParts {
		value, err := strconv.Atoi(versionPart)
		if err != nil {
			return nil, fmt.Errorf("can't convert %q to int (version = %q): %s", versionPart, version, err)
		}
		versionSlice[i] = value
	}

	// remove last zeros from slice. For example: [1 2 0 0] --> [1 2]
	sliceLen := len(versionSlice)
	for i := sliceLen - 1; i >= 0; i-- {
		if versionSlice[i] != 0 {
			break
		} else {
			sliceLen--
		}
	}

	return versionSlice[:sliceLen], nil
}

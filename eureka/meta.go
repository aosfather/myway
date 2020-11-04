package eureka

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/clbanning/x2j"
	"io"
	"strconv"
	"time"
)

// EurekaURLSlugs is a map of resource names->Eureka URLs.
var EurekaURLSlugs = map[string]string{
	"Apps":                        "apps",
	"Instances":                   "instances",
	"InstancesByVIPAddress":       "vips",
	"InstancesBySecureVIPAddress": "svips",
}

// EurekaConnection is the settings required to make Eureka requests.
type EurekaConnection struct {
	ServiceUrls    []string
	ServicePort    int
	ServerURLBase  string
	Timeout        time.Duration
	PollInterval   time.Duration
	PreferSameZone bool
	Retries        int
	DNSDiscovery   bool
	DiscoveryZone  string
	discoveryTtl   chan struct{}
	UseJson        bool
}

// GetAppsResponseJson lets us deserialize the eureka/v2/apps response JSON—a wrapped GetAppsResponse.
type GetAppsResponseJson struct {
	Response *GetAppsResponse `json:"applications"`
}

// GetAppsResponse lets us deserialize the eureka/v2/apps response XML.
type GetAppsResponse struct {
	Applications  []*Application `xml:"application" json:"application"`
	AppsHashcode  string         `xml:"apps__hashcode" json:"apps__hashcode"`
	VersionsDelta int            `xml:"versions__delta" json:"versions__delta"`
}

// GetAppResponseJson wraps an Application for deserializing from Eureka JSON.
type GetAppResponseJson struct {
	Application Application `json:"application"`
}

// Application deserializeable from Eureka XML.
type Application struct {
	Name      string      `xml:"name" json:"name"`
	Instances []*Instance `xml:"instance" json:"instance"`
}

// StatusType is an enum of the different statuses allowed by Eureka.
type StatusType string

// Supported statuses
const (
	UP           StatusType = "UP"
	DOWN         StatusType = "DOWN"
	STARTING     StatusType = "STARTING"
	OUTOFSERVICE StatusType = "OUT_OF_SERVICE"
	UNKNOWN      StatusType = "UNKNOWN"
)

// Datacenter names
const (
	Amazon = "Amazon"
	MyOwn  = "MyOwn"
)

// RegisterInstanceJson lets us serialize the eureka/v2/apps/<ins> request JSON—a wrapped Instance.
type RegisterInstanceJson struct {
	Instance *Instance `json:"instance"`
}

// Instance [de]serializeable [to|from] Eureka [XML|JSON].
type Instance struct {
	InstanceId       string `xml:"instanceId" json:"instanceId"`
	HostName         string `xml:"hostName" json:"hostName"`
	App              string `xml:"app" json:"app"`
	IPAddr           string `xml:"ipAddr" json:"ipAddr"`
	VipAddress       string `xml:"vipAddress" json:"vipAddress"`
	SecureVipAddress string `xml:"secureVipAddress" json:"secureVipAddress"`

	Status           StatusType `xml:"status" json:"status"`
	Overriddenstatus StatusType `xml:"overriddenstatus" json:"overriddenstatus"`

	Port              int  `xml:"-" json:"-"`
	PortEnabled       bool `xml:"-" json:"-"`
	SecurePort        int  `xml:"-" json:"-"`
	SecurePortEnabled bool `xml:"-" json:"-"`

	HomePageUrl    string `xml:"homePageUrl" json:"homePageUrl"`
	StatusPageUrl  string `xml:"statusPageUrl" json:"statusPageUrl"`
	HealthCheckUrl string `xml:"healthCheckUrl" json:"healthCheckUrl"`

	CountryId      int64          `xml:"countryId" json:"countryId"`
	DataCenterInfo DataCenterInfo `xml:"dataCenterInfo" json:"dataCenterInfo"`

	LeaseInfo LeaseInfo        `xml:"leaseInfo" json:"leaseInfo"`
	Metadata  InstanceMetadata `xml:"metadata" json:"metadata"`
	//V2 新增
	IsCoordinatingDiscoveryServer string `xml:"isCoordinatingDiscoveryServer,omitempty" json:"isCoordinatingDiscoveryServer,omitempty"`
	LastUpdatedTimestamp          string `xml:"lastUpdatedTimestamp,omitempty" json:"lastUpdatedTimestamp,omitempty"`
	LastDirtyTimestamp            string `xml:"lastDirtyTimestamp,omitempty" json:"lastDirtyTimestamp,omitempty"`
	ActionType                    string `xml:"actionType,omitempty" json:"actionType,omitempty"`

	UniqueID func(i Instance) string `xml:"-" json:"-"`
}

// InstanceMetadata represents the eureka metadata, which is arbitrary XML.
// See metadata.go for more info.
type InstanceMetadata struct {
	Raw    []byte `xml:",innerxml" json:"-"`
	parsed map[string]interface{}
}

// AmazonMetadataType is information about AZ's, AMI's, and the AWS instance.
// <xsd:complexType name="amazonMetdataType">
// from http://docs.amazonwebservices.com/AWSEC2/latest/DeveloperGuide/index.html?AESDG-chapter-instancedata.html
type AmazonMetadataType struct {
	AmiLaunchIndex   string `xml:"ami-launch-index" json:"ami-launch-index"`
	LocalHostname    string `xml:"local-hostname" json:"local-hostname"`
	AvailabilityZone string `xml:"availability-zone" json:"availability-zone"`
	InstanceID       string `xml:"instance-id" json:"instance-id"`
	PublicIpv4       string `xml:"public-ipv4" json:"public-ipv4"`
	PublicHostname   string `xml:"public-hostname" json:"public-hostname"`
	AmiManifestPath  string `xml:"ami-manifest-path" json:"ami-manifest-path"`
	LocalIpv4        string `xml:"local-ipv4" json:"local-ipv4"`
	HostName         string `xml:"hostname" json:"hostname"`
	AmiID            string `xml:"ami-id" json:"ami-id"`
	InstanceType     string `xml:"instance-type" json:"instance-type"`
}

// DataCenterInfo indicates which type of data center hosts this instance
// and conveys details about the instance's environment.
type DataCenterInfo struct {
	// Name indicates which type of data center hosts this instance.
	Name string
	// Class indicates the Java class name representing this structure in the Eureka server,
	// noted only when encoding communication with JSON.
	//
	// When registering an instance, if the name is neither "Amazon" nor "MyOwn", this field's
	// value is used. Otherwise, a suitable default value will be supplied to the server. This field
	// is available for specifying custom data center types other than the two built-in ones, for
	// which no suitable default value could be known.
	Class string
	// Metadata provides details specific to an Amazon data center,
	// populated and honored when the Name field's value is "Amazon".
	Metadata AmazonMetadataType
	// AlternateMetadata provides details specific to a data center other than Amazon,
	// populated and honored when the Name field's value is not "Amazon".
	AlternateMetadata map[string]string
}

// LeaseInfo tells us about the renewal from Eureka, including how old it is.
type LeaseInfo struct {
	RenewalIntervalInSecs int32 `xml:"renewalIntervalInSecs" json:"renewalIntervalInSecs"`
	DurationInSecs        int32 `xml:"durationInSecs" json:"durationInSecs"`
	RegistrationTimestamp int64 `xml:"registrationTimestamp" json:"registrationTimestamp"`
	LastRenewalTimestamp  int64 `xml:"lastRenewalTimestamp" json:"lastRenewalTimestamp"`
	EvictionTimestamp     int64 `xml:"evictionTimestamp" json:"evictionTimestamp"`
	ServiceUpTimestamp    int64 `xml:"serviceUpTimestamp" json:"serviceUpTimestamp"`
}

// ParseAllMetadata iterates through all instances in an application
func (a *Application) ParseAllMetadata() error {
	for _, instance := range a.Instances {
		err := instance.Metadata.parse()
		if err != nil {
			log.Errorf("Failed parsing metadata for Instance=%s of Application=%s: %s",
				instance.HostName, a.Name, err.Error())
			return err
		}
	}
	return nil
}

// SetMetadataString for a given instance before register
func (ins *Instance) SetMetadataString(key, value string) {
	if ins.Metadata.parsed == nil {
		ins.Metadata.parsed = map[string]interface{}{}
	}
	ins.Metadata.parsed[key] = value
}

func (im *InstanceMetadata) parse() error {
	if len(im.Raw) == 0 {
		im.parsed = make(map[string]interface{})
		return nil
	}
	metadataLog.Debugf("InstanceMetadata.parse: %s", im.Raw)

	if len(im.Raw) > 0 && im.Raw[0] == '{' {
		// JSON
		err := json.Unmarshal(im.Raw, &im.parsed)
		if err != nil {
			log.Errorf("Error unmarshalling: %s", err.Error())
			return fmt.Errorf("error unmarshalling: %s", err.Error())
		}
	} else {
		// XML: wrap in a BS xml tag so all metadata tags are pulled
		fullDoc := append(append([]byte("<d>"), im.Raw...), []byte("</d>")...)
		parsedDoc, err := x2j.ByteDocToMap(fullDoc, true)
		if err != nil {
			log.Errorf("Error unmarshalling: %s", err.Error())
			return fmt.Errorf("error unmarshalling: %s", err.Error())
		}
		im.parsed = parsedDoc["d"].(map[string]interface{})
	}
	return nil
}

// GetMap returns a map of the metadata parameters for this instance
func (im *InstanceMetadata) GetMap() map[string]interface{} {
	return im.parsed
}

func (im *InstanceMetadata) getItem(key string) (interface{}, bool, error) {
	err := im.parse()
	if err != nil {
		return "", false, fmt.Errorf("parsing error: %s", err.Error())
	}
	val, present := im.parsed[key]
	return val, present, nil
}

// GetString pulls a value cast as a string. Swallows panics from type
// assertion and returns empty string + an error if conversion fails
func (im *InstanceMetadata) GetString(key string) (s string, err error) {
	defer func() {
		if r := recover(); r != nil {
			s = ""
			err = fmt.Errorf("failed to cast interface to string")
		}
	}()
	v, prs, err := im.getItem(key)
	if !prs {
		return "", err
	}
	return v.(string), err
}

// GetInt pulls a value cast as int. Swallows panics from type assertion and
// returns 0 + an error if conversion fails
func (im *InstanceMetadata) GetInt(key string) (i int, err error) {
	defer func() {
		if r := recover(); r != nil {
			i = 0
			err = fmt.Errorf("failed to cast interface to int")
		}
	}()
	v, err := im.GetFloat64(key)
	return int(v), err
}

// GetFloat32 pulls a value cast as float. Swallows panics from type assertion
// and returns 0.0 + an error if conversion fails
func (im *InstanceMetadata) GetFloat32(key string) (f float32, err error) {
	defer func() {
		if r := recover(); r != nil {
			f = 0.0
			err = fmt.Errorf("failed to cast interface to float32")
		}
	}()
	v, err := im.GetFloat64(key)
	return float32(v), err
}

// GetFloat64 pulls a value cast as float. Swallows panics from type assertion
// and returns 0.0 + an error if conversion fails
func (im *InstanceMetadata) GetFloat64(key string) (f float64, err error) {
	defer func() {
		if r := recover(); r != nil {
			f = 0.0
			err = fmt.Errorf("failed to cast interface to float64")
		}
	}()
	v, prs, err := im.getItem(key)
	if !prs {
		return 0.0, err
	}
	return v.(float64), err
}

// GetBool pulls a value cast as bool.  Swallows panics from type assertion and
// returns false + an error if conversion fails
func (im *InstanceMetadata) GetBool(key string) (b bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			b = false
			err = fmt.Errorf("failed to cast interface to bool")
		}
	}()
	v, prs, err := im.getItem(key)
	if !prs {
		return false, err
	}
	return v.(bool), err
}

func intFromJSONNumberOrString(jv interface{}, description string) (int, error) {
	switch v := jv.(type) {
	case float64:
		return int(v), nil
	case string:
		n, err := strconv.Atoi(v)
		if err != nil {
			return 0, err
		}
		return n, nil
	default:
		return 0, fmt.Errorf("unexpected %s: %[2]v (type %[2]T)", description, jv)
	}
}

// UnmarshalJSON is a custom JSON unmarshaler for GetAppsResponse to deal with
// sometimes non-wrapped Application arrays when there is only a single Application item.
func (r *GetAppsResponse) UnmarshalJSON(b []byte) error {
	marshalLog.Debugf("GetAppsResponse.UnmarshalJSON b:%s\n", string(b))
	resolveDelta := func(d interface{}) (int, error) {
		return intFromJSONNumberOrString(d, "versions delta")
	}

	// Normal array case
	type getAppsResponse GetAppsResponse
	auxArray := struct {
		*getAppsResponse
		VersionsDelta interface{} `json:"versions__delta"`
	}{
		getAppsResponse: (*getAppsResponse)(r),
	}
	var err error
	if err = json.Unmarshal(b, &auxArray); err == nil {
		marshalLog.Debugf("GetAppsResponse.UnmarshalJSON array:%+v\n", auxArray)
		r.VersionsDelta, err = resolveDelta(auxArray.VersionsDelta)
		return err
	}

	// Bogus non-wrapped case
	auxSingle := struct {
		Application   *Application `json:"application"`
		AppsHashcode  string       `json:"apps__hashcode"`
		VersionsDelta interface{}  `json:"versions__delta"`
	}{}
	if err := json.Unmarshal(b, &auxSingle); err != nil {
		return err
	}
	marshalLog.Debugf("GetAppsResponse.UnmarshalJSON single:%+v\n", auxSingle)
	if r.VersionsDelta, err = resolveDelta(auxSingle.VersionsDelta); err != nil {
		return err
	}
	r.Applications = make([]*Application, 1, 1)
	r.Applications[0] = auxSingle.Application
	r.AppsHashcode = auxSingle.AppsHashcode
	return nil
}

// Temporary structs used for Application unmarshalling
type applicationArray Application
type applicationSingle struct {
	Name     string    `json:"name"`
	Instance *Instance `json:"instance"`
}

// UnmarshalJSON is a custom JSON unmarshaler for Application to deal with
// sometimes non-wrapped Instance array when there is only a single Instance item.
func (a *Application) UnmarshalJSON(b []byte) error {
	marshalLog.Debugf("Application.UnmarshalJSON b:%s\n", string(b))
	var err error

	// Normal array case
	var aa applicationArray
	if err = json.Unmarshal(b, &aa); err == nil {
		marshalLog.Debugf("Application.UnmarshalJSON aa:%+v\n", aa)
		*a = Application(aa)
		return nil
	}

	// Bogus non-wrapped case
	var as applicationSingle
	if err = json.Unmarshal(b, &as); err == nil {
		marshalLog.Debugf("Application.UnmarshalJSON as:%+v\n", as)
		a.Name = as.Name
		a.Instances = make([]*Instance, 1, 1)
		a.Instances[0] = as.Instance
		return nil
	}
	return err
}

func stringAsBool(s string) bool {
	return s == "true"
}

// UnmarshalJSON is a custom JSON unmarshaler for Instance, transcribing the two composite port
// specifications up to top-level fields.
func (i *Instance) UnmarshalJSON(b []byte) error {
	// Preclude recursive calls to MarshalJSON.
	type instance Instance
	// inboundJSONFormatPort describes an instance's network port, including whether its registrant
	// considers the port to be enabled or disabled.
	//
	// Example JSON encoding:
	//
	//   Eureka versions 1.2.1 and prior:
	//     "port":{"@enabled":"true", "$":"7101"}
	//
	//   Eureka version 1.2.2 and later:
	//     "port":{"@enabled":"true", "$":7101}
	//
	// Note that later versions of Eureka write the port number as a JSON number rather than as a
	// decimal-formatted string. We accept it as either an integer or a string. Strangely, the
	// "@enabled" field remains a string.
	type inboundJSONFormatPort struct {
		Number  interface{} `json:"$"`
		Enabled bool        `json:"@enabled,string"`
	}
	aux := struct {
		*instance
		Port       inboundJSONFormatPort `json:"port"`
		SecurePort inboundJSONFormatPort `json:"securePort"`
	}{
		instance: (*instance)(i),
	}
	if err := json.Unmarshal(b, &aux); err != nil {
		return err
	}
	resolvePort := func(port interface{}) (int, error) {
		return intFromJSONNumberOrString(port, "port number")
	}
	var err error
	if i.Port, err = resolvePort(aux.Port.Number); err != nil {
		return err
	}
	i.PortEnabled = aux.Port.Enabled
	if i.SecurePort, err = resolvePort(aux.SecurePort.Number); err != nil {
		return err
	}
	i.SecurePortEnabled = aux.SecurePort.Enabled
	return nil
}

// MarshalJSON is a custom JSON marshaler for Instance, adapting the top-level raw port values to
// the composite port specifications.
func (i *Instance) MarshalJSON() ([]byte, error) {
	// Preclude recursive calls to MarshalJSON.
	type instance Instance
	// outboundJSONFormatPort describes an instance's network port, including whether its registrant
	// considers the port to be enabled or disabled.
	//
	// Example JSON encoding:
	//
	//   "port":{"@enabled":"true", "$":"7101"}
	//
	// Note that later versions of Eureka write the port number as a JSON number rather than as a
	// decimal-formatted string. We emit the port number as a string, not knowing the Eureka
	// server's version. Strangely, the "@enabled" field remains a string.
	type outboundJSONFormatPort struct {
		Number  int  `json:"$,string"`
		Enabled bool `json:"@enabled,string"`
	}
	aux := struct {
		*instance
		Port       outboundJSONFormatPort `json:"port"`
		SecurePort outboundJSONFormatPort `json:"securePort"`
	}{
		(*instance)(i),
		outboundJSONFormatPort{i.Port, i.PortEnabled},
		outboundJSONFormatPort{i.SecurePort, i.SecurePortEnabled},
	}
	return json.Marshal(&aux)
}

// xmlFormatPort describes an instance's network port, including whether its registrant considers
// the port to be enabled or disabled.
//
// Example XML encoding:
//
//     <port enabled="true">7101</port>
type xmlFormatPort struct {
	Number  int  `xml:",chardata"`
	Enabled bool `xml:"enabled,attr"`
}

// UnmarshalXML is a custom XML unmarshaler for Instance, transcribing the two composite port
// specifications up to top-level fields.
func (i *Instance) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type instance Instance
	aux := struct {
		*instance
		Port       xmlFormatPort `xml:"port"`
		SecurePort xmlFormatPort `xml:"securePort"`
	}{
		instance: (*instance)(i),
	}
	if err := d.DecodeElement(&aux, &start); err != nil {
		return err
	}
	i.Port = aux.Port.Number
	i.PortEnabled = aux.Port.Enabled
	i.SecurePort = aux.SecurePort.Number
	i.SecurePortEnabled = aux.SecurePort.Enabled
	return nil
}

// startLocalName creates a start-tag of an XML element with the given local name and no namespace name.
func startLocalName(local string) xml.StartElement {
	return xml.StartElement{Name: xml.Name{Space: "", Local: local}}
}

// MarshalXML is a custom XML marshaler for Instance, adapting the top-level raw port values to
// the composite port specifications.
func (i *Instance) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	type instance Instance
	aux := struct {
		*instance
		Port       xmlFormatPort `xml:"port"`
		SecurePort xmlFormatPort `xml:"securePort"`
	}{
		instance:   (*instance)(i),
		Port:       xmlFormatPort{i.Port, i.PortEnabled},
		SecurePort: xmlFormatPort{i.SecurePort, i.SecurePortEnabled},
	}
	return e.EncodeElement(&aux, startLocalName("instance"))
}

// UnmarshalJSON is a custom JSON unmarshaler for InstanceMetadata to handle squirreling away
// the raw JSON for later parsing.
func (i *InstanceMetadata) UnmarshalJSON(b []byte) error {
	i.Raw = b
	// TODO(cq) could actually parse Raw here, and in a parallel UnmarshalXML as well.
	return nil
}

// MarshalJSON is a custom JSON marshaler for InstanceMetadata.
func (i *InstanceMetadata) MarshalJSON() ([]byte, error) {
	if i.parsed != nil {
		return json.Marshal(i.parsed)
	}

	if i.Raw == nil {
		i.Raw = []byte("{}")
	}

	return i.Raw, nil
}

// MarshalXML is a custom XML marshaler for InstanceMetadata.
func (i InstanceMetadata) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	tokens := []xml.Token{start}

	if i.parsed != nil {
		for key, value := range i.parsed {
			t := startLocalName(key)
			tokens = append(tokens, t, xml.CharData(value.(string)), xml.EndElement{Name: t.Name})
		}
	}
	tokens = append(tokens, xml.EndElement{Name: start.Name})

	for _, t := range tokens {
		err := e.EncodeToken(t)
		if err != nil {
			return err
		}
	}

	// flush to ensure tokens are written
	return e.Flush()
}

type metadataMap map[string]string

// MarshalXML is a custom XML marshaler for metadataMap, mapping each metadata name/value pair to a
// correspondingly named XML element with the pair's value as character data content.
func (m metadataMap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if err := e.EncodeToken(start); err != nil {
		return err
	}

	for k, v := range m {
		if err := e.EncodeElement(v, startLocalName(k)); err != nil {
			return err
		}
	}

	return e.EncodeToken(start.End())
}

// UnmarshalXML is a custom XML unmarshaler for metadataMap, mapping each XML element's name and
// character data content to a corresponding metadata name/value pair.
func (m metadataMap) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	for {
		t, err := d.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if k, ok := t.(xml.StartElement); ok {
			if err := d.DecodeElement(&v, &k); err != nil {
				return err
			}
			m[k.Name.Local] = v
		}
	}
	return nil
}

func metadataValue(i *DataCenterInfo) interface{} {
	if i.Name == Amazon {
		return i.Metadata
	}
	return metadataMap(i.AlternateMetadata)
}

var (
	startName     = startLocalName("name")
	startMetadata = startLocalName("metadata")
)

// MarshalXML is a custom XML marshaler for DataCenterInfo, writing either Metadata or AlternateMetadata
// depending on the type of data center indicated by the Name.
func (i *DataCenterInfo) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if err := e.EncodeToken(start); err != nil {
		return err
	}

	if err := e.EncodeElement(i.Name, startName); err != nil {
		return err
	}
	if err := e.EncodeElement(metadataValue(i), startMetadata); err != nil {
		return err
	}

	return e.EncodeToken(start.End())
}

type preliminaryDataCenterInfo struct {
	Name     string      `xml:"name" json:"name"`
	Class    string      `xml:"-" json:"@class"`
	Metadata metadataMap `xml:"metadata" json:"metadata"`
}

func bindValue(dst *string, src map[string]string, k string) bool {
	if v, ok := src[k]; ok {
		*dst = v
		return true
	}
	return false
}

func populateAmazonMetadata(dst *AmazonMetadataType, src map[string]string) {
	bindValue(&dst.AmiLaunchIndex, src, "ami-launch-index")
	bindValue(&dst.LocalHostname, src, "local-hostname")
	bindValue(&dst.AvailabilityZone, src, "availability-zone")
	bindValue(&dst.InstanceID, src, "instance-id")
	bindValue(&dst.PublicIpv4, src, "public-ipv4")
	bindValue(&dst.PublicHostname, src, "public-hostname")
	bindValue(&dst.AmiManifestPath, src, "ami-manifest-path")
	bindValue(&dst.LocalIpv4, src, "local-ipv4")
	bindValue(&dst.HostName, src, "hostname")
	bindValue(&dst.AmiID, src, "ami-id")
	bindValue(&dst.InstanceType, src, "instance-type")
}

func adaptDataCenterInfo(dst *DataCenterInfo, src *preliminaryDataCenterInfo) {
	dst.Name = src.Name
	dst.Class = src.Class
	if src.Name == Amazon {
		populateAmazonMetadata(&dst.Metadata, src.Metadata)
	} else {
		dst.AlternateMetadata = src.Metadata
	}
}

// UnmarshalXML is a custom XML unmarshaler for DataCenterInfo, populating either Metadata or AlternateMetadata
// depending on the type of data center indicated by the Name.
func (i *DataCenterInfo) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	p := preliminaryDataCenterInfo{
		Metadata: make(map[string]string, 11),
	}
	if err := d.DecodeElement(&p, &start); err != nil {
		return err
	}
	adaptDataCenterInfo(i, &p)
	return nil
}

// MarshalJSON is a custom JSON marshaler for DataCenterInfo, writing either Metadata or AlternateMetadata
// depending on the type of data center indicated by the Name.
func (i *DataCenterInfo) MarshalJSON() ([]byte, error) {
	type named struct {
		Name  string `json:"name"`
		Class string `json:"@class"`
	}
	if i.Name == Amazon {
		return json.Marshal(struct {
			named
			Metadata AmazonMetadataType `json:"metadata"`
		}{
			named{i.Name, "com.netflix.appinfo.AmazonInfo"},
			i.Metadata,
		})
	}
	class := "com.netflix.appinfo.MyDataCenterInfo"
	if i.Name != MyOwn {
		class = i.Class
	}
	return json.Marshal(struct {
		named
		Metadata map[string]string `json:"metadata,omitempty"`
	}{
		named{i.Name, class},
		i.AlternateMetadata,
	})
}

func jsonValueAsString(i interface{}) string {
	switch v := i.(type) {
	case string:
		return v
	case float64:
		return fmt.Sprintf("%.f", v)
	case bool:
		return strconv.FormatBool(v)
	case []interface{}, map[string]interface{}:
		// Don't bother trying to decode these.
		return ""
	case nil:
		return ""
	default:
		panic("type of unexpected value")
	}
}

// UnmarshalJSON is a custom JSON unmarshaler for DataCenterInfo, populating either Metadata or AlternateMetadata
// depending on the type of data center indicated by the Name.
func (i *DataCenterInfo) UnmarshalJSON(b []byte) error {
	// The Eureka server will mistakenly convert metadata values that look like numbers to JSON numbers.
	// Convert them back to strings.
	aux := struct {
		*preliminaryDataCenterInfo
		PreliminaryMetadata map[string]interface{} `json:"metadata"`
	}{
		PreliminaryMetadata: make(map[string]interface{}, 11),
	}
	if err := json.Unmarshal(b, &aux); err != nil {
		return err
	}
	metadata := make(map[string]string, len(aux.PreliminaryMetadata))
	for k, v := range aux.PreliminaryMetadata {
		metadata[k] = jsonValueAsString(v)
	}
	aux.Metadata = metadata
	adaptDataCenterInfo(i, aux.preliminaryDataCenterInfo)
	return nil
}

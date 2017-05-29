package data

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"

	"time"

	"github.com/clbanning/checkxml"
	"github.com/fatih/structs"
	"github.com/wcerfgba/nationstates-xlsx/util"
)

type Provider20170429 struct {
	config struct {
		nation string
	}
}

func (p *Provider20170429) Get() (res Result) {
	// log.Printf("Fetching %v", nation)
	url := fmt.Sprintf("https://www.nationstates.net/cgi-bin/api.cgi?nation=%v&q=deaths+gdp+publicsector+govt+income+sectors+freedomscores", p.config.nation)
	apiRes, err := http.Get(url)
	if err != nil {
		// log.Fatal(err)
	}

	// log.Println("Reading response")
	xmlStr, err := ioutil.ReadAll(apiRes.Body)
	apiRes.Body.Close()
	if err != nil {
		// log.Fatal(err)
	}

	ref := nation20170429{}
	err = xml.Unmarshal(xmlStr, &ref)
	extra, _, _ := checkxml.UnknownXMLTags(xmlStr, ref)
	extra = remove([]string{
		"DEATHS",
		"-id",
	}, extra)
	data := structs.Map(ref)
	data["DEATHS"] = flattenDeaths(data["DEATHS"].([]interface{}))

	// The source data does not have a timestamp, so we inject one.
	data["_timestamp"] = time.Now().Format(time.RFC3339)

	data = restructureData(data)

	res = NewResult()
	buildResult(data, res)

	return
}

func restructureData(in map[string]interface{}) (out map[string]interface{}) {

	gdpInt, err := strconv.ParseInt(in["GDP"].(string), 10, 64)
	if err != nil {
		return
	}
	gdpBnsStr := fmt.Sprintf("%.3f", float64(gdpInt)/math.Pow10(9))
	gdpBnsFloat, err := strconv.ParseFloat(gdpBnsStr, 64)
	if err != nil {
		return
	}
	publicSectorFloat, err := strconv.ParseFloat(in["PUBLICSECTOR"].(string), 64)
	if err != nil {
		return
	}
	publicExpenditureBns := (publicSectorFloat / 100) * gdpBnsFloat
	publicExpenditureBnsStr := fmt.Sprintf("%.3f", publicExpenditureBns)

	out = map[string]interface{}{
		"Causes of death": in["DEATHS"],
		"Government expenditure": map[string]interface{}{
			"Expenditure (billion)": publicExpenditureBnsStr,
			"% of GDP":              in["PUBLICSECTOR"],
			"Education":             in["GOVT"].(map[string]interface{})["EDUCATION"],
			"Environment":           in["GOVT"].(map[string]interface{})["ENVIRONMENT"],
			"Healthcare":            in["GOVT"].(map[string]interface{})["HEALTHCOARE"],
			"Commerce":              in["GOVT"].(map[string]interface{})["COMMERCE"],
			"International aid":     in["GOVT"].(map[string]interface{})["INTERNATIONALAID"],
			"Law and Order":         in["GOVT"].(map[string]interface{})["LAWANDORDER"],
			"Public Transport":      in["GOVT"].(map[string]interface{})["PUBLICTRANSPORT"],
			"Social Equality":       in["GOVT"].(map[string]interface{})["SOCIALEQUALITY"],
			"Spirituality":          in["GOVT"].(map[string]interface{})["SPIRITUALITY"],
			"Welfare":               in["GOVT"].(map[string]interface{})["WELFARE"],
		},
		"Economy": map[string]interface{}{
			"GDP (billion)":        gdpBnsStr,
			"Ave. wage":            in["INCOME"],
			"Government":           in["SECTORS"].(map[string]interface{})["GOVERNMENT"],
			"State-owned Industry": in["SECTORS"].(map[string]interface{})["PUBLIC"],
			"Private Industry":     in["SECTORS"].(map[string]interface{})["INDUSTRY"],
			"Black Market":         in["SECTORS"].(map[string]interface{})["BLACKMARKET"],
		},
		"Rights": map[string]interface{}{
			"Civil Rights":      in["FREEDOMSCORES"].(map[string]interface{})["CIVILRIGHTS"],
			"Economy":           in["FREEDOMSCORES"].(map[string]interface{})["ECONOMY"],
			"Political Freedom": in["FREEDOMSCORES"].(map[string]interface{})["POLITICALFREEDOM"],
		},
	}

	for _, submap := range out {
		submap.(map[string]interface{})["Timestamp"] = in["_timestamp"]
	}

	return
}

func buildResult(in map[string]interface{}, res Result) {
	for k, v := range in {
		switch v.(type) {
		case string:
			res.SetChildValue(k, v.(string))
		case map[string]interface{}:
			child := res.AddChild(k)
			buildResult(v.(map[string]interface{}), child)
		}
	}
	return
}

func (p *Provider20170429) Configure(c util.Configuration) {
	p.config.nation = c["nation"].(string)
}

// nation20170429 is one particular InputSpec. The struct itself is used by
// xml.Unmarshal, and also acts as a receiver for the InputSpec interface
// functions.
type nation20170429 struct {
	GDP,
	INCOME,
	PUBLICSECTOR string
	GOVT struct {
		ADMINISTRATION,
		DEFENCE,
		EDUCATION,
		ENVIRONMENT,
		HEALTHCARE,
		COMMERCE,
		INTERNATIONALAID,
		LAWANDORDER,
		PUBLICTRANSPORT,
		SOCIALEQUALITY,
		SPIRITUALITY,
		WELFARE string
	}
	SECTORS struct {
		BLACKMARKET,
		GOVERNMENT,
		INDUSTRY,
		PUBLIC string
	}
	FREEDOMSCORES struct {
		CIVILRIGHTS,
		ECONOMY,
		POLITICALFREEDOM string
	}
	DEATHS []causeOfDeath20170429 `xml:"DEATHS>CAUSE"`
}

// causeOfDeath20170429 is used to store the causes of death from the XML. We
// define it as a separate struct as we have to extract both an attribute and
// the chardata.
type causeOfDeath20170429 struct {
	Type   string `xml:"type,attr"`
	Amount string `xml:",chardata"`
}

// remove is a helper function to filter out values from a []string.
func remove(excludes, from []string) (out []string) {
	out = []string{}
search:
	for _, v := range from {
		for _, ex := range excludes {
			if v == ex {
				continue search
			}
		}
		out = append(out, v)
	}
	return out
}

// flattenDeaths takes a list of maps with "Type" and "Amount" keys and builds
// a map from the type to the amount.
func flattenDeaths(in []interface{}) (out map[string]interface{}) {
	out = map[string]interface{}{}
	for _, death := range in {
		castDeath := death.(map[string]interface{})
		out[castDeath["Type"].(string)] = castDeath["Amount"].(string)
	}
	return
}

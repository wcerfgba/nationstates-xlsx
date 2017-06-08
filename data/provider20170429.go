package data

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"

	"time"

	"github.com/clbanning/checkxml"
	"github.com/fatih/structs"
	"github.com/wcerfgba/nationstates-xlsx/util"
)

// ProviderV1 is a production Provider which retrieves data via the official
// NationStates REST API. It is built to provide a Result consisting of a
// particlar set of statistics, arranged in a way suitable for storing in a
// spreadsheet.
type ProviderV1 struct {
	config struct {
		nation        string
		resultFactory func() Result
	}
}

func (p *ProviderV1) Get() (res Result) {
	log.Printf("Fetching %v", p.config.nation)
	url := fmt.Sprintf("https://www.nationstates.net/cgi-bin/api.cgi?nation=%v&q=deaths+gdp+publicsector+govt+income+sectors+freedomscores", p.config.nation)
	apiRes, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Reading response")
	xmlStr, err := ioutil.ReadAll(apiRes.Body)
	apiRes.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	// Since the NationStates API returns XML, we can model the expected result
	// as a struct and unmarshal the data into it.
	ref := nationV1{}
	err = xml.Unmarshal(xmlStr, &ref)

	// Check for any tags that are unknown (not in the struct) or missing
	// (not in the XML).
	unknown, _, _ := checkxml.UnknownXMLTags(xmlStr, ref)
	unknown = remove([]string{
		"DEATHS",
		"-id",
	}, unknown)
	if len(unknown) > 0 {
		log.Printf("Unknown tags: %v", unknown)
	}
	unmatched, _, _ := checkxml.MissingXMLTags(xmlStr, ref)
	if len(unmatched) > 0 {
		log.Printf("Unmatched tags: %v", unmatched)
	}

	// Convert struct to a map.
	data := structs.Map(ref)

	// Array of deaths can be a map.
	data["DEATHS"] = flattenDeaths(data["DEATHS"].([]interface{}))

	// The source data does not have a timestamp, so we inject one.
	data["_timestamp"] = time.Now().Format(time.RFC3339)

	// Re-arrange data for spreadsheet.
	sheetData := restructureData(data)

	// Construct Result
	res = p.config.resultFactory()
	buildResult(sheetData, res)

	return
}

// restructureData takes in the 'raw' map of parsed data from the API and
// returns a relatively structured map of map[string]map[string]string. The
// first key is the sheet name, the second key is the column title.
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

	// Although it would be nice to type this as `map[string]map[string]string`,
	// golang will complain when trying to cast to `map[string]interface{}`,
	// which we need to do later when we recurse over the map, so we don't
	// bother. #typesystemwoes
	out = map[string]interface{}{
		"Causes of death": in["DEATHS"].(map[string]string),
		"Government expenditure": map[string]string{
			"Expenditure (billion)": publicExpenditureBnsStr,
			"% of GDP":              in["PUBLICSECTOR"].(string),
			"Education":             in["GOVT"].(map[string]interface{})["EDUCATION"].(string),
			"Environment":           in["GOVT"].(map[string]interface{})["ENVIRONMENT"].(string),
			"Healthcare":            in["GOVT"].(map[string]interface{})["HEALTHCOARE"].(string),
			"Commerce":              in["GOVT"].(map[string]interface{})["COMMERCE"].(string),
			"International aid":     in["GOVT"].(map[string]interface{})["INTERNATIONALAID"].(string),
			"Law and Order":         in["GOVT"].(map[string]interface{})["LAWANDORDER"].(string),
			"Public Transport":      in["GOVT"].(map[string]interface{})["PUBLICTRANSPORT"].(string),
			"Social Equality":       in["GOVT"].(map[string]interface{})["SOCIALEQUALITY"].(string),
			"Spirituality":          in["GOVT"].(map[string]interface{})["SPIRITUALITY"].(string),
			"Welfare":               in["GOVT"].(map[string]interface{})["WELFARE"].(string),
		},
		"Economy": map[string]string{
			"GDP (billion)":        gdpBnsStr,
			"Ave. wage":            in["INCOME"].(string),
			"Government":           in["SECTORS"].(map[string]interface{})["GOVERNMENT"].(string),
			"State-owned Industry": in["SECTORS"].(map[string]interface{})["PUBLIC"].(string),
			"Private Industry":     in["SECTORS"].(map[string]interface{})["INDUSTRY"].(string),
			"Black Market":         in["SECTORS"].(map[string]interface{})["BLACKMARKET"].(string),
		},
		"Rights": map[string]string{
			"Civil Rights":      in["FREEDOMSCORES"].(map[string]interface{})["CIVILRIGHTS"].(string),
			"Economy":           in["FREEDOMSCORES"].(map[string]interface{})["ECONOMY"].(string),
			"Political Freedom": in["FREEDOMSCORES"].(map[string]interface{})["POLITICALFREEDOM"].(string),
		},
	}

	// Add the timestamp to each sheet.
	for _, submap := range out {
		submap.(map[string]string)["Timestamp"] = in["_timestamp"].(string)
	}

	return
}

// buildResult recurses over an arbitrarily-nested map with string keys and
// values, and constructs a Result.
func buildResult(in map[string]interface{}, res Result) {
	for k, v := range in {
		switch v.(type) {
		case string:
			res.AddOrGetChild(k).SetValue(v.(string))
		case map[string]interface{}:
			child := res.AddOrGetChild(k)
			buildResult(v.(map[string]interface{}), child)
		}
	}
	return
}

func (p *ProviderV1) Configure(c util.Configuration) {
	p.config.nation = c["nation"].(string)
	p.config.resultFactory = c["resultFactory"].(func() Result)
}

// nationV1 is one particular input spec. The struct itself is used by
// xml.Unmarshal, and also acts as a receiver for the interface functions.
type nationV1 struct {
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

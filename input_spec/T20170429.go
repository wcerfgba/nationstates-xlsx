package input_spec

import (
	"encoding/xml"
	"fmt"

	"github.com/fatih/structs"
)

type Nation_20170429 struct {
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
	DEATHS []Death_20170429
}

type Death_20170429 struct {
	Type   string `xml:",attr"`
	Amount string
}

func (s *Nation_20170429) BuildRequestUrl(nation string) (url string) {
	url = fmt.Sprintf("https://www.nationstates.net/cgi-bin/api.cgi?nation=%v&q=deaths+gdp+publicsector+govt+income+sectors+freedomscores", nation)
	return
}

func (s *Nation_20170429) Parse(xmlStr []byte) (data InputData, err error) {
	ref := Nation_20170429{}
	err = xml.Unmarshal(xmlStr, &ref)
	data = structs.Map(ref)
	return
}

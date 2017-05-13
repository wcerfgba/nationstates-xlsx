package input_spec

type InputData map[string]interface{}

type InputSpec interface {
	BuildRequestUrl(nation string) (url string)
	Parse(xml []byte) (data InputData, extra []string, err error)
}

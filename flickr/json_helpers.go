package flickr

import "encoding/json"

type Content string

type contentObj struct {
	Content string `json:"_content"`
}

func (c *Content) UnmarshalJSON(b []byte) error {
	str := ""
	if b[0] == '"' {
		str = string(b[1 : len(b)-2])
	} else {
		o := contentObj{}
		if err := json.Unmarshal(b, &o); err != nil {
			return err
		}
		str = o.Content
	}
	*c = Content(str)
	return nil
}

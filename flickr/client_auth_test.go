package flickr

import (
	"encoding/json"
	"testing"
)

func TestUnmarshalUser(t *testing.T) {
	b := []byte(`{"auth":{"token":{"_content":"72157653336438823-XYZ"},"perms":{"_content":"delete"},"user":{"nsid":"59132923@N08","username":"jpillora","fullname":"Jaime Pillora"}},"stat":"ok"}`)

	resp := TokenResp{}
	if err := json.Unmarshal(b, &resp); err != nil {
		t.Fatal(err)
	}

	if resp.Auth.Token.Content != "72157653336438823-XYZ" {
		t.Error("Missing token")
	}
}

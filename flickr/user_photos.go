package flickr

import "encoding/json"

type PhotoSize struct {
	Label  string          `json:"label"`
	Width  json.RawMessage `json:"width"`
	Height json.RawMessage `json:"height"`
	Source string          `json:"source"`
	URL    string          `json:"url"`
	Media  string          `json:"media"`
}

type PhotoSizes struct {
	CanBlog     int         `json:"canblog"`
	CanPrint    int         `json:"canprint"`
	CanDownload int         `json:"candownload"`
	Sizes       []PhotoSize `json:"size"`
}

func (u *User) PhotoSizes(photoID string) (*PhotoSizes, error) {
	b, err := u.Do("flickr.photos.getSizes", Args{"photo_id": photoID})
	if err != nil {
		return nil, err
	}
	resp := &struct {
		Sizes PhotoSizes `json:"sizes"`
	}{}
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}
	return &resp.Sizes, nil
}

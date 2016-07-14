package flickr

import (
	"encoding/json"
	"strconv"
)

type PhotoSetsResponse struct {
	Data struct {
		Cancreate int         `json:"cancreate"`
		Page      int         `json:"page"`
		Pages     int         `json:"pages"`
		Perpage   int         `json:"perpage"`
		Total     int         `json:"total"`
		Photosets []*PhotoSet `json:"photoset"`
	} `json:"photosets"`
}

type PhotoSet struct {
	ID                  string   `json:"id"`
	Title               Content  `json:"title"`
	Description         Content  `json:"description"`
	Owner               string   `json:"owner"`
	Primary             string   `json:"primary"`
	Photos              []*Photo `json:"photo"`
	Farm                int      `json:"farm"`
	Secret              string   `json:"secret"`
	Server              string   `json:"server"`
	Username            string   `json:"username"`
	CanComment          int      `json:"can_comment"`
	CountComments       string   `json:"count_comments"`
	CountPhotos         int      `json:"count_photos"`
	CountVideos         string   `json:"count_videos"`
	CountViews          string   `json:"count_views"`
	CoverphotoFarm      int      `json:"coverphoto_farm"`
	CoverphotoServer    string   `json:"coverphoto_server"`
	DateCreate          string   `json:"date_create"`
	DateUpdate          string   `json:"date_update"`
	Videos              string   `json:"videos"`
	NeedsInterstitial   int      `json:"needs_interstitial"`
	VisibilityCanSeeSet int      `json:"visibility_can_see_set"`
}

func (u *User) Photoset(id string) (*PhotoSet, error) {
	b, err := u.Do("flickr.photosets.getInfo", Args{"photoset_id": id})
	if err != nil {
		return nil, err
	}
	resp := PhotosetResponse{}
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}
	set := resp.Data
	if err := u.GetPhotos(set); err != nil {
		return nil, err
	}
	return set, nil
}

func (u *User) Photosets() ([]*PhotoSet, error) {
	b, err := u.Do("flickr.photosets.getList", nil)
	if err != nil {
		return nil, err
	}
	resp := PhotoSetsResponse{}
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}
	return resp.Data.Photosets, nil
}

type GetPhotos struct {
	ID        string   `json:"id"`
	Primary   string   `json:"primary"`
	Owner     string   `json:"owner"`
	Ownername string   `json:"ownername"`
	Photos    []*Photo `json:"photo"`
	Page      string   `json:"page"`
	PerPage   int      `json:"per_page"`
	Perpage   int      `json:"perpage"`
	Pages     int      `json:"pages"`
	Total     string   `json:"total"`
	Title     string   `json:"title"`
}

func (u *User) GetPhotos(set *PhotoSet) error {
	page := 0
	pages := 1
	for page < pages {
		page++
		b, err := u.Do("flickr.photosets.getPhotos", Args{
			"photoset_id": set.ID,
			"extras":      "url_o,original_format,date_taken,media,path_alias",
			"page":        strconv.Itoa(page),
		})
		if err != nil {
			return err
		}
		resp := struct {
			GetPhotos `json:"photoset"`
		}{}
		if err := json.Unmarshal(b, &resp); err != nil {
			return err
		}
		getPhotos := resp.GetPhotos
		set.Photos = append(set.Photos, getPhotos.Photos...)
		pages = getPhotos.Pages
	}
	return nil
}

type PhotosetResponse struct {
	Data *PhotoSet `json:"photoset"`
}

// 	Page      int      `json:"page"`
// 	Pages     int      `json:"pages"`
// 	PerPage   int      `json:"per_page"`
// 	Perpage   int      `json:"perpage"`
// 	Primary   string   `json:"primary"`
// 	Title     string   `json:"title"`
// 	Total     string   `json:"total"`

type Photo struct {
	Farm           int    `json:"farm"`
	ID             string `json:"id"`
	Isfamily       int    `json:"isfamily"`
	Isfriend       int    `json:"isfriend"`
	Isprimary      string `json:"isprimary"`
	Ispublic       int    `json:"ispublic"`
	Secret         string `json:"secret"`
	Server         string `json:"server"`
	Title          string `json:"title"`
	DateTaken      string `json:"datetaken"`
	OriginalURL    string `json:"url_o"`
	OriginalSecret string `json:"originalsecret"`
	OriginalFormat string `json:"originalformat"`
	Media          string `json:"media"`
	MediaState     string `json:"media_status"`
}

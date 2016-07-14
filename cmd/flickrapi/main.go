package main

import (
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/jpillora/flickr-bot/flickr"
	"github.com/jpillora/opts"
)

type Config struct {
	API    string
	Secret string
	Token  string
	Cmd    string `type:"cmdname"`
	Run    struct {
		Method string   `type:"arg"`
		Args   []string `type:"args"`
	}
	GetToken struct{}
	PhotoSet struct {
		ID string `type:"arg"`
	}
	PhotoSets struct{}
	Photos    struct {
		Set string `type:"arg"`
	}
	PhotoList struct {
		Set string `type:"arg"`
	}
	PhotoSizes struct {
		ID string `type:"arg"`
	}
	Merge struct {
		Dest    string   `type:"arg"`
		Sources []string `type:"args" min:"1"`
	}
	MergeMatch struct {
		Dest    string `type:"arg"`
		Matcher string `type:"arg"`
	}
}

func main() {
	conf := Config{
		API:    "",
		Secret: "",
		Token:  "",
	}
	opts.Parse(&conf)

	client := flickr.New(conf.API, conf.Secret)

	user, err := client.UserFromToken(conf.Token)
	if conf.Token != "" && err != nil {
		log.Fatal(err)
	}

	log.Printf("Running %s...", conf.Cmd)

	switch conf.Cmd {
	//===================================
	case "run":
		args := flickr.Args{}
		for _, a := range conf.Run.Args {
			kv := strings.SplitN(a, "=", 2)
			if len(kv) != 2 {
				log.Fatalf("Invalid arg '%s'", a)
			}
			args[kv[0]] = kv[1]
		}
		if user == nil {
			client.Test(conf.Run.Method, args)
		} else {
			user.Test(conf.Run.Method, args)
		}
		//===================================
	case "get-token":
		frob, url, err := client.GetAuthURL()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Please visit\n\n=> %s\n\nand allow access\n", url)
		var user *flickr.User = nil
		for i := 0; i < 30; i++ {
			user, err = client.UserFromFrob(frob)
			if err != nil {
				log.Printf("Attempt #%02d: %s...", i+1, err)
				time.Sleep(3 * time.Second)
				continue
			}
			break
		}
		if user == nil {
			log.Fatal("Timeout")
		}
		log.Printf("Authenticated as %s with token:\n\n=> %s", user.Fullname, user.Token)
		//===================================
	case "photo-sets":
		if user == nil {
			log.Fatal("Token required")
		}
		sets, err := user.Photosets()
		check(err)
		for i, s := range sets {
			log.Printf("%s #%03d %50s (%04d photos)", s.ID, i+1, s.Title, s.CountPhotos)
		}
		//===================================
	case "photo-set":
		if user == nil {
			log.Fatal("Token required")
		}
		set, err := user.Photoset(conf.PhotoSet.ID)
		check(err)
		log.Printf("%+v", set)
	//===================================
	case "photos":
		if user == nil {
			log.Fatal("Token required")
		}
		user.Test("flickr.photosets.getInfo", flickr.Args{"photoset_id": conf.Photos.Set})
		//===================================
	case "photo-sizes":
		if user == nil {
			log.Fatal("Token required")
		}
		// go run main.go photo-sizes 14238614047
		user.Test("flickr.photos.getSizes", flickr.Args{"photo_id": conf.PhotoSizes.ID})
		//===================================
	case "photo-list":
		if user == nil {
			log.Fatal("Token required")
		}
		s, err := user.Photoset(conf.PhotoList.Set)
		check(err)
		err = user.GetPhotos(s)
		check(err)

		photos := s.Photos
		if len(photos) > 3 {
			photos = photos[:3]
		}
		for _, p := range photos {
			log.Printf("photo: %+v", p)
			log.Printf("source: https://farm%d.staticflickr.com/%s/%s_%s_o.jpg", p.Farm, p.Server, p.ID, p.Secret)
			log.Printf("get sizes...")
			sizes, err := user.PhotoSizes(p.ID)
			check(err)
			for _, s := range sizes.Sizes {
				log.Printf("%s = %s", s.Label, s.Source)
			}
		}
	//===================================
	case "merge":
		if user == nil {
			log.Fatal("Token required")
		}

		dest, err := user.Photoset(conf.Merge.Dest)
		check(err)
		for _, srcID := range conf.Merge.Sources {
			src, err := user.Photoset(srcID)
			check(err)
			log.Printf("Moving photos from %s...", src.Title)
			for i, photo := range src.Photos {
				_, err = user.Do("flickr.photosets.addphoto", flickr.Args{"photoset_id": dest.ID, "photo_id": photo.ID})
				if err != nil && err.Error() != "Photo already in set" {
					check(err)
				}
				_, err = user.Do("flickr.photosets.removephoto", flickr.Args{"photoset_id": src.ID, "photo_id": photo.ID})
				check(err)
				log.Printf("Moved photo %d", i+1)
			}
		}
		//===================================
	case "merge-match":
		if user == nil {
			log.Fatal("Token required")
		}

		re, err := regexp.Compile(conf.MergeMatch.Matcher)
		check(err)

		dest, err := user.Photoset(conf.MergeMatch.Dest)
		check(err)
		log.Printf("Destination photoset %s", dest.Title)

		sets, err := user.Photosets()
		check(err)
		for _, src := range sets {
			if !re.MatchString(string(src.Title)) {
				continue
			}
			err = user.GetPhotos(src)
			check(err)
			log.Printf("Moving photos from %s...", src.Title)
			for i, photo := range src.Photos {

				log.Printf("Moving photo %d (%s)...", i+1, photo.ID)
				_, err = user.Do("flickr.photosets.addphoto", flickr.Args{"photoset_id": dest.ID, "photo_id": photo.ID})
				if err != nil && err.Error() != "Photo already in set" {
					check(err)
				}
				_, err = user.Do("flickr.photosets.removephoto", flickr.Args{"photoset_id": src.ID, "photo_id": photo.ID})
				check(err)
			}
		}
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

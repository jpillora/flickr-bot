package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/jpillora/backoff"
	"github.com/jpillora/flickr-bot/flickr"
	"github.com/tj/go-dropbox"
)

const t = ``

const dest = "/Flickr"

var d = dropbox.New(dropbox.NewConfig(t))

var existing = map[string]string{}

type tx struct {
	id, url, path, ext string
	taken              time.Time
}

const dequeuers = 4

var queue = make(chan tx)

func main() {
	log.Println("start")
	t := time.Now()
	loadexisting()
	wg := sync.WaitGroup{}
	for i := 0; i < dequeuers; i++ {
		wg.Add(1)
		go dequeue(i, &wg)
	}
	enqueue()
	wg.Wait()
	log.Printf("done (%s)", time.Now().Sub(t))
}

func loadexisting() {
	fmt.Printf("dropbox load existing..")
	resp, err := d.Files.ListFolder(&dropbox.ListFolderInput{
		Path:      dest,
		Recursive: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	for {
		//add to map
		for _, e := range resp.Entries {
			path := e.PathDisplay
			ext := filepath.Ext(path)
			path = strings.TrimSuffix(path, ext)
			existing[path] = ext
		}
		//load more?
		if !resp.HasMore {
			break
		}
		resp, err = d.Files.ListFolderContinue(&dropbox.ListFolderContinueInput{
			Cursor: resp.Cursor,
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf(".")
	}
	fmt.Printf("\ndropbox directory %s already has %d files\n", dest, len(existing))
}

func enqueue() {
	client := flickr.New("", "")
	user, err := client.UserFromToken("-")
	if err != nil {
		log.Fatal(err)
	}
	sets, err := user.Photosets()
	if err != nil {
		log.Fatal(err)
	}
	for i, set := range sets {
		log.Printf("[%d/%d] Loading photo set '%s' ", i+1, len(sets), set.Title)
		if err := user.GetPhotos(set); err != nil {
			log.Fatal(err)
		}
		log.Printf("[%d/%d] Transferring %d photos...", i+1, len(sets), len(set.Photos))
		for _, photo := range set.Photos {
			if photo.OriginalSecret == "" {
				fmt.Printf("\n[ERROR] %s: HAS NO SECRET\n", photo.ID)
				continue
			}
			t := tx{
				id:   photo.ID,
				path: fmt.Sprintf("%s/%s/%s", dest, set.Title, photo.ID),
			}
			//path check
			//original photo/video url
			if photo.Media == "video" {
				t.url = fmt.Sprintf("https://www.flickr.com/photos/%s/%s/play/orig/%s/", user.ID, photo.ID, photo.OriginalSecret)
				if ext, ok := existing[t.path]; ok && ext != ".jpg" {
					fmt.Print("_")
					continue
				}
			} else {
				t.url = photo.OriginalURL
				t.ext = filepath.Ext(photo.OriginalURL)
				if ext, ok := existing[t.path]; ok && ext == t.ext {
					fmt.Print("_")
					continue
				}
			}
			//calculate timestamp for file
			dateTaken, err := time.Parse("2006-01-02 15:04:05", photo.DateTaken)
			if err != nil {
				fmt.Printf("\n[ERROR] %s: INVALID-DATE: %s\n", photo.ID, err)
				continue
			}
			t.taken = dateTaken.Add(-10 * time.Hour) //manual AEST fix
			//enqueue!
			fmt.Print(".")
			queue <- t
		}
		fmt.Print("\n")
	}
	close(queue)
}

func dequeue(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	b := backoff.Backoff{Min: time.Second}
	for tx := range queue {
	attempt:
		if err := transfer(tx); err != nil {
			if b.Attempt() == 10 {
				fmt.Printf("\n[ERROR] %s: TRANSFER '%s' (give up)\n", tx.id, err)
				return
			}
			d := b.Duration()
			fmt.Printf("\n[ERROR] %s: TRANSFER '%s' (retrying in %s)\n", tx.id, err, d)
			time.Sleep(d)
			goto attempt
		}
		b.Reset()
	}
}

func transfer(t tx) error {
	// log.Printf("downloading %s", t.url)
	//download photo
	resp, err := http.Get(t.url)
	if err != nil {
		return fmt.Errorf("flickr GET error: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("flickr GET error: status %s", resp.Status)
	}
	//calculate filename
	if t.ext == "" {
		t.ext = filepath.Ext(resp.Header.Get("Content-Disposition"))
	}
	// log.Printf("uploading %s", t.path+t.ext)
	//and upload result
	_, err = d.Files.Upload(&dropbox.UploadInput{
		Path:           t.path + t.ext,
		Mode:           dropbox.WriteModeAdd,
		AutoRename:     true,
		Mute:           false,
		ClientModified: t.taken,
		Reader:         resp.Body,
	})
	if err != nil {
		return fmt.Errorf("dropbox POST error: %s", err)
	}
	return nil
}

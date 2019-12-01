package saver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"strconv"

	"github.com/dghubble/go-twitter/twitter"
)

type TweetSaver struct {
	SaveDir   string
	SaveJson  bool
	SaveMedia bool
}

func (ts *TweetSaver) Save(tweet twitter.Tweet) error {
	if ts.SaveJson {
		if err := ts.saveJson(tweet); err != nil {
			return err
		}
	}

	if ts.SaveMedia && hasMedia(tweet) {
		if err := ts.saveMedia(tweet); err != nil {
			return err
		}
	}

	return nil
}

func (ts *TweetSaver) saveJson(tweet twitter.Tweet) error {
	tweetId := strconv.FormatInt(tweet.ID, 10)

	bytes, err := json.MarshalIndent(tweet, "", "    ")
	if err != nil {
		return fmt.Errorf("could not marshal tweet JSON: %s", err)
	}

	if err := ioutil.WriteFile(path.Join(ts.SaveDir, tweetId+".json"), bytes, 0644); err != nil {
		return fmt.Errorf("could not write JSON file: %s", err)
	}

	return nil
}

func (ts *TweetSaver) saveMedia(tweet twitter.Tweet) error {
	tweetId := strconv.FormatInt(tweet.ID, 10)

	num := 1
	for _, media := range extractMedia(tweet.ExtendedEntities.Media) {
		ext, err := getExtensionFromURL(media)
		if err != nil {
			return fmt.Errorf("could not save tweet ID %s media: %s", tweetId, media)
		}

		saveMediaFromURL(media, path.Join(ts.SaveDir, tweetId+"-"+strconv.Itoa(num)+ext))
		num++
	}

	return nil
}

func hasMedia(tweet twitter.Tweet) bool {
	return (tweet.ExtendedEntities != nil && len(tweet.ExtendedEntities.Media) > 0)
}

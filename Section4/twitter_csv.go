package main

import (
	"encoding/json"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"io/ioutil"
	"net/url"
	//"os"
)

var(
	consumerKey = "zh8CBpJvq4plTsoY2ERkV45Be"
	consumerSecret = "wj3gDQ3VDEjL5OWf2MMLSmgLHoUPujkh8589yGDUXUNGamFna0"
	accessToken = "1214744656442994688-01PNTehLqJBhFykIYYcuAt48tBQUYY"
	accessSecret = "wfFFbrTMiHutcqK7z7rvUeQ3Kn3aZEs5Pa4iqSiPS7RTi"
	//consumerKey = os.Getenv("CONSUMER_KEY_TWITTER")
	//consumerSecret = os.Getenv("CONSUMER_SECRET_TWITTER")
	//accessToken = os.Getenv("ACCESS_KEY_TWITTER")
	//accessSecret = os.Getenv("CONSUMER_SECRET_TWITTER")
)

type cleanTweet struct {
	Id string
	Text string
	Likes int
	Retweets int
	Language string
	URL string
}

var CleanTweet []cleanTweet

func getCleanTweet(tweet anaconda.Tweet) cleanTweet{
	var t = cleanTweet{
		Id:       tweet.IdStr,
		Text:     tweet.Text,
		Likes:    tweet.FavoriteCount,
		Retweets: tweet.RetweetCount,
		Language: tweet.Lang,
		URL:      "www.twitter.com/i/web/status/" + tweet.IdStr,
	}
	return t
}

func PrettyPrintTweet(tweet anaconda.Tweet){
	t := getCleanTweet(tweet)
	tweetJSON, _ := json.MarshalIndent(t, "", "\t")
	fmt.Println(string(tweetJSON))
}

func SaveTweetJSON(TweetsJSON []cleanTweet) error{
	tweetJSON, _ := json.MarshalIndent(TweetsJSON, "", "\t")
	err := ioutil.WriteFile("tweets.json", tweetJSON, 0644)
	if err != nil{
		return err
	}
	return nil
}

func LoadTweetsJSON() ([]cleanTweet, error){
	fileData, err := ioutil.ReadFile("tweets.json")

	if err != nil{
		return CleanTweet, err
	}
	err = json.Unmarshal(fileData, &CleanTweet)
	if err != nil{
		return CleanTweet, err
	}
	return CleanTweet, nil
}

func main() {
	Tweets, err:= LoadTweetsJSON()
	if err == nil{
		fmt.Println("Loading Tweets First ..")
		for _, t := range Tweets{
			tweetJSON, _ := json.MarshalIndent(t, "", "\t")
			fmt.Println(string(tweetJSON))
		}
	}

	api := anaconda.NewTwitterApiWithCredentials(
		accessToken, accessSecret, consumerKey, consumerSecret)
	fmt.Println("Started the API...")

	searchResult, _ := api.GetSearch("Falcon 9",
		url.Values{"result_type": []string{"popular"}})


	fmt.Printf("Retrieved %v tweets\n",
		len(searchResult.Statuses))

	var TweetsForfile []cleanTweet
	for _, tweet := range searchResult.Statuses{
		if !tweet.Retweeted && tweet.FavoriteCount > 900 && tweet.RetweetCount > 500 {
			TweetsForfile = append(TweetsForfile, getCleanTweet(tweet))
		}else{
			fmt.Println("Skipping tweet")
			fmt.Printf("\n")
			//PrettyPrintTweet(tweet)
		}
	}

	err = SaveTweetJSON(TweetsForfile)
	if err != nil {
		fmt.Println("Error in saving tweets")
	} else {
		fmt.Println("Successfully saved popular tweets!")
	}
}

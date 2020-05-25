package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"io/ioutil"
	"net/url"
	"os"
	"strconv"
)

var(
	consumerKey = os.Getenv("CONSUMER_KEY_TWITTER")
	consumerSecret = os.Getenv("CONSUMER_SECRET_TWITTER")
	accessToken = os.Getenv("ACCESS_KEY_TWITTER")
	accessSecret = os.Getenv("CONSUMER_SECRET_TWITTER")
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

func SaveTweetsCSV(tweets []cleanTweet) error{
	file, err := os.Create("tweets.csv")

	defer file.Close()
	if err != nil{
		panic(err)
	}

	w := csv.NewWriter(file)
	defer w.Flush()

	err = w.Write([]string{
		"Index", "ID", "Likes", "Retweets", "Language", "URL", "Text"})
	if err != nil{
		return err
	}

	for idx, tweet := range tweets{
		stringData := []string{strconv.Itoa(idx), tweet.Id,
			strconv.Itoa(tweet.Likes),
			strconv.Itoa(tweet.Retweets),
			tweet.Language, tweet.URL, tweet.Text}
		err = w.Write(stringData)
		if err != nil{
			return err
		}
	}
	return nil
}

func LoadTweetCSV() ([]cleanTweet, error){
	csvFile, err := os.Open("tweets.csv")

	defer csvFile.Close()
	if err != nil{
		return CleanTweet, err
	}

	lines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil{
		return CleanTweet, err
	}
	for idx, line := range lines{
		if idx == 0{
			continue
		}
		likecounts, _ := strconv.Atoi(line[3])
		retweetcounts, _ := strconv.Atoi(line[3])
		lineData := cleanTweet{
			Id:       line[1],
			Text:     line[2],
			Likes:    likecounts,
			Retweets: retweetcounts,
			Language: line[5],
			URL:      line[6],
		}
		CleanTweet = append(CleanTweet, lineData)
	}
	return CleanTweet, err
}


func main() {
	Tweets, err:= LoadTweetCSV()
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

	err = SaveTweetsCSV(TweetsForfile)
	if err != nil {
		fmt.Println("Error in saving tweets")
	} else {
		fmt.Println("Successfully saved popular tweets!")
	}
}

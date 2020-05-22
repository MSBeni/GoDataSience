package main

import(
	"encoding/json"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"net/url"
	"os"
)

var(
	consumerKey = os.Getenv("CONSUMER_KEY_TWITTER")
	consumerSecret = os.Getenv("CONSUMER_SECRET_TWITTER")
	accessToken = os.Getenv("ACCESS_KEY_TWITTER")
	accessSecret = os.Getenv("CONSUMER_SECRET_TWITTER")
)
func PrettyPrintTweet(tweet anaconda.Tweet){
	type cleanTweet struct {
		Id string
		Text string
		Likes int
		Retweets int
		Language string
		URL string
	}
	var t = cleanTweet{
		Id:       tweet.IdStr,
		Text:     tweet.Text,
		Likes:    tweet.FavoriteCount,
		Retweets: tweet.RetweetCount,
		Language: tweet.Lang,
		URL:      "www.twitter.com/i/web/status/" + tweet.IdStr,
	}
	tweetJSON, _ := json.MarshalIndent(t, "", "    ")
	fmt.Println(string(tweetJSON))
}

func main(){
	api := anaconda.NewTwitterApiWithCredentials(
		accessToken, accessSecret, consumerKey, consumerSecret)
	fmt.Println("Started the api ...")

	searchResult, _ := api.GetSearch("deep learning",
		url.Values{"result_type": []string{"popular"}})

	fmt.Printf("Retrieved %v tweets\n",
		len(searchResult.Statuses))

	for _, tweet := range searchResult.Statuses{
		if tweet.FavoriteCount > 90000 && tweet.RetweetCount > 50000 {
			_, err := api.Retweet(tweet.Id, false)
			if err != nil {
				fmt.Println("Error in Retweeting")
				continue
			}
		}else{
			fmt.Printf("Skeeping tweet")
			PrettyPrintTweet(tweet)
		}
	}
}
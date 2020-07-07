package main

import (
	R "GoDataScience/Section8/References"
	//"encoding/csv"
	"fmt"
	//"math"
	//"os"
	//"sort"
	//"strconv"
	//"strings"
)

func main() {
	_, allUsers := R.LoadData()

	thisUserId := 2
	sortedUsers := R.SimilarUsers(thisUserId)
	otherUserId := sortedUsers[1].UserId
	fmt.Println("Other user is", otherUserId)

	for i, rating := range R.Users[thisUserId].RatingsVector {
		if R.Users[otherUserId].RatingsVector[i] > 0 && rating > 0 {
			fmt.Printf("This user rates %s movie %0.1f stars, while other user rates it %0.1f stars\n",
				R.MoviesSlice[i].Title, rating, R.Users[otherUserId].RatingsVector[i])
		}
	}

	total := 0
	fmt.Println("-------")
	fmt.Println("Let's find movies that this user might like using user-user CF...")
	for _, otherUser := range sortedUsers {
		for _, k := range allUsers[otherUser.UserId].LikedMovies {
			_, ok := allUsers[thisUserId].Ratings[k.Title]
			if total == 20 {
				break
			}
			if !ok {
				fmt.Println(k.Title)
				total = total + 1
			}
		}
	}
	fmt.Println("-------")
	fmt.Println("Let's find movies using Content Recommendation System...")
	total = 0
	for _, m := range allUsers[thisUserId].LikedMovies {
		sortedMovies := R.GetMovieRecommendations(m.Title)
		for i, n := range sortedMovies {
			_, ok := allUsers[thisUserId].Ratings[n.Title]
			if total == 20 {
				break
			}
			if !ok {
				fmt.Println(n.Title)
				total = total + 1
			}
			if i == 1 {
				break
			}
		}
	}
}
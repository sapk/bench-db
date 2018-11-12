package generator

//Tweet generic struct
type Tweet struct {
	user string
	text string
}

//NewTweet generate a tweet
func NewTweet() Tweet {
	return Tweet{
		"me",
		"some random text",
	}
}

//NewTweetArray generate a tweet list
func NewTweetArray(n int) []Tweet {
	res := make([]Tweet, n)
	for i := 0; i < n; i++ {
		res[i] = NewTweet()
	}
	return res
}

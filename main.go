package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

const (
	BASE_URL = "https://www.reddit.com/api/v1/authorize?client_id=%s&response_type=%s&redirect_uri=%s&scope=*&state=1234"
	CLIENT_ID = "ohXpoqrZYub1kg"
	RESPONSE_TYPE = "token"
)

func main() {
	done := make(chan bool, 1)
	fmt.Println("Which subreddit do you want broadcast?")
	reader := bufio.NewReader(os.Stdin)
	subreddit, _ := reader.ReadString('\n')
	subreddit = strings.TrimSuffix(subreddit, "\n")
	fmt.Printf("We will make a post request subreddit: %v", subreddit)
	fmt.Println("What is the title?\n")
	reader = bufio.NewReader(os.Stdin)
	title, _ := reader.ReadString('\n')
	title = strings.TrimSuffix(title, "\n")
	fmt.Printf("Nice title: %v\n", title)

	go func() {
		<-time.After(2 * time.Second)
		browser(fmt.Sprintf(BASE_URL, CLIENT_ID, RESPONSE_TYPE, "http://localhost:65010/callback"))
	}()


	server := &http.Server{Addr: ":65010"}


	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Content-Type", "text/html")
		_, _ = writer.Write([]byte("Wait we are redirecting to another page. <script>window.location = 'http://localhost:65010/access_key?access_key=' + window.location.hash.substring(1);</script>"))
	})

	mux.HandleFunc("/access_key", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		accessKeyString := request.URL.Query().Get("access_key")
		parts := strings.Split(accessKeyString, "=")
		fmt.Printf("We got token %v, now lets make a post request in order to get authorization\n", parts[1])
		req, err := http.NewRequest("POST", "https://strapi.reddit.com/r/" + subreddit + "/broadcasts?title=" + url.QueryEscape(title), strings.NewReader(""))
		if err != nil {
			fmt.Println("Unable to create a new request")
			fmt.Println(err)
			return
		}
		req.Header.Add("User-Agent", "xxx/0.1")
		req.Header.Add("Authorization", "Bearer " + parts[1])
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("An error happen %v", err)
			return
		}
		respBody := ResponseBody{}
		buf := []byte{}
		_, _ = resp.Body.Read(buf)
		_ = resp.Body.Close()
		_ = json.Unmarshal(buf, &respBody)
		fmt.Println("We got token in order to be use in other software")
		fmt.Printf("Streamer Key: %v", respBody.Data.StreamerKey)
		fmt.Println("Streamer Endpoint: rtmp://ingest.redd.it/inbound/")
		done<-true
	})

	fmt.Println("Listening...")

	server.Handler = mux

	go func() {
		if err := server.ListenAndServe(); err != nil {
			fmt.Println("An error happend: %v", err)
		}
	}()

	<-done
	if err := server.Shutdown(context.TODO()); err != nil {
		fmt.Println("An error happend!")
	}

	fmt.Println("=====> end. Good Luck!")
}


func browser(url string) error {
	var err error

	switch runtime.GOOS {
		case "linux":
			err = exec.Command("xdg-open", url).Start()
		case "windows":
			err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
		case "darwin":
			err = exec.Command("open", url).Start()
		default:
			err = fmt.Errorf("unsupported platform")
	}
	return err
}

type ResponseBody struct {
	Status        string `json:"status"`
	StatusMessage string `json:"status_message"`
	Data          struct {
		VideoID     string `json:"video_id"`
		StreamerKey string `json:"streamer_key"`
		RtmpURL     string `json:"rtmp_url"`
		HlsURL      string `json:"hls_url"`
		Post        struct {
			Typename       string      `json:"__typename"`
			CreatedAt      string      `json:"createdAt"`
			CrosspostCount int         `json:"crosspostCount"`
			Domain         string      `json:"domain"`
			EditedAt       interface{} `json:"editedAt"`
			Flair          struct {
				Richtext  string `json:"richtext"`
				Text      string `json:"text"`
				Type      string `json:"type"`
				TextColor string `json:"textColor"`
				Template  struct {
					ID              string `json:"id"`
					Type            string `json:"type"`
					Text            string `json:"text"`
					Richtext        string `json:"richtext"`
					TextColor       string `json:"textColor"`
					BackgroundColor string `json:"backgroundColor"`
				} `json:"template"`
			} `json:"flair"`
			Awardings             []interface{} `json:"awardings"`
			ID                    string        `json:"id"`
			IsArchived            bool          `json:"isArchived"`
			IsContestMode         bool          `json:"isContestMode"`
			IsHidden              bool          `json:"isHidden"`
			IsLocked              bool          `json:"isLocked"`
			IsNsfw                bool          `json:"isNsfw"`
			IsOriginalContent     bool          `json:"isOriginalContent"`
			IsSaved               bool          `json:"isSaved"`
			IsScoreHidden         bool          `json:"isScoreHidden"`
			IsSelfPost            bool          `json:"isSelfPost"`
			IsSpoiler             bool          `json:"isSpoiler"`
			IsStickied            bool          `json:"isStickied"`
			IsVisited             bool          `json:"isVisited"`
			LiveCommentsWebsocket string        `json:"liveCommentsWebsocket"`
			ModerationInfo        interface{}   `json:"moderationInfo"`
			OutboundLink          struct {
				ExpiresAt string `json:"expiresAt"`
				URL       string `json:"url"`
			} `json:"outboundLink"`
			Permalink            string  `json:"permalink"`
			Score                float64 `json:"score"`
			SuggestedCommentSort string  `json:"suggestedCommentSort"`
			Title                string  `json:"title"`
			URL                  string  `json:"url"`
			WhitelistStatus      string  `json:"whitelistStatus"`
			VoteState            string  `json:"voteState"`
			AuthorInfo           struct {
				Typename string `json:"__typename"`
				ID       string `json:"id"`
				Name     string `json:"name"`
			} `json:"authorInfo"`
			AuthorOnlyInfo struct {
				ContentMode            string `json:"contentMode"`
				IsReceivingPostReplies bool   `json:"isReceivingPostReplies"`
			} `json:"authorOnlyInfo"`
			CommentCount    float64     `json:"commentCount"`
			Content         interface{} `json:"content"`
			DistinguishedAs interface{} `json:"distinguishedAs"`
			IsCrosspostable bool        `json:"isCrosspostable"`
			IsMediaOnly     bool        `json:"isMediaOnly"`
			IsPollIncluded  bool        `json:"isPollIncluded"`
			Media           struct {
				TypeHint  string      `json:"typeHint"`
				Streaming interface{} `json:"streaming"`
			} `json:"media"`
			PostEventInfo interface{} `json:"postEventInfo"`
			Thumbnail     interface{} `json:"thumbnail"`
			UpvoteRatio   float64     `json:"upvoteRatio"`
			ViewCount     interface{} `json:"viewCount"`
			Subreddit     struct {
				Typename string `json:"__typename"`
				ID       string `json:"id"`
				Styles   struct {
					LegacyIcon struct {
						Dimensions struct {
							Width  int `json:"width"`
							Height int `json:"height"`
						} `json:"dimensions"`
						URL string `json:"url"`
					} `json:"legacyIcon"`
					PrimaryColor string `json:"primaryColor"`
					Icon         string `json:"icon"`
				} `json:"styles"`
				Name              string  `json:"name"`
				Subscribers       float64 `json:"subscribers"`
				Title             string  `json:"title"`
				Type              string  `json:"type"`
				Path              string  `json:"path"`
				IsNSFW            bool    `json:"isNSFW"`
				IsQuarantined     bool    `json:"isQuarantined"`
				Wls               string  `json:"wls"`
				PrefixedName      string  `json:"prefixedName"`
				PostFlairSettings struct {
					Position  string `json:"position"`
					IsEnabled bool   `json:"isEnabled"`
				} `json:"postFlairSettings"`
				OriginalContentCategories  []interface{} `json:"originalContentCategories"`
				IsThumbnailsEnabled        bool          `json:"isThumbnailsEnabled"`
				IsFreeFormReportingAllowed bool          `json:"isFreeFormReportingAllowed"`
			} `json:"subreddit"`
		} `json:"post"`
		ShareLink string `json:"share_link"`
	} `json:"data"`
}
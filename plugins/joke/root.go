package joke

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"servon/core"
	"time"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

type Joke struct {
	ID     string `json:"id"`
	Joke   string `json:"joke"`
	Status int    `json:"status"`
}

type SearchResult struct {
	Results    json.RawMessage `json:"results"`
	SearchTerm string          `json:"search_term"`
	Status     int             `json:"status"`
	TotalJokes int             `json:"total_jokes"`
}

// JokeCmd represents the joke command
var JokeCmd = &cobra.Command{
	Use:   "joke",
	Short: color.Blue.Render("输出一条笑话"),
	Long:  color.Success.Render("\r\n输出一条笑话。\r\n可以随机输出一条笑话，或者通过关键字搜索特定主题的笑话。"),
	Run: func(cmd *cobra.Command, args []string) {
		jokeTerm, _ := cmd.Flags().GetString("term")

		if jokeTerm != "" {
			getRandomJokeWithTerm(jokeTerm)
		} else {
			getRandomJoke()
		}
	},
}

func init() {
	JokeCmd.PersistentFlags().String("term", "", color.Blue.Render("要查询的关键字"))
}

// Setup 注册joke插件到应用程序
func Setup(app *core.App) {
	app.GetRootCommand().AddCommand(JokeCmd)
}

func getRandomJoke() {
	url := "https://icanhazdadjoke.com/"
	responseBytes := getJokeData(url)
	joke := Joke{}

	if err := json.Unmarshal(responseBytes, &joke); err != nil {
		fmt.Printf("Could not unmarshal responseBytes. %v", err)
		return
	}

	fmt.Println(string(joke.Joke))
}

func getRandomJokeWithTerm(jokeTerm string) {
	total, results := getJokeDataWithTerm(jokeTerm)
	randomizeJokeList(total, results)
}

func getJokeData(baseAPI string) []byte {
	request, err := http.NewRequest(
		http.MethodGet, //method
		baseAPI,        //url
		nil,            //body
	)

	if err != nil {
		log.Printf("Could not request a dadjoke. %v", err)
		return nil
	}

	request.Header.Add("Accept", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Printf("Could not make a request. %v", err)
		return nil
	}
	defer response.Body.Close()

	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("Could not read response body. %v", err)
		return nil
	}

	return responseBytes
}

func getJokeDataWithTerm(jokeTerm string) (totalJokes int, jokeList []Joke) {
	url := fmt.Sprintf("https://icanhazdadjoke.com/search?term=%s", jokeTerm)
	responseBytes := getJokeData(url)

	jokeListRaw := SearchResult{}

	if err := json.Unmarshal(responseBytes, &jokeListRaw); err != nil {
		log.Printf("Could not unmarshal responseBytes. %v", err)
		return 0, nil
	}

	jokes := []Joke{}
	if err := json.Unmarshal(jokeListRaw.Results, &jokes); err != nil {
		log.Printf("Could not unmarshal responseBytes. %v", err)
		return 0, nil
	}

	return jokeListRaw.TotalJokes, jokes
}

func randomizeJokeList(length int, jokeList []Joke) {
	if length <= 0 {
		fmt.Println("没有找到相关的笑话")
		return
	}

	rand.Seed(time.Now().Unix())
	randomNum := rand.Intn(length)
	fmt.Println(jokeList[randomNum].Joke)
}

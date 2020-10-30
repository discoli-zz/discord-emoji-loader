package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/remeh/sizedwaitgroup"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	var token string
	fmt.Println("Welcome to /discord-emoji-loader/ — steal emojis without the hassle.")
	fmt.Print("Enter Discord Token: ")
	_, _ = fmt.Scanln(&token)

	s, err := discordgo.New(token)

	if err != nil {
		fmt.Println("There was an error with your token:", err)
		exit()
	}

	guilds, err := s.UserGuilds(100, "", "")
	if err != nil {
		fmt.Println("There was an error loading your guilds:", err)
		exit()
	}

	for _, g := range guilds {
		folder := g.Name + " (" + g.ID + ")"
		_ = os.Mkdir(folder, os.ModePerm)

		fmt.Println("Created the directory", folder)

		emojis, err := s.GuildEmojis(g.ID)
		if err != nil {
			fmt.Println("There was an error loading emojis for guild:", g.Name)
			continue
		}

		swg := sizedwaitgroup.New(30)
		for _, e := range emojis {
			swg.Add()
			go func(e *discordgo.Emoji) {
				defer swg.Done()

				extension := ".png"
				if e.Animated {
					extension = ".gif"
				}

				fetch := "https://cdn.discordapp.com/emojis/" + e.ID + extension
				resp, err := http.Get(fetch)
				if err != nil {
					fmt.Println("There was an error loading emoji:", e.Name)
					return
				}

				body, _ := ioutil.ReadAll(resp.Body)
				_ = ioutil.WriteFile(folder+"/"+e.Name+extension, body, os.ModePerm)
				fmt.Println("Downloaded the emoji:", e.Name)
			}(e)
		}
		swg.Wait()
	}

	fmt.Println("The process has completed.")
	exit()
}

func exit() {
	_, _ = fmt.Scanln()
	os.Exit(0)
}

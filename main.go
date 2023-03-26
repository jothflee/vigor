package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gokitcloud/ginkit"
	"github.com/jothflee/vigor/openai"
	log "github.com/sirupsen/logrus"
)

type CPostInput struct {
	ID       string           `json:"id"`
	Messages []openai.Message `json:"messages"`
}

func main() {
	r := ginkit.Default()
	r.GET("/", ginkit.H{
		"message": "pong",
	})

	r.POST("/api/d", func(c *gin.Context) {
		var req CPostInput
		// json unmarshal body to messages
		err := c.BindJSON(&req)
		if err != nil {
			log.Error(err)
			c.JSON(500, ginkit.H{
				"error": err.Error(),
			})
			return
		}

		// create a text file with the chat history
		chatHistoryBuilder := make([]string, len(req.Messages))
		for i, m := range req.Messages {
			chatHistoryBuilder[i] = fmt.Sprintf("%s: %s", m.Role, m.Content)
		}
		chatHistory := strings.Join(chatHistoryBuilder, "\n")
		// return chatHistory as a text file
		c.Writer.Header().Set("Content-Type", "text/plain")
		c.Writer.Header().Set("Content-Disposition", "attachment; filename=chat-history.txt")
		c.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", len(chatHistory)))
		c.Writer.Header().Set("Connection", "close")
		c.Writer.Header().Set("Expires", "0")
		c.Writer.Header().Set("Cache-Control", "must-revalidate, post-check=0, pre-check=0")
		c.Writer.Header().Set("Pragma", "public")
		c.Writer.Write([]byte(chatHistory))
	})

	r.GET("/api/ical/:id", func(c *gin.Context) {
		id := c.Param("id")
		ical, err := ioutil.ReadFile(fmt.Sprintf("calendar-%s.ics", id))
		if err != nil {
			log.Error(err)
			c.JSON(500, ginkit.H{
				"error": err.Error(),
			})
			return
		}

		c.Writer.Header().Set("Content-Type", "text/calendar")
		c.Writer.Header().Set("Content-Disposition", "attachment; filename=calendar.ics")
		c.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", len(ical)))
		c.Writer.Header().Set("Connection", "close")
		c.Writer.Header().Set("Expires", "0")
		c.Writer.Header().Set("Cache-Control", "must-revalidate, post-check=0, pre-check=0")
		c.Writer.Header().Set("Pragma", "public")
		c.Writer.Write(ical)
	})

	r.POST("/api/c", func(c *gin.Context) {
		var req CPostInput
		// json unmarshal body to messages
		err := c.BindJSON(&req)
		if err != nil {
			log.Error(err)
			c.JSON(500, ginkit.H{
				"error": err.Error(),
			})
			return
		}

		m, err := ChatWithChatGPT(req.Messages)
		if err != nil {
			log.Error(err)
			c.JSON(500, ginkit.H{
				"error": err.Error(),
			})
			return
		}

		message := openai.Message{
			Role:    "assistant",
			Content: m,
		}

		if strings.Contains(m, "BEGIN:VCALENDAR") {

			icalStr := m[strings.Index(m, "BEGIN:VCALENDAR"):(strings.Index(m, "END:VCALENDAR") + len("END:VCALENDAR"))]

			ioutil.WriteFile(fmt.Sprintf("calendar-%s.ics", req.ID), []byte(icalStr), 0644)

			message.Content = ""
			message.File = fmt.Sprintf("/api/ical/%s", req.ID)
			message.Mime = "text/calendar"

		}

		c.JSON(200, message)

	})
	r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

var CACHE = map[string]string{}

// ChatWithChatGPT
// inputs: userID string, messages []string
// outputs: response string, err error
func ChatWithChatGPT(messages []openai.Message) (response string, err error) {
	if len(messages) < 6 {
		if resp, ok := CACHE[messages[len(messages)-1].Content]; ok {
			return resp, nil
		}
	}

	ms := []openai.Message{
		{
			Role:    "system",
			Content: fmt.Sprintf("Todays date is: %s", time.Now().Format("2006-01-02")),
		},
		{
			Role:    "system",
			Content: "Act as an encouraging life coach that only helps users develop schedules to complete personal goals when the user may not know the steps to complete the task. Your tone is a blend of Carl Rogers and Bob Ross.",
		},
		{
			Role:    "system",
			Content: "Break down the goal into smaller sub-tasks. Prompt the user with specific details. Do not wait for the user to complete a task before prompting for the next task.",
		},
		{
			Role:    "system",
			Content: "Write each sub-task to contain a goal, a short description of the task, and a timeframe for completion.",
		},
		{
			Role:    "system",
			Content: "Be concise, use lists instead of paragraphs, ask at most 1 question per response",
		},
		{
			Role:    "system",
			Content: "Present one sub-task at a time until the sub task is fully planned.",
		},
		{
			Role:    "system",
			Content: "You can only generate iCal (.ics) formatted calendar files as text in the response. You do not upload the calendar anywhere.",
		},
		{
			Role:    "system",
			Content: "When all sub-tasks are planned and accepted by the user, offer to generate an ical.",
		},
	}

	// append the messages and filer out the file messages
	for _, m := range messages {
		if m.File == "" && m.Content != "" {
			ms = append(ms, m)
		}
	}

	// create new openai client using the environment variable OPENAI_API_KEY
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"), "")
	log.Info("sending message to chatgpt")
	resp, err := client.CreateChats(context.Background(), openai.CreateChatsRequest{
		Model:       "gpt-3.5-turbo",
		Messages:    ms,
		Temperature: .1,
	})

	// handle the error
	if err != nil {
		// log the error
		log.Errorf("error creating chat: %s", err.Error())

		return "", err
	}

	output := ""

	if len(resp.Choices) > 0 {
		output = resp.Choices[0].Message.Content

		log.Infof("cost: $%.6f\ntokens: %d\ninput: %d\noutput: %d", (float32(resp.Usage.TotalTokens)/1000)*.002, resp.Usage.TotalTokens, resp.Usage.PromptTokens, resp.Usage.CompletionTokens)
		CACHE[ms[len(ms)-1].Content] = output
	} else {
		log.Errorf("no response from chatgpt")
		log.Infof("%v", resp)
	}

	return output, nil
}

package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/jothflee/vigor/openai"
)

func TestSnip(t *testing.T) {
	m := `ddd
	
	BEGIN:VCALENDAR
test
END:VCALENDAR

xxxxx`
	if strings.Contains(m, "BEGIN:VCALENDAR") {

		icalStr := m[strings.Index(m, "BEGIN:VCALENDAR"):(strings.Index(m, "END:VCALENDAR") + len("END:VCALENDAR"))]

		fmt.Println(icalStr)

	}

}

// create a test for ChatWithChatGPT
func TestChatWithChatGPT(t *testing.T) {
	output, err := ChatWithChatGPT([]openai.Message{
		{
			Role:    "user",
			Content: "I'm Joe.",
		},
		{
			Role:    "assistant",
			Content: "Hello Joe! How can I assist you with your creative project? What type of project are you interested in completing and by when do you need it done?",
		},
		{
			Role:    "user",
			Content: "I am looking to write a book.",
		},
		{
			Role:    "assistant",
			Content: "That's awesome, Joe! Writing a book can be a very rewarding experience. Do you have a specific genre in mind for your book? And by when do you want to complete your book?",
		},
		{
			Role:    "user",
			Content: "A fiction book, a sci-fi action type story. I would like to have it done by the end of the year.",
		},
		{
			Role: "assistant",
			Content: `Great choice! A science fiction action book can be very exciting to write. Let's break down your goal into steps so that we can make a plan to achieve your goal. 

					The first step can be to create an outline for your book. The outline will help you to organize your thoughts and create a roadmap for your book. 
					
					Can you make a goal for creating the outline? For example, you could set a goal to have a detailed outline completed within the next week.`,
		},
		{
			Role:    "user",
			Content: "That sounds good.",
		},
		{
			Role:    "assistant",
			Content: "Great! So, your first step is to create a detailed outline for your book. When do you want to complete this step? Is there a specific day you want to finish creating your outline?",
		},
		{
			Role:    "user",
			Content: "in a few weeks.",
		},
	})
	// panic if error
	if err != nil {
		panic(err)
	}
	fmt.Println(output)
}

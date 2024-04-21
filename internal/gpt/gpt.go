package gpt

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

type GeneratedCourseInformation struct {
	Title string `json:"title"`
	Units []Unit `json:"units"`
}

type Unit struct {
	Title    string    `json:"title"`
	Chapters []Chapter `json:"chapters"`
}

type Chapter struct {
	ChapterTitle       string `json:"chapter_title"`
	YouTubeSearchQuery string `json:"youtube_search_query"`
	Summary            string `json:"summary"`
}
type ImageSearchTerm struct {
	SearchTerm string `json:"image_search_term"`
}

func GenerateCourseTitleAndUnitChapters(client *openai.Client, title string, unitValues map[string]string) (GeneratedCourseInformation, error) {
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "You are an AI capable of curating course content, coming up with relevent chapter titles, and finding relevant youtube videos for each chapter. Provide the search term in the JSON format as shown:{ title: \"title of the unit\", units: \"an array of units \"}. Each unit should have { chapters: \"an array of chapters\", title: \"title of the chapter\" }. Each chapter should have a youtube_search_query and chapter_title and summary (3 paragraphs long) key in the JSON object. Return at max 4 chapters per unit but make sure to generate chapters for every unit inputted by the user. Make sure to reword the course titles, unit titles, and chapter titles into better ones.",
		},
	}
	for _, name := range unitValues {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: fmt.Sprintf("I am creating a course about, %s. One of my units are named %s, for this unit create relevant chapters. For each chapter, provide a detailed youtube search query that can be used to find an informative educational video.", title, name),
		})
	}
	resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model:    openai.GPT3Dot5Turbo,
		Messages: messages,
	})
	if err != nil {
		return GeneratedCourseInformation{}, err
	}

	var courseTitleAndChapters GeneratedCourseInformation
	if err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &courseTitleAndChapters); err != nil {
		return GeneratedCourseInformation{}, err
	}
	return courseTitleAndChapters, nil
}

type GeneratedQuestion struct {
	Question string   `json:"question"`
	Answer   string   `json:"answer"`
	Options  []string `json:"options"`
}

func GenerateQuestionFromChapter(client *openai.Client, chapterName string) (GeneratedQuestion, error) {
	resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: `You are an AI capable of generating a question, answer, and 3 options (no more, no less) for a given chapter title. Provide these in the JSON format as shown: { "question": "generated question here", "answer": "generated answer here", "options": ["Option 1","Option 2","Option 3"] }`,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: fmt.Sprintf("I have a chapter titled, %s.", chapterName),
			},
		},
	})
	_ = resp
	if err != nil {
		return GeneratedQuestion{}, err
	}
	var generatedQuestion GeneratedQuestion
	if err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &generatedQuestion); err != nil {
		return generatedQuestion, err
	}
	return generatedQuestion, nil
}

func GetImageSearchTermFromTitle(client *openai.Client, title string) (ImageSearchTerm, error) {
	resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are an AI capable of suggesting an image search term for a given course title. Provide the search term in the JSON format as shown: { image_search_term: \"search term here\" }.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: fmt.Sprintf("I have a course titled, %s. Provide a good image search term. This search term will be fed into the unsplash API, so make sure it is a good search term that will return good results.", title),
			},
		},
	})
	if err != nil {
		return ImageSearchTerm{}, err
	}

	var imageSearchTerm ImageSearchTerm
	if err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &imageSearchTerm); err != nil {
		return ImageSearchTerm{}, err
	}
	return imageSearchTerm, nil
}

package quests

import (
	"encoding/json"
	"ethno/internal/models"
	"fmt"
)

var contentParsers = map[models.StepType]func(json.RawMessage) (models.StepContent, error){
	models.StepTypeIntro: func(raw json.RawMessage) (models.StepContent, error) {
		var c models.IntroContent
		return c, json.Unmarshal(raw, &c)
	},
	models.StepTypeQuiz: func(raw json.RawMessage) (models.StepContent, error) {
		var c models.QuizContent
		return c, json.Unmarshal(raw, &c)
	},
	models.StepTypeBuilder: func(raw json.RawMessage) (models.StepContent, error) {
		var c models.BuilderContent
		return c, json.Unmarshal(raw, &c)
	},
}

func ParseContent(stepType models.StepType, raw json.RawMessage) (models.StepContent, error) {
	parser, ok := contentParsers[stepType]
	if !ok {
		return nil, fmt.Errorf("unknown step type: %s", stepType)
	}
	return parser(raw)
}

func ParseStepContent(step *models.Step) error {
	parser, ok := contentParsers[step.Type]
	if !ok {
		return fmt.Errorf("unknown step type: %s", step.Type)
	}
	content, err := parser(step.RawContent)
	if err != nil {
		return err
	}
	step.Content = content
	return nil
}
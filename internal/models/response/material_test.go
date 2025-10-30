package response

import (
	"encoding/json"
	"testing"
)

func TestMaterialSummaryResponse_JSON(t *testing.T) {
	response := MaterialSummaryResponse{
		Sections: []SummarySection{
			{
				Title:                "Test Section",
				Content:              "Test content",
				Difficulty:           "basic",
				EstimatedTimeMinutes: 5,
				Order:                1,
			},
		},
		Glossary: []GlossaryTerm{
			{
				Term:       "Test Term",
				Definition: "Test definition",
				Order:      1,
			},
		},
		ReflectionQuestions: []string{"Test question?"},
		ProcessingMetadata: ProcessingMetadata{
			NLPProvider: "openai",
			Model:       "gpt-4",
			TokensUsed:  100,
		},
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var unmarshaled MaterialSummaryResponse
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if len(unmarshaled.Sections) != 1 {
		t.Errorf("Expected 1 section, got %d", len(unmarshaled.Sections))
	}
	if unmarshaled.Sections[0].Title != "Test Section" {
		t.Errorf("Expected 'Test Section', got '%s'", unmarshaled.Sections[0].Title)
	}
}

func TestAssessmentResponse_JSON(t *testing.T) {
	response := AssessmentResponse{
		Title:            "Test Quiz",
		Description:      "Test description",
		TotalQuestions:   2,
		TotalPoints:      100,
		PassingScore:     70,
		TimeLimitMinutes: 15,
		Questions: []Question{
			{
				ID:         "q1",
				Text:       "Test question?",
				Type:       "multiple_choice",
				Difficulty: "basic",
				Points:     50,
				Order:      1,
				Options: []QuestionOption{
					{ID: "a", Text: "Option A"},
					{ID: "b", Text: "Option B"},
				},
			},
		},
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var unmarshaled AssessmentResponse
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if unmarshaled.TotalQuestions != 2 {
		t.Errorf("Expected 2 questions, got %d", unmarshaled.TotalQuestions)
	}
	if len(unmarshaled.Questions) != 1 {
		t.Errorf("Expected 1 question, got %d", len(unmarshaled.Questions))
	}
	if len(unmarshaled.Questions[0].Options) != 2 {
		t.Errorf("Expected 2 options, got %d", len(unmarshaled.Questions[0].Options))
	}
}

func TestAttemptResultResponse_JSON(t *testing.T) {
	response := AttemptResultResponse{
		Score:       85.5,
		TotalPoints: 100,
		Passed:      true,
		DetailedFeedback: []QuestionFeedback{
			{
				QuestionID:      "q1",
				IsCorrect:       true,
				YourAnswer:      "a",
				FeedbackMessage: "Correct!",
			},
		},
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var unmarshaled AttemptResultResponse
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if unmarshaled.Score != 85.5 {
		t.Errorf("Expected score 85.5, got %f", unmarshaled.Score)
	}
	if len(unmarshaled.DetailedFeedback) != 1 {
		t.Errorf("Expected 1 feedback, got %d", len(unmarshaled.DetailedFeedback))
	}
	if !unmarshaled.DetailedFeedback[0].IsCorrect {
		t.Error("Expected is_correct to be true")
	}
}

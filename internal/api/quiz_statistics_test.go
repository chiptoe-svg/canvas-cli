package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestQuizStatisticsService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/1/quizzes/2/statistics" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(QuizStatisticsResponse{
			QuizStatistics: []QuizStatistics{
				{
					ID:     1,
					QuizID: 2,
					SubmissionStatistics: &QuizSubmissionStatistics{
						UniqueCount:  15,
						ScoreAverage: 82.5,
					},
				},
			},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewQuizStatisticsService(client)

	stats, err := svc.List(context.Background(), 1, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(stats) != 1 {
		t.Fatalf("expected 1 statistics entry, got %d", len(stats))
	}
	if stats[0].QuizID != 2 {
		t.Errorf("expected QuizID 2, got %d", stats[0].QuizID)
	}
	if stats[0].SubmissionStatistics == nil {
		t.Fatal("expected SubmissionStatistics, got nil")
	}
	if stats[0].SubmissionStatistics.UniqueCount != 15 {
		t.Errorf("expected UniqueCount 15, got %d", stats[0].SubmissionStatistics.UniqueCount)
	}
}

package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

func TestPollListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list polls successfully",
			Args: []string{},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/polls": cmdtest.NewMockResponse(`{"polls":[
					{"id":1,"question":"What is 2+2?","user_id":10},
					{"id":2,"question":"Favourite colour?","user_id":10}
				]}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "What is 2+2?") {
					t.Error("expected question in output")
				}
			},
		},
		{
			Name: "list polls - empty",
			Args: []string{},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/polls": cmdtest.NewMockResponse(`{"polls":[]}`),
			},
			ExpectError:  false,
			ExpectOutput: "No polls found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newPollListCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestPollGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get poll successfully",
			Args: []string{"1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/polls/1": cmdtest.NewMockResponse(`{"polls":[{"id":1,"question":"What is 2+2?","user_id":10}]}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "What is 2+2?") {
					t.Error("expected question in output")
				}
			},
		},
		{
			Name:        "get poll - missing ID",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newPollGetCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestPollCreateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "create poll successfully",
			Args: []string{"--question", "Which language?"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/polls": cmdtest.NewMockResponse(`{"polls":[{"id":5,"question":"Which language?","user_id":10}]}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Which language?") {
					t.Error("expected question in output")
				}
			},
		},
		{
			Name:        "create poll - missing question",
			Args:        []string{"--description", "some desc"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newPollCreateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestPollUpdateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "update poll successfully",
			Args: []string{"1", "--question", "Updated question"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/polls/1": cmdtest.NewMockResponse(`{"polls":[{"id":1,"question":"Updated question","user_id":10}]}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Updated question") {
					t.Error("expected updated question in output")
				}
			},
		},
		{
			Name:        "update poll - missing ID",
			Args:        []string{"--question", "x"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newPollUpdateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestPollDeleteCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "delete poll with force",
			Args: []string{"1", "--force"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/polls/1": {StatusCode: 204, Body: ""},
			},
			ExpectError: false,
		},
		{
			Name:        "delete poll - missing ID",
			Args:        []string{"--force"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newPollDeleteCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestPollChoiceListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list poll choices successfully",
			Args: []string{"--poll-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/polls/1/poll_choices": cmdtest.NewMockResponse(`{"poll_choices":[{"id":10,"poll_id":1,"text":"Option A","is_correct":true}]}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Option A") {
					t.Error("expected 'Option A' in output")
				}
			},
		},
		{
			Name:        "list poll choices - missing poll-id",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newPollChoiceListCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestPollChoiceGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get poll choice successfully",
			Args: []string{"10", "--poll-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/polls/1/poll_choices/10": cmdtest.NewMockResponse(`{"poll_choices":[{"id":10,"poll_id":1,"text":"Option A","is_correct":true}]}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Option A") {
					t.Error("expected 'Option A' in output")
				}
			},
		},
		{
			Name:        "get poll choice - missing choice ID",
			Args:        []string{"--poll-id", "1"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newPollChoiceGetCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestPollChoiceCreateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "create poll choice successfully",
			Args: []string{"--poll-id", "1", "--text", "Paris"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/polls/1/poll_choices": cmdtest.NewMockResponse(`{"poll_choices":[{"id":20,"poll_id":1,"text":"Paris","is_correct":true}]}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Paris") {
					t.Error("expected 'Paris' in output")
				}
			},
		},
		{
			Name:        "create poll choice - missing text",
			Args:        []string{"--poll-id", "1"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newPollChoiceCreateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestPollChoiceDeleteCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "delete poll choice with force",
			Args: []string{"10", "--poll-id", "1", "--force"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/polls/1/poll_choices/10": {StatusCode: 204, Body: ""},
			},
			ExpectError: false,
		},
		{
			Name:        "delete poll choice - missing poll-id",
			Args:        []string{"10", "--force"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newPollChoiceDeleteCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestPollSessionListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list poll sessions successfully",
			Args: []string{"--poll-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/polls/1/poll_sessions": cmdtest.NewMockResponse(`{"poll_sessions":[{"id":100,"poll_id":1,"course_id":999,"is_published":true}]}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "100") {
					t.Error("expected session ID in output")
				}
			},
		},
		{
			Name:        "list poll sessions - missing poll-id",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newPollSessionListCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestPollSessionCreateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "create poll session successfully",
			Args: []string{"--poll-id", "1", "--course-id", "999"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/polls/1/poll_sessions": cmdtest.NewMockResponse(`{"poll_sessions":[{"id":101,"poll_id":1,"course_id":999,"is_published":false}]}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "101") {
					t.Error("expected session ID in output")
				}
			},
		},
		{
			Name:        "create poll session - missing course-id",
			Args:        []string{"--poll-id", "1"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newPollSessionCreateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestPollSessionOpenCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "open poll session successfully",
			Args: []string{"100", "--poll-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/polls/1/poll_sessions/100/open": cmdtest.NewMockResponse(`{"poll_sessions":[{"id":100,"poll_id":1,"course_id":999,"is_published":true}]}`),
			},
			ExpectError: false,
		},
		{
			Name:        "open poll session - missing poll-id",
			Args:        []string{"100"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newPollSessionOpenCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestPollSessionCloseCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "close poll session successfully",
			Args: []string{"100", "--poll-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/polls/1/poll_sessions/100/close": cmdtest.NewMockResponse(`{"poll_sessions":[{"id":100,"poll_id":1,"course_id":999,"is_published":false}]}`),
			},
			ExpectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newPollSessionCloseCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestPollSessionListOpenedCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list opened poll sessions successfully",
			Args: []string{},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/poll_sessions/opened": cmdtest.NewMockResponse(`{"poll_sessions":[{"id":100,"poll_id":1,"course_id":999,"is_published":true}]}`),
			},
			ExpectError: false,
		},
		{
			Name: "list opened poll sessions - empty",
			Args: []string{},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/poll_sessions/opened": cmdtest.NewMockResponse(`{"poll_sessions":[]}`),
			},
			ExpectError:  false,
			ExpectOutput: "No open poll sessions found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newPollSessionListOpenedCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestPollSessionListClosedCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list closed poll sessions successfully",
			Args: []string{},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/poll_sessions/closed": cmdtest.NewMockResponse(`{"poll_sessions":[{"id":50,"poll_id":2,"course_id":888,"is_published":false}]}`),
			},
			ExpectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newPollSessionListClosedCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestPollSubmissionGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get poll submission successfully",
			Args: []string{"200", "--poll-id", "1", "--session-id", "100"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/polls/1/poll_sessions/100/poll_submissions/200": cmdtest.NewMockResponse(`{"poll_submissions":[{"id":200,"poll_choice_id":10,"user_id":42}]}`),
			},
			ExpectError: false,
		},
		{
			Name:        "get poll submission - missing session-id",
			Args:        []string{"200", "--poll-id", "1"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newPollSubmissionGetCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestPollSubmissionCreateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "create poll submission successfully",
			Args: []string{"--poll-id", "1", "--session-id", "100", "--choice-id", "10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/polls/1/poll_sessions/100/poll_submissions": cmdtest.NewMockResponse(`{"poll_submissions":[{"id":201,"poll_choice_id":10,"user_id":42}]}`),
			},
			ExpectError: false,
		},
		{
			Name:        "create poll submission - missing choice-id",
			Args:        []string{"--poll-id", "1", "--session-id", "100"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newPollSubmissionCreateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

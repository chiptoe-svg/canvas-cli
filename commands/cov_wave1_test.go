package commands

import (
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

// Coverage for the polls, folders, and appointment-groups command groups
// (Wave 1 resources that earlier coverage passes skipped).

func TestCovWave1_PollsList(t *testing.T) {
	cmdtest.RunCommandTest(t, newPollListCmd(), cmdtest.CommandTestCase{
		Name: "polls list",
		Args: []string{},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/polls": cmdtest.NewMockResponse(`{"polls":[{"id":1,"question":"Q?"}]}`),
		},
	})
}

func TestCovWave1_PollsGet(t *testing.T) {
	cmdtest.RunCommandTest(t, newPollGetCmd(), cmdtest.CommandTestCase{
		Name: "polls get",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/polls/1": cmdtest.NewMockResponse(`{"polls":[{"id":1,"question":"Q?"}]}`),
		},
	})
}

func TestCovWave1_PollsGetError(t *testing.T) {
	cmdtest.RunCommandTest(t, newPollGetCmd(), cmdtest.CommandTestCase{
		Name: "polls get API error",
		Args: []string{"9"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/polls/9": cmdtest.NewErrorResponse(404, "not found"),
		},
		ExpectError: true,
	})
}

func TestCovWave1_PollChoicesList(t *testing.T) {
	cmdtest.RunCommandTest(t, newPollChoiceListCmd(), cmdtest.CommandTestCase{
		Name: "poll choices list",
		Args: []string{"--poll-id", "1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/polls/1/poll_choices": cmdtest.NewMockResponse(`{"poll_choices":[{"id":3,"text":"A"}]}`),
		},
	})
}

func TestCovWave1_AppointmentGroupsList(t *testing.T) {
	cmdtest.RunCommandTest(t, newAppointmentGroupListCmd(), cmdtest.CommandTestCase{
		Name: "appointment-groups list",
		Args: []string{},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/appointment_groups": cmdtest.NewMockResponse(`[{"id":1,"title":"Office Hours"}]`),
		},
	})
}

func TestCovWave1_AppointmentGroupsGet(t *testing.T) {
	cmdtest.RunCommandTest(t, newAppointmentGroupGetCmd(), cmdtest.CommandTestCase{
		Name: "appointment-groups get",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/appointment_groups/1": cmdtest.NewMockResponse(`{"id":1,"title":"Office Hours"}`),
		},
	})
}

func TestCovWave1_AppointmentGroupsListError(t *testing.T) {
	cmdtest.RunCommandTest(t, newAppointmentGroupListCmd(), cmdtest.CommandTestCase{
		Name: "appointment-groups list API error",
		Args: []string{},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/appointment_groups": cmdtest.NewErrorResponse(500, "server error"),
		},
		ExpectError: true,
	})
}

func TestCovWave1_FoldersGet(t *testing.T) {
	cmdtest.RunCommandTest(t, newFoldersGetCmd(), cmdtest.CommandTestCase{
		Name: "folders get",
		Args: []string{"--folder-id", "5"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/folders/5": cmdtest.NewMockResponse(`{"id":5,"name":"Course Files"}`),
		},
	})
}

func TestCovWave1_FoldersListCourse(t *testing.T) {
	cmdtest.RunCommandTest(t, newFoldersListCmd(), cmdtest.CommandTestCase{
		Name: "folders list for course",
		Args: []string{"--course-id", "1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":         cmdtest.NewMockResponse(`{"id":1,"name":"Bio"}`),
			"/api/v1/courses/1/folders": cmdtest.NewMockResponse(`[{"id":5,"name":"Root"}]`),
		},
	})
}

func TestCovWave1_FoldersGetError(t *testing.T) {
	cmdtest.RunCommandTest(t, newFoldersGetCmd(), cmdtest.CommandTestCase{
		Name: "folders get API error",
		Args: []string{"--folder-id", "9"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/folders/9": cmdtest.NewErrorResponse(404, "not found"),
		},
		ExpectError: true,
	})
}

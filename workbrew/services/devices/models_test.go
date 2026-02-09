package devices

import (
	"encoding/json"
	"testing"
	"time"
)

func TestTimeOrStatus_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name       string
		jsonInput  string
		wantStatus string
		wantTime   bool
		wantErr    bool
	}{
		{
			name:       "Never status",
			jsonInput:  `"Never"`,
			wantStatus: "Never",
			wantTime:   false,
			wantErr:    false,
		},
		{
			name:       "Not Started status",
			jsonInput:  `"Not Started"`,
			wantStatus: "Not Started",
			wantTime:   false,
			wantErr:    false,
		},
		{
			name:       "Not Finished status",
			jsonInput:  `"Not Finished"`,
			wantStatus: "Not Finished",
			wantTime:   false,
			wantErr:    false,
		},
		{
			name:       "Valid RFC3339 time",
			jsonInput:  `"2023-11-01T12:34:56.000Z"`,
			wantStatus: "",
			wantTime:   true,
			wantErr:    false,
		},
		{
			name:       "Invalid time format",
			jsonInput:  `"invalid-time"`,
			wantStatus: "",
			wantTime:   false,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result TimeOrStatus
			err := json.Unmarshal([]byte(tt.jsonInput), &result)

			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			if result.Status != tt.wantStatus {
				t.Errorf("Status = %v, want %v", result.Status, tt.wantStatus)
			}

			if (result.Time != nil) != tt.wantTime {
				t.Errorf("Time present = %v, want %v", result.Time != nil, tt.wantTime)
			}
		})
	}
}

func TestTimeOrStatus_MarshalJSON(t *testing.T) {
	tests := []struct {
		name       string
		input      TimeOrStatus
		wantOutput string
	}{
		{
			name: "Never status",
			input: TimeOrStatus{
				Status: "Never",
			},
			wantOutput: `"Never"`,
		},
		{
			name: "Not Started status",
			input: TimeOrStatus{
				Status: "Not Started",
			},
			wantOutput: `"Not Started"`,
		},
		{
			name: "Not Finished status",
			input: TimeOrStatus{
				Status: "Not Finished",
			},
			wantOutput: `"Not Finished"`,
		},
		{
			name: "Time value",
			input: TimeOrStatus{
				Time: func() *time.Time {
					t, _ := time.Parse(time.RFC3339, "2023-11-01T12:34:56Z")
					return &t
				}(),
			},
			wantOutput: `"2023-11-01T12:34:56Z"`,
		},
		{
			name: "Nil time defaults to Never",
			input: TimeOrStatus{
				Time: nil,
			},
			wantOutput: `"Never"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := json.Marshal(tt.input)
			if err != nil {
				t.Errorf("MarshalJSON() error = %v", err)
				return
			}

			if string(result) != tt.wantOutput {
				t.Errorf("MarshalJSON() = %v, want %v", string(result), tt.wantOutput)
			}
		})
	}
}

func TestTimeOrStatus_HelperMethods(t *testing.T) {
	tests := []struct {
		name            string
		input           TimeOrStatus
		wantIsNever     bool
		wantIsNotStarted bool
		wantIsNotFinished bool
		wantHasTime     bool
		wantString      string
	}{
		{
			name: "Never status",
			input: TimeOrStatus{
				Status: "Never",
			},
			wantIsNever:     true,
			wantIsNotStarted: false,
			wantIsNotFinished: false,
			wantHasTime:     false,
			wantString:      "Never",
		},
		{
			name: "Not Started status",
			input: TimeOrStatus{
				Status: "Not Started",
			},
			wantIsNever:     false,
			wantIsNotStarted: true,
			wantIsNotFinished: false,
			wantHasTime:     false,
			wantString:      "Not Started",
		},
		{
			name: "Not Finished status",
			input: TimeOrStatus{
				Status: "Not Finished",
			},
			wantIsNever:     false,
			wantIsNotStarted: false,
			wantIsNotFinished: true,
			wantHasTime:     false,
			wantString:      "Not Finished",
		},
		{
			name: "Time value",
			input: TimeOrStatus{
				Time: func() *time.Time {
					t, _ := time.Parse(time.RFC3339, "2023-11-01T12:34:56Z")
					return &t
				}(),
			},
			wantIsNever:     false,
			wantIsNotStarted: false,
			wantIsNotFinished: false,
			wantHasTime:     true,
			wantString:      "2023-11-01T12:34:56Z",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.input.IsNever() != tt.wantIsNever {
				t.Errorf("IsNever() = %v, want %v", tt.input.IsNever(), tt.wantIsNever)
			}
			if tt.input.IsNotStarted() != tt.wantIsNotStarted {
				t.Errorf("IsNotStarted() = %v, want %v", tt.input.IsNotStarted(), tt.wantIsNotStarted)
			}
			if tt.input.IsNotFinished() != tt.wantIsNotFinished {
				t.Errorf("IsNotFinished() = %v, want %v", tt.input.IsNotFinished(), tt.wantIsNotFinished)
			}
			if tt.input.HasTime() != tt.wantHasTime {
				t.Errorf("HasTime() = %v, want %v", tt.input.HasTime(), tt.wantHasTime)
			}
			if tt.input.String() != tt.wantString {
				t.Errorf("String() = %v, want %v", tt.input.String(), tt.wantString)
			}
		})
	}
}

func TestTimeOrNever_BackwardCompatibility(t *testing.T) {
	tests := []struct {
		name       string
		jsonInput  string
		wantNever  bool
		wantTime   bool
		wantErr    bool
	}{
		{
			name:       "Never status",
			jsonInput:  `"Never"`,
			wantNever:  true,
			wantTime:   false,
			wantErr:    false,
		},
		{
			name:       "Valid RFC3339 time",
			jsonInput:  `"2023-11-01T12:34:56.000Z"`,
			wantNever:  false,
			wantTime:   true,
			wantErr:    false,
		},
		{
			name:       "Invalid time format",
			jsonInput:  `"invalid-time"`,
			wantNever:  false,
			wantTime:   false,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result TimeOrNever
			err := json.Unmarshal([]byte(tt.jsonInput), &result)

			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			if result.Never != tt.wantNever {
				t.Errorf("Never = %v, want %v", result.Never, tt.wantNever)
			}

			if (result.Time != nil) != tt.wantTime {
				t.Errorf("Time present = %v, want %v", result.Time != nil, tt.wantTime)
			}
		})
	}
}

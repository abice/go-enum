//go:build example
// +build example

package example

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type annotationTestData struct {
	Status AnnotationStatus `json:"status"`
	Color  AnnotationColor  `json:"color"`
	Number AnnotationNumber `json:"number"`
}

func TestAnnotationStatus(t *testing.T) {
	// Test prefix "My" was applied
	assert.Equal(t, MyAnnotationStatusPending, AnnotationStatus("pending"))
	assert.Equal(t, MyAnnotationStatusRunning, AnnotationStatus("running"))
	assert.Equal(t, MyAnnotationStatusCompleted, AnnotationStatus("completed"))
	assert.Equal(t, MyAnnotationStatusFailed, AnnotationStatus("failed"))

	// Test String()
	assert.Equal(t, "pending", MyAnnotationStatusPending.String())
	assert.Equal(t, "running", MyAnnotationStatusRunning.String())

	// Test IsValid()
	assert.True(t, MyAnnotationStatusPending.IsValid())
	assert.True(t, MyAnnotationStatusRunning.IsValid())
	assert.False(t, AnnotationStatus("invalid").IsValid())

	// Test Parse
	parsed, err := ParseAnnotationStatus("pending")
	assert.NoError(t, err)
	assert.Equal(t, MyAnnotationStatusPending, parsed)

	_, err = ParseAnnotationStatus("invalid")
	assert.Error(t, err)
	assert.Equal(t, "invalid is not a valid AnnotationStatus", err.Error())

	// Test Marshal/Unmarshal
	jsonData := `{"status":"pending"}`
	var data struct {
		Status AnnotationStatus `json:"status"`
	}
	err = json.Unmarshal([]byte(jsonData), &data)
	assert.NoError(t, err)
	assert.Equal(t, MyAnnotationStatusPending, data.Status)

	marshaled, err := json.Marshal(data)
	assert.NoError(t, err)
	assert.JSONEq(t, jsonData, string(marshaled))

	// Test AppendText (method has pointer receiver)
	status := MyAnnotationStatusPending
	text, err := status.AppendText(nil)
	assert.NoError(t, err)
	assert.Equal(t, "pending", string(text))
}

func TestAnnotationColor(t *testing.T) {
	// Test noprefix - no "AnnotationColor" prefix
	assert.Equal(t, AnnotationRed, AnnotationColor("annotation_red"))
	assert.Equal(t, AnnotationGreen, AnnotationColor("annotation_green"))
	assert.Equal(t, AnnotationBlue, AnnotationColor("annotation_blue"))

	// Test nocase - case insensitive parsing
	parsed, err := ParseAnnotationColor("ANNOTATION_RED")
	assert.NoError(t, err)
	assert.Equal(t, AnnotationRed, parsed)

	parsed, err = ParseAnnotationColor("annotation_red")
	assert.NoError(t, err)
	assert.Equal(t, AnnotationRed, parsed)

	parsed, err = ParseAnnotationColor("AnNoTaTiOn_ReD")
	assert.NoError(t, err)
	assert.Equal(t, AnnotationRed, parsed)

	// Test invalid
	_, err = ParseAnnotationColor("invalid")
	assert.Error(t, err)
	assert.Equal(t, "invalid is not a valid AnnotationColor", err.Error())

	// Test String()
	assert.Equal(t, "annotation_red", AnnotationRed.String())

	// Test IsValid()
	assert.True(t, AnnotationRed.IsValid())
	assert.False(t, AnnotationColor("invalid").IsValid())

	// Note: No marshal methods for AnnotationColor (not specified)
}

func TestAnnotationNumber(t *testing.T) {
	// Test constants
	assert.Equal(t, AnnotationNumberOne, AnnotationNumber(0))
	assert.Equal(t, AnnotationNumberTwo, AnnotationNumber(1))
	assert.Equal(t, AnnotationNumberThree, AnnotationNumber(2))

	// Test String()
	assert.Equal(t, "one", AnnotationNumberOne.String())
	assert.Equal(t, "two", AnnotationNumberTwo.String())
	assert.Equal(t, "three", AnnotationNumberThree.String())

	// Test IsValid()
	assert.True(t, AnnotationNumberOne.IsValid())
	assert.False(t, AnnotationNumber(999).IsValid())

	// Test Parse
	parsed, err := ParseAnnotationNumber("one")
	assert.NoError(t, err)
	assert.Equal(t, AnnotationNumberOne, parsed)

	_, err = ParseAnnotationNumber("invalid")
	assert.Error(t, err)
	assert.Equal(t, "invalid is not a valid AnnotationNumber", err.Error())

	// Test Marshal/Unmarshal
	jsonData := `{"number":"one"}`
	var data struct {
		Number AnnotationNumber `json:"number"`
	}
	err = json.Unmarshal([]byte(jsonData), &data)
	assert.NoError(t, err)
	assert.Equal(t, AnnotationNumberOne, data.Number)

	marshaled, err := json.Marshal(data)
	assert.NoError(t, err)
	assert.JSONEq(t, jsonData, string(marshaled))

	// Test SQL Scan/Value (basic test)
	var numScan AnnotationNumber
	err = numScan.Scan("one")
	assert.NoError(t, err)
	assert.Equal(t, AnnotationNumberOne, numScan)

	val, err := AnnotationNumberOne.Value()
	assert.NoError(t, err)
	assert.Equal(t, "one", val)

	// Test AppendText (method has pointer receiver)
	numAppend := AnnotationNumberOne
	text, err := numAppend.AppendText(nil)
	assert.NoError(t, err)
	assert.Equal(t, "one", string(text))
}

func TestAnnotationSQL(t *testing.T) {
	// Test AnnotationNumber SQL (enabled)
	var num AnnotationNumber

	// Scan from string
	err := num.Scan("two")
	assert.NoError(t, err)
	assert.Equal(t, AnnotationNumberTwo, num)

	// Scan from int
	err = num.Scan(1)
	assert.NoError(t, err)
	assert.Equal(t, AnnotationNumberTwo, num)

	// Value returns string
	val, err := num.Value()
	assert.NoError(t, err)
	assert.Equal(t, "two", val)

	// Test AnnotationStatus SQL (disabled - should not have Scan/Value methods)
	// We can't test absence directly, but we can verify that the type doesn't implement
	// driver.Valuer and sql.Scanner for AnnotationStatus (they're not generated)
}

func TestAnnotationMarshalCombined(t *testing.T) {
	// Test all three together
	jsonData := `{"status":"completed","color":"annotation_green","number":"three"}`
	var data annotationTestData
	err := json.Unmarshal([]byte(jsonData), &data)
	assert.NoError(t, err)
	assert.Equal(t, MyAnnotationStatusCompleted, data.Status)
	assert.Equal(t, AnnotationGreen, data.Color)
	assert.Equal(t, AnnotationNumberThree, data.Number)

	marshaled, err := json.Marshal(data)
	assert.NoError(t, err)
	assert.JSONEq(t, jsonData, string(marshaled))
}

func BenchmarkAnnotationParse(b *testing.B) {
	knownItems := []string{
		"pending",
		"annotation_red",
		"one",
	}

	var err error
	for _, item := range knownItems {
		b.Run(item, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				// Try to parse with appropriate parser
				switch {
				case strings.Contains(item, "annotation_"):
					_, err = ParseAnnotationColor(item)
				case item == "one" || item == "two" || item == "three":
					_, err = ParseAnnotationNumber(item)
				default:
					_, err = ParseAnnotationStatus(item)
				}
				assert.NoError(b, err)
			}
		})
	}
}

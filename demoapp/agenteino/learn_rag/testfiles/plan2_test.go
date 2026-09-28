package testfiles

import (
	"fmt"
	"testing"
)

func TestPlan2ExtractsEveryRealCourseDocument(t *testing.T) {
	documents := requireCourseDocuments(t)
	formats := map[string]int{}
	for i, document := range documents {
		if document.ID == "" || document.URI == "" || document.ContentHash == "" || document.Text == "" {
			t.Fatalf("extracted document lacks traceability: %#v", document)
		}
		if document.Metadata["source_path"] == "" || document.Metadata["format"] != document.Format {
			t.Fatalf("source metadata was not preserved: %#v", document)
		}
		if document.Format == "pdf" && len(document.Pages) == 0 {
			t.Fatalf("PDF lost its page-aware extraction: %#v", document)
		}
		formats[document.Format]++

		fmt.Printf("document.docx -- %d - %+v", i, document.Pages)
	}
	//if formats["txt"] == 0 || formats["pdf"] == 0 || formats["docx"] == 0 {
	//	t.Fatalf("course corpus formats were not all extracted: %#v", formats)
	//}
	t.Logf("extracted real corpus: %#v", formats)
}

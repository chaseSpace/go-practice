package testfiles

import "testing"

func TestPlan1RAGUsesTraceableEvidence(t *testing.T) {
	document := requireCourseDocuments(t)[0]
	chunks, err := chunkDocument(document, ChunkPolicy{MaxRunes: 500, OverlapRunes: 50})
	if err != nil {
		t.Fatalf("chunk real source %s: %v", document.URI, err)
	}
	result := RetrievedChunk{Chunk: chunks[0], Score: 1}

	if result.DocumentID != document.ID || result.URI == "" {
		t.Fatalf("retrieved evidence lost its source: %#v", result)
	}
	if result.DocumentID != document.ID || result.Text == "" {
		t.Fatalf("RAG evidence must originate from the loaded corpus: %#v", result)
	}

	for i, ck := range chunks {
		t.Logf("loop print: %+v", ck)
		if i >= 5 {
			break
		}
	}
}

package testfiles

import "testing"

func TestPlan3ChunkingKeepsStableSourceAndPDFPages(t *testing.T) {
	document := requireCourseDocumentFormat(t, "pdf")

	chunks, err := chunkDocument(document, ChunkPolicy{MaxRunes: 24, OverlapRunes: 4})
	if err != nil {
		t.Fatalf("chunk document: %v", err)
	}
	if len(chunks) < 2 {
		t.Fatalf("real PDF should produce multiple chunks, got %d", len(chunks))
	}
	for index, chunk := range chunks {
		if chunk.Index != index || chunk.DocumentID != document.ID || chunk.URI != document.URI {
			t.Fatalf("chunk %d lost stable provenance: %#v", index, chunk)
		}
		if chunk.PageStart != chunk.PageEnd || chunk.PageStart < 1 {
			t.Fatalf("chunk %d lost PDF page: %#v", index, chunk)
		}
	}

	again, err := chunkDocument(document, ChunkPolicy{MaxRunes: 24, OverlapRunes: 4})
	if err != nil {
		t.Fatalf("chunk document again: %v", err)
	}
	for index := range chunks {
		if chunks[index].ID != again[index].ID {
			t.Fatalf("chunk ID changed on identical input: %q vs %q", chunks[index].ID, again[index].ID)
		}
	}
}

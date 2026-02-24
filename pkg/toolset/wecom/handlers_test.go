package wecom

import (
	"encoding/base64"
	"strings"
	"testing"

	wecombot "github.com/futuretea/go-wecom-bot"
)

// --- getBot tests ---

func TestGetBot_Valid(t *testing.T) {
	bot := wecombot.New("test-key")
	result, err := getBot(bot)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result != bot {
		t.Fatalf("expected same bot instance")
	}
}

func TestGetBot_Nil(t *testing.T) {
	_, err := getBot(nil)
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestGetBot_WrongType(t *testing.T) {
	_, err := getBot("not a bot")
	if err == nil {
		t.Fatal("expected error for wrong type")
	}
}

// --- stringParam tests ---

func TestStringParam_Exists(t *testing.T) {
	params := map[string]any{"key": "value"}
	if got := stringParam(params, "key"); got != "value" {
		t.Fatalf("expected 'value', got %q", got)
	}
}

func TestStringParam_Missing(t *testing.T) {
	params := map[string]any{}
	if got := stringParam(params, "key"); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestStringParam_WrongType(t *testing.T) {
	params := map[string]any{"key": 123}
	if got := stringParam(params, "key"); got != "" {
		t.Fatalf("expected empty string for non-string type, got %q", got)
	}
}

// --- stringSliceParam tests ---

func TestStringSliceParam_AnySlice(t *testing.T) {
	params := map[string]any{"tags": []any{"a", "b", "c"}}
	got := stringSliceParam(params, "tags")
	if len(got) != 3 || got[0] != "a" || got[1] != "b" || got[2] != "c" {
		t.Fatalf("unexpected result: %v", got)
	}
}

func TestStringSliceParam_StringSlice(t *testing.T) {
	params := map[string]any{"tags": []string{"x", "y"}}
	got := stringSliceParam(params, "tags")
	if len(got) != 2 || got[0] != "x" || got[1] != "y" {
		t.Fatalf("unexpected result: %v", got)
	}
}

func TestStringSliceParam_Missing(t *testing.T) {
	params := map[string]any{}
	if got := stringSliceParam(params, "tags"); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestStringSliceParam_WrongType(t *testing.T) {
	params := map[string]any{"tags": "not-a-slice"}
	if got := stringSliceParam(params, "tags"); got != nil {
		t.Fatalf("expected nil for wrong type, got %v", got)
	}
}

// --- mapSliceParam tests ---

func TestMapSliceParam_Valid(t *testing.T) {
	params := map[string]any{
		"items": []any{
			map[string]any{"k": "v1"},
			map[string]any{"k": "v2"},
		},
	}
	got := mapSliceParam(params, "items")
	if len(got) != 2 {
		t.Fatalf("expected 2 items, got %d", len(got))
	}
}

func TestMapSliceParam_Missing(t *testing.T) {
	params := map[string]any{}
	if got := mapSliceParam(params, "items"); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

// --- mapParam tests ---

func TestMapParam_Valid(t *testing.T) {
	inner := map[string]any{"a": "b"}
	params := map[string]any{"obj": inner}
	got := mapParam(params, "obj")
	if got == nil || got["a"] != "b" {
		t.Fatalf("unexpected result: %v", got)
	}
}

func TestMapParam_Missing(t *testing.T) {
	params := map[string]any{}
	if got := mapParam(params, "obj"); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

// --- parseSource tests ---

func TestParseSource_Valid(t *testing.T) {
	params := map[string]any{
		"source": map[string]any{
			"icon_url": "https://example.com/icon.png",
			"desc":     "Test Source",
		},
	}
	iconURL, desc, ok := parseSource(params)
	if !ok {
		t.Fatal("expected ok to be true")
	}
	if iconURL != "https://example.com/icon.png" {
		t.Fatalf("unexpected iconURL: %s", iconURL)
	}
	if desc != "Test Source" {
		t.Fatalf("unexpected desc: %s", desc)
	}
}

func TestParseSource_Missing(t *testing.T) {
	params := map[string]any{}
	_, _, ok := parseSource(params)
	if ok {
		t.Fatal("expected ok to be false for missing source")
	}
}

// --- parseCardActionURL tests ---

func TestParseCardActionURL_Valid(t *testing.T) {
	params := map[string]any{
		"card_action": map[string]any{
			"url": "https://example.com",
		},
	}
	url, err := parseCardActionURL(params)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if url != "https://example.com" {
		t.Fatalf("unexpected url: %s", url)
	}
}

func TestParseCardActionURL_MissingAction(t *testing.T) {
	params := map[string]any{}
	_, err := parseCardActionURL(params)
	if err == nil {
		t.Fatal("expected error for missing card_action")
	}
}

func TestParseCardActionURL_MissingURL(t *testing.T) {
	params := map[string]any{
		"card_action": map[string]any{},
	}
	_, err := parseCardActionURL(params)
	if err == nil {
		t.Fatal("expected error for missing url in card_action")
	}
}

// --- Handler validation tests ---
// These test parameter validation only; they do not call the WeCom API.

func TestHandleSendText_EmptyContent(t *testing.T) {
	bot := wecombot.New("test-key")
	_, err := handleSendText(bot, map[string]any{})
	if err == nil || !strings.Contains(err.Error(), "content is required") {
		t.Fatalf("expected 'content is required' error, got %v", err)
	}
}

func TestHandleSendText_ContentTooLong(t *testing.T) {
	bot := wecombot.New("test-key")
	longContent := strings.Repeat("a", maxTextContentBytes+1)
	_, err := handleSendText(bot, map[string]any{"content": longContent})
	if err == nil || !strings.Contains(err.Error(), "exceeds maximum size") {
		t.Fatalf("expected size limit error, got %v", err)
	}
}

func TestHandleSendText_InvalidClient(t *testing.T) {
	_, err := handleSendText("not-a-bot", map[string]any{"content": "hello"})
	if err == nil {
		t.Fatal("expected error for invalid client")
	}
}

func TestHandleSendMarkdown_EmptyContent(t *testing.T) {
	bot := wecombot.New("test-key")
	_, err := handleSendMarkdown(bot, map[string]any{})
	if err == nil || !strings.Contains(err.Error(), "content is required") {
		t.Fatalf("expected 'content is required' error, got %v", err)
	}
}

func TestHandleSendMarkdown_ContentTooLong(t *testing.T) {
	bot := wecombot.New("test-key")
	longContent := strings.Repeat("a", maxMarkdownContentBytes+1)
	_, err := handleSendMarkdown(bot, map[string]any{"content": longContent})
	if err == nil || !strings.Contains(err.Error(), "exceeds maximum size") {
		t.Fatalf("expected size limit error, got %v", err)
	}
}

func TestHandleSendImage_MissingParams(t *testing.T) {
	bot := wecombot.New("test-key")
	_, err := handleSendImage(bot, map[string]any{})
	if err == nil || !strings.Contains(err.Error(), "base64 is required") {
		t.Fatalf("expected 'base64 is required' error, got %v", err)
	}

	_, err = handleSendImage(bot, map[string]any{"base64": "abc"})
	if err == nil || !strings.Contains(err.Error(), "md5 is required") {
		t.Fatalf("expected 'md5 is required' error, got %v", err)
	}
}

func TestHandleSendNews_EmptyArticles(t *testing.T) {
	bot := wecombot.New("test-key")
	_, err := handleSendNews(bot, map[string]any{})
	if err == nil || !strings.Contains(err.Error(), "articles is required") {
		t.Fatalf("expected 'articles is required' error, got %v", err)
	}
}

func TestHandleSendNews_TooManyArticles(t *testing.T) {
	bot := wecombot.New("test-key")
	articles := make([]any, maxNewsArticles+1)
	for i := range articles {
		articles[i] = map[string]any{"title": "t", "url": "u"}
	}
	_, err := handleSendNews(bot, map[string]any{"articles": articles})
	if err == nil || !strings.Contains(err.Error(), "must not exceed") {
		t.Fatalf("expected max articles error, got %v", err)
	}
}

func TestHandleSendNews_ArticleMissingTitle(t *testing.T) {
	bot := wecombot.New("test-key")
	articles := []any{map[string]any{"url": "https://example.com"}}
	_, err := handleSendNews(bot, map[string]any{"articles": articles})
	if err == nil || !strings.Contains(err.Error(), "must have a title") {
		t.Fatalf("expected 'must have a title' error, got %v", err)
	}
}

func TestHandleSendNews_ArticleMissingURL(t *testing.T) {
	bot := wecombot.New("test-key")
	articles := []any{map[string]any{"title": "Test"}}
	_, err := handleSendNews(bot, map[string]any{"articles": articles})
	if err == nil || !strings.Contains(err.Error(), "must have a url") {
		t.Fatalf("expected 'must have a url' error, got %v", err)
	}
}

func TestHandleSendTextNoticeCard_MissingTitle(t *testing.T) {
	bot := wecombot.New("test-key")
	_, err := handleSendTextNoticeCard(bot, map[string]any{})
	if err == nil || !strings.Contains(err.Error(), "main_title is required") {
		t.Fatalf("expected 'main_title is required' error, got %v", err)
	}
}

func TestHandleSendTextNoticeCard_MissingCardAction(t *testing.T) {
	bot := wecombot.New("test-key")
	_, err := handleSendTextNoticeCard(bot, map[string]any{"main_title": "Test"})
	if err == nil || !strings.Contains(err.Error(), "card_action is required") {
		t.Fatalf("expected 'card_action is required' error, got %v", err)
	}
}

func TestHandleSendNewsNoticeCard_MissingFields(t *testing.T) {
	bot := wecombot.New("test-key")
	_, err := handleSendNewsNoticeCard(bot, map[string]any{})
	if err == nil || !strings.Contains(err.Error(), "main_title is required") {
		t.Fatalf("expected 'main_title is required' error, got %v", err)
	}

	_, err = handleSendNewsNoticeCard(bot, map[string]any{"main_title": "Test"})
	if err == nil || !strings.Contains(err.Error(), "card_image_url is required") {
		t.Fatalf("expected 'card_image_url is required' error, got %v", err)
	}
}

func TestHandleUploadFile_MissingParams(t *testing.T) {
	bot := wecombot.New("test-key")
	_, err := handleUploadFile(bot, map[string]any{})
	if err == nil || !strings.Contains(err.Error(), "filename is required") {
		t.Fatalf("expected 'filename is required' error, got %v", err)
	}

	_, err = handleUploadFile(bot, map[string]any{"filename": "test.txt"})
	if err == nil || !strings.Contains(err.Error(), "base64_data is required") {
		t.Fatalf("expected 'base64_data is required' error, got %v", err)
	}
}

func TestHandleUploadFile_InvalidBase64(t *testing.T) {
	bot := wecombot.New("test-key")
	_, err := handleUploadFile(bot, map[string]any{
		"filename":    "test.txt",
		"base64_data": "not-valid-base64!!!",
	})
	if err == nil || !strings.Contains(err.Error(), "failed to decode base64 data") {
		t.Fatalf("expected base64 decode error, got %v", err)
	}
}

func TestHandleUploadFile_TooLarge(t *testing.T) {
	bot := wecombot.New("test-key")
	// Create data just over 20MB
	largeData := make([]byte, maxUploadFileBytes+1)
	encoded := base64.StdEncoding.EncodeToString(largeData)
	_, err := handleUploadFile(bot, map[string]any{
		"filename":    "large.bin",
		"base64_data": encoded,
	})
	if err == nil || !strings.Contains(err.Error(), "file size exceeds maximum") {
		t.Fatalf("expected file size error, got %v", err)
	}
}

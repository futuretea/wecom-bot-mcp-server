package wecom

import (
	"encoding/base64"
	"fmt"

	wecombot "github.com/futuretea/go-wecom-bot"
	"github.com/futuretea/go-wecom-bot/image"
	"github.com/futuretea/go-wecom-bot/markdown"
	"github.com/futuretea/go-wecom-bot/news"
	"github.com/futuretea/go-wecom-bot/templatecard"
	"github.com/futuretea/go-wecom-bot/text"
)

// WeCom API limits
const (
	maxTextContentBytes     = 2048
	maxMarkdownContentBytes = 4096
	maxNewsArticles         = 8
	maxUploadFileBytes      = 20 * 1024 * 1024 // 20MB

	// defaultCardImageAspectRatio is the default aspect ratio for news notice card images.
	defaultCardImageAspectRatio = 2.35
)

// getBot validates and returns the WeCom bot client from the generic client.
func getBot(client any) (*wecombot.Bot, error) {
	bot, ok := client.(*wecombot.Bot)
	if !ok || bot == nil {
		return nil, fmt.Errorf("weCom bot client is not configured")
	}
	return bot, nil
}

// stringParam extracts a string parameter from the params map.
func stringParam(params map[string]any, key string) string {
	value, exists := params[key]
	if !exists {
		return ""
	}
	str, _ := value.(string)
	return str
}

// stringSliceParam extracts a string slice parameter from the params map.
func stringSliceParam(params map[string]any, key string) []string {
	value, exists := params[key]
	if !exists {
		return nil
	}

	switch arr := value.(type) {
	case []any:
		result := make([]string, 0, len(arr))
		for _, item := range arr {
			if str, ok := item.(string); ok {
				result = append(result, str)
			}
		}
		return result
	case []string:
		return arr
	default:
		return nil
	}
}

// mapSliceParam extracts a slice of maps from the params map.
func mapSliceParam(params map[string]any, key string) []map[string]any {
	value, exists := params[key]
	if !exists {
		return nil
	}

	arr, ok := value.([]any)
	if !ok {
		return nil
	}

	result := make([]map[string]any, 0, len(arr))
	for _, item := range arr {
		if entry, ok := item.(map[string]any); ok {
			result = append(result, entry)
		}
	}
	return result
}

// mapParam extracts a map parameter from the params map.
func mapParam(params map[string]any, key string) map[string]any {
	value, exists := params[key]
	if !exists {
		return nil
	}
	m, _ := value.(map[string]any)
	return m
}

// parseSource extracts source fields (icon_url, desc) from the "source" param.
func parseSource(params map[string]any) (iconURL, desc string, ok bool) {
	source := mapParam(params, "source")
	if source == nil {
		return "", "", false
	}
	return stringParam(source, "icon_url"), stringParam(source, "desc"), true
}

// parseCardActionURL extracts and validates the card_action.url from params.
func parseCardActionURL(params map[string]any) (string, error) {
	cardAction := mapParam(params, "card_action")
	if cardAction == nil {
		return "", fmt.Errorf("card_action is required")
	}
	url := stringParam(cardAction, "url")
	if url == "" {
		return "", fmt.Errorf("card_action.url is required")
	}
	return url, nil
}

// handleSendText handles the send_text tool call.
func handleSendText(client any, params map[string]any) (string, error) {
	bot, err := getBot(client)
	if err != nil {
		return "", err
	}

	content := stringParam(params, "content")
	if content == "" {
		return "", fmt.Errorf("content is required")
	}
	if len(content) > maxTextContentBytes {
		return "", fmt.Errorf("content exceeds maximum size of %d bytes", maxTextContentBytes)
	}

	msg := text.New(content)

	// Add mentions if specified
	if mentionedList := stringSliceParam(params, "mentioned_list"); len(mentionedList) > 0 {
		msg.WithMention(mentionedList...)
	}
	if mentionedMobileList := stringSliceParam(params, "mentioned_mobile_list"); len(mentionedMobileList) > 0 {
		msg.WithMentionMobile(mentionedMobileList...)
	}

	if err := bot.Send(msg); err != nil {
		return "", fmt.Errorf("failed to send text message: %w", err)
	}

	return "Text message sent successfully", nil
}

// handleSendMarkdown handles the send_markdown tool call.
func handleSendMarkdown(client any, params map[string]any) (string, error) {
	bot, err := getBot(client)
	if err != nil {
		return "", err
	}

	content := stringParam(params, "content")
	if content == "" {
		return "", fmt.Errorf("content is required")
	}
	if len(content) > maxMarkdownContentBytes {
		return "", fmt.Errorf("content exceeds maximum size of %d bytes", maxMarkdownContentBytes)
	}

	msg := markdown.New(content)
	if err := bot.Send(msg); err != nil {
		return "", fmt.Errorf("failed to send markdown message: %w", err)
	}

	return "Markdown message sent successfully", nil
}

// handleSendImage handles the send_image tool call.
func handleSendImage(client any, params map[string]any) (string, error) {
	bot, err := getBot(client)
	if err != nil {
		return "", err
	}

	encodedImage := stringParam(params, "base64")
	if encodedImage == "" {
		return "", fmt.Errorf("base64 is required")
	}

	md5Hash := stringParam(params, "md5")
	if md5Hash == "" {
		return "", fmt.Errorf("md5 is required")
	}

	msg := image.New(encodedImage, md5Hash)
	if err := bot.Send(msg); err != nil {
		return "", fmt.Errorf("failed to send image message: %w", err)
	}

	return "Image message sent successfully", nil
}

// handleSendNews handles the send_news tool call.
func handleSendNews(client any, params map[string]any) (string, error) {
	bot, err := getBot(client)
	if err != nil {
		return "", err
	}

	articles := mapSliceParam(params, "articles")
	if len(articles) == 0 {
		return "", fmt.Errorf("articles is required and must not be empty")
	}
	if len(articles) > maxNewsArticles {
		return "", fmt.Errorf("articles must not exceed %d items, got %d", maxNewsArticles, len(articles))
	}

	msg := news.New()
	for _, article := range articles {
		title := stringParam(article, "title")
		if title == "" {
			return "", fmt.Errorf("each article must have a title")
		}
		url := stringParam(article, "url")
		if url == "" {
			return "", fmt.Errorf("each article must have a url")
		}
		msg.AddArticle(title, stringParam(article, "description"), url, stringParam(article, "picurl"))
	}

	if err := bot.Send(msg); err != nil {
		return "", fmt.Errorf("failed to send news message: %w", err)
	}

	return fmt.Sprintf("News message sent successfully with %d article(s)", len(articles)), nil
}

// handleSendTextNoticeCard handles the send_text_notice_card tool call.
func handleSendTextNoticeCard(client any, params map[string]any) (string, error) {
	bot, err := getBot(client)
	if err != nil {
		return "", err
	}

	mainTitle := stringParam(params, "main_title")
	if mainTitle == "" {
		return "", fmt.Errorf("main_title is required")
	}

	card := templatecard.NewTextNotice().
		WithMainTitle(mainTitle, stringParam(params, "main_title_desc"))

	if iconURL, desc, ok := parseSource(params); ok {
		// Use blue for text notice cards as a default visual distinction
		card.WithSource(iconURL, desc, templatecard.SourceDescColorBlue)
	}

	if subTitle := stringParam(params, "sub_title"); subTitle != "" {
		card.WithSubTitle(subTitle)
	}

	if emphasis := mapParam(params, "emphasis_content"); emphasis != nil {
		card.WithEmphasisContent(stringParam(emphasis, "title"), stringParam(emphasis, "desc"))
	}

	for _, content := range mapSliceParam(params, "horizontal_content_list") {
		card.AddHorizontalContent(stringParam(content, "keyname"), stringParam(content, "value"), templatecard.HorizontalContentTypeText)
	}

	for _, jump := range mapSliceParam(params, "jump_list") {
		card.AddJump(templatecard.JumpTypeURL, stringParam(jump, "title"), stringParam(jump, "url"))
	}

	actionURL, err := parseCardActionURL(params)
	if err != nil {
		return "", err
	}
	card.WithCardAction(templatecard.ActionTypeURL, actionURL)

	if err := bot.Send(card); err != nil {
		return "", fmt.Errorf("failed to send text notice card: %w", err)
	}

	return "Text notice card sent successfully", nil
}

// handleSendNewsNoticeCard handles the send_news_notice_card tool call.
func handleSendNewsNoticeCard(client any, params map[string]any) (string, error) {
	bot, err := getBot(client)
	if err != nil {
		return "", err
	}

	mainTitle := stringParam(params, "main_title")
	if mainTitle == "" {
		return "", fmt.Errorf("main_title is required")
	}

	cardImageURL := stringParam(params, "card_image_url")
	if cardImageURL == "" {
		return "", fmt.Errorf("card_image_url is required")
	}

	card := templatecard.NewNewsNotice().
		WithMainTitle(mainTitle, stringParam(params, "main_title_desc")).
		WithCardImage(cardImageURL, defaultCardImageAspectRatio)

	if iconURL, desc, ok := parseSource(params); ok {
		// Use green for news notice cards as a default visual distinction
		card.WithSource(iconURL, desc, templatecard.SourceDescColorGreen)
	}

	actionURL, err := parseCardActionURL(params)
	if err != nil {
		return "", err
	}
	card.WithCardAction(templatecard.ActionTypeURL, actionURL)

	if err := bot.Send(card); err != nil {
		return "", fmt.Errorf("failed to send news notice card: %w", err)
	}

	return "News notice card sent successfully", nil
}

// handleUploadFile handles the upload_file tool call.
func handleUploadFile(client any, params map[string]any) (string, error) {
	bot, err := getBot(client)
	if err != nil {
		return "", err
	}

	filename := stringParam(params, "filename")
	if filename == "" {
		return "", fmt.Errorf("filename is required")
	}

	base64Data := stringParam(params, "base64_data")
	if base64Data == "" {
		return "", fmt.Errorf("base64_data is required")
	}

	// Decode base64 data
	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 data: %w", err)
	}
	if len(data) > maxUploadFileBytes {
		return "", fmt.Errorf("file size exceeds maximum of %d bytes", maxUploadFileBytes)
	}

	media, err := bot.UploadMedia(filename, data)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	return fmt.Sprintf("File uploaded successfully. media_id: %s, type: %s, created_at: %s",
		media.MediaID, media.Type, media.CreatedAt), nil
}

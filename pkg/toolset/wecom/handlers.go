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

// getBot validates and returns the WeCom bot client from the generic client.
func getBot(client any) (*wecombot.Bot, error) {
	bot, ok := client.(*wecombot.Bot)
	if !ok || bot == nil {
		return nil, fmt.Errorf("WeCom bot client is not configured")
	}
	return bot, nil
}

// getStringParam extracts a string parameter from the map.
func getStringParam(m map[string]any, key string) string {
	v, ok := m[key]
	if !ok {
		return ""
	}
	s, _ := v.(string)
	return s
}

// getStringSliceParam extracts a string slice parameter from the map.
func getStringSliceParam(m map[string]any, key string) []string {
	v, ok := m[key]
	if !ok {
		return nil
	}

	switch arr := v.(type) {
	case []any:
		result := make([]string, 0, len(arr))
		for _, item := range arr {
			if s, ok := item.(string); ok {
				result = append(result, s)
			}
		}
		return result
	case []string:
		return arr
	default:
		return nil
	}
}

// getMapSliceParam extracts a slice of maps from the map.
func getMapSliceParam(m map[string]any, key string) []map[string]any {
	v, ok := m[key]
	if !ok {
		return nil
	}

	arr, ok := v.([]any)
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

// getMapParam extracts a map parameter from the map.
func getMapParam(m map[string]any, key string) map[string]any {
	v, ok := m[key]
	if !ok {
		return nil
	}
	sub, _ := v.(map[string]any)
	return sub
}

// parseSource extracts source fields (icon_url, desc) from the "source" param.
func parseSource(params map[string]any) (iconURL, desc string, ok bool) {
	source := getMapParam(params, "source")
	if source == nil {
		return "", "", false
	}
	return getStringParam(source, "icon_url"), getStringParam(source, "desc"), true
}

// parseCardActionURL extracts and validates the card_action.url from params.
func parseCardActionURL(params map[string]any) (string, error) {
	cardAction := getMapParam(params, "card_action")
	if cardAction == nil {
		return "", fmt.Errorf("card_action is required")
	}
	url := getStringParam(cardAction, "url")
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

	content := getStringParam(params, "content")
	if content == "" {
		return "", fmt.Errorf("content is required")
	}

	msg := text.New(content)

	// Add mentions if specified
	if mentionedList := getStringSliceParam(params, "mentioned_list"); len(mentionedList) > 0 {
		msg.WithMention(mentionedList...)
	}
	if mentionedMobileList := getStringSliceParam(params, "mentioned_mobile_list"); len(mentionedMobileList) > 0 {
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

	content := getStringParam(params, "content")
	if content == "" {
		return "", fmt.Errorf("content is required")
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

	b64 := getStringParam(params, "base64")
	if b64 == "" {
		return "", fmt.Errorf("base64 is required")
	}

	md5Hash := getStringParam(params, "md5")
	if md5Hash == "" {
		return "", fmt.Errorf("md5 is required")
	}

	msg := image.New(b64, md5Hash)
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

	articles := getMapSliceParam(params, "articles")
	if len(articles) == 0 {
		return "", fmt.Errorf("articles is required and must not be empty")
	}

	msg := news.New()
	for _, article := range articles {
		title := getStringParam(article, "title")
		if title == "" {
			return "", fmt.Errorf("each article must have a title")
		}
		url := getStringParam(article, "url")
		if url == "" {
			return "", fmt.Errorf("each article must have a url")
		}
		msg.AddArticle(title, getStringParam(article, "description"), url, getStringParam(article, "picurl"))
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

	mainTitle := getStringParam(params, "main_title")
	if mainTitle == "" {
		return "", fmt.Errorf("main_title is required")
	}

	card := templatecard.NewTextNotice().
		WithMainTitle(mainTitle, getStringParam(params, "main_title_desc"))

	if iconURL, desc, ok := parseSource(params); ok {
		card.WithSource(iconURL, desc, templatecard.SourceDescColorBlue)
	}

	if subTitle := getStringParam(params, "sub_title"); subTitle != "" {
		card.WithSubTitle(subTitle)
	}

	if emphasis := getMapParam(params, "emphasis_content"); emphasis != nil {
		card.WithEmphasisContent(getStringParam(emphasis, "title"), getStringParam(emphasis, "desc"))
	}

	for _, hc := range getMapSliceParam(params, "horizontal_content_list") {
		card.AddHorizontalContent(getStringParam(hc, "keyname"), getStringParam(hc, "value"), templatecard.HorizontalContentTypeText)
	}

	for _, jump := range getMapSliceParam(params, "jump_list") {
		card.AddJump(templatecard.JumpTypeURL, getStringParam(jump, "title"), getStringParam(jump, "url"))
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

	mainTitle := getStringParam(params, "main_title")
	if mainTitle == "" {
		return "", fmt.Errorf("main_title is required")
	}

	cardImageURL := getStringParam(params, "card_image_url")
	if cardImageURL == "" {
		return "", fmt.Errorf("card_image_url is required")
	}

	card := templatecard.NewNewsNotice().
		WithMainTitle(mainTitle, getStringParam(params, "main_title_desc")).
		WithCardImage(cardImageURL, 2.35)

	if iconURL, desc, ok := parseSource(params); ok {
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

	filename := getStringParam(params, "filename")
	if filename == "" {
		return "", fmt.Errorf("filename is required")
	}

	base64Data := getStringParam(params, "base64_data")
	if base64Data == "" {
		return "", fmt.Errorf("base64_data is required")
	}

	// Decode base64 data
	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 data: %w", err)
	}

	media, err := bot.UploadMedia(filename, data)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	return fmt.Sprintf("File uploaded successfully. media_id: %s, type: %s, created_at: %s",
		media.MediaID, media.Type, media.CreatedAt), nil
}

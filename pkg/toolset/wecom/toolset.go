package wecom

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/wecom-bot-mcp-server/pkg/toolset"
)

// Toolset provides WeCom bot messaging tools.
type Toolset struct{}

// GetName returns the name of the toolset.
func (t *Toolset) GetName() string {
	return "wecom"
}

// GetDescription returns the description of the toolset.
func (t *Toolset) GetDescription() string {
	return "WeCom (WeChat Work) bot messaging tools"
}

// GetTools returns all WeCom bot tools.
func (t *Toolset) GetTools(_ any) []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("send_text",
				mcp.WithDescription("Send a text message through a WeCom bot webhook. Supports @mentioning users by ID or mobile number."),
				mcp.WithString("content",
					mcp.Required(),
					mcp.Description("The text content to send. Maximum 2048 bytes."),
				),
				mcp.WithArray("mentioned_list",
					mcp.Description("List of user IDs to @mention. Use \"@all\" to mention everyone."),
					mcp.Items(map[string]any{"type": "string"}),
				),
				mcp.WithArray("mentioned_mobile_list",
					mcp.Description("List of mobile numbers to @mention. Use \"@all\" to mention everyone."),
					mcp.Items(map[string]any{"type": "string"}),
				),
			),
			Handler: handleSendText,
		},
		{
			Tool: mcp.NewTool("send_markdown",
				mcp.WithDescription("Send a Markdown message through a WeCom bot webhook. Supports headings, bold, links, quotes, etc."),
				mcp.WithString("content",
					mcp.Required(),
					mcp.Description("The markdown content to send. Maximum 4096 bytes."),
				),
			),
			Handler: handleSendMarkdown,
		},
		{
			Tool: mcp.NewTool("send_image",
				mcp.WithDescription("Send an image (JPG/PNG, base64-encoded) through a WeCom bot webhook."),
				mcp.WithString("base64",
					mcp.Required(),
					mcp.Description("Base64-encoded image content. Max image size: 2MB. Supported formats: JPG, PNG."),
				),
				mcp.WithString("md5",
					mcp.Required(),
					mcp.Description("MD5 hash of the original image content (before base64 encoding)."),
				),
			),
			Handler: handleSendImage,
		},
		{
			Tool: mcp.NewTool("send_news",
				mcp.WithDescription("Send a news message (article list) through a WeCom bot webhook. Accepts 1-8 articles."),
				mcp.WithArray("articles",
					mcp.Required(),
					mcp.Description("Array of news articles (1-8 items)."),
					mcp.Items(map[string]any{
						"type": "object",
						"properties": map[string]any{
							"title": map[string]any{
								"type":        "string",
								"description": "Article title (required).",
							},
							"description": map[string]any{
								"type":        "string",
								"description": "Article description (optional).",
							},
							"url": map[string]any{
								"type":        "string",
								"description": "Article link URL (required).",
							},
							"picurl": map[string]any{
								"type":        "string",
								"description": "Article cover image URL (optional).",
							},
						},
						"required": []string{"title", "url"},
					}),
				),
			),
			Handler: handleSendNews,
		},
		{
			Tool: mcp.NewTool("send_text_notice_card",
				mcp.WithDescription("Send a text notice template card through a WeCom bot webhook. Supports highlighted content, key-value pairs, and links."),
				mcp.WithString("main_title",
					mcp.Required(),
					mcp.Description("Main title of the card."),
				),
				mcp.WithString("main_title_desc",
					mcp.Description("Description text below the main title."),
				),
				mcp.WithString("sub_title",
					mcp.Description("Subtitle text displayed in the card body."),
				),
				mcp.WithObject("source",
					mcp.Description("Source information displayed at the top of the card."),
					mcp.Properties(map[string]any{
						"icon_url": map[string]any{
							"type":        "string",
							"description": "URL of the source icon.",
						},
						"desc": map[string]any{
							"type":        "string",
							"description": "Source description text.",
						},
					}),
				),
				mcp.WithObject("emphasis_content",
					mcp.Description("Emphasized content area (large text)."),
					mcp.Properties(map[string]any{
						"title": map[string]any{
							"type":        "string",
							"description": "Emphasis title (displayed in large font).",
						},
						"desc": map[string]any{
							"type":        "string",
							"description": "Emphasis description.",
						},
					}),
				),
				mcp.WithArray("horizontal_content_list",
					mcp.Description("Key-value pairs displayed horizontally."),
					mcp.Items(map[string]any{
						"type": "object",
						"properties": map[string]any{
							"keyname": map[string]any{
								"type":        "string",
								"description": "Key name (label).",
							},
							"value": map[string]any{
								"type":        "string",
								"description": "Value text.",
							},
						},
						"required": []string{"keyname"},
					}),
				),
				mcp.WithArray("jump_list",
					mcp.Description("Jump links displayed at the bottom of the card."),
					mcp.Items(map[string]any{
						"type": "object",
						"properties": map[string]any{
							"title": map[string]any{
								"type":        "string",
								"description": "Jump link title.",
							},
							"url": map[string]any{
								"type":        "string",
								"description": "Jump link URL.",
							},
						},
						"required": []string{"title", "url"},
					}),
				),
				mcp.WithObject("card_action",
					mcp.Required(),
					mcp.Description("Card click action. Defines what happens when the card is clicked."),
					mcp.Properties(map[string]any{
						"url": map[string]any{
							"type":        "string",
							"description": "URL to open when the card is clicked (required).",
						},
					}),
				),
			),
			Handler: handleSendTextNoticeCard,
		},
		{
			Tool: mcp.NewTool("send_news_notice_card",
				mcp.WithDescription("Send a news notice template card with a cover image through a WeCom bot webhook."),
				mcp.WithString("main_title",
					mcp.Required(),
					mcp.Description("Main title of the card."),
				),
				mcp.WithString("main_title_desc",
					mcp.Description("Description text below the main title."),
				),
				mcp.WithString("card_image_url",
					mcp.Required(),
					mcp.Description("URL of the card cover image."),
				),
				mcp.WithObject("source",
					mcp.Description("Source information displayed at the top of the card."),
					mcp.Properties(map[string]any{
						"icon_url": map[string]any{
							"type":        "string",
							"description": "URL of the source icon.",
						},
						"desc": map[string]any{
							"type":        "string",
							"description": "Source description text.",
						},
					}),
				),
				mcp.WithObject("card_action",
					mcp.Required(),
					mcp.Description("Card click action. Defines what happens when the card is clicked."),
					mcp.Properties(map[string]any{
						"url": map[string]any{
							"type":        "string",
							"description": "URL to open when the card is clicked (required).",
						},
					}),
				),
			),
			Handler: handleSendNewsNoticeCard,
		},
		{
			Tool: mcp.NewTool("upload_file",
				mcp.WithDescription("Upload a file to the WeCom server (up to 20MB). Returns a media_id you can use to send file messages."),
				mcp.WithString("filename",
					mcp.Required(),
					mcp.Description("Name of the file to upload."),
				),
				mcp.WithString("base64_data",
					mcp.Required(),
					mcp.Description("Base64-encoded file content. Max file size: 20MB."),
				),
			),
			Handler: handleUploadFile,
		},
	}
}

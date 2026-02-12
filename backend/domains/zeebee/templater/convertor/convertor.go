package convertor

import (
	"bytes"
	"context"
	"html/template"

	"zeebee-backend/app/pdf_generator/client"
)

type Convertor interface {
	ConvertToHtml(ctx context.Context, name string, templateContent string, data map[string]interface{}) (content []byte, contentType string, err error)
	ConvertToPdf(ctx context.Context, name string, templateContent string, data map[string]interface{}) (content []byte, contentType string, err error)
}

type DefaultConvertor struct {
	// pdfGeneratorUrl string

	pdfClient *client.Client
}

func NewDefaultConvertor(pdfGeneratorUrl string) *DefaultConvertor {
	return &DefaultConvertor{
		// pdfGeneratorUrl: pdfGeneratorUrl,
		pdfClient: client.NewClient(pdfGeneratorUrl),
	}
}

func (c *DefaultConvertor) ConvertToHtml(ctx context.Context, name string, templateContent string, data map[string]interface{}) (content []byte, contentType string, err error) {
	contentType = "text/html"
	tpl, err := template.New(name).Parse(templateContent)
	if err != nil {
		return nil, contentType, err
	}
	var buf bytes.Buffer
	err = tpl.Execute(&buf, data)
	if err != nil {
		return nil, contentType, err
	}
	return buf.Bytes(), contentType, nil
}
func (c *DefaultConvertor) ConvertToPdf(ctx context.Context, name string, templateContent string, data map[string]interface{}) (content []byte, contentType string, err error) {
	contentType = "application/pdf"
	content, _, err = c.ConvertToHtml(ctx, name, templateContent, data)
	if err != nil {
		return nil, contentType, err
	}

	req := client.GeneratePdfRequest{
		Type: client.SourceTypeText,
		Text: string(content),
	}

	result, err := c.pdfClient.GeneratePdf(ctx, req)
	if err != nil {
		return nil, contentType, err
	}
	return result, contentType, nil

	// client := &http.Client{}
	// payload := GeneratePdfRequest{
	// 	Type: "Text",
	// 	Text: string(content),
	// }
	// requestBody, err := json.Marshal(payload)
	// if err != nil {
	// 	return nil, contentType, err
	// }
	//
	// req, err := http.NewRequestWithContext(
	// 	ctx,
	// 	http.MethodPost,
	// 	c.pdfGeneratorUrl,
	// 	bytes.NewReader(requestBody),
	// )
	// if err != nil {
	// 	return nil, contentType, err
	// }
	// req.Header.Set("Content-Type", "application/json")
	//
	// resp, err := client.Do(req)
	// if err != nil {
	// 	return nil, contentType, err
	// }
	// defer resp.Body.Close()
	// var generatePdfResponse GeneratePdfResponse
	// err = json.NewDecoder(resp.Body).Decode(&generatePdfResponse)
	// if err != nil {
	// 	return nil, contentType, err
	// }
	// if generatePdfResponse.Code != http.StatusOK {
	// 	return nil, contentType, fmt.Errorf(generatePdfResponse.Message)
	// }
	// if generatePdfResponse.Data.Data == nil {
	// 	return nil, contentType, fmt.Errorf("No data returned")
	// }
	// return generatePdfResponse.Data.Data, contentType, nil
}

type GeneratePdfRequest struct {
	Type string `json:"type"`
	Url  string `json:"url"`
	Text string `json:"text"`
}

type GeneratePdfResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Data []byte `json:"data"`
	} `json:"data"`
}

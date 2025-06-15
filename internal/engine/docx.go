package engine

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"baliance.com/gooxml/color"
	"baliance.com/gooxml/common"
	"baliance.com/gooxml/document"
	"baliance.com/gooxml/measurement"
	"baliance.com/gooxml/schema/soo/wml"
)

// DocxTemplate represents a DOCX template
type DocxTemplate struct {
	doc *document.Document
}

// DocxOptions provides configuration for DOCX rendering
type DocxOptions struct {
	// PreserveFormatting determines if placeholder formatting should be preserved
	PreserveFormatting bool
	// ImageDir specifies the directory to look for images
	ImageDir string
	// DefaultImageFormat is the format to use when no format is specified
	DefaultImageFormat string
}

// DefaultDocxOptions returns the default options for DOCX rendering
func DefaultDocxOptions() *DocxOptions {
	return &DocxOptions{
		PreserveFormatting: true,
		ImageDir:           "images",
		DefaultImageFormat: "png",
	}
}

// LoadDocxTemplate loads a DOCX template file
func LoadDocxTemplate(path string) (*DocxTemplate, error) {
	doc, err := document.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening DOCX file: %w", err)
	}

	return &DocxTemplate{doc: doc}, nil
}

// RenderDocxTemplate renders a DOCX template with the given data
func RenderDocxTemplate(templatePath string, data map[string]string, opts *DocxOptions) error {
	if opts == nil {
		opts = DefaultDocxOptions()
	}

	// Load template
	template, err := LoadDocxTemplate(templatePath)
	if err != nil {
		return err
	}
	doc := template.doc

	// Process paragraphs
	for _, para := range doc.Paragraphs() {
		if err := processParagraph(doc, para, data, opts); err != nil {
			return err
		}
	}

	// Process tables
	for _, table := range doc.Tables() {
		if err := processTable(doc, table, data, opts); err != nil {
			return err
		}
	}

	// Save the rendered document
	outputPath := strings.TrimSuffix(templatePath, filepath.Ext(templatePath)) + "_rendered.docx"
	return template.SaveDocxTemplate(outputPath)
}

// processParagraph processes a paragraph for placeholders
func processParagraph(doc *document.Document, para document.Paragraph, data map[string]string, opts *DocxOptions) error {
	// Check for conditional content
	if isConditionalParagraph(para) {
		if !shouldKeepParagraph(para, data) {
			// Remove paragraph by clearing its runs
			for _, run := range para.Runs() {
				run.Clear()
			}
			return nil
		}
	}

	// Process runs
	for _, run := range para.Runs() {
		text := run.Text()
		if text == "" {
			continue
		}

		// Check for image placeholder
		if isImagePlaceholder(text) {
			if err := replaceImagePlaceholder(doc, run, text, data, opts); err != nil {
				return err
			}
			continue
		}

		// Replace text placeholders
		if err := replaceTextPlaceholder(run, text, data, opts); err != nil {
			return err
		}
	}

	return nil
}

// processTable processes a table for placeholders
func processTable(doc *document.Document, table document.Table, data map[string]string, opts *DocxOptions) error {
	// Check for row repetition
	if isRepeatingRow(table) {
		return processRepeatingRow(doc, table, data, opts)
	}

	// Process each cell
	for _, row := range table.Rows() {
		for _, cell := range row.Cells() {
			for _, para := range cell.Paragraphs() {
				if err := processParagraph(doc, para, data, opts); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// isConditionalParagraph checks if a paragraph contains conditional content
func isConditionalParagraph(para document.Paragraph) bool {
	var text string
	for _, run := range para.Runs() {
		text += run.Text()
	}
	return strings.Contains(text, "{{#if") || strings.Contains(text, "{{#unless")
}

// shouldKeepParagraph determines if a conditional paragraph should be kept
func shouldKeepParagraph(para document.Paragraph, data map[string]string) bool {
	var text string
	for _, run := range para.Runs() {
		text += run.Text()
	}

	// Extract condition
	re := regexp.MustCompile(`{{#(if|unless)\s+(\w+)}}`)
	matches := re.FindStringSubmatch(text)
	if len(matches) != 3 {
		return true
	}

	conditionType := matches[1]
	conditionKey := matches[2]
	value, exists := data[conditionKey]

	if conditionType == "if" {
		return exists && value != "" && value != "false" && value != "0"
	} else { // unless
		return !exists || value == "" || value == "false" || value == "0"
	}
}

// isImagePlaceholder checks if text is an image placeholder
func isImagePlaceholder(text string) bool {
	return strings.HasPrefix(text, "{{image:") && strings.HasSuffix(text, "}}")
}

// replaceImagePlaceholder replaces an image placeholder with an actual image
func replaceImagePlaceholder(doc *document.Document, run document.Run, text string, data map[string]string, opts *DocxOptions) error {
	// Extract image key
	key := strings.TrimPrefix(text, "{{image:")
	key = strings.TrimSuffix(key, "}}")
	key = strings.TrimSpace(key)

	// Get image path from data
	imagePath, exists := data[key]
	if !exists {
		return fmt.Errorf("image data not found for key: %s", key)
	}

	// Create image struct
	img := common.Image{
		Path:   imagePath,
		Format: filepath.Ext(imagePath)[1:], // Remove the dot from extension
	}

	// Add image to document
	imgRef, err := doc.AddImage(img)
	if err != nil {
		return fmt.Errorf("error adding image: %w", err)
	}

	// Add drawing to run
	drawing, err := run.AddDrawingInline(imgRef)
	if err != nil {
		return fmt.Errorf("error adding drawing: %w", err)
	}

	// Set image size
	drawing.SetSize(2*measurement.Inch, 2*measurement.Inch)
	return nil
}

// isRepeatingRow checks if a table row should be repeated
func isRepeatingRow(table document.Table) bool {
	if len(table.Rows()) < 2 {
		return false
	}

	// Check if first row contains {{#each}}
	firstRow := table.Rows()[0]
	for _, cell := range firstRow.Cells() {
		for _, para := range cell.Paragraphs() {
			var text string
			for _, run := range para.Runs() {
				text += run.Text()
			}
			if strings.Contains(text, "{{#each") {
				return true
			}
		}
	}
	return false
}

// processRepeatingRow processes a repeating table row
func processRepeatingRow(doc *document.Document, table document.Table, data map[string]string, opts *DocxOptions) error {
	if len(table.Rows()) < 2 {
		return nil
	}

	// Get template row
	templateRow := table.Rows()[1]

	// Extract array key from first row
	var arrayKey string
	firstRow := table.Rows()[0]
	for _, cell := range firstRow.Cells() {
		for _, para := range cell.Paragraphs() {
			var text string
			for _, run := range para.Runs() {
				text += run.Text()
			}
			if strings.Contains(text, "{{#each") {
				re := regexp.MustCompile(`{{#each\s+(\w+)}}`)
				matches := re.FindStringSubmatch(text)
				if len(matches) == 2 {
					arrayKey = matches[1]
				}
			}
		}
	}

	if arrayKey == "" {
		return fmt.Errorf("no array key found in repeating row")
	}

	// Get array from data
	arrayStr, exists := data[arrayKey]
	if !exists {
		return fmt.Errorf("array data not found for key: %s", arrayKey)
	}

	// Parse array data
	arrayData := parseJSONArray(arrayStr)
	if len(arrayData) == 0 {
		return nil
	}

	// Clear template row content
	for _, cell := range templateRow.Cells() {
		for _, para := range cell.Paragraphs() {
			for _, run := range para.Runs() {
				run.Clear()
			}
		}
	}

	// Add rows for each array item
	for _, item := range arrayData {
		newRow := table.AddRow()
		for _, cell := range templateRow.Cells() {
			newCell := newRow.AddCell()
			for _, para := range cell.Paragraphs() {
				newPara := newCell.AddParagraph()
				for _, run := range para.Runs() {
					newRun := newPara.AddRun()
					if opts.PreserveFormatting {
						copyRunFormatting(run, newRun)
					}
					if err := replaceTextPlaceholder(newRun, run.Text(), item, opts); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

// copyRunFormatting copies formatting from one run to another
func copyRunFormatting(src, dst document.Run) {
	dst.Properties().SetBold(src.Properties().IsBold())
	dst.Properties().SetItalic(src.Properties().IsItalic())

	// Underline and color
	if src.Properties().X().U != nil {
		val := src.Properties().X().U.ValAttr
		if val == wml.ST_UnderlineSingle {
			dst.Properties().SetUnderline(wml.ST_UnderlineSingle, color.Auto)
		}
	}

	// Font size
	if src.Properties().X().Sz != nil {
		// Set default size to 12pt since we can't get the value directly
		dst.Properties().SetSize(measurement.Point * 12)
	}

	// Font color
	if src.Properties().X().Color != nil {
		dst.Properties().SetColor(color.FromHex(src.Properties().X().Color.ValAttr.String()))
	}
}

// replaceTextPlaceholder replaces a text placeholder with its value
func replaceTextPlaceholder(run document.Run, text string, data map[string]string, opts *DocxOptions) error {
	// Find all placeholders in the text
	re := regexp.MustCompile(`{{([^}]+)}}`)
	matches := re.FindAllStringSubmatch(text, -1)

	if len(matches) == 0 {
		return nil
	}

	// Replace each placeholder
	for _, match := range matches {
		if len(match) != 2 {
			continue
		}

		placeholder := match[0]
		key := strings.TrimSpace(match[1])

		// Skip special placeholders
		if strings.HasPrefix(key, "#") || strings.HasPrefix(key, "/") || strings.HasPrefix(key, "image:") {
			continue
		}

		// Get value from data
		value, exists := data[key]
		if !exists {
			return fmt.Errorf("data not found for key: %s", key)
		}

		// Replace placeholder with value
		text = strings.Replace(text, placeholder, value, 1)
	}

	// Update run text
	run.Clear()
	run.AddText(text)

	return nil
}

// parseJSONArray parses a JSON array string into a slice of maps
func parseJSONArray(str string) []map[string]string {
	var result []map[string]string
	if err := json.Unmarshal([]byte(str), &result); err != nil {
		return nil
	}
	return result
}

// SaveDocxTemplate saves the DOCX template to a file
func (t *DocxTemplate) SaveDocxTemplate(outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer file.Close()

	if err := t.doc.Save(file); err != nil {
		return fmt.Errorf("error saving document: %w", err)
	}

	return nil
}

// ExtractDocxPlaceholders extracts all placeholders from a DOCX template
func ExtractDocxPlaceholders(templatePath string) ([]string, error) {
	doc, err := document.Open(templatePath)
	if err != nil {
		return nil, fmt.Errorf("error opening DOCX file: %w", err)
	}

	placeholders := make(map[string]bool)

	// Extract from paragraphs
	for _, para := range doc.Paragraphs() {
		var text string
		for _, run := range para.Runs() {
			text += run.Text()
		}
		extractPlaceholders(text, placeholders)
	}

	// Extract from tables
	for _, table := range doc.Tables() {
		for _, row := range table.Rows() {
			for _, cell := range row.Cells() {
				for _, para := range cell.Paragraphs() {
					var text string
					for _, run := range para.Runs() {
						text += run.Text()
					}
					extractPlaceholders(text, placeholders)
				}
			}
		}
	}

	// Convert map to slice
	var result []string
	for placeholder := range placeholders {
		result = append(result, placeholder)
	}

	return result, nil
}

// extractPlaceholders extracts placeholders from text
func extractPlaceholders(text string, placeholders map[string]bool) {
	re := regexp.MustCompile(`{{([^}]+)}}`)
	matches := re.FindAllStringSubmatch(text, -1)

	for _, match := range matches {
		if len(match) != 2 {
			continue
		}

		key := strings.TrimSpace(match[1])
		if !strings.HasPrefix(key, "#") && !strings.HasPrefix(key, "/") {
			placeholders[key] = true
		}
	}
}

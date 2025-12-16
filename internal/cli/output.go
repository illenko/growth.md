package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorGray   = "\033[90m"
)

func PrintTable(data interface{}) error {
	if data == nil {
		return nil
	}

	v := reflect.ValueOf(data)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Slice:
		if v.Len() == 0 {
			fmt.Println("No results found.")
			return nil
		}
		return printSliceAsTable(v)
	case reflect.Struct:
		return printStructAsTable(v)
	default:
		fmt.Println(data)
		return nil
	}
}

func printSliceAsTable(slice reflect.Value) error {
	if slice.Len() == 0 {
		return nil
	}

	first := slice.Index(0)
	if first.Kind() == reflect.Ptr {
		first = first.Elem()
	}

	if first.Kind() != reflect.Struct {
		for i := 0; i < slice.Len(); i++ {
			fmt.Println(slice.Index(i).Interface())
		}
		return nil
	}

	headers, widths := getTableHeaders(first.Type())

	printTableHeader(headers, widths)
	printTableSeparator(widths)

	for i := 0; i < slice.Len(); i++ {
		item := slice.Index(i)
		if item.Kind() == reflect.Ptr {
			item = item.Elem()
		}
		printTableRow(item, headers, widths)
	}

	return nil
}

func printStructAsTable(s reflect.Value) error {
	headers, widths := getTableHeaders(s.Type())

	printTableHeader(headers, widths)
	printTableSeparator(widths)
	printTableRow(s, headers, widths)

	return nil
}

func getTableHeaders(t reflect.Type) ([]string, []int) {
	var headers []string
	var widths []int

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if field.Tag.Get("yaml") == "-" || field.Name == "Body" {
			continue
		}

		name := field.Name
		yamlTag := field.Tag.Get("yaml")
		if yamlTag != "" && yamlTag != "-" {
			parts := strings.Split(yamlTag, ",")
			if parts[0] != "" {
				name = parts[0]
			}
		}

		headers = append(headers, strings.ToUpper(name))
		widths = append(widths, max(len(name), 10))
	}

	return headers, widths
}

func printTableHeader(headers []string, widths []int) {
	for i, header := range headers {
		fmt.Printf("%-*s  ", widths[i], header)
	}
	fmt.Println()
}

func printTableSeparator(widths []int) {
	for _, width := range widths {
		fmt.Print(strings.Repeat("-", width) + "  ")
	}
	fmt.Println()
}

func printTableRow(item reflect.Value, headers []string, widths []int) {
	t := item.Type()
	headerIdx := 0

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if field.Tag.Get("yaml") == "-" || field.Name == "Body" {
			continue
		}

		fieldValue := item.Field(i)
		value := formatFieldValue(fieldValue)

		if len(value) > widths[headerIdx] {
			value = value[:widths[headerIdx]-3] + "..."
		}

		fmt.Printf("%-*s  ", widths[headerIdx], value)
		headerIdx++
	}
	fmt.Println()
}

func formatFieldValue(v reflect.Value) string {
	if !v.IsValid() {
		return ""
	}

	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			return ""
		}
		return formatFieldValue(v.Elem())
	case reflect.Slice, reflect.Array:
		if v.Len() == 0 {
			return ""
		}
		parts := make([]string, v.Len())
		for i := 0; i < v.Len() && i < 3; i++ {
			parts[i] = fmt.Sprint(v.Index(i).Interface())
		}
		if v.Len() > 3 {
			return strings.Join(parts, ",") + "..."
		}
		return strings.Join(parts, ",")
	case reflect.Struct:
		if timeValue, ok := v.Interface().(interface{ Format(string) string }); ok {
			return timeValue.Format("2006-01-02")
		}
		return fmt.Sprint(v.Interface())
	default:
		return fmt.Sprint(v.Interface())
	}
}

func PrintJSON(data interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func PrintYAML(data interface{}) error {
	encoder := yaml.NewEncoder(os.Stdout)
	encoder.SetIndent(2)
	return encoder.Encode(data)
}

func PrintSuccess(message string) {
	fmt.Printf("%s✓%s %s\n", colorGreen, colorReset, message)
}

func PrintError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s✗ Error:%s %s\n", colorRed, colorReset, err.Error())
	}
}

func PrintWarning(message string) {
	fmt.Printf("%s⚠%s  %s\n", colorYellow, colorReset, message)
}

func PrintInfo(message string) {
	fmt.Printf("%sℹ%s  %s\n", colorBlue, colorReset, message)
}

func Print(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func Println(args ...interface{}) {
	fmt.Println(args...)
}

func PrintOutput(data interface{}, format string) error {
	if format == "" {
		format = "table"
	}

	switch format {
	case "json":
		return PrintJSON(data)
	case "yaml":
		return PrintYAML(data)
	default:
		return PrintTable(data)
	}
}

func PrintOutputWithConfig(data interface{}) error {
	format := "table"
	if config != nil && config.Display.OutputFormat != "" {
		format = config.Display.OutputFormat
	}
	if outputFormat != "" {
		format = outputFormat
	}
	return PrintOutput(data, format)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

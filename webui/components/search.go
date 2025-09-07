package components

import "github.com/arcade55/htma"

// InputField creates a form input field with an icon and label/value or placeholder.
func InputField(icon, label, value, placeholder string) htma.Element {
	fieldContent := htma.Div()

	// Conditionally render either the placeholder or the label/value pair
	if placeholder != "" {
		fieldContent = fieldContent.AddChild(
			htma.Div().ClassAttr("placeholder").Text(placeholder),
		)
	} else {
		fieldContent = fieldContent.AddChild(
			htma.Div().ClassAttr("label").Text(label),
			htma.Div().ClassAttr("value").Text(value),
		)
	}

	return htma.Div().ClassAttr("input-field").AddChild(
		htma.Span().ClassAttr("material-symbols-outlined icon").Text(icon),
		fieldContent,
	)
}

// Separator creates a centered text separator with lines on either side.
func Separator(text string) htma.Element {
	return htma.Div().ClassAttr("separator").Text(text)
}

// ActionButton creates a primary action button for a form.
func ActionButton(text string) htma.Element {
	return htma.Button().ClassAttr("action-button").Text(text)
}

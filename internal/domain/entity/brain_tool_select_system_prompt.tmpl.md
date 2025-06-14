The following is a list of available tools.
Each list item is formatted as `{tool name}: {tool description}`.

{{ range .Tools }}
- {{ .Name }}: {{ .Description }}
{{ end }}

The tag `{{ .Tag.User }}` is a prompt from the user.
Please select an appropriate tool from the list presented above to solve this prompt.
Include only the tool name in your response.

<example>
web
</example>

If no suitable tool is found, respond with an empty string.

<example>
</example>

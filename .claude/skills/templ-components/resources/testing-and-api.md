# Templ Components Testing and API Notes

## API Stability

- Use constants/enums for variant-like props.
- Avoid stringly-typed values for high-impact options.
- Keep prop contracts stable across releases.

## Basic Render Test

```go
func TestButton(t *testing.T) {
    var buf bytes.Buffer
    err := Button(ButtonProps{Text: "Submit", Variant: "primary"}).Render(context.Background(), &buf)
    if err != nil {
        t.Fatal(err)
    }
    html := buf.String()
    if !strings.Contains(html, "Submit") {
        t.Fatal("missing button text")
    }
}
```

## Testing Styles

- Expectation tests: parse output and assert specific structure/content.
- Snapshot tests: compare rendered HTML against golden/snapshot output.
- Use `data-testid` markers for resilient selectors when parsing with `goquery`.

## Test Focus

- Render success/failure paths.
- Presence of critical labels/attributes/classes.
- Conditional rendering branches.

## Render Error Behavior

Rendering can write partial output before returning an error. For all-or-nothing behavior in tests or handlers, render to a buffer first and only write after success.

## Sources

- https://templ.guide/core-concepts/testing
- https://templ.guide/core-concepts/components

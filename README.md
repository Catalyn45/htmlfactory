# htmlfactory

Very simple tool to statically do templating with your html files.

## Usage

1. Write your html files, and where you want to insert a template, just use: `<factory src="./template.html">`.
example:
```html
<html>
    <head>
        <factory src="./templates/part.html">
    </head>
</html>
```

2. Run the tool on your html file
```sh
.\htmlfactory.exe ./page.html`
```

3. This will produce your newly html file, `page.compiled.html`:
```html
<html>
    <head>
        <title>hey</title>
    </head>
</html>
```

# Advantages
- Statically compiling your templates into your pages, 0 runtime cost
- You don't need to learn yet another new templating language, just the one you already now, html and js:
    - use  `<factory src="./template.html">` to insert a fragment
    - use `<factory src="./template.html" my_var="helo">` to send variables to templates
    - use `${my_var}` inside your fragments to use the variable. (ex. `<title>${my_var}</title>`)
- It's not a giant website generation tool with required file names and directory hierarchy, you just compile templates into your pages to avoid repetition

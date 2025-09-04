# htmlfactory

Very simple tool to statically do templating with your html files.

## Usage

```sh
Usage: .\htmlfactory.exe <input_file_or_directory> [options]
Options:
  --help
        Show usage
  --out string
        Path to output directory (default ".")
```

1. Write your html files, and where you want to insert a template, just use: `<fragment src="./template.html">`.
example:
```html
<html>
    <head>
        <fragment src="./templates/part.html">
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
    - use  `<fragment src="./template.html">` to insert a fragment
    - use `<fragment src="./template.html" my_var="helo">` to send variables to templates
    - use `${my_var}` inside your fragments to use the variable. (ex. `<title>${my_var}</title>`)
- It's not a giant website generation tool with required file names and directory hierarchy, you just compile templates into your pages to avoid repetition

# htmlfactory

Very simple tool to statically do templating with your html files.

## Advantages
- Statically compiling your templates into your pages, 0 runtime cost
- You don't need to learn yet another new templating language, just the ones you already now, html and js:
    - use  `<fragment src="./template.html"></fragment>` to insert a fragment
    - use `<fragment src="./template.html" my_var="helo"></fragment>` to send variables to fragments
    - use `${my_var}` inside your fragments to use the variable. (ex. `<title>${my_var}</title>`)
    - use `<any_tag content="fragment"> </any_tag>` in your fragment to insert the content of `<fragment>` tag in your fragment
- It's not a giant website generation tool with required file names and directory hierarchy, you just compile templates into your pages to avoid repetition

## Limitations
- You can't put fragments directly into the `<html>` or `<head>` tags, they need to be in `<body>`.
- The output html file will be correct but not formatted, this is how the `html.Render()` function makes it. A prettier renderer may be added to this project in the future.

## Usage

```sh
Usage: .\htmlfactory.exe <input_file_or_directory> [options]
Options:
  --help
        Show usage
  --out string
        Path to output directory (default ".")
  --watch
        keep watching for modified files and template them
```

1. Write your html files, and where you want to insert a template, just use: `<fragment src="./template.html"></fragment>`.
example:
```html
<html>
    <head>
        <fragment src="./templates/part.html"></fragment>
    </head>
</html>
```

2. Run the tool on your html file
```sh
.\htmlfactory.exe ./page.html
```

3. This will produce your newly html file, `page.compiled.html`:
```html
<html>
    <head>
        <title>hey</title>
    </head>
</html>
```

## Features

### Variables
In your main file you can have
```html
<html>
    <head>
        <fragment src="./templates/part.html" abc="10"></fragment>
    </head>
</html>
```

Then in the fragment:
```html
<h1> hello ${abc} </h1>
```

The result will be 
```html
<html>
    <head>
        <h1> hello 10 </h1>
    </head>
</html>
```

<hr>

### Block replacement
In your main file you can have
```html
<html>
    <head>
        <fragment src="./templates/part.html">
          <p> Hello </p>
        </fragment>
    </head>
</html>
```

Then in the fragment:
```html
<div content="fragment">
</div>
```

The result will be 
```html
<html>
    <head>
        <div>
          <p> Hello </p>
        </div>
    </head>
</html>
```

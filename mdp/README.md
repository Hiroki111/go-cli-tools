# MDP (Markdown Preview) Tool

## Dependencies

If you use Linux, make sure `xdg-utils` is installed. (Run `sudo apt-get install --reinstall xdg-utils` if you use Ubuntu and `xdg-open` isn't installed yet.)

## How to use

```bash
# Convert the MD file to HTML with a default template in main.go, show it on a browser, and delete the converted HTML
go run main.go -file <my-md-file>.md

# Convert the MD file to HTML with a template-fmt.html.tmpl (custom template that you're free to update), show it on a browser, and delete the converted HTML
go run main.go -file <my-md-file>.md -t template-fmt.html.tmpl

# Convert the MD file to HTML with a default template in main.go, show it on a browser, and save the converted HTML in /tmp (Linux) or C:\Users\<user>\AppData\Local\Temp (Windows)
go run main.go -file <my-md-file>.md -s
```
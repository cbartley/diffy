
## Diffy, an HTML Diff Tool

*Disclaimer: Diffy is still very much a work in progress.*

**Diffy** &mdash; rhymes with "spiffy" &mdash; is an HTML diff tool.  It compares two input files and generates HTML which will show a side-by-side "diff" of the two files when viewed in a web browser.

Diffy can output HTML to *stdout* like any good Unix tool and it can also output the HTML to a temporary file and then open that file with a web browser using a specific command that you provide.

### A diffy HTML snippet pasted into the Gmail Compose window

One of the key design goals for Diffy was to generate HTML which could easily be copied and pasted (in whole or in part) into email messages while retaining the formatting.  Of course it has to generate some pretty awful HTML to make this possible, but I figure it's worth it.

It works pretty well for Gmail.  I haven't tried any other email apps yet.

![](docs/gmail-compose-example.png?raw=true)

Note that you can insert comments between lines in the diff.

## Using Diffy

### Writing an HTML diff to a file

`diffy left.js right.js > differences.html`

### Displaying an HTML diff with Google Chrome on Mac OS X

`diffy --open-with 'open -a "Google Chrome.app"' left.js right.js`

## Supported Platforms

I have only built and tested diffy on Mac OS X so far.  It *should* work on other platforms without any changes, but I haven't had the opportunity to do any cross-platform testing yet.

## Licensing

Diffy is provided under the MIT License.
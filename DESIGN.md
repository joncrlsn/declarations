# Declarations Design

There are two components to this application, both served from the same Golang process.

- The first component is an API using the standard Golang HTTP framework that serves read-only text-based declarations of God's promises and who I am in Christ.
- The second component is an HTTP web client, also served with Golang HTTP. This HTML and Javascript interacts with the API
- Both components run in the same Golang process and listen on port 8080

## API

The Declarations API is a read-only JSON REST API that provides the following endpoints:

- GET /api/v1/declarations - provides a list of declarations
- GET /api/v1/declarations/random - returns a random declaration 
- GET /api/v1/declarations/{id} - returns the declaration with the given id, if it exists.  Otherwise, HTTP 404 is returned.
- GET /api/v1/declarations/label/{label} - provides list of declarations that have the given label
- GET /api/v1/labels - provides a list of labels that are used by declarations.  Included in the output is the count of declarations for each label.
- GET /api/v1/bible-esv/{reference} - provides the Bible text for the given reference from the ESV Bible API.  If the Bible text is returned in the response body.
- GET /api/v1/health - this API response includes the number of declarations currently stored

## UI

The web client provides a user interface for interacting with the API. It provides the following features:

- The UI is a single HTML page that is served from the same Golang process as the API.
  - GET / - provides the default page
  - The common header for each view should have a status line that shows the number of declarations currently stored and the timestamp (in the local timezone) of when the declarations were last loaded from the bucket. The timestamp should include the timezone and be formatted like this: `8:44 PM MST` (No date or seconds)

- **Line:** Each line is a declaration.
- **Comment:** If the line starts with a hash, then the rest of the line is a comment and is not processed except when sorting.
- **Labels:** If the line starts with a word (no spaces) surrounded by colons, then whatever is inside the colons is a label. After a label is found, there may be additional labels also separated by colons. Labels are optional.  If the line does not start with a colon, there are no labels. Examples:
  - `:Promise: I *will* look upon the goodness of the LORD in the land of the living. - Psalm 27:13` - one label
  - `:Promise:Blessing: I *will* look upon the goodness of the LORD in the land of the living - Psalm 27:13` - two labels
  - `I *will* look upon the goodness of the LORD in the land of the living. - Psalm 27:13` - No Label
- **Declaration:** The rest of the line is the declaration with a bible reference at the end.  Declarations can end with a period `.`, an exclamation point `!`, a ending square bracket `]`, an ending parenthesis `)`, or none of those.
- **Reference:** There are two possible types of references at the end of the line. Bible references and person references.  References are separated from the declaration text by a space-dash-space ` - ` followed by one or more references.
  - Bible reference examples:
    - `Gal 4:7` - single verse reference
    - `Gal 4:7,8` - single reference with multiple verses
    - `1 John 5:14-15` - single reference with multiple verses
    - `Gal 4:7,8; 1 John 5:14-15` - multiple references separated by semicolon
    - `1 Peter 2:9, Eph 2:6` - multiple references separated by comma
  - Person reference examples:
    - `Smith Wigglesworth` - name of a person
    - `Joe J. Smith` - name of a person
    - `Jim John Doe` - name of a person

### Declarations file example

```text
I am blessed in everything I undertake - Deut 28:8
:Promise: I *will* look upon the goodness of the LORD in the land of the living. - Psalm 27:13
:Promise:Blessing: I *will* look upon the goodness of the LORD in the land of the living - Psalm 27:13
```

The curent list of declarations is in the file declarations.txt in the root directory of the project.

## Bible API

- The Bible API is at <https://api.esv.org>
- The token header is like this: `Authorization: Token {token}`
- The URL for finding the text of a bible reference is: <https://api.esv.org/v3/passage/text/?q={verseReference}>.  This returns the Bible text for the given verse reference.
- When running locally, the API token is in the file .esv-api-token in the root directory of the project.  It must be left in that file and never added to the Golang code.  If it does not exist, the application will write a message to standard error and not start.
- When running in Google Cloud Run, the API token is in a secret in Google Secret Manager.  The secret is called `esv-api-token`

## The source code

- The UI code is placed in it's own Golang file called ui.go.
- The API code is placed in it's own Golang file called api.go.
- The main.go file is the entry point for the application.  It initializes the API and UI and starts the HTTP server.

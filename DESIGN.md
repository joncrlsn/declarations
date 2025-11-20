# Declarations Design

There are two components to this application, both running in the same Golang process.
- The first component is an API using the standard Golang HTTP framework that stores and edits text-based declarations of who I am in Christ. 
- The second component is a web client, also using Golang HTTP for interacting with the API
- Both components run in the same Golang process and listen on port 8080

## API

The API is a JSON REST API that provides the following endpoints: 
- GET /api/v1/declarations - provides a JSON list of declarations
- POST /api/v1/declarations - creates a new declaration
- PUT /api/v1/declarations/{id} - updates an existing declaration
- DELETE /api/v1/declarations/{id} - deletes an existing declaration
- GET /api/v1/declarations/random - provides a random declaration on it's own page
- GET /api/v1/health - this API response includes the number of declarations currently stored

## UI

The web client provides a user interface for interacting with the API. It provides the following features:
- A list of declarations in a compact table format.  The declaration id is not shown to the user.
- A form for creating new declarations
- A form for editing existing declarations
- A form for deleting declarations
- A page for viewing a random declaration from the API.  This is the default page.  When this page is first loaded, a random declaration is pre-selected so the user does not have to click a button to view a random declaration.  As new declarations are requested with a button click, they are appended to the bottom of the page.
- Only dark mode is supported.  In dark mode, the background is dark and the font color is light.
- Pages that look good on either a desktop or mobile device.


## Storage

Each declaration is a text-based description of who the Bible says I am or what I have been given in Christ. Each line of the file is a declaration.  It is stored in a flat file and can be edited.  The file is sorted by a bible reference that occurs at the end of each line, where books are sorted in the order of their appearance in the Bible.  The file needs to stay sorted by bible reference.  The code is written to read and write the file in the same format as the current declarations.txt file.

The current declarations file has very little formatting.  If the line starts with a word surrounded by colons, then whatever is inside the colons is a label. After a label is found, there may be additional labels separated by colons. Labels are optional.  If the line does not start with a colon, there are no labels.  The rest of the line is the declaration with a bible reference at the end.  The bible reference at the end of the line is denoted by a space, then a dash and a book of the bible with a chapter and verse separated by a colon. i.e. ` - Gal 4:7`.

File Example:

```
I am blessed in everything I undertake - Deut 28:8
:Promise: I *will* look upon the goodness of the LORD in the land of the living. - Psalm 27:13
:Promise:Blessing: I *will* look upon the goodness of the LORD in the land of the living - Psalm 27:13
```

The curent list of declarations is in the file declarations.txt in the root directory of the project.

## The code

- The UI code is placed in it's own Golang file called ui.go.
- The API code is placed in it's own Golang file called api.go.
- The main.go file is the entry point for the application.  It initializes the API and UI and starts the HTTP server.
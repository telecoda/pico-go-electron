# About the docs

The pico-go in app documentation is powered by hugo.

To maintain the in app docs, edit the markdown files here in the [docs/content](./docs/content) folder.

## Generating the docs

To generate new docs type:

    go generate ./...

In the root dir of the project

## Editing the docs

You can use `hugo` to edit the docs dynamically whilst developing content.

In the `docs` dir type:

    hugo server

The docs will be temporarily hosted on http://localhost:1313 and automatically refreshed as you edit them

## oddities

Due to the docs being hosted inside and Electron app the docs need to be compatible with NodeJS.  There are not being hosted inside a webserver but being fetched via file:// urls.  

Therefore all urls need to be fairly explicit.  For this reasons you cannot have any folder level links in the docs such as http://mysite/my-content they MUST include a full filepath. eg.  http://mysite/my-content/index.html

Otherwise user will click on a link like this and end up in a black screen dead end.

The `uglyurls = true` setting in `config.toml` helps to resolve this issue.

Also templates have been edited to force home links to `./index.html` instead of just `/`.

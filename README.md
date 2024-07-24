# OCIDex Api

A simple golang application to get OCI image information from a container image registry.

Generally manifest files are available through the api itself; however, some information, like the image config, 
are just stored in blobs and not easily handled by a front end application to display. For now, this only handles the
manifest and config, with hopes to extend this to image signatures, attestations, etc.

# API Docs

API docs are generated with [swag](https://github.com/swaggo/swag)
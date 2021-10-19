## Problem

Goal: Create an HTTP proxy in Go which grayscales images.

The origin server should be configurable at start. Example: `https://maps.wikimedia.org/`
The proxy should accept requests for any path and attempt to request that same path from the origin. Example: Requesting a path of `/osm-intl/1/0/0.png` will cause the proxy to in turn make a request to `https://maps.wikimedia.org/osm-intl/1/0/0.png`.
Assume that most requests will be for PNG images but the server should also support JPEG.
Assume that most of the response images will be paletted but non-paletted representations should also be supported.
The request should fail if the origin server does not respond within 5 seconds.
The image should be returned in the same format as the original image.
External resources can be used to complete this, similar to how you might solve a problem in the real world. Though direct copy/pasting of external code is not within the spirit of the assessment.

# A Go wrapper around the OpenCL 1.1 API

This somewhat higher than a raw wrapper attempts to take some of the tediousness
out of using OpenCL. A couple of examples are when querying platforms it
automatically fills in all platform and device information, and maps buffers and
images to Go types.

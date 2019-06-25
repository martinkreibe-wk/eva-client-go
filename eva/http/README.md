# HTTP Binding Implementation

This package is the implementation for communicating with the eva-client-service over http.

<!-- toc -->

- [Usage](#usage)

<!-- tocstop -->

To use this package include it as a blind dependency:

```go
import (
	
	// ... other dependencies.
	
	"github.com/Workiva/eva-client-go/eva"
	_ "github.com/Workiva/eva-client-go/eva/http"
)
```

## Usage

This package is used to facilitate http interaction with eva, and is not intended to be used directly. That said, users
will need to configure the eva package to use this http implementation.

```json
{
    "source": {
        "type":    "http",                // required to enable the http implementation
        "server":  "<server>[?:<port>]",  // required for this package
        "retries": "<tries>[?@<pause>]",  // optional retry logic for connections to the eva client service
        "mime":    <serializer-type>,     // optional way to set the serializer. See the eva package for details.
        
        // optional certificate to call eva with, the eva client service will need to know how to resolve this cert.
        "cert": "-----BEGIN CERTIFICATE-----\n ... cert ...  \n-----END CERTIFICATE-----"
    },

  "category": "<category>"  // required by the eva package
}
```

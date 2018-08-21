# Elastic Search stream reader flogo activity
This activity allows your flogo application to query Elastic Search using a Elastic stream query

WORK IN PROGRSS - DO NOT USE FOR PRODUCTION!!!

## Installation

```bash
flogo install github.com/abasse/flogoelasticstream
```

## Schema
Inputs and Outputs:

```json
  { "inputs":[
        {
          "name": "basicAuthUser",
          "type": "string",
          "required": false
        },
        {
          "name": "basicAuthPassword",
          "type": "string",
          "required": false
        },
        {
          "name": "elasticbaseURL",
          "type": "string",
          "required": true
        },
        {
          "name": "elasticQuery",
          "type": "string",
          "required": true
        }
      ],
      "outputs": [
        {
            "name": "result",
            "type": "any"
        },
        {
            "name": "hits",
            "type": "number"
        }
      ]
  }
```

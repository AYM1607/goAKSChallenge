# goAKSChallenge

## Challenge description

### Requirements

- Build a Golang RESTful API server for application metadata.
- An endpoint to persist application metadata (In memory is fine). The API must support YAML as a valid payload format.
- An endpoint to search application metadata and retrieve a list that matches the query parameters.
- Include tests if you feel it’s appropriate.

We've provided example yaml data payloads. Two that should persist, and two that should error due to missing fields.

### "Rules"

Use golang for the server, but any other software or open source libraries are fair game to help you solve this problem. The response from the server as well as the structure of the query endpoint is intentionally vague to allow latitude in your solution.

### Advice

This exercise is an opportunity to show off your passion and the craftsmanship of your solution. Optimize your solution for quality and reliability. If you feel your solution is missing a cool feature and you have time, have fun and add it. Make the solution your own, and show off your skills.

### What about the database?

It's recommended that you don't use a database. Integrating with a database driver or ORM gives you less room to shine, and us less ability to evaluate your work.

### Example payloads

All fields in the payload are required. For illustration purposes, we have a few example payloads. One example payload where the maintainer email is not a valid email and another where the version is missing that should fail on submit and two that should be valid.

#### Invalid Payloads

```yaml
title: App w/ Invalid maintainer email
version: 1.0.1
maintainers:
  - name: Firstname Lastname
    email: apptwohotmail.com
company: Upbound Inc.
website: https://upbound.io
source: https://github.com/upbound/repo
license: Apache-2.0
description: |
  ### blob of markdown
  More markdown
```

```yaml
title: App w/ missing version
maintainers:
  - name: first last
    email: email@hotmail.com
  - name: first last
    email: email@gmail.com
company: Company Inc.
website: https://website.com
source: https://github.com/company/repo
license: Apache-2.0
description: |
  ### blob of markdown
  More markdown
```

#### Valid Payloads

```yaml
title: Valid App 1
version: 0.0.1
maintainers:
  - name: firstmaintainer app1
    email: firstmaintainer@hotmail.com
  - name: secondmaintainer app1
    email: secondmaintainer@gmail.com
company: Random Inc.
website: https://website.com
source: https://github.com/random/repo
license: Apache-2.0
description: |
  ### Interesting Title
  Some application content, and description
```

```yaml
title: Valid App 2
version: 1.0.1
maintainers:
  - name: AppTwo Maintainer
    email: apptwo@hotmail.com
company: Upbound Inc.
website: https://upbound.io
source: https://github.com/upbound/repo
license: Apache-2.0
description: |
  ### Why app 2 is the best
  Because it simply is...
```

## Solution

### Libraries used

- [go-yaml](https://github.com/goccy/go-yaml) to parse the yaml payloads.
- [gorilla mux](https://github.com/gorilla/mux) as the http router.
- [bleve](https://github.com/blevesearch/bleve) as a full text search index implemented in go (only used to search through the description).
- [ulid](https://github.com/oklog/ulid) for internal indexing purposes.

### Implementation details

#### Assumptions

Since the description was intentionally vague, I took the time to define some constraints that would allow for a clearer direction while implementing the code and also while testing. The following assumptions are true for this solution:
- All the requests are performed using JSON. The yaml that represents the metadata is always transferred encoded as a string.
- All the fields are searchable individually and join queries can be used to be more granular.
- All fields are indexed as an exact match, meaning that the exact field value must be provided for its corresponding record to be returned. The exception is the description field which is indexed for full text search and thus the queries can be more flexible.
- Most of the internal apis have to be tested with unit tests.
- Due to the amount of search combinations and request validations, integration teststing is helpful and thus a set of them is included.

#### Request formats

A post request to the /records endpoint is required to add a metadata record to the system. The required schema is the following:
```json
{
  "record": "<a yaml document encoded as a string>" 
}
```


A post request to the /records/search endpoint is required to query the existing records. The required schema is the following:
```json
{
  "joinMethod": "or",
  "searchTerms": [
    {
      "field": "title",
      "query": "Valid App 1"
    }
  ]
}
```

The options for join method are:
- or
- and

The options for fields are the following:
- title
- version
- maintainerEmail
- maintainerName
- company
- website
- source
- license
- description

As mentioned previously, only description supports full text search but can be combined with "or" or "and" joins with other search terms.

#### Architecture

All of the fields are indexed separately. An internal index interface has implementations for both exact match and fts indexing:
- The exact match implementation uses a simple golang map.
- The fts implementation uses the bleve library, which is overkill for this purpose but I really wanted to have fts at least for the description field.

When a request to add a new record is received, the server populates all of the indexes with the record data and fails if any of the fields are not indexed successfully.
When a request to search for records is received, all the indexes are queried concurrently and the results are merged afterwards depending on the join method. If a single index query fails, it doesn't fail the whole request.
The requests are protected by a RW lock thus, multiple concurrent reads are performant but we're still protected against race conditions.

#### Testing

Unit tests are provided for the functions that I though were more error prone but had I had more time I would've definitely increased the coverage.
A set of integration tests is provided which goes through the most common scenarions (both fails and successes)

There's a make file that can be used to run the whole test suite or to run the server and test manually with a tool like postman.
Valid make commands:

- `make test` run all the tests
- `make serve` start the server on port :8888

#### Areas of improvement

This list of things were not included due to lack of time but would be nice to have:
- More efficient indexing: there are some parts of the algorith that are linear in time complexity and could cause problems if the queries get too large.
- More robust concurrent search: The current implementation uses a single unbuffered channel that could be causing a bottleneck. I didn't profile the code but If I were to do so I could put together a different and more efficient concurrent appraoch.
- More robust testing for invalid inputs: While I do tests for missing and invalid inputs in the yaml documents, I would like for the errors to be more robuts. I do return a string that ends up in the API response and indicates what fields are missing/invalid, this is not easily testable. A custom error that contains the fields would allow better tests.
- Interact with env vars for server configuration.


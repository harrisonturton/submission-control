# Continuous Integration Service
This service is what powers the automated testing of student projects.

It is an entirely separate service from the Submission Control backend. You send requests, and it will
run student code, and return the test results.

## API

### Test Requests
The test spec is sent in a JSON POST request. I've written the example in YAML, because it's easier to read,
and these configs should be written in it. The equivalent JSON is at the bottom of this README.

```
version: 1
environment:
  image: haskell
  vars:
    user: harrisonturton
    CC: gcc4.8
defualts:
  failOnError: false
tests:
  - name: Check correct files exist
    message: Missing the required files. Please check the deliverables page.
    command: ls
    expect:
      containsAll:
        - Student_Packages
        - report.pdf
```


## Testing

Run `make test` to run all unit and integration tests.

Run `make test-unit` for unit tests only, and `make test-integration` for integration tests.

The integration testing environment is specified in `docker-compose.test.yml`.

## Running

The entire service is run as a set of containers, as defined in `docker-compose.yml`. It has multiple executables that
need to be managed, so this makes it really easy.


```
{
    // Version of the configuration file. Defaults to the latest, use "1" if unsure
    "version": "1",
    
    // Configuration of the testing environment inside the container
    "environment": {
        // The base image to run the code in. Usually one image per assignment.
        "image": "haskell",
        
        // Environment variables to set inside the container
        "vars": {
            "user": "harrisonturton",
            "CC": "gcc4.8,
        }
    },
    
    // Default values of each test case
    "defaults": {
        failOnError: false,
    }
    
    // The different test cases
    "tests": [
        {
            "name": "Check the right files exist",
            "command": "ls",
            "expect": {
                "containsAll": [
                    "Student_Packages/",
                    "report.pdf"
                ]
            }
        }
    ]
}
```

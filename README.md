# HomeVision Backend Take Home Exercise

This is a take home interview for HomeVision that focuses primarily on writing clean code that accomplishes a very
practical task. We have a simple paginated API that returns a list of houses along with some metadata. Your
challenge is to write a script that meets the requirements.

**Note:** this is a *flaky* API! That means that it will likely fail with a non-200 response code. Your code *must*
handle these errors correctly so that all photos are downloaded.

## API Endpoint

You can request the data using the following endpoint:

```bash
http://app-homevision-staging.herokuapp.com/api_project/houses
```

This route by itself will respond with a default list of houses (or a server error!). You can use the following URL parameters:

- `page`: the page number you want to retrieve (default is *1*)
- `per_page`: the number of houses per page (default is *10*)

## Requirements

- Requests the first 10 pages of results from the API
- Parses the JSON returned by the API
- Downloads the photo for each house and saves it in a file with the name formatted as:

  `[id]-[address].[ext]`

- Downloading photos is slow so please optimize them and make use of concurrency

## Bonus Points

- Write tests
- Write your code in a strongly typed language
- Structure your code as if you were planning to evolve it to production quality

## Managing your time

Include “TODO:” comments if there are things that you might want to address but that would take too much time for this exercise. That will help us understand items that you are considering but aren’t putting into this implementation. We can talk about what the improvements might look like during the interview that would get the code to final production quality.

## Submitting

- Create a  GitHub repo with clear readme instructions for running your code on MacOS or Linux
- Send us a zip of the files, or a link to the public repo containing your submission (latter is preferred).

## Usage
- Download dependencies
```bash
make deps
``` 
- Build executable
```bash
make build
``` 
- Run application
```bash
make run
``` 
- Run unit tests
```bash
make test
``` 
- Remove application and saved photos
```bash
make clean
``` 

## Configuration
The following configuration variables are used by the application and are stored in the app.env file:
- **HOUSE_SERVICE_URL**: This contains the URL of the House Service that is used to retrieve the house data
- **NUM_PAGES**: This is the number of pages of house data to request from the House Service
- **NUM_PER_PAGE**: This is the number of houses per page of data requested from the House Service

**NOTE**: These variables can be modified in the app.env file or overridden using environment variables with the same
names.

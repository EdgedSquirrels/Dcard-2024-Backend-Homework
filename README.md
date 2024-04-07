# Dcard 2024 Intern Backend Homework
## Spec Description
[Decard 2024 Backend Homework Spec](https://drive.google.com/file/d/1dnDiBDen7FrzOAJdKZMDJg479IC77_zT/view)

## Prerequisite
* Go (1.22.1)
* Docker

## Quick Start
* Set up postgreSQL Docker image
```bash
$ docker run \
--rm --name postgres \
-e POSTGRES_USER=postgres \
-e POSTGRES_PASSWORD=postgres \
-e POSTGRES_DB=ad \
-p 5432:5432 \
-d postgres:latest
```

* Run the service:
```bash
$ go run main.go
```

* Testing
```bash
$ go test
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

PASS
ok      dcard2024       2.553s
```
## API Usages
* Post new advertisement
```bash
$ curl -X POST -H "Content-Type: application/json" \
  "http://localhost:8080/api/v1/ad" \
  --data '{
     "title": "AD 55",
     "startAt": "2023-12-10T03:00:00.000Z",
     "endAt": "2023-12-31T16:00:00.000Z", 
     "conditions": {
        "ageStart": 20,
        "country": ["TW", "JP"],
        "platform": ["android", "ios"]
     }
  }'
# Response (information of the added advertisement):
{"Title":"AD 55","StartAt":"2023-12-10T03:00:00Z","EndAt":"2023-12-31T16:00:00Z","Conditions":{"AgeStart":20,"AgeEnd":null,"Gender":null,"Country":["TW","JP"],"Platform":["android","ios"]}}
```

* Get advertisements
```bash
$ curl -X GET -H "Content-Type: application/json" \
"http://localhost:8080/api/v1/ad?offset=0&limit=2&age=24&gender=F&country=TW&platform=ios"

# Response:
{"items":[{"title":"AD 1","endAt":"2023-12-22T01:00:00.000Z"},{"title":"AD 31","endAt":"2023-12-30T12:00:00.000Z"}]}
```

## File Structure
```
├── go.mod
├── go.sum
├── internal
│   ├── get_ads
│   │   └── get_ads.go (GET endpoint implementation)
│   └── post_ads
│       └── post_ads.go (POST endpoint implementation)
├── main.go
├── main_test.go
├── README.md
└── test
    ├── (json files of unit test cases)
```

## Idea & Design Choice
### API Format
Generally, I applied the same API format as the API example in the spec. Most of the parameters are verified via `binding`. In addition, I also checked that `AgeStart > AgeEnd` should be regarded as an invalid request.
```go
type Advertisement struct {
	Title      string    `binding:"required"`
	StartAt    time.Time `binding:"required"`
	EndAt      time.Time `binding:"required"`
	Conditions struct {
		AgeStart null.Int
		AgeEnd   null.Int
		Gender   []string `binding:"dive,oneof= M F"`
		Country  []string `binding:"dive,iso3166_1_alpha2"`
		Platform []string `binding:"dive,oneof= android ios web"`
	} `binding:"required"`
}
```
### Database
Since the advertisement format is determined. I did not chose a NoSQL database for better performance. Since I have used MySQL before, I decided to use `PostgreSQL` as a trial. Due to the tedious database configuration, I used docker image for the database.

### Testing
* Correctness: It contains some test cases to verify the correctness of the AD Assigner. Those includes query formats and AD responses.
* Performance: I tried to send 1000 requests at the end of the test, but it took about 2.5 seconds. It showed that it can only handles roughly 400 requests per second. This homework aims to reach 10,000 requests per second. Perhaps cache storage is needed for this homework to reach the target performance.

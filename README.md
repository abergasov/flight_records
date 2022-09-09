## track flight records

### Local run
```bash
make run
# or run on custom port
make port=8081 run
```

endpoint for calculate source and destination in list of flights
```go
POST /api/v1/track
```
example request
```shell
curl --request POST \
  --url http://127.0.0.1:8080/api/v1/track \
  --header 'Content-Type: application/json' \
  --data '[["IND", "EWR"], ["SFO", "ATL"], ["GSO", "IND"], ["ATL", "GSO"]]'
```
sample response:
```json
["SFO","EWR"]
```

## Comments

### 1. Airport duplicates
Current task description does not give answer - is it possible that the user visited the same airport 2 or more times.

As we assume that it is related to the real world, so the answer is `yes`, the user can visit airport 2 times.

So i make 2 realizations:
```
internal/service/tracker/tracker.go:36 - user route without duplicates
internal/service/tracker/tracker.go:71 - user route with possible duplicates
```

Nevertheless, the 2nd realization is not so fast. You can look at the difference in benchmarks.  
```shell
make bench

# my laptop results
cpu: Intel(R) Core(TM) i7-10510U CPU @ 1.80GHz
# realization without duplicates
BenchmarkRouteTracker_CalculateRouteOverMaps-8                   3617112          322.8 ns/op
# realization with duplicates
BenchmarkRouteTracker_CalculateRouteWithDuplicates-8              596617          1970 ns/op
# realization without duplicates for route with 2 points
BenchmarkRouteTracker_CalculateRouteOverMaps_2Values-8           5407621          218.1 ns/op
# direct realization `if - else if - else` for 2 and 3 points routes
BenchmarkRouteTracker_CalculateRouterFor2Pairs-8                70915405          16.02 ns/op
BenchmarkRouteTracker_CalculateRouterFor3Pairs-8                65531049          16.79 ns/op
```
If we write code for real usage in production, i would like to suggest to do some research and try to get answers for several questions:
1. how often does user visit same airport more than 1 time in his route.
2. how long is average route for user, how many points in average route.

In my opinion - scenario, when user visit same airport 2 time not very often case. 
Based on this assumption i make my realization:
1. we try to calculate route, like there are no duplicates (using `CalculateRouteOverMaps`)
2. if build route is fail - than try to calculate route with method #2 (using `CalculateRouteWithDuplicates`)

for most cases - it will work quite fast. But for several cases - not so good. 

#### Possible improvements
Try to guess if duplicates exist before parse route. 
##### Codogeneration
I make some methods, which use very naive approach - `if - else if - else` for 2 and 3 points. 
Benchmark for these methods shows best results.

If it is possible to collect statistic - how long is average route for user, then other improvements are possible.

For example, 90% of users fly using 3-4 point. 
So we can generate code for these cases - make all possible pairs using combinatorics and build route. 
For 3-4 points it will be 6-24 combinations.

#### 2. Dual solutions
Some cases can generate more than 1 correct answer. For example:
```
# user fly in circle
[a, b]
[b, c]
[c, a]
```

for this case we can have several correct answers:
```
a -> b, b -> c, c -> a
b -> c, c -> a, a -> b
c -> a, a -> b, b -> c
```
As we do not know start point, we can not say which one is correct. In this case i return first correct answer. 
But probably, we can return all correct answers, let service consumer decide.

Or make additional endpoint, which also accept start point and based on this info service can decide, which route is correct.

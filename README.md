Compare the performance of [distroless](https://github.com/GoogleContainerTools/distroless) with alpine image.

# Result
The peformance is not so different.
```
> docker build -f Dockerfile-alpine -t distrolesstrial:alpine . && docker run --rm -i -t distrolesstrial:alpine
Allocs: 0
time taken: 561464600

> docker build -f Dockerfile-distroless -t distrolesstrial:distroless . && docker run --rm -i -t distrolesstrial:distroless
Allocs: 0
time taken: 58299090
```



# Appendix
- I used the script from [this article](https://chris124567.github.io/2021-06-21-go-performance/)

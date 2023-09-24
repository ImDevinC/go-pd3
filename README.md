# pd3
This is a POC of how to call the challenges API in Payday 3. You will need to inspect your game traffic to get the bearer token, which is outside the scope of this project.

Once you have the token, you can either set an environment variable named `NEBULA_BEARER_TOKEN` or use the `-token` flag when calling `go run main.go`

The results are printed to the screen, or you can use the `-outputFile` flag to save the results to a file in JSON format.

This project is very much a POC and has no real-world applications at this time.
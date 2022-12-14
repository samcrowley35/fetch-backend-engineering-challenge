# fetch-backend-engineering-challenge

To run my version of the solution to this challenge, follow these steps:

First, ```go mod init gfg``` to create the module for this project

Second, ```go get github.com/google/uuid``` , I used this package to generate uuid's for the receipts

Third, to build the app with docker run this command: ```docker build -t my-go-app .```

Fourth, ```docker run -p 8080:8081 -it my-go-app```

Fifth, navigate to http://localhost:8080/receipts/process in your browser

To get the ID for the receipt, simply enter in the json's full name including ".json" at the end, if adding a new one, be sure to place it in the examples folder 

To get the total number of points for the receipt, navigate to http://localhost:8080/receipts/points/ followed by the id previously displayed.

To see another receipt's id and points, navigate back to http://localhost:8080/receipts/process and start the process over.

Notes:
- Some of the solutions may be off by one due to the math.Round() function not behaving the way I wanted it too, for example, the result of target-receipt.json will be 24 
instead of 25 because 2.45 doesn't round up to 3 with math.Round()

- Also as outlined in the challenge description, data for these receipts does not need to persist, so if you use an old id, this program won't work

I appreciate you all taking the time to evaluate my solution, if I messed something up you can reach out to me at samcrowley35@gmail.com




<samp>

# Go_Mongo_Template
<h4>This is a basic Go Mongo template inspired from <a href="https://dev.to/hackmamba/build-a-rest-api-with-golang-and-mongodb-gin-gonic-version-269m"><img align="center" height="30" src="https://img.shields.io/badge/dev.to-0A0A0A?style=for-the-badge&logo=devdotto&logoColor=white">article</a></h4>

## prerequisites

Add `.env` in the root folder of your project. 
add a variable of `MONGOURI`

<ins>Example:</ins> MONGOURI=mongodb+srv://<YOUR USERNAME HERE>:<YOUR PASSWORD HERE>@cluster0.e5akf.mongodb.net/myFirstDatabese?retryWrites=true&w=majority

## Run the App
to run the project, first open your terminal or command prompt in the root folder. paste
```
go run main.go
```

## issues may occur
With Mongodb atlas there is a chance of getting timedout. So tweak timeout on   `configs/setup.go` file.


<hr>
<h3>

![withLove](https://img.shields.io/badge/Made%20with-Go-1f425f.svg) Feel free to contribute or contact me. </h3>
</samp>

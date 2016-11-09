package main

var html = `
<html>
  <head>
    <meta http-equiv="refresh" content="5">
  </head>
  <body>
    <h1>Here we go!</h1>
    <iframe src="http://localhost:8880"></iframe>
    <iframe src="http://localhost:8881"></iframe>
    <iframe src="http://localhost:8882"></iframe>
    <iframe src="http://localhost:8883"></iframe>
    <iframe src="http://localhost:8884"></iframe>
    <iframe src="http://localhost:8885"></iframe>
  </body>
	<hr />
	Refreshed at %s
</html>
`

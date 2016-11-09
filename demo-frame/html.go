package main

var html = `
<html>
  <head>
    <meta http-equiv="refresh" content="5">
  </head>
  <body>
    <h1>Here we go!</h1>
    <iframe src="http://localhost:8880" height=275 width=275></iframe>
    <iframe src="http://localhost:8881" height=275 width=275></iframe>
    <iframe src="http://localhost:8882" height=275 width=275></iframe>
    <iframe src="http://localhost:8883" height=275 width=275></iframe>
    <iframe src="http://localhost:8884" height=275 width=275></iframe>
    <iframe src="http://localhost:8885" height=275 width=275></iframe>
  </body>
	<hr />
	Refreshed at %s
</html>
`

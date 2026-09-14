package main

import (
	"fmt"
	"net/http"
	"nightveil-demon/api"
	"nightveil-demon/media"
	"nightveil-demon/middleware"
	"os"
	"os/exec"
	"time"
)

const (
	appName = "NightVeil Demon"
	version = "v0.1.0"
	port    = 5226
)

var server *http.Server

func uiHandler(w http.ResponseWriter, r *http.Request) {
	html := fmt.Sprintf(`
<!doctype html>
<html>
<head>
<title>%s</title>
<meta charset="utf-8">
<style>
body{
	background:#0f0f14;
	color:white;
	font-family:Arial;
	display:flex;
	height:100vh;
	justify-content:center;
	align-items:center;
}
.box{
	background:#1c1c25;
	padding:30px;
	border-radius:12px;
	width:400px;
	text-align:center;
}
button{
	margin-top:15px;
	padding:10px 15px;
	border:none;
	border-radius:8px;
	cursor:pointer;
	color:white;
}
.test{
	background:#4f46e5;
}
.quit{
	background:#ef4444;
}
.status{
	margin-top:15px;
	color:#22c55e;
}
</style>
</head>
<body>
<div class="box">
<h1>🌙 %s</h1>
<p>Version: %s</p>
<p class="status" id="status">
Checking FFmpeg...
</p>
<button class="test" onclick="testFFmpeg()">
Test FFmpeg
</button>
<button class="quit" onclick="shutdown()">
⛔ Shut down
</button>
<pre id="output"></pre>
</div>
<script>
async function testFFmpeg(){
	const res=await fetch('/test-ffmpeg');
	const data=await res.text();
	document.getElementById('output').innerText=data;
}
async function shutdown(){
	document.getElementById('status').innerText="Shutting down...";
	await fetch('/quit');
	setTimeout(()=>{
		document.body.innerHTML="<h2 style='color:white'>Server stopped</h2>";
	},500);
}
fetch('/health')
.then(r=>r.text())
.then(t=>{
	document.getElementById('status').innerText=t;
});
</script>
</body>
</html>
`, appName, appName, version)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}
func healthHandler(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command(
		media.FFmpegPath(),
		"-version",
	)
	if err := cmd.Run(); err != nil {
		w.Write([]byte("❌ FFmpeg NOT working"))
		return
	}
	w.Write([]byte("✅ FFmpeg ready"))
}
func testFFmpegHandler(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command(
		media.FFmpegPath(),
		"-f",
		"lavfi",
		"-i",
		"testsrc=duration=1:size=1280x720:rate=30",
		"-f",
		"null",
		"-",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		w.Write(out)
		return
	}
	w.Write([]byte("FFmpeg OK"))
}
func quitHandler(w http.ResponseWriter, r *http.Request) {
	go func() {
		time.Sleep(
			200 * time.Millisecond,
		)
		if server != nil {
			server.Close()
		}
		os.Exit(0)
	}()
	w.Write([]byte("bye"))
}
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(
		"/",
		uiHandler,
	)
	mux.HandleFunc(
		"/health",
		healthHandler,
	)
	mux.HandleFunc(
		"/test-ffmpeg",
		testFFmpegHandler,
	)
	mux.HandleFunc(
		"/render",
		renderHandler,
	)
	mux.HandleFunc(
		"/quit",
		quitHandler,
	)
	mux.HandleFunc(
		"/ws",
		wsHandler,
	)
	api.RegisterAPIRoutes(mux)
	server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: middleware.CorsMiddleware(mux),
	}
	fmt.Printf("%s running on :%d\n", appName, port)
	err := LoadIdentity()
	if err != nil {
		fmt.Println(err)
		return
	}
	InitDemonState()
	StartMonitor()
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		fmt.Println(err)
	}
}

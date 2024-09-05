package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
)

var bffURL = "http://bff.bff-server.svc.cluster.local:8080" // BFF 서비스의 도메인 이름 및 포트

func main() {
	// HTTP 핸들러 등록
	http.HandleFunc("/", renderHome)
	http.HandleFunc("/test1", handleTest1)
	http.HandleFunc("/test2", handleTest2)

	// 서버 시작
	log.Println("SSR Frontend Server started at :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

// 홈 페이지 렌더링 함수
func renderHome(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("index").Parse(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<title>SSR Frontend</title>
		</head>
		<body>
			<h1>SSR Frontend</h1>
			<p>Click below to test BFF communication:</p>
			<a href="/test1">Test 1 (BFF Test1)</a><br>
			<a href="/test2">Test 2 (BFF Test2)</a>
		</body>
		</html>
	`)
	if err != nil {
		http.Error(w, "Failed to parse template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// /test1 핸들러 - BFF와 통신하여 Test1 요청
func handleTest1(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(fmt.Sprintf("%s/api/v1/bff/test1", bffURL))
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to communicate with BFF", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	// 결과를 HTML로 렌더링
	tmpl, err := template.New("result").Parse(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<title>SSR Frontend - Test1</title>
		</head>
		<body>
			<h1>Test 1 Result</h1>
			<p>{{.}}</p>
			<a href="/">Go back</a>
		</body>
		</html>
	`)
	if err != nil {
		http.Error(w, "Failed to parse template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, string(body))
}

// /test2 핸들러 - BFF와 통신하여 Test2 요청
func handleTest2(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(fmt.Sprintf("%s/api/v1/bff/test2", bffURL))
	if err != nil {
		http.Error(w, "Failed to communicate with BFF", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	// 결과를 HTML로 렌더링
	tmpl, err := template.New("result").Parse(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<title>SSR Frontend - Test2</title>
		</head>
		<body>
			<h1>Test 2 Result</h1>
			<p>{{.}}</p>
			<a href="/">Go back</a>
		</body>
		</html>
	`)
	if err != nil {
		http.Error(w, "Failed to parse template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, string(body))
}

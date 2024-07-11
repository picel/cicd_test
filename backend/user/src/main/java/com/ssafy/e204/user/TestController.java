package com.ssafy.e204.user;

import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class TestController {

  public TestController() {
    System.out.println("TestController()");
  }

  @GetMapping
  public ResponseEntity<String> test() {
    System.out.println("asdf");
    return ResponseEntity.ok("Hello, World!");
  }

}

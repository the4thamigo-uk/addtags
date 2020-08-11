# addtags
_Add arbitrary struct field tags to golang code_

## Introduction

This code-generation tool allows you to add struct field tags to golang structs that are defined in a configuration file.

To see how it works run the following command :

```
$ go run ./main.go -t ./examples/gotags.yml -d ./examples | git diff ./examples/
diff --git a/examples/examples.go b/examples/examples.go
index 759fee1..5ca5445 100644
--- a/examples/examples.go
+++ b/examples/examples.go
@@ -1,10 +1,10 @@
 package examples
 
 type MyStruct1 struct {
-       X int
-       Y int `existing:"tag"`
+       X int `new:"tag1"`
+       Y int `existing:"tag" new:"tag2"`
 }
 
 type MyStruct2 struct {
-       Z int
+       Z int `new:"tag3" othernew:"tag4"`
 }
```

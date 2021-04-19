# The Monkey programming language

An interpreter for the Monkey programming language.

# Run
```
$ go run main.go
Hello michael.cuffaro! This is the Monkey programming language!
Feel free to type in commands
>> let a = 5;
>> a
5
>> let add = fn(x, y) { x + y }
>> add(a, 3)
8
>> 5 > 1
true
>> let hello = "Hello, world!"
>> hello
Hello, world!
>> len(hello)
13
```

# TODO

- [ ] Documentation
- [ ] Read and execute *.monkey files
- [ ] More builtin functions

## Credit:
https://interpreterbook.com/

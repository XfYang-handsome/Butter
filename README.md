# Butter
A Go-made programming language that is easy to learn.

Butter has already had some basic tools which can realize some basic algorithms.

You can write a Butter code like this:

```
run main()
  print("Your name:")
  butter n string = readln() //Enter your name.
  welcome(n)
  
  print("Factorial of 15:")
  println(fact(15))

  print("Calculate the sum from 1 to 10:")
  butter i, sum int = 1, 0
  for compareTo(11, i, 0)
    sum = add(sum, i)
    i = add(i, 1)
  /for
  println(sum)

//Calculate the factorial of an integer.
func fact(n int)
  if equalTo(n, 0)
    return 1
  /if
  return mul(n, fact(sub(n, 1)))

//Welcome user.
func welcome(name string)
  println("Hello,", name, "!")
```

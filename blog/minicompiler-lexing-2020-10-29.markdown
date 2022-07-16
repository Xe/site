---
title: "Minicompiler: Lexing"
date: 2020-10-29
series: rust
tags:
 - rust
 - templeos
 - compiler
---

I've always wanted to make my own compiler. Compilers are an integral part of
my day to day job and I use the fruits of them constantly. A while ago while I
was browsing through the TempleOS source code I found
[MiniCompiler.HC][minicompiler] in the `::/Demos/Lectures` folder and I was a
bit blown away. It implements a two phase compiler from simple math expressions
to AMD64 bytecode (complete with bit-banging it to an array that the code later
jumps to) and has a lot to teach about how compilers work. For those of you that
don't have a TempleOS VM handy, here is a video of MiniCompiler.HC in action:

[minicompiler]: https://github.com/Xe/TempleOS/blob/master/Demo/Lectures/MiniCompiler.HC

<video controls width="100%">
    <source src="https://cdn.xeiaso.net/file/christine-static/img/minicompiler/tmp.YDcgaHSb3z.webm"
            type="video/webm">
    <source src="https://cdn.xeiaso.net/file/christine-static/img/minicompiler/tmp.YDcgaHSb3z.mp4"
            type="video/mp4">
    Sorry, your browser doesn't support embedded videos.
</video>

You put in a math expression, the compiler builds it and then spits out a bunch
of assembly and runs it to return the result. In this series we are going to be
creating an implementation of this compiler that targets [WebAssembly][wasm].
This compiler will be written in Rust and will use only the standard library for
everything but the final bytecode compilation and execution phase. There is a
lot going on here, so I expect this to be at least a three part series. The
source code will be in [Xe/minicompiler][Xemincompiler] in case you want to read
it in detail. Follow along and let's learn some Rust on the way!

[wasm]: https://webassembly.org/
[Xemincompiler]: https://github.com/Xe/minicompiler

[Compilers for languages like C are built on top of the fundamentals here, but
they are _much_ more complicated.](conversation://Mara/hacker)

## Description of the Language

This language uses normal infix math expressions on whole numbers. Here are a
few examples:

- `2 + 2`
- `420 * 69`
- `(34 + 23) / 38 - 42`
- `(((34 + 21) / 5) - 12) * 348`

Ideally we should be able to nest the parentheses as deep as we want without any
issues. 

Looking at these values we can notice a few patterns that will make parsing this
a lot easier:

- There seems to be only 4 major parts to this language:
  - numbers
  - math operators
  - open parentheses
  - close parentheses
- All of the math operators act identically and take two arguments
- Each program is one line long and ends at the end of the line

Let's turn this description into Rust code:

## Bringing in Rust

Make a new project called `minicompiler` with a command that looks something
like this:

```console
$ cargo new minicompiler
```

This will create a folder called `minicompiler` and a file called `src/main.rs`.
Open that file in your editor and copy the following into it:

```rust
// src/main.rs

/// Mathematical operations that our compiler can do.
#[derive(Debug, Eq, PartialEq)]
enum Op {
    Mul,
    Div,
    Add,
    Sub,
}

/// All of the possible tokens for the compiler, this limits the compiler
/// to simple math expressions.
#[derive(Debug, Eq, PartialEq)]
enum Token {
    EOF,
    Number(i32),
    Operation(Op),
    LeftParen,
    RightParen,
}
```

[In compilers, "tokens" refer to the individual parts of the language you are
working with. In this case every token represents every possible part of a
program.](conversation://Mara/hacker)

And then let's start a function that can turn a program string into a bunch of
tokens:

```rust
// src/main.rs

fn lex(input: &str) -> Vec<Token> {
    todo!("implement this");
}
```

[Wait, what do you do about bad input such as things that are not math expressions?
Shouldn't this function be able to fail?](conversation://Mara/hmm)

You're right! Let's make a little error type that represents bad input. For
creativity's sake let's call it `BadInput`:

```rust
// src/main.rs

use std::error::Error;
use std::fmt;

/// The error that gets returned on bad input. This only tells the user that it's
/// wrong because debug information is out of scope here. Sorry.
#[derive(Debug, Eq, PartialEq)]
struct BadInput;

// Errors need to be displayable.
impl fmt::Display for BadInput {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "something in your input is bad, good luck")
    }
}

// The default Error implementation will do here.
impl Error for BadInput {}
```

And then let's adjust the type of `lex()` to compensate for this:

```rust
// src/main.rs

fn lex(input: &str) -> Result<Vec<Token>, BadInput> {
    todo!("implement this");
}
```

So now that we have the function type we want, let's start implementing `lex()`
by setting up the result and a loop over the characters in the input string:

```rust
// src/main.rs

fn lex(input: &str) -> Result<Vec<Token>, BadInput> {
    let mut result: Vec<Token> = Vec::new();
    
    for character in input.chars() {
        todo!("implement this");
    }

    Ok(result)
}
```

Looking at the examples from earlier we can start writing some boilerplate to
turn characters into tokens:

```rust
// src/main.rs

// ...

for character in input.chars() {
    match character {
        // Skip whitespace
        ' ' => continue,

        // Ending characters
        ';' | '\n' => {
            result.push(Token::EOF);
            break;
        }

        // Math operations
        '*' => result.push(Token::Operation(Op::Mul)),
        '/' => result.push(Token::Operation(Op::Div)),
        '+' => result.push(Token::Operation(Op::Add)),
        '-' => result.push(Token::Operation(Op::Sub)),

        // Parentheses
        '(' => result.push(Token::LeftParen),
        ')' => result.push(Token::RightParen),

        // Numbers
        '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9' => {
            todo!("implement number parsing")
        }

        // Everything else is bad input
        _ => return Err(BadInput),
    }
}

// ...
```

[Ugh, you're writing `Token::` and `Op::` a lot. Is there a way to simplify
that?](conversation://Mara/hmm)

Yes! enum variants can be shortened to their names with a `use` statement like
this:

```rust
// src/main.rs

// ...

use Op::*;
use Token::*;

match character {
    // ...

    // Math operations
    '*' => result.push(Operation(Mul)),
    '/' => result.push(Operation(Div)),
    '+' => result.push(Operation(Add)),
    '-' => result.push(Operation(Sub)),

    // Parentheses
    '(' => result.push(LeftParen),
    ')' => result.push(RightParen),

    // ...
}
    
// ...
```

Which looks a _lot_ better.

[You can use the `use` statement just about anywhere in your program. However to
keep things flowing nicer, the `use` statement is right next to where it is
needed in these examples.](conversation://Mara/hacker)

Now we can get into the fun that is parsing numbers. When he wrote MiniCompiler,
Terry Davis used an approach that is something like this (spacing added for readability):

```c
case '0'...'9':
  i = 0;
  do {
    i = i * 10 + *src - '0';
    src++;
  } while ('0' <= *src <= '9');
  *num=i;
```

This sets an intermediate variable `i` to 0 and then consumes characters from
the input string as long as they are between `'0'` and `'9'`. As a neat side
effect of the numbers being input in base 10, you can conceptualize `40` as `(4 *
10) + 2`. So it multiplies the old digit by 10 and then adds the new digit to
the resulting number. Our setup doesn't let us get that fancy as easily, however
we can emulate it with a bit of stack manipulation according to these rules:

- If `result` is empty, push this number to result and continue lexing the
  program
- Pop the last item in `result` and save it as `last`
- If `last` is a number, multiply that number by 10 and add the current number
  to it
- Otherwise push the node back into `result` and push the current number to
  `result` as well
  
Translating these rules to Rust, we get this:

```rust
// src/main.rs

// ...

// Numbers
'0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9' => {
    let num: i32 = (character as u8 - '0' as u8) as i32;
    if result.len() == 0 {
        result.push(Number(num));
        continue;
    }

    let last = result.pop().unwrap();

    match last {
        Number(i) => {
            result.push(Number((i * 10) + num));
        }
        _ => {
            result.push(last);
            result.push(Number(num));
        }
    }
}
            
// ...
```

[This is not the most robust number parsing code in the world, however it will
suffice for now. Extra credit if you can identify the edge
cases!](conversation://Mara/hacker)

This should cover the tokens for the language. Let's write some tests to be sure
everything is working the way we think it is!

## Testing

Rust has a [robust testing
framework](https://doc.rust-lang.org/book/ch11-00-testing.html) built into the
standard library. We can use it here to make sure we are generating tokens
correctly. Let's add the following to the bottom of `main.rs`:

```rust
#[cfg(test)] // tells the compiler to only build this code when tests are being run
mod tests {
    use super::{Op::*, Token::*, *};

    // registers the following function as a test function
    #[test]
    fn basic_lexing() {
        assert!(lex("420 + 69").is_ok());
        assert!(lex("tacos are tasty").is_err());

        assert_eq!(
            lex("420 + 69"),
            Ok(vec![Number(420), Operation(Add), Number(69)])
        );
        assert_eq!(
            lex("(30 + 560) / 4"),
            Ok(vec![
                LeftParen,
                Number(30),
                Operation(Add),
                Number(560),
                RightParen,
                Operation(Div),
                Number(4)
            ])
        );
    }
}
```

This test can and probably should be expanded on, but when we run `cargo test`:

```console
$ cargo test
   Compiling minicompiler v0.1.0 (/home/cadey/code/Xe/minicompiler)

    Finished test [unoptimized + debuginfo] target(s) in 0.22s
     Running target/debug/deps/minicompiler-03cad314858b0419

running 1 test
test tests::basic_lexing ... ok

test result: ok. 1 passed; 0 failed; 0 ignored; 0 measured; 0 filtered out
```

And hey presto! We verified that all of the parsing is working correctly. Those
test cases should be sufficient to cover all of the functionality of the
language. 

---

This is it for part 1. We covered a lot today. Next time we are going to run a
validation pass on the program, convert the infix expressions to reverse polish
notation and then also get started on compiling that to WebAssembly. This has
been fun so far and I hope you were able to learn from it.

Special thanks to the following people for reviewing this post:
 - Steven Weeks
 - sirpros
 - Leonora Tindall
 - Chetan Conikee
 - Pablo
 - boopstrap
 - ash2x3

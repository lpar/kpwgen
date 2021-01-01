# kpwgen

## Command line password generator producing more readable/pronounceable passwords

This is a simple random password generator, based on an earlier Ruby script I
used to use. That came about because I found the standard Linux `pwgen` produced
passwords which were very difficult to transcribe manually. This code tries to
generate more readable and pronounceable passwords, so that they are easier to
transcribe -- though they have to be a bit longer to get equivalent security,
see below.

You may want to use this algorithm if you need a password that a user can write 
down or type manually. If you're using a password manager that does the typing 
and remembering for you, this is the wrong tool.

Instead of using random characters for the passwords, this code uses the
[koremutake syllables](http://shorl.com/koremutake.php). I exclude a few which
might cause the system to generate very rude English words, because I use the 
algorithm for automatically generating passwords for user accounts at work.

## Security analysis

_Trigger warning: math._

The filtered list has 213 possible syllables. The number of bits of entropy per
symbol, therefore, is log₂(213), or 7.7347. That compares to 4.7 bits of
entropy for randomly chosen letters of the alphabet with no syllable clustering.

For a desired password entropy H (in bits), choosing at random from N symbols
to generate the password, the number of symbols in your password needs to be L,
where L = H / log2(N) rounded up.

See: <https://en.wikipedia.org/wiki/Password_strength#Random_passwords>

Given that formula, we can assemble a table showing the number of symbols you
need in order to get various levels of password entropy, for both a-z letters
(alphabet symbols) and koremutake symbols. 

There are fewer koremutake symbols needed, of course, because there are 216 of
them rather than 26. However, the symbols are longer -- the average length of a
syllable in characters is 2.35. So the "koremutake length" column lists the number of
_characters_ needed, on average, for a koremutake password of the appropriate strength.

The final column is the bottom line: how many more characters you need in a
koremutake password, rather than an alphabetic one, for the same entropy.

| Entropy | Alphabet symbols | Koremutake symbols | Koremutake length | Increase |
|---------|------------------|--------------------|-------------------|----------|
| 16 | 4 | 3 | 5 | 1.25x |
| 24 | 6 | 4 | 8 | 1.33x |
| 32 | 7 | 5 | 10 | 1.43x |
| 40 | 9 | 6 | 13 | 1.44x |
| 48 | 11 | 7 | 15 | 1.36x |
| 56 | 12 | 8 | 18 | 1.5x |
| 64 | 14 | 9 | 20 | 1.43x |
| 72 | 16 | 10 | 22 | 1.38x |
| 80 | 18 | 11 | 25 | 1.39x |
| 88 | 19 | 12 | 27 | 1.42x |
| 96 | 21 | 13 | 30 | 1.43x |
| 104 | 23 | 14 | 32 | 1.39x |
| 112 | 24 | 15 | 35 | 1.46x |
| 120 | 26 | 16 | 37 | 1.42x |
| 128 | 28 | 17 | 39 | 1.39x |

As you can see, your koremutake passwords need to be longer than your regular
random alphabetic ones by a factor of up to 1.5, mostly around 1.4, to have the
same number of bits of entropy and be as unguessable.

That's not great, but there are some mitigating factors. Firstly, this assumes
that the attacker knows you're using koremutake syllables. Also, we can use
some of the usual password strengthening tricks -- adding a number, a capital
letter or a random punctuation symbol somewhere. This utility offers all three
options. The modifications are carried out after the password is generated.

## The bottom line

As a rule of thumb, multiply your password lengths by 1.4-1.5x when using
this utility. That may seem bad, but which would you rather have to transcribe:
`traobtrotitromen` or `jqfsylgvyp`? `pytrarotit0atro` or `cq4dkpwxdm`?

Personally, I can keep the longer koremutake passwords in my short term memory
for long enough to retype them. I can't do that with the shorter alphabetic
passwords.

## Other things of note

 * The utility has a `--help` command to list options. 
 * The option to add a punctuation character excludes backslash (because
   quoting rules) and single quote (because badly written code often disallows it
   as a half-assed attempt at stopping SQL injection).
 * Yes, the code uses `[crypto/rand](https://golang.org/pkg/crypto/rand/)`.
 * It's pure Go, so you can compile it for Windows, Mac or Linux.


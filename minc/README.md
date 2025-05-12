# minc

The official compiler for Minimal.

## Usage

```
minc <source> <output>
```

Arguments:

- `source`: is the path to the Minimal source file (file extension is `.mn`);
- `output`: is the path to the Minimal bytecode file to be created (file extension is `.mnb`).

### Special values

Some arguments have special values used to refer to the _standard input_, _standard output_ and _standard error_.

- `<source>`: can be replaced with `*stdin` to refer to the _standard input_;
- `<output>`: can be replaced with `*stdout` to refer to the _standard output_, or `*stderr` to refer to the _standard error_.

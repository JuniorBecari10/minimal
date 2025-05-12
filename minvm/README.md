# minvm

The VM for Minimal.

## Usage

```
minvm <bytecode>
```

Argument:

- `bytecode`: is the path to the Minimal bytecode file to be read (file extension is `.mnb`);

### Special values

The argument have some special values used to refer to the _standard input_, _standard output_ and _standard error_.

- `*stdin` to refer to the _standard input_;
- `*stdout` to refer to the _standard output_;
- `*stderr` to refer to the _standard error_.

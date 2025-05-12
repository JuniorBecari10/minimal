# min

The official tool for working with Minimal. <br />
It has several tools for compiling, running and disassembling Minimal code.

## Usage

```
min
    build <source> <output> - compile the specified source file into the specified output file
    disasm <source>         - compile and disassemble the specified source file
    disasmb <bytecode>      - disassemble the specified bytecode file
    execute <bytecode>      - runs the specified bytecode file
    run <source>            - compile and run the specified source file
```

### Special values

Some arguments have special values used to refer to the _standard input_, _standard output_ and _standard error_.

- `<source>` and `<bytecode>`: can be replaced with `*stdin` to refer to the _standard input_;
- `<output>`: can be replaced with `*stdout` to refer to the _standard output_, or `*stderr` to refer to the _standard error_.

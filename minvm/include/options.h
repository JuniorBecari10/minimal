#ifndef OPTIONS_H
#define OPTIONS_H

// Some compiler options for the VM.

// When defined, type checking (at runtime) is added to the VM's code.
// If the compiler's type checker is sufficiently safe to not cause UB at runtime,
// this may be disabled.
#define ENABLE_TYPE_CHECK

#endif

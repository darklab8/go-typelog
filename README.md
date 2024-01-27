Static typed structured logging

# Description

This project araised from the need to log backend applications, aws lambdas and other stuff in modern cloud ecosystem. Logging systems today are able easily parsing JSON format out of the box.
Static typing approach brings here consistent way to define key values to final msg, as well as easier following Domain Driven Design, where logs consistently describe what they log. Static typed logging brings easy refactoring to any present logs.

# Features

- Accepts static typed components as optional params
  - has shortcut WithFields, to make clone of the logger with default logging fields
- Easy to turn on/off parameters by environment variables
  - Ability to define different log levss for different created loggers
- Easier turning complex objects into structured logging
  - accepts maps and structs as its params. It will parse them on their own.

[See folder examples](./examples)

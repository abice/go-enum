CHANGELOG
=========

# unreleased

* Del: Remove extra package `now`.
* Add: Adds `UsageFn` for customizing usage.
* Mod: Replaces `NeedArgs` with `NumArg`.

# v0.0.1 (2016-05-21)

Initializes the repository, and implements many important features.

* Add: Prases flags base on golang tag.
* Add: Supports almost all basic types, slice and map as a flag.
* Add: Fuzzy matching for suggestion.
* Add: Supports command tree.
* Add: Pretty output for usage of command.
* Add: Supports any type which implement `cli.Decoder` as a flag.
* Add: Parser for flag.
* Add: Supports editor like `git commit`.
* Add: Rich examples to help other quick start.

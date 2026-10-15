# Changelog

## v1.2.0 — 2026-10-15

### Added
- `shaktest` package with test helpers: `Require`, `Assert`, `RequireValidation` and `AssertValidation`. They report every failing rule against any `testing.TB`; the `Require` variants then stop the test.

### Changed
- `EqualToAny` and `NotEqualToAny` use `slices.Contains`. No behaviour change.

## v1.1.0 — 2026-03-15

### Changed
- **Breaking:** `RunValidation`, `RunFullValidation` and `SatisfiesValidation` now take a `Validatable` instead of a `validation.Validation`. Pass the struct itself, e.g. `shak.RunValidation(order)` instead of `shak.RunValidation(order.Validation())`.

## v1.0.0 — 2026-03-15

- Initial release: rules for ordered, comparable, string, bool, pointer, time, slice and map values; rule composition with `And`, `Or`, `If`, `Else`, `Ref` and `Deref`; nested struct validation with field paths.

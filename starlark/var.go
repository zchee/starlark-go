package starlark

// A Variable represents an addressable variable.
//
// Addressable variables are used only in Stargo,
// which enables resolve.AllowAddressing.
//
// An addressable variable supports three operations.
// It contains a value, retrieved using the Value method.
// It has an address, obtained using the Address method.
// An address is a value that denotes the identity of the variable; in
// Go, addresses are pointers.
// And its contents may be updated, using the SetValue method.
//
// Ordinary Starlark variables are not addressable, but they contain
// only references. To perform an update such as x[i].f = 1, the
// expression x[i] is evaluated, which yields a reference; then the
// operation "set field .f to 1" is applied to the reference.
// Similarly x.f[i] = 2 is executed by evaluating x.f to obtain a
// reference, then applying the operation "set element i to 2" to it.
// One may introduce a temporary variable for the reference without
// changing the meaning of the statement, for example:
//    tmp = x.f; tmp[i] = 2.
//
// By contrast, in Go, an expression such as x[i].f, where x is variable
// of type array-of-structs, has a dual meaning. When it appears in an
// ordinary expression, it means: find the variable x[i].f within x and
// retrieve its value. But when it appears on the left side of an
// assignment, as in x[i].f = 1, it means find the variable x[i].f
// within x and update its value. It cannot be decomposed into two
// operations tmp = x[i]; tmp.f = 2 without changing its meaning because
// the first operation would make a copy of the variable x[i] and the
// second would mutate the copy, not the original. It can be
// decomposed only by using pointers, for example:
//    ptr = &x[i]; ptr.f = 2.
// A Go compiler implicitly does this decomposition using pointers when
// it generates code for x[i].f on the left side of an assignment, or
// as the operand of a &-operator. This is called "l-mode" (l for left),
// opposed to r-mode code generation, in which the expression's value,
// not its address, is needed.
//
// In order to support these operations on Go variables with the usual
// semantics, in AllowAddressing mode the Starlark compiler, like a Go
// compiler, generates different code for sequences of operations such
// as x[i].f based on whether they appear on the left or right side of
// an assignment (l-mode or r-mode).
//
// In both cases the compiler generates a sequence of calls to
// Indexable.Index and HasAttrs.Attr for all but the last field/index
// operation. The sequence is then followed by a call to one of the
// following.
// 1. Variable.Value, for an r-mode expression. If the operand is not a
// Variable, it is assumed to be an ordinary Starlark value and is left
// unchanged.
// 2. SetIndexable.SetIndex, for an element update.
// 3. HasSetField.SetField, for a field update;
// 4. Variable.Address, for an explicit & operation.
//
// An concrete type that satisifies Variable could be represented by a
// non-nil Go pointer, p: Address would return p; Value would return *p;
// and SetValue(x) would execute *p=x. But the Variable is logically an
// abstraction of the variable *p, not the pointer itself.
//
type Variable interface {
	Value
	Address() Value
	Value() Value
	SetValue(Value) error
}
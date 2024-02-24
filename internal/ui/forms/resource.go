package forms

// I need to know what I am attaching
// I want to let the state do that...

// I can tell ResourceForm which I am attaching it to
// based off the route ID, that becomes the child.
// Save state will need to create that association
// by grabbing the item again from the DB and then
// adding it there. Once the form receives a "StateSavedMsg"
// we will go back.

type ResourceFormModel struct {
}

# Supporting undoable actions with the Command pattern

This plan describes the implementation of the Command pattern to support
reversable filesystem operations in `odc`. The goal is to provide a "safety
net" for destructive actions like deletions or overwrites.

Currently, once a file is deleted or moved, there is no built-in way to revert
the action without manual intervention. By encapsulating each operation as an
"Undoable Command" object, `odc` can track the necessary steps to roll back
changes if requested by the user or in the event of an error.

## Objectives

The primary goals for this refactoring include:

*   **Reversibility:** Provide a clear path for undoing file operations (e.g.,
    restoring a deleted file from a temporary "trash" location).
*   **Atomic batches:** Group multiple operations into a single logical unit
    that can be rolled back entirely if any part fails.
*   **Consistent history:** Maintain a log of executed commands that can be
    inspected or replayed.
*   **Safety:** Reduce the risk of data loss during complex batch operations
    like `rm -r` or `mv *.txt`.

## Proposed infrastructure (`internal/fs/commands`)

The infrastructure defines an interface for commands that can be executed and
undone.

```go
type UndoableCommand interface {
    Execute(ctx context.Context) error
    Undo(ctx context.Context) error
    Name() string
}
```

## Potential undoable commands

Several common operations can be refactored into undoable commands:

*   **`DeleteCommand`**: Moves a file to a temporary trash location instead of
    performing a permanent deletion.
*   **`MoveCommand`**: Records the original source path to allow for moving the
    item back to its initial location.
*   **`OverwriteCommand`**: Backs up an existing file before it is replaced by
    a new version.

## Implementation steps

1.  **Define the command interface:** Create the `UndoableCommand` interface
    and a `CommandHistory` manager in `internal/fs/commands/`.
2.  **Implement core commands:** Refactor `Delete`, `Move`, and `Copy` into
    dedicated command objects that store their revert logic.
3.  **Implement a "Trash" mechanism:** Create a temporary storage area (either
    local or remote) for items that are marked for deletion.
4.  **Update CLI handlers:** Modify command handlers to create and execute
    these command objects, recording them in the history manager.

## Next steps

After establishing the command pattern, we can introduce a new `undo` CLI
command that allows users to revert their most recent actions directly from
the terminal.

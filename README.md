# wdgo
My work day in Go

# Elements

* __board__: represents a kind of Kanban board
    * __stage__: every board can have different stages. The different stages have a fixed position on the board.
        * __card__: a card is pinned on a position on a stage. The cards can move around the board.

## Events

All events are text-based.

```
[Time]|[ID]|[Action]|[Value]|[[Action]|[Value]|...]
2020.01.24 12:13:04|28d...8ca1|Title|New Title
2020.01.24 12:13:14|28d...8ca1|AddComment|3d2...23f1
2020.01.24 12:13:15|3d2...23f1|Title|some text|Bar|some value
```

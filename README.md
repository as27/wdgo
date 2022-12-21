# About wdgo

wdgo is for organizing a work day as a simple kanban-board in the terminal. Every change inside the app is implemented as event. All events are stored inside simple text files.

I wrote this tool as a prototype for getting started with the great library [tview](https://github.com/rivo/tview). I used that tool active at work to record my time. 

At the moment I don't have really much time left to implement more features, but for my use cases it works like it is. But support or further suggestions are welcome.

Waring: I tested the tool and it works for me, but there is no warranty, if something breaks. If you are using wdgo you are using it at your own risk.

## Different elements of wdgo

The tool abstracts everything into different elements:

* __board__: represents a kind of Kanban board
    * __stage__: every board can have different stages. The different stages have a fixed position on the board.
        * __card__: a card is pinned on a position on a stage. The cards can move around the board. 
        * __session__: a session is used to record the time

## Getting started

If you have Go installed on your machine, just use `go install` for download and installing.

```
go install github.com/as27/wdgo/cmd/wdgo
```

Then create a folder, where all the data will be stored. Inside that folder create a additional subfolder `events`. 

Then start the tool. At the first start there are no boards, where you can put any cards. So go to "New Board". Give your board a name and save that name. 

After that your new boar shows up at the main menu. You can select your new board. After that you will see the window "New Stage". Press `Ctrl+e` to edit this stage. Now you can change the name of the stage. You could name it _Backlog_.

To add new stages you have to press `Ctrl+a`. So you could add a stage _Doing_ and _Done_.

To add cards to a stage press `Ctrl+n`. Now you can add a title description, support nr and a customer. When selecting "ToDo" that card will also show up inside the ToDo-List. "Archived" is used, when the card should not apear inside any stage.

Inside a card you can start and stop a session. With this feature you can track the time you are working on that topic. For fast start/stop you can use `Ctrl+j`. 

To go back you can always press `ESC`.




## Events

All events are text-based.

```
[Time]|[ID]|[Action]|[Value]|[[Action]|[Value]|...]
2020.01.24 12:13:04|28d...8ca1|Title|New Title
2020.01.24 12:13:14|28d...8ca1|AddComment|3d2...23f1
2020.01.24 12:13:15|3d2...23f1|Title|some text|Bar|some value
```

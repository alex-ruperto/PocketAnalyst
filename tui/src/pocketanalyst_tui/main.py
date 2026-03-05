from textual.app import App, ComposeResult
from textual.containers import HorizontalGroup, VerticalScroll
from textual.widgets import Button, Digits, Footer, Header

class TimeDisplay(Digits):
    """A widget to display elapsed time"""

class Stopwatch(HorizontalGroup):
    """A stopwatch widget."""
    def compose(self) -> ComposeResult:
        """Create child widgets of a stopwatch."""
        yield Button("Start", id="start", variant="success")
        yield Button("Stop", id="stop", variant="error")
        yield Button("Reset", id="reset")
        yield TimeDisplay("00:00:00.00")


class StopwatchApp(App):
    """A textual app to manage stopwatches"""

    # Bindings is a list of tuples that maps keys to actions in the app.
    # First value = key
    # Second value = name of the action
    # Third value = short description
    CSS_PATH ="stopwatch.tcss"
    BINDINGS = [
        ("d", "toggle_dark", "Toggle dark mode"),
    ]

    def compose(self) -> ComposeResult:
        """Create child widgets for the app. This is where the user interface is created with widgets."""
        # Textual will mount widgets in the order you yield them.
        # yield makes compose() a generator function.
        # Calling those functions with yield returns a generator object. Iterating the generator produces widget instances like Header and Footer one at a time.
        yield Header()
        yield Footer()
        yield VerticalScroll(Stopwatch(), Stopwatch(), Stopwatch())

    def action_toggle_dark(self) -> None:
        """An action to toggle dark mode."""
        # This is an action method. Actions are methods that begin with action_ followed by the name of the action. The BINDINGS list tells Textual to run
        # this action when the user hits the D key.
        self.theme = (
            "textual-dark" if self.theme == "textual-light" else "textual-light"
        )

def main():
    app = StopwatchApp()
    app.run()

if __name__ == "__main__":
    main()


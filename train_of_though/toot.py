#!/usr/bin/python3
# program purpose: simple amusement and fiddle around with 
# ANSI sequences, coloring and text formatting

from random import randrange
from shutil import which
import subprocess
import textwrap


def random_color():
    "returns a random foreground ANSI color with bold attribute set"
    fg = randrange(31, 39)
    return f"\x1b[1;{fg}m"


def has_fortune():
    "detects when fortune program is installed"
    return which('fortune') is not None


def get_fortune_quote():
    "run external program fortune, and formats the output to fit in 48 char width"
    result = subprocess.check_output('fortune')
    return textwrap.wrap(result.decode('utf-8'), width=48)


def insert_string_fixed_length(original_string, insert_string):
    """Inserts a string into another string, keeping the total length the same.
    Args:
      original_string: The string to insert into.
      insert_string: The string to be inserted.
    Returns:
      A new string with the inserted string, or None if the insertion is not possible.
    """
    if len(insert_string) > len(original_string):
        return None  # Cannot insert if the new string is longer
    # Calculate the middle of the string
    midpoint = len(original_string) // 2
    # Calculate the starting position for the insert string
    start_pos = 9+midpoint - len(insert_string) // 2
    # Calculate how many characters need to be removed from the original string
    chars_to_remove = len(insert_string)
    # Construct the new string
    new_string = (
        original_string[: start_pos]
        + insert_string
        + original_string[(start_pos + chars_to_remove):]
    )
    return new_string


COLOR_RESET = "\x1b[0m"

if not has_fortune():
    print("Sorry you need to install the fortune package!")
    print("$ sudo zypper install fortune")
    exit(1)

# ASCII train art from https://www.asciiart.eu/vehicles/trains
train = [
  "ToOT : THE Train OOF Thought - SUSE HackWeek 2024",
  "                 _-====-__-======-__-========-_____-============-__",
  "               _(                                                 _)",
  "            OO(                                                   )_",
  "           0  (_                                                   _)",
  "         o0     (_                                                _)",
  "        o         '=-___-===-_____-========-___________-===-___-='",
  "      .o                                _________",
  "     . ______          ______________  |         |      _____",
  "   _()_||__|| ________ |            |  |_________|   __||___||__",
  "  (BNSF 1995| |      | |            | __Y______00_| |_         _|",
  ' /-OO----OO""="OO--OO"="OO--------OO"="OO-------OO"="OO-------OO"=P',
  "#####################################################################",
  "Train art by Donovan Bake"]

# only accepts 'short' quotes that fits in the train bubble
while True:
    quotes = get_fortune_quote()
    if len(quotes) < 5:
        break

# embed the output of fortune inside the blank bubble space
for i, q in enumerate(quotes):
    train[i+2] = insert_string_fixed_length(train[i+2], q)

print(random_color())
print("\n".join(train))
print(COLOR_RESET)

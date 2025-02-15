# Bible-cli
A simple command-line interface (CLI) tool for accessing and reading the Bible, written in Go. Designed for personal use. Supports colored and plain text output.

<img alt="Bible CLI Tool" src="https://raw.githubusercontent.com/ButbkaDrug/bible/refs/heads/master/bible.gif" width="600" />


## Table of Contents

- [Introduction](#introduction)
- [Features](#features)
- [Usage](#usage)
- [Configuration](#configuration)
- [Database](#database)

## Introduction

`bible-cli` provides a quick and easy way to read the Bible from your terminal. It's designed for simple verse lookups, chapter reading, range selections, and keyword searches.  It's particularly useful for integration with text editors like Vim/Neovim for sermon preparation and study.

## Features

* **Verse Lookup:** Retrieve specific verses (e.g., `bible john 3:16`).
* **Chapter Reading:** Read entire chapters (e.g., `bible john 3`).
* **Verse Range Selection:** Specify a range of verses across chapters (e.g., `bible john 3:12-4:12`).
* **Keyword Search:** Search for phrases within the Bible (e.g., `bible search "love your neighbor"`).
* **Colored and Plain Text Output:** Choose between colored output for readability or plain text for simpler displays via environment variable.
* **Go Implementation:** Built for performance and cross-platform compatibility.
* **NVIM Integration:** Designed for seamless integration with NVIM for quick verse lookups and pasting into the editor.

## Usage
```bash

bible john 3:16     # Read John 3:16
bible john 3        # Read John chapter 3
bible john 3:12-4:12 # Read John 3:12 through John 4:12
bible love your neighbor # Search for the phrase "love your neighbor"
bible john 3:16     # Read John 3:16 (colored output by default)
BIBLE_ENV=plain bible john 3:16 # Set env variable for plain output
```

Note: Mixed requests like bible john 3:16, Luke 4:5-10 are not yet implemented.

## Configuration

**bible-cli** looks for an SQLite3 database in $HOME/.config/bible-cli/. The default database name is ESV.SQLite3.

Database Location: You can change the database directory by setting the BIBLECLI environment variable to the path containing your database files. For example:

```bash
export BIBLECLI=/path/to/my/bible/dbs
```

Database Selection: You can specify a different database file (e.g., a different translation) by setting the TRANSLATION environment variable. For example:

```bash
TRANSLATION=NIV bible john 3:16
```

Plain Text Output: Set the BIBLE_ENV environment variable to plain to force plain text output. If this variable is not set, output will be colored.

```Bash
BIBLE_ENV=plain bible 1 John 1:10 # Enable plain text output
```
## Database

The Bible databases are not distributed with the CLI. You can [Download Bible Databases here](https://www.ph4.ru/b4_1.php?l=en&q=).  Place the downloaded .SQLite3 files in the directory specified by $BIBLECLI (or the default $HOME/.config/bible-cli/).

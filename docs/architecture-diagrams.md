# Architecture Diagrams

This page contains visual diagrams that illustrate the architecture and data flows within `govman`.

## High-Level Component Diagram

This diagram shows the main components of `govman` and their relationships.

```mermaid
graph TD
    subgraph User Interaction
        CLI[CLI Layer<br>(cobra)]
    end

    subgraph Core Logic
        Manager[Manager Layer]
    end

    subgraph Services
        Config[Config<br>(viper)]
        Downloader
        Shell
        GoReleases[Go Releases API]
        Symlink
        Logger
    end

    CLI -->|Executes| Manager

    Manager -->|Uses| Config
    Manager -->|Uses| Downloader
    Manager -->|Uses| Shell
    Manager -->|Uses| GoReleases
    Manager -->|Uses| Symlink
    Manager -->|Uses| Logger

    style CLI fill:#d4e6f1,stroke:#333,stroke-width:2px
    style Manager fill:#d1f2eb,stroke:#333,stroke-width:2px
    style Services fill:#fdebd0,stroke:#333,stroke-width:1px
```

## `install` Command Sequence Diagram

This diagram illustrates the sequence of events when a user runs `govman install <version>`.

```mermaid
sequenceDiagram
    actor User
    participant CLI
    participant Manager
    participant GoReleases
    participant Downloader

    User->>CLI: govman install 1.25.1
    CLI->>Manager: Install("1.25.1")
    Manager->>GoReleases: GetDownloadURL("1.25.1")
    GoReleases-->>Manager: Returns URL & Checksum
    Manager->>Downloader: Download(url, checksum)
    Downloader-->>Manager: Success/Failure
    Manager-->>CLI: Result
    CLI-->>User: Display output
```

## `use` Command Sequence Diagram

This diagram shows the process for activating a Go version with `govman use <version> --default`.

```mermaid
sequenceDiagram
    actor User
    participant CLI
    participant Manager
    participant Config
    participant Symlink

    User->>CLI: govman use 1.25.1 --default
    CLI->>Manager: Use("1.25.1", setDefault=true)
    Manager->>Config: SetDefaultVersion("1.25.1")
    Config-->>Manager: Success
    Manager->>Symlink: CreateLink("1.25.1")
    Symlink-->>Manager: Success
    Manager-->>CLI: Result
    CLI-->>User: Display output
```

## Auto-Switching Sequence Diagram (on `cd`)

This diagram illustrates how shell integration works when changing directories.

```mermaid
sequenceDiagram
    actor User
    participant Shell
    participant govman

    User->>Shell: cd my-project
    Shell->>Shell: Hook triggered, finds .govman-version
    Shell->>govman: Executes `govman use 1.22.4`
    govman->>Shell: Returns new PATH
    Shell->>Shell: Updates current session's PATH
```
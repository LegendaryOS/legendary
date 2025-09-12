#!/usr/bin/env python3

import argparse
import subprocess
import sys
from pathlib import Path

from rich.console import Console
from rich.panel import Panel
from rich.progress import (
    Progress,
    SpinnerColumn,
    BarColumn,
    TextColumn,
    MofNCompleteColumn,
    TimeElapsedColumn,
    TimeRemainingColumn,
    DownloadColumn,
    TaskID,
)
from rich.table import Table
from rich.text import Text
from rich.prompt import Prompt
from rich.live import Live
from rich.align import Align
from rich.markdown import Markdown

# Configuration
APT_PATH = Path("/usr/lib/legendary/apt")
if not APT_PATH.exists():
    console = Console()
    console.print(Panel("[bold red]Error: Legendary APT wrapper not found at /usr/lib/legendary/apt[/bold red]", border_style="red", title="Initialization Error"))
    sys.exit(1)

console = Console()

def print_banner():
    """Display an enhanced banner with styling."""
    banner_text = Text.assemble(
        Text("Legendary Package Manager\n", style="bold magenta underline"),
        Text("For LegendaryOS (based on Ubuntu non-LTS)", style="italic dim white")
    )
    console.print(Panel(Align.center(banner_text), border_style="bold magenta", expand=False, padding=(1, 2)))

def run_apt_command(args, command_type="general"):
    """Run the apt command with advanced progress bar simulation."""
    full_cmd = [str(APT_PATH)] + args
    cmd_str = " ".join(full_cmd)
    
    progress_columns = [
        SpinnerColumn(style="bold blue"),
        TextColumn("[progress.description]{task.description}", style="cyan"),
        BarColumn(bar_width=40, style="bar.back", complete_style="bar.complete", finished_style="bar.finished"),
        MofNCompleteColumn(),
        TextColumn("[progress.percentage]{task.percentage:>3.1f}%", style="green"),
        "•",
        DownloadColumn(),
        "•",
        TimeElapsedColumn(),
        "•",
        TimeRemainingColumn(),
    ]
    
    with Progress(*progress_columns, console=console, transient=False) as progress_instance:
        main_task = progress_instance.add_task(f"[bold cyan]Executing command: {command_type}[/bold cyan]", total=100)
        
        # Simulate stages for better progress visualization
        stages = ["Preparing", "Fetching data", "Processing", "Applying changes", "Cleaning up"]
        stage_tasks = {stage: progress_instance.add_task(f"[italic dim]{stage}...", total=20, visible=False) for stage in stages}
        
        try:
            # Start subprocess
            process = subprocess.Popen(full_cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True, bufsize=1, universal_newlines=True)
            
            # Simulate progress while reading output
            output_lines = []
            error_lines = []
            current_stage = 0
            progress_instance.update(main_task, advance=0)
            progress_instance.update(stage_tasks[stages[current_stage]], visible=True)
            
            while process.poll() is None:
                line = process.stdout.readline().strip()
                if line:
                    output_lines.append(line)
                    # Simulate progress based on output or time
                    progress_instance.update(stage_tasks[stages[current_stage]], advance=1)
                    if progress_instance.tasks[stage_tasks[stages[current_stage]]].completed >= 20:
                        progress_instance.update(stage_tasks[stages[current_stage]], visible=False)
                        current_stage += 1
                        if current_stage < len(stages):
                            progress_instance.update(stage_tasks[stages[current_stage]], visible=True)
                    progress_instance.update(main_task, advance=1)
            
            # Read remaining
            output_lines.extend(process.stdout.readlines())
            error_lines.extend(process.stderr.readlines())
            
            # Complete progress
            while not progress_instance.finished:
                progress_instance.update(main_task, advance=1)
            
            if process.returncode != 0:
                raise subprocess.CalledProcessError(process.returncode, full_cmd, "\n".join(output_lines), "\n".join(error_lines))
            
            # Display output in a styled panel
            if output_lines:
                output_md = Markdown("\n".join(output_lines))
                console.print(Panel(output_md, title="[green]Command Output[/green]", border_style="green"))
            if error_lines:
                error_md = Markdown("\n".join(error_lines))
                console.print(Panel(error_md, title="[yellow]Command Warnings[/yellow]", border_style="yellow"))
        
        except subprocess.CalledProcessError as e:
            progress_instance.update(main_task, description="[bold red]Command Failed[/bold red]")
            if e.stdout:
                console.print(Panel(Markdown(e.stdout), title="[red]Output[/red]", border_style="red"))
            if e.stderr:
                console.print(Panel(Markdown(e.stderr), title="[red]Errors[/red]", border_style="red"))
            sys.exit(1)

def show_help():
    """Display enhanced help with tables and styling."""
    print_banner()
    
    # Main commands table
    main_table = Table(title="Core Commands", show_header=True, header_style="bold magenta", border_style="bright_blue")
    main_table.add_column("Command", style="bold cyan", no_wrap=True, width=20)
    main_table.add_column("Description", style="white", width=60)
    main_table.add_column("Usage Example", style="italic dim green", width=40)

    commands = [
        ("install <packages>", "Install specified packages. Consider using 'legendary isolator' for isolated environments.", "legendary install vim curl"),
        ("remove <packages>", "Remove specified packages.", "legendary remove vim"),
        ("autoremove", "Remove unused automatically installed dependencies.", "legendary autoremove"),
        ("autoclean", "Clean up retrieved package files from local repository.", "legendary autoclean"),
        ("update", "Update the package index files.", "legendary update"),
        ("upgrade", "Upgrade all upgradable packages.", "legendary upgrade"),
        ("info <package>", "Display detailed information about a package.", "legendary info vim"),
        ("help / ?", "Display this help information.", "legendary help"),
    ]

    for cmd, desc, ex in commands:
        main_table.add_row(cmd, desc, ex)

    console.print(Align.center(main_table))
    
    # Additional tips table
    tips_table = Table(title="Additional Tips", show_header=True, header_style="bold yellow", border_style="yellow")
    tips_table.add_column("Tip", style="bold yellow", width=80)

    tips = [
        "For isolated package management, use 'legendary isolator' as an alternative during installations.",
        "Always run 'update' before 'upgrade' to ensure the latest package information.",
        "Use 'info' to check package details before installing.",
    ]

    for tip in tips:
        tips_table.add_row(tip)

    console.print(Align.center(tips_table))

def show_info(package):
    """Display package info in a styled format."""
    print_banner()
    console.print(Text(f"Package Information for: {package}", style="bold underline cyan"))
    
    # Run apt info and capture
    try:
        result = subprocess.run([str(APT_PATH), "info", package], capture_output=True, text=True, check=True)
        info_output = result.stdout
        
        # Attempt to parse into table for better display
        info_table = Table(title="Package Details", show_header=False, border_style="bright_green")
        info_table.add_column("Field", style="bold magenta", width=20)
        info_table.add_column("Value", style="white", width=60)
        
        lines = info_output.splitlines()
        for line in lines:
            if ":" in line:
                field, value = line.split(":", 1)
                info_table.add_row(field.strip(), value.strip())
        
        console.print(Align.center(info_table))
        
        if result.stderr:
            console.print(Panel(Markdown(result.stderr), title="[yellow]Warnings[/yellow]", border_style="yellow"))
    
    except subprocess.CalledProcessError as e:
        console.print(Panel(f"[bold red]Error fetching info: {e}[/bold red]", border_style="red"))

def confirm_action(action, packages=None):
    """Prompt for confirmation with styled prompt."""
    if packages:
        pkg_list = ", ".join(packages)
        question = f"Confirm {action} for {pkg_list}?"
    else:
        question = f"Confirm {action}?"
    
    response = Prompt.ask(Text(question, style="bold yellow"), choices=["y", "n"], default="n")
    return response.lower() == "y"

def main():
    parser = argparse.ArgumentParser(
        description="Legendary Package Manager - Advanced CLI Tool for LegendaryOS",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="Enhanced interface with Rich library for superior terminal experience."
    )
    parser.add_argument("command", choices=["install", "remove", "autoremove", "autoclean", "update", "upgrade", "info", "help", "?"], nargs="?", default="help")
    parser.add_argument("packages", nargs="*", help="Package names for relevant commands")

    args = parser.parse_args()

    command = args.command
    packages = args.packages

    if command in ["help", "?"]:
        show_help()
        return

    if command == "info":
        if not packages:
            console.print(Panel("[bold red]Error: 'info' requires a package name.[/bold red]", border_style="red"))
            sys.exit(1)
        show_info(packages[0])  # Support single package for info
        return

    print_banner()

    if command in ["install", "remove"]:
        if not packages:
            console.print(Panel(f"[bold red]Error: '{command}' requires package names.[/bold red]", border_style="red"))
            sys.exit(1)
        if command == "install":
            console.print(Panel(
                Text("Note: For isolated environments, consider using 'legendary isolator'.", style="italic yellow"),
                title="Installation Note",
                border_style="yellow",
                padding=(1, 2)
            ))
        if not confirm_action(command, packages):
            console.print(Panel("[bold yellow]Action cancelled.[/bold yellow]", border_style="yellow"))
            return
        run_apt_command([command] + packages, command_type=command.capitalize())

    else:
        # For autoremove, autoclean, update, upgrade
        if command in ["autoremove", "autoclean"]:
            if not confirm_action(command):
                console.print(Panel("[bold yellow]Action cancelled.[/bold yellow]", border_style="yellow"))
                return
        run_apt_command([command], command_type=command.capitalize())

    # Success message
    success_msgs = {
        "install": "Packages installed successfully.",
        "remove": "Packages removed successfully.",
        "update": "Package index updated.",
        "upgrade": "System upgraded.",
        "autoremove": "Unused dependencies removed.",
        "autoclean": "Cache cleaned."
    }
    if command in success_msgs:
        console.print(Panel(Text(success_msgs[command], style="bold green"), border_style="green", title="Success"))

if __name__ == "__main__":
    main()

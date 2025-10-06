use clap::{Parser, Subcommand};
use colored::*;
use std::fs;
use std::process::{Command, Stdio};
use std::io::{self, Write};

#[derive(Parser)]
#[clap(name = "legendary", about = "A colorful CLI tool for managing LegendaryOS", version = "1.0.0")]
struct Cli {
    #[clap(subcommand)]
    command: Option<Commands>,
}

#[derive(Subcommand)]
enum Commands {
    /// Installs a package using pacman, falls back to yay if not found
    Install {
        #[clap(value_parser)]
        package: String,
    },
    /// Updates package lists (pacman -Sy)
    Update,
    /// Upgrades all packages (pacman -Syu)
    Upgrade,
    /// Removes a package (pacman -R)
    Remove {
        #[clap(value_parser)]
        package: String,
    },
    /// Rolls back to a Btrfs snapshot using snapper
    Rollback,
    /// Displays system info and ASCII art from /usr/share/ascii
    About,
    /// Shows available commands (same as running without arguments)
    Help,
}

fn run_command(program: &str, args: &[&str]) -> bool {
    let output = Command::new(program)
        .args(args)
        .stdout(Stdio::inherit())
        .stderr(Stdio::inherit())
        .output();

    match output {
        Ok(output) => output.status.success(),
        Err(e) => {
            eprintln!("{}", format!("Error running {}: {}", program, e).red());
            false
        }
    }
}

fn install_package(package: &str) {
    println!("{}", format!("Installing package: {}", package).green());
    let pacman_success = run_command("/usr/lib/LegendaryOS/pacman", &["-S", package, "--noconfirm"]);
    if !pacman_success {
        println!("{}", format!("Package {} not found in pacman repos, trying yay...", package).yellow());
        let yay_success = run_command("yay", &["-S", package, "--noconfirm"]);
        if !yay_success {
            eprintln!("{}", format!("Failed to install package {} with yay.", package).red());
        } else {
            println!("{}", format!("Package {} installed successfully with yay!", package).cyan());
        }
    } else {
        println!("{}", format!("Package {} installed successfully with pacman!", package).cyan());
    }
}

fn update_system() {
    println!("{}", "Updating package lists...".green());
    if run_command("/usr/lib/LegendaryOS/pacman", &["-Sy", "--noconfirm"]) {
        println!("{}", "Package lists updated successfully!".cyan());
    } else {
        eprintln!("{}", "Failed to update package lists.".red());
    }
}

fn upgrade_system() {
    println!("{}", "Upgrading system...".green());
    if run_command("/usr/lib/LegendaryOS/pacman", &["-Syu", "--noconfirm"]) {
        println!("{}", "System upgraded successfully!".cyan());
    } else {
        eprintln!("{}", "Failed to upgrade system.".red());
    }
}

fn remove_package(package: &str) {
    println!("{}", format!("Removing package: {}", package).green());
    if run_command("/usr/lib/LegendaryOS/pacman", &["-R", package, "--noconfirm"]) {
        println!("{}", format!("Package {} removed successfully!", package).cyan());
    } else {
        eprintln!("{}", format!("Failed to remove package {}.", package).red());
    }
}

fn rollback_snapshot() {
    println!("{}", "Rolling back to a previous snapshot...".green());
    if run_command("snapper", &["undochange", "0..1"]) {
        println!("{}", "Snapshot rollback completed successfully!".cyan());
    } else {
        eprintln!("{}", "Failed to rollback snapshot.".red());
    }
}

fn display_about() {
    println!("{}", "LegendaryOS System Information".purple().bold());
    println!("{}", "-----------------------------".purple());
    if let Ok(ascii_content) = fs::read_to_string("/usr/share/ascii") {
        println!("{}", ascii_content.blue());
    } else {
        eprintln!("{}", "Failed to read /usr/share/ascii".red());
    }
    println!("{}", "System: LegendaryOS".cyan());
    println!("{}", "Tool: legendary v1.0.0".cyan());
    println!("{}", "Description: A colorful CLI tool for managing packages and snapshots".cyan());
}

fn show_help() {
    println!("{}", "LegendaryOS CLI Tool - Available Commands".purple().bold());
    println!("{}", "---------------------------------------".purple());
    println!("{}", "help                - Show this help message".yellow());
    println!("{}", "install <package>   - Install a package (falls back to yay if not found)".yellow());
    println!("{}", "update              - Update package lists (pacman -Sy)".yellow());
    println!("{}", "upgrade             - Upgrade all packages (pacman -Syu)".yellow());
    println!("{}", "remove <package>    - Remove a package (pacman -R)".yellow());
    println!("{}", "rollback            - Rollback to a previous Btrfs snapshot".yellow());
    println!("{}", "about               - Display system info and ASCII art".yellow());
}

fn main() {
    let cli = Cli::parse();

    match cli.command {
        Some(Commands::Install { package }) => install_package(&package),
        Some(Commands::Update) => update_system(),
        Some(Commands::Upgrade) => upgrade_system(),
        Some(Commands::Remove { package }) => remove_package(&package),
        Some(Commands::Rollback) => rollback_snapshot(),
        Some(Commands::About) => display_about(),
        Some(Commands::Help) => show_help(),
        None => show_help(),
    }
}

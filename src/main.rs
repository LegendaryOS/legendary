use clap::{Parser, Subcommand};
use colored::*;
use std::fs;
use std::process::{Command, Stdio};
use std::io::{self, Write};

#[derive(Parser)]
#[clap(
    name = "legendary",
    about = "A vibrant CLI tool for managing LegendaryOS with style".bold().cyan(),
    version = "1.0.0"
)]
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
    /// Launches the graphical UI for package and snapshot management
    Ui,
    /// Searches for a package in pacman and yay repositories
    Search {
        #[clap(value_parser)]
        query: String,
    },
    /// Lists all installed packages
    List,
    /// Cleans the package cache (pacman -Sc)
    Clean,
    /// Displays system status and snapshot information
    Status,
}

fn run_command(program: &str, args: &[&str]) -> bool {
    println!("{}", format!("Executing: {} {}", program, args.join(" ")).dimmed());
    let output = Command::new(program)
        .args(args)
        .stdout(Stdio::inherit())
        .stderr(Stdio::inherit())
        .output();
    match output {
        Ok(output) => {
            if output.status.success() {
                println!("{}", format!("Command {} completed successfully!", program).green().bold());
                true
            } else {
                eprintln!("{}", format!("Command {} failed.", program).red().bold());
                false
            }
        }
        Err(e) => {
            eprintln!("{}", format!("Error running {}: {}", program, e).red().bold());
            false
        }
    }
}

fn install_package(package: &str) {
    println!("{}", format!("Installing package: {}", package).green().bold());
    let pacman_success = run_command("/usr/lib/LegendaryOS/pacman", &["-S", package, "--noconfirm"]);
    if !pacman_success {
        println!("{}", format!("Package {} not found in pacman repos, trying yay...", package).yellow().bold());
        let yay_success = run_command("yay", &["-S", package, "--noconfirm"]);
        if !yay_success {
            eprintln!("{}", format!("Failed to install package {} with yay.", package).red().bold());
        } else {
            println!("{}", format!("Package {} installed successfully with yay!", package).cyan().bold());
        }
    } else {
        println!("{}", format!("Package {} installed successfully with pacman!", package).cyan().bold());
    }
}

fn update_system() {
    println!("{}", "Updating package lists...".green().bold());
    if run_command("/usr/lib/LegendaryOS/pacman", &["-Sy", "--noconfirm"]) {
        println!("{}", "Package lists updated successfully!".cyan().bold());
    } else {
        eprintln!("{}", "Failed to update package lists.".red().bold());
    }
}

fn upgrade_system() {
    println!("{}", "Upgrading system packages...".green().bold());
    if run_command("/usr/lib/LegendaryOS/pacman", &["-Syu", "--noconfirm"]) {
        println!("{}", "System upgraded successfully!".cyan().bold());
    } else {
        eprintln!("{}", "Failed to upgrade system.".red().bold());
    }
}

fn remove_package(package: &str) {
    println!("{}", format!("Removing package: {}", package).green().bold());
    if run_command("/usr/lib/LegendaryOS/pacman", &["-R", package, "--noconfirm"]) {
        println!("{}", format!("Package {} removed successfully!", package).cyan().bold());
    } else {
        eprintln!("{}", format!("Failed to remove package {}.", package).red().bold());
    }
}

fn rollback_snapshot() {
    println!("{}", "Initiating rollback to previous snapshot...".green().bold());
    if run_command("snapper", &["undochange", "0..1"]) {
        println!("{}", "Snapshot rollback completed successfully!".cyan().bold());
    } else {
        eprintln!("{}", "Failed to rollback snapshot.".red().bold());
    }
}

fn display_about() {
    println!("{}", "LegendaryOS System Information".purple().bold());
    println!("{}", "-----------------------------".purple().underline());
    if let Ok(ascii_content) = fs::read_to_string("/usr/share/ascii") {
        println!("{}", ascii_content.blue().bold());
    } else {
        eprintln!("{}", "Failed to read /usr/share/ascii".red().bold());
    }
    println!("{}", "System: LegendaryOS".cyan().bold());
    println!("{}", "Tool: legendary v1.0.0".cyan().bold());
    println!("{}", "Description: A vibrant CLI tool for managing packages and snapshots".cyan().bold());
    println!("{}", "Developed by: LegendaryOS Team".magenta().bold());
}

fn show_help() {
    println!("{}", "Legendary CLI Tool - Available Commands".purple().bold());
    println!("{}", "---------------------------------------".purple().underline());
    println!("{}", "help           - Show this help message".yellow().bold());
    println!("{}", "install <pkg>  - Install a package (falls back to yay)".yellow().bold());
    println!("{}", "update         - Update package lists (pacman -Sy)".yellow().bold());
    println!("{}", "upgrade        - Upgrade all packages (pacman -Syu)".yellow().bold());
    println!("{}", "remove <pkg>   - Remove a package (pacman -R)".yellow().bold());
    println!("{}", "rollback       - Rollback to a previous Btrfs snapshot".yellow().bold());
    println!("{}", "about          - Display system info and ASCII art".yellow().bold());
    println!("{}", "ui             - Launch graphical UI for package management".yellow().bold());
    println!("{}", "search <query> - Search for packages in repositories".yellow().bold());
    println!("{}", "list           - List all installed packages".yellow().bold());
    println!("{}", "clean          - Clean package cache (pacman -Sc)".yellow().bold());
    println!("{}", "status         - Show system and snapshot status".yellow().bold());
}

fn launch_ui() {
    println!("{}", "Launching LegendaryOS graphical interface...".green().bold());
    if run_command("legendary-ui", &[]) {
        println!("{}", "Graphical UI launched successfully!".cyan().bold());
    } else {
        eprintln!("{}", "Failed to launch graphical UI. Is legendary-ui installed?".red().bold());
    }
}

fn search_package(query: &str) {
    println!("{}", format!("Searching for package: {}", query).green().bold());
    let pacman_success = run_command("/usr/lib/LegendaryOS/pacman", &["-Ss", query]);
    if pacman_success {
        println!("{}", "Search completed in pacman repositories!".cyan().bold());
    } else {
        println!("{}", "No results in pacman repos, trying yay...".yellow().bold());
        if run_command("yay", &["-Ss", query]) {
            println!("{}", "Search completed in yay repositories!".cyan().bold());
        } else {
            eprintln!("{}", "Failed to search for packages.".red().bold());
        }
    }
}

fn list_packages() {
    println!("{}", "Listing all installed packages...".green().bold());
    if run_command("/usr/lib/LegendaryOS/pacman", &["-Q"]) {
        println!("{}", "Installed packages listed successfully!".cyan().bold());
    } else {
        eprintln!("{}", "Failed to list installed packages.".red().bold());
    }
}

fn clean_cache() {
    println!("{}", "Cleaning package cache...".green().bold());
    if run_command("/usr/lib/LegendaryOS/pacman", &["-Sc", "--noconfirm"]) {
        println!("{}", "Package cache cleaned successfully!".cyan().bold());
    } else {
        eprintln!("{}", "Failed to clean package cache.".red().bold());
    }
}

fn display_status() {
    println!("{}", "LegendaryOS System Status".purple().bold());
    println!("{}", "------------------------".purple().underline());
    println!("{}", "Checking system status...".green().bold());
    if run_command("/usr/lib/LegendaryOS/pacman", &["-Qdtq"]) {
        println!("{}", "No orphaned packages found.".cyan().bold());
    } else {
        println!("{}", "Orphaned packages detected.".yellow().bold());
    }
    println!("{}", "Checking snapshot status...".green().bold());
    if run_command("snapper", &["list"]) {
        println!("{}", "Snapshots listed successfully!".cyan().bold());
    } else {
        eprintln!("{}", "Failed to list snapshots.".red().bold());
    }
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
        Some(Commands::Ui) => launch_ui(),
        Some(Commands::Search { query }) => search_package(&query),
        Some(Commands::List) => list_packages(),
        Some(Commands::Clean) => clean_cache(),
        Some(Commands::Status) => display_status(),
        None => show_help(),
    }
}

#!/usr/bin/env node

const { execSync } = require('child_process');
const fs = require('fs');
const path = require('path');
const chalk = require('chalk');

const PACMAN_PATH = '/usr/lib/LegendaryOS/pacman';
const UPGRADE_SCRIPT = '/usr/share/LegendaryOS/SCRIPTS/legendaryos-upgrade';

function runCommand(command, options = {}) {
    try {
        return execSync(command, { stdio: 'inherit', ...options }).toString().trim();
    } catch (error) {
        if (options.silent) {
            return '';
        }
        console.error(chalk.red(`Error executing command: ${command}`));
        process.exit(1);
    }
}

function packageExistsInPacman(pkg) {
    const output = runCommand(`${PACMAN_PATH} -Ss ${pkg}`, { silent: true });
    return output !== '';
}

function packageExistsInYay(pkg) {
    try {
        const output = runCommand(`yay -Ss ${pkg}`, { silent: true });
        return output !== '';
    } catch (error) {
        return false; // If yay is not installed, skip
    }
}

function packageExistsInFlatpak(pkg) {
    try {
        const output = runCommand(`flatpak search ${pkg}`, { silent: true });
        return output !== '';
    } catch (error) {
        return false; // If flatpak is not installed, skip
    }
}

function installPackage(pkg) {
    if (packageExistsInPacman(pkg)) {
        console.log(chalk.blue(`Installing ${pkg} via pacman...`));
        runCommand(`sudo ${PACMAN_PATH} -S ${pkg} --noconfirm`);
    } else if (packageExistsInYay(pkg)) {
        console.log(chalk.blue(`Installing ${pkg} via yay...`));
        runCommand(`yay -S ${pkg} --noconfirm`);
    } else if (packageExistsInFlatpak(pkg)) {
        console.log(chalk.blue(`Installing ${pkg} via flatpak...`));
        runCommand(`flatpak install ${pkg} -y`);
    } else {
        console.error(chalk.red(`Package ${pkg} not found in pacman, yay, or flatpak.`));
        process.exit(1);
    }
    console.log(chalk.green(`Package ${pkg} installed successfully.`));
}

function removePackage(pkg) {
    console.log(chalk.blue(`Removing ${pkg}...`));
    runCommand(`sudo ${PACMAN_PATH} -R ${pkg} --noconfirm`);
    console.log(chalk.green(`Package ${pkg} removed successfully.`));
}

function update() {
    console.log(chalk.blue('Running system update...'));
    runCommand(`sudo ${PACMAN_PATH} -Syu --noconfirm`);
    console.log(chalk.green('System updated successfully.'));
}

function upgrade() {
    if (!fs.existsSync(UPGRADE_SCRIPT)) {
        console.error(chalk.red(`Upgrade script not found at ${UPGRADE_SCRIPT}`));
        process.exit(1);
    }
    console.log(chalk.blue('Running LegendaryOS upgrade script...'));
    runCommand(`bash ${UPGRADE_SCRIPT}`);
    console.log(chalk.green('Upgrade completed.'));
}

function refresh() {
    console.log(chalk.blue('Refreshing package databases...'));
    runCommand(`sudo ${PACMAN_PATH} -Sy`);
    console.log(chalk.green('Packages refreshed.'));
}

function listInstalled() {
    console.log(chalk.blue('Listing installed packages:'));
    runCommand(`${PACMAN_PATH} -Q`);
}

function showVersion() {
    console.log(chalk.yellow('Legendary CLI Tool'));
    console.log(chalk.yellow('Version: 1.0.0'));
    console.log(chalk.yellow('For LegendaryOS based on Arch Linux'));
}

function showHelp() {
    console.log(chalk.bold.magenta(`
Legendary CLI Tool for LegendaryOS
`));
    console.log(chalk.cyan('Available commands:'));
    console.log(chalk.white('  install <package>   - Install a package (searches pacman, then yay, then flatpak)'));
    console.log(chalk.white('  remove <package>    - Remove a package'));
    console.log(chalk.white('  update              - Update system (sudo pacman -Syu)'));
    console.log(chalk.white('  upgrade             - Run LegendaryOS upgrade script'));
    console.log(chalk.white('  refresh             - Refresh package databases (sudo pacman -Sy)'));
    console.log(chalk.white('  list                - List installed packages'));
    console.log(chalk.white('  version             - Show version information'));
    console.log(chalk.white('  help                - Show this help message'));
    console.log(chalk.white('  info                - Show information about LegendaryOS'));
    console.log(chalk.gray('\nNote: Multiple packages can be specified for install/remove, separated by spaces.'));
}

function showInfo() {
    console.log(chalk.bold.magenta(`
LegendaryOS Information:
`));
    console.log(chalk.cyan('- Based on Arch Linux'));
    console.log(chalk.cyan('- Custom package manager integration with pacman, yay, and flatpak'));
    console.log(chalk.cyan('- Designed for efficient system management'));
    console.log(chalk.cyan('- Version: 1.0 (initial release)'));
    console.log(chalk.cyan('- Developed to enhance user experience on LegendaryOS'));
}

const args = process.argv.slice(2);
if (args.length === 0) {
    showHelp();
    process.exit(0);
}

const command = args[0];
const params = args.slice(1);

switch (command) {
    case 'install':
        if (params.length === 0) {
            console.error(chalk.red('Please provide at least one package name to install.'));
            process.exit(1);
        }
        params.forEach(installPackage);
        break;
    case 'remove':
        if (params.length === 0) {
            console.error(chalk.red('Please provide at least one package name to remove.'));
            process.exit(1);
        }
        params.forEach(removePackage);
        break;
    case 'update':
        update();
        break;
    case 'upgrade':
        upgrade();
        break;
    case 'refresh':
        refresh();
        break;
    case 'list':
        listInstalled();
        break;
    case 'version':
        showVersion();
        break;
    case 'help':
        showHelp();
        break;
    case 'info':
        showInfo();
        break;
    default:
        console.error(chalk.red(`Unknown command: ${command}`));
        showHelp();
        process.exit(1);
}

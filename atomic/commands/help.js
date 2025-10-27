const chalk = require('chalk');

module.exports = () => {
  console.log(chalk.blue('Legendary CLI Commands:'));
  console.log(chalk.yellow('install <package>') + ' - Install a package transactionally');
  console.log(chalk.yellow('remove <package>') + ' - Remove a package transactionally');
  console.log(chalk.yellow('update') + ' - Update all packages transactionally');
  console.log(chalk.yellow('upgrade') + ' - Run LegendaryOS upgrade script');
  console.log(chalk.yellow('clean') + ' - Clean old snapshots and pacman cache');
  console.log(chalk.yellow('info') + ' - Display current system info');
  console.log(chalk.yellow('help') + ' - Display this help');
  console.log(chalk.green('\nAdditional features:'));
  console.log(' - Systemd timer for auto-updates can be set up separately.');
  console.log(' - For rollback, use TUI: fzf to select snapshot, then set as default with btrfs.');
  console.log(' - Automatic rollback on boot failure via bootloader watchdog (configure in grub).');
};

#!/usr/bin/env node

const { Command } = require('commander');
const chalk = require('chalk');
const program = new Command();

const install = require('./commands/install');
const remove = require('./commands/remove');
const update = require('./commands/update');
const upgrade = require('./commands/upgrade');
const clean = require('./commands/clean');
const info = require('./commands/info');
const helpCmd = require('./commands/help');

program.name('legendary').description('Transactional package manager using Btrfs').version('1.0.0');

program.command('install')
  .description('Install a package transactionally')
  .argument('<package>', 'Package to install')
  .action(install);

program.command('remove')
  .description('Remove a package transactionally')
  .argument('<package>', 'Package to remove')
  .action(remove);

program.command('update')
  .description('Update all packages transactionally')
  .action(update);

program.command('upgrade')
  .description('Run the LegendaryOS upgrade script')
  .action(upgrade);

program.command('clean')
  .description('Clean old snapshots and pacman cache')
  .action(clean);

program.command('info')
  .description('Display current system info')
  .action(info);

program.command('help')
  .description('Display help')
  .action(helpCmd);

program.parse(process.argv);

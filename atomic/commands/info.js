const { execSync } = require('child_process');
const fs = require('fs');
const path = require('path');
const chalk = require('chalk');

module.exports = () => {
  try {
    // Current active snapshot
    const current = execSync('btrfs subvolume get-default /').toString().trim();
    console.log(chalk.blue(`Current active snapshot ID: ${current}`));

    // List snapshots
    const snapBase = '/var/lib/legendary/snapshots/';
    const snapshots = fs.readdirSync(snapBase).filter(dir => fs.statSync(path.join(snapBase, dir)).isDirectory());
    console.log(chalk.blue(`Number of snapshots: ${snapshots.length}`));
    snapshots.sort().forEach(snap => {
      console.log(chalk.yellow(`Snapshot: ${snap}`));
    });

    // Last upgrade - assume from last snapshot date
    if (snapshots.length > 0) {
      const lastSnap = snapshots[snapshots.length - 1];
      console.log(chalk.blue(`Last upgrade/snapshot: ${lastSnap}`));
    }
  } catch (error) {
    console.error(chalk.red(`Error retrieving info: ${error.message}`));
  }
};

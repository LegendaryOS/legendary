const { execSync } = require('child_process');
const fs = require('fs');
const path = require('path');
const chalk = require('chalk');

const PACMAN_PATH = '/usr/lib/LegendaryOS/pacman';

module.exports = () => {
  console.log(chalk.green('Cleaning old snapshots and pacman cache'));

  try {
    // Clean pacman cache
    execSync(`${PACMAN_PATH} -Scc --noconfirm`);

    // Clean old snapshots > 30 days
    const snapBase = '/var/lib/legendary/snapshots/';
    const snapshots = fs.readdirSync(snapBase).filter(dir => fs.statSync(path.join(snapBase, dir)).isDirectory());
    const thirtyDaysAgo = new Date(Date.now() - 30 * 24 * 60 * 60 * 1000);

    snapshots.forEach(snap => {
      const snapDate = new Date(snap.substring(0, 4), snap.substring(4, 6) - 1, snap.substring(6, 8), snap.substring(8, 10), snap.substring(10, 12), snap.substring(12, 14));
      if (snapDate < thirtyDaysAgo) {
        const snapRoot = path.join(snapBase, snap, '@');
        execSync(`btrfs subvolume delete ${snapRoot}`);
        fs.rmdirSync(path.join(snapBase, snap));
        console.log(chalk.yellow(`Deleted old snapshot: ${snap}`));
      }
    });

    console.log(chalk.green('Cleaning complete.'));
  } catch (error) {
    console.error(chalk.red(`Error during clean: ${error.message}`));
  }
};

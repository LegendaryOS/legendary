const { execSync } = require('child_process');
const fs = require('fs');
const path = require('path');
const chalk = require('chalk');
const { createSnapshot, prepareChroot, runInChroot, cleanupChroot, deploySnapshot } = require('../utils/btrfs');

const PACMAN_PATH = '/usr/lib/LegendaryOS/pacman';

module.exports = (pkg) => {
  console.log(chalk.green(`Installing package: ${pkg}`));

  const timestamp = new Date().toISOString().replace(/[:.-]/g, '');
  const snapDir = `/var/lib/legendary/snapshots/${timestamp}`;
  const snapRoot = path.join(snapDir, '@');

  try {
    fs.mkdirSync(snapDir, { recursive: true });
    createSnapshot(snapRoot);
    prepareChroot(snapRoot);
    runInChroot(snapRoot, `${PACMAN_PATH} -S --noconfirm ${pkg}`);
    cleanupChroot(snapRoot);
    deploySnapshot(snapRoot);
    console.log(chalk.green('Installation complete. Reboot to apply changes.'));
  } catch (error) {
    console.error(chalk.red(`Error during install: ${error.message}`));
    // Cleanup failed snapshot
    execSync(`btrfs subvolume delete ${snapRoot}`);
    fs.rmdirSync(snapDir);
  }
};
